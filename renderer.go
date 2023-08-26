package startprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/terminalcode"
	"github.com/yetsing/startprompt/token"
)

// cCompletionMenu 辅助补全菜单的渲染
type cCompletionMenu struct {
	screen        *Screen
	completeState *cCompletionState
	maxHeight     int
}

func newCompleteMenu(screen *Screen, completeState *cCompletionState, maxHeight int) *cCompletionMenu {
	return &cCompletionMenu{
		screen:        screen,
		completeState: completeState,
		maxHeight:     maxHeight,
	}
}

// 返回光标的位置坐标
func (c *cCompletionMenu) getOrigin() Coordinate {
	return c.screen.getCursorCoordinate(
		c.completeState.originalDocument.CursorPositionRow(),
		c.completeState.originalDocument.CursorPositionCol())
}

// getDrawCoordinate 返回菜单渲染位置坐标（因为是从左上角开始，所以这个就是左上角的坐标）
func (c *cCompletionMenu) getDrawCoordinate(menuWidth int) Coordinate {
	coordinate := c.getOrigin()
	x := coordinate.X
	y := coordinate.Y
	y++
	//    这里 x - 1 是因为前面会加个空格
	x = maxInt(0, x-1)
	//    补全项前后总共有 4 个空格
	if x+menuWidth+4 > c.screen.Width() {
		x -= (x + menuWidth + 4) - c.screen.Width() + 1
	}
	return Coordinate{
		X: x,
		Y: y,
	}
}

// getMenuWidth 计算补全菜单的宽度
func (c *cCompletionMenu) getMenuWidth() int {
	menuWidth := 0
	for _, completion := range c.completeState.currentCompletions {
		w := runewidth.StringWidth(completion.Display)
		if w > menuWidth {
			menuWidth = w
		}
	}
	return menuWidth
}

// 将菜单写入 screen 里面
func (c *cCompletionMenu) write() {
	completions := c.completeState.currentCompletions
	index := c.completeState.completeIndex

	//    计算补全菜单的宽度
	menuWidth := c.getMenuWidth()

	//    决定从哪个补全项开始展示
	sliceFrom := 0
	//    补全项多于最大高度并且当前选择项在下半部分位置，需要向上移动补全菜单
	//    尽可能地让选中的补全项位于菜单中上部分
	if len(completions) > c.maxHeight && index != -1 && index > c.maxHeight/2 {
		sliceFrom = minInt(
			index-c.maxHeight/2,          // 将选择项移到中间位置
			len(completions)-c.maxHeight, // 最后一个补全在最底部
		)
	}

	sliceTo := minInt(sliceFrom+c.maxHeight, len(completions))

	//    获取菜单的位置坐标
	coordinate := c.getDrawCoordinate(menuWidth)
	//    写入补全到 screen
	for i, completion := range completions[sliceFrom:sliceTo] {
		var tokenType token.TokenType
		var button token.Token
		//    判断补全项是否已选中
		if i+sliceFrom == index {
			tokenType = token.CompletionMenuCurrentCompletion
			button = token.Token{
				Type:    token.CompletionProgressButton,
				Literal: " ",
			}
		} else {
			tokenType = token.CompletionMenuCompletion
			button = token.Token{
				Type:    token.CompletionProgressBar,
				Literal: " ",
			}
		}

		c.screen.WriteTokensAtPos(coordinate.X, coordinate.Y+i, []token.Token{
			{
				Type:    token.Unspecific,
				Literal: " ",
			},
			{
				Type:    tokenType,
				Literal: fmt.Sprintf(" %s", ljustWidth(completion.Display, menuWidth)),
			},
			button,
			{
				Type:    token.Unspecific,
				Literal: " ",
			},
		})

	}
}

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
	//    光标在文本的行，用于将光标移动到文本第一行
	cursorRow int
}

type _Size struct {
	width  int
	height int
}

func (r *rRenderer) getSize() _Size {
	width, height := getSize(int(os.Stdin.Fd()))
	return _Size{
		width:  width,
		height: height,
	}
}

func (r *rRenderer) getNewScreen(renderContext *RenderContext) *Screen {
	screen := NewScreen(r.schema, r.getSize())

	//    写入提示符
	prompts := renderContext.prompt.GetPrompt()
	screen.WriteTokens(prompts, false)

	//    设置后续行前缀函数
	screen.setSecondLinePrefix(func() []token.Token {
		return renderContext.prompt.GetSecondLinePrefix()
	})

	//    写入分词后的用户输入
	screen.WriteTokens(renderContext.code.GetTokens(), true)
	screen.saveInputPos()

	screen.setSecondLinePrefix(nil)

	//    写入补全菜单
	if renderContext.completeState != nil {
		newCompleteMenu(screen, renderContext.completeState, 7).write()
	}

	return screen
}

func (r *rRenderer) renderToStr(renderContext *RenderContext, abort bool, accept bool) string {
	var buf bytes.Buffer

	//    移动光标到输入的左上方
	if r.cursorRow > 0 {
		buf.WriteString(terminalcode.CursorUp(r.cursorRow))
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
		r.cursorRow = 0
		buf.WriteString(terminalcode.CRLF)
	} else {
		// 移动光标到正确位置
		cursorCoordinate := screen.getCursorCoordinate(
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
		r.cursorRow = cursorCoordinate.Y
	}
	return buf.String()
}

func (r *rRenderer) render(renderContext *RenderContext, abort bool, accept bool) {
	out := r.renderToStr(renderContext, abort, accept)
	r.write(out)
	r.flush()
}

// renderCompletions 将补全选项一行行打印出来
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

	r.cursorRow = 0
}

// inColumns 将词语按行自适应排列， marginLeft 左边空格数量
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

// clear 清空屏幕，移动到左上角
func (r *rRenderer) clear() {
	r.write(terminalcode.EraseScreen)
	r.write(terminalcode.CursorGoto(0, 0))
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
