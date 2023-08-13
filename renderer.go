package startprompt

import (
	"bufio"
	"bytes"
	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/terminalcode"
	"github.com/yetsing/startprompt/token"
	"os"
)

func newRender(schema Schema) *rRenderer {
	return &rRenderer{
		writer:    bufio.NewWriter(os.Stdout),
		schema:    schema,
		cursorRow: 0,
	}
}

type rRenderer struct {
	writer *bufio.Writer
	schema Schema
	// 光标在文本的行，用于将光标移动到文本第一行
	cursorRow int
}

func (r *rRenderer) getWidth() int {
	_, width := getSize(int(os.Stdout.Fd()))
	return width
}

func (r *rRenderer) getNewScreen(renderContext *RenderContext) *Screen {
	screen := NewScreen(r.schema, r.getWidth())

	// write prompt
	prompts := renderContext.prompt.GetPrompt()
	screen.WriteTokens(prompts, false)

	// set second line prefix
	screen.setSecondLinePrefix(func() []token.Token {
		return renderContext.prompt.GetSecondLinePrefix()
	})

	// write code object
	screen.WriteTokens(renderContext.code.GetTokens(), true)
	screen.saveInputPos()

	return screen
}

func (r *rRenderer) renderToStr(renderContext *RenderContext) string {
	var buf bytes.Buffer

	// 移动光标到输入的左上方
	if r.cursorRow > 0 {
		buf.WriteString(terminalcode.CursorUp(r.cursorRow))
	}
	buf.WriteString(terminalcode.CarriageReturn)
	// 删除当行到屏幕下方
	buf.WriteString(terminalcode.EraseDown)

	// 生成屏幕输出
	screen := r.getNewScreen(renderContext)
	o, lastCoordinate := screen.Output()
	buf.WriteString(o)

	// 用户输入完毕或者放弃输入，另起一行
	if renderContext.accept || renderContext.abort {
		buf.WriteString(terminalcode.CRLF)
	} else {
		// 移动光标到正确位置
		cursorCoordinate := screen.getCursorCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		if lastCoordinate.Y > cursorCoordinate.Y {
			buf.WriteString(terminalcode.CursorUp(lastCoordinate.Y - cursorCoordinate.Y))
		}
		if lastCoordinate.X > cursorCoordinate.X {
			buf.WriteString(terminalcode.CursorBackward(lastCoordinate.X - cursorCoordinate.X))
		}
		if lastCoordinate.X < cursorCoordinate.X {
			buf.WriteString(terminalcode.CursorForward(cursorCoordinate.X - lastCoordinate.X))
		}
		r.cursorRow = cursorCoordinate.Y
	}
	return buf.String()
}

func (r *rRenderer) render(renderContext *RenderContext) {
	out := r.renderToStr(renderContext)
	r.write(out)
	r.flush()
}

// 将补全选项一行行打印出来
func (r *rRenderer) renderCompletions(completions []*Completion) {
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
}

// marginLeft 左边空格数量
func (r *rRenderer) inColumns(items []string, marginLeft int) []string {
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
	termWidth := r.getWidth() - marginLeft
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

func (r *rRenderer) write(s string) {
	_, err := r.writer.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (r *rRenderer) flush() {
	err := r.writer.Flush()
	if err != nil {
		panic(err)
	}
}
