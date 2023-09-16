package startprompt

import (
	"bufio"
	"bytes"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/terminalcode"
	"github.com/yetsing/startprompt/token"
)

func newRenderer(schema Schema, promptFactory PromptFactory) *Renderer {
	return &Renderer{
		writer:        bufio.NewWriter(os.Stdout),
		schema:        schema,
		promptFactory: promptFactory,
	}
}

type Renderer struct {
	writer *bufio.Writer
	schema Schema
	//    光标在输入文本中的坐标（这是一个相对于输入文本左上角的坐标）
	cursorCoordinate Coordinate
	promptFactory    PromptFactory
}

type _Size struct {
	width  int
	height int
}

func (r *Renderer) getSize() _Size {
	width, height := getSize(int(os.Stdin.Fd()))
	return _Size{
		width:  width,
		height: height,
	}
}

func (r *Renderer) getNewScreen(renderContext *RenderContext) *Screen {
	screen := NewScreen(r.schema, r.getSize())

	//    写入提示符
	prompt := r.promptFactory(renderContext.code)
	prompts := prompt.GetPrompt()
	screen.WriteTokens(prompts, false)
	//    设置后续行前缀函数
	screen.setSecondLinePrefix(func() []token.Token {
		return prompt.GetSecondLinePrefix()
	})

	//    写入分词后的用户输入
	screen.WriteTokens(renderContext.code.GetTokens(), true)
	screen.saveInputPos()

	screen.setSecondLinePrefix(nil)

	//    写入补全菜单
	if renderContext.completeState != nil {
		newCompletionMenu(screen, renderContext.completeState, 7).write()
	}

	return screen
}

func (r *Renderer) renderToStr(renderContext *RenderContext, abort bool, accept bool) string {
	var buf bytes.Buffer

	//    移动光标到输入的左上方
	if r.cursorCoordinate.Y > 0 {
		buf.WriteString(terminalcode.CursorUp(r.cursorCoordinate.Y))
	}
	buf.WriteString(terminalcode.CarriageReturn)
	//    删除当前行到屏幕下方
	buf.WriteString(terminalcode.EraseDown)

	//    生成屏幕输出
	screen := r.getNewScreen(renderContext)
	o, lastCoordinate := screen.Output()
	buf.WriteString(o)

	//    用户输入完毕或者放弃输入或者退出，另起一行
	if accept || abort {
		r.cursorCoordinate = Coordinate{0, 0}
		buf.WriteString(terminalcode.CRLF)
	} else {
		// 移动光标到正确位置
		cursorCoordinate := screen.getCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		if lastCoordinate.Y > cursorCoordinate.Y {
			buf.WriteString(terminalcode.CursorUp(lastCoordinate.Y - cursorCoordinate.Y))
		}
		// 当光标的坐标刚好是终端宽度时，这个时候用偏移量计算会有 1 的偏差
		if lastCoordinate.X >= r.getSize().width {
			buf.WriteString(terminalcode.CarriageReturn)
			buf.WriteString(terminalcode.CursorForward(cursorCoordinate.X))
		} else if lastCoordinate.X > cursorCoordinate.X {
			buf.WriteString(terminalcode.CursorBackward(lastCoordinate.X - cursorCoordinate.X))
		} else if lastCoordinate.X < cursorCoordinate.X {
			buf.WriteString(terminalcode.CursorForward(cursorCoordinate.X - lastCoordinate.X))
		}
		r.cursorCoordinate = cursorCoordinate
	}
	return buf.String()
}

func (r *Renderer) render(renderContext *RenderContext, abort bool, accept bool) {
	out := r.renderToStr(renderContext, abort, accept)
	r.write(out)
	r.flush()
}

// renderCompletions 将补全选项一行行打印出来
func (r *Renderer) renderCompletions(completions []*Completion) {
	r.write(terminalcode.CRLF)
	items := make([]string, len(completions))
	for i, completion := range completions {
		items[i] = completion.Display
	}

	for _, line := range r.inColumns(items, 0) {
		r.write(line)
		r.write(terminalcode.CRLF)
	}
	r.flush()

	r.cursorCoordinate = Coordinate{0, 0}
}

// inColumns 将词语按行自适应排列， marginLeft 左边空格数量
func (r *Renderer) inColumns(items []string, marginLeft int) []string {
	// 计算最宽的选项
	maxWidth := 0
	for _, item := range items {
		w := runewidth.StringWidth(item)
		if w > maxWidth {
			// 需要一个空格作为分割
			maxWidth = w + 1
		}
	}

	// 每行打印几个单词
	termWidth := r.getSize().width - marginLeft
	wordsPerLine := termWidth / maxWidth
	if wordsPerLine == 0 {
		wordsPerLine = 1
	}

	var lines []string
	margin := repeatByte(' ', marginLeft)
	var buf bytes.Buffer
	buf.WriteString(margin)
	for i, item := range items {
		buf.WriteString(item)

		// 到达这行最后一个单词
		if (i+1)%wordsPerLine == 0 {
			lines = append(lines, buf.String())
			buf.Reset()
		} else {
			// 加上单词之间的空格
			buf.WriteString(repeatByte(' ', maxWidth-runewidth.StringWidth(item)))
		}
	}
	if buf.Len() > 0 {
		lines = append(lines, buf.String())
	}
	return lines
}

// erase 清空当前输出，移动光标到第一行
func (r *Renderer) erase() {
	r.write(terminalcode.CursorBackward(r.cursorCoordinate.X))
	r.write(terminalcode.CursorUp(r.cursorCoordinate.Y))
	r.write(terminalcode.EraseDown)
	r.write(terminalcode.ResetAttributes)
	r.flush()
	r.reset()
}

// Clear 清空屏幕，移动到左上角
func (r *Renderer) Clear() {
	r.write(terminalcode.EraseScreen)
	r.write(terminalcode.CursorGoto(0, 0))
	r.flush()
}

func (r *Renderer) write(s string) {
	_, err := r.writer.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) flush() {
	err := r.writer.Flush()
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) reset() {

}
