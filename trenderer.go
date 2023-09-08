package startprompt

import (
	"github.com/gdamore/tcell/v2"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type TRenderer struct {
	tscreen       tcell.Screen
	schema        Schema
	promptFactory PromptFactory

	//    保存至今为止全部的输出
	//    {y: {x: Char}}
	totalBuffer map[int]map[int]*Char
	//    当前输入在 totalBuffer 的坐标（输入左上角）
	bufferCoordinate Coordinate
	//    渲染 totalBuffer 中 >= bufferOffsetY 的内容
	bufferOffsetY int
	//    当前输入在窗口中的坐标（输入左上角）
	cursorCoordinate Coordinate
	//    在窗口中显示的光标坐标
	showCursorCoordinate Coordinate
}

func newTRenderer(tscreen tcell.Screen, schema Schema, promptFactory PromptFactory) *TRenderer {
	return &TRenderer{
		totalBuffer:   map[int]map[int]*Char{},
		tscreen:       tscreen,
		schema:        schema,
		promptFactory: promptFactory,
	}
}

func (tr *TRenderer) getSize() _Size {
	width, height := tr.tscreen.Size()
	return _Size{
		width:  width,
		height: height,
	}
}

func (tr *TRenderer) getNewScreen(renderContext *RenderContext) *Screen {
	screen := NewScreen(tr.schema, tr.getSize())

	//    写入提示符
	prompt := tr.promptFactory(renderContext.code)
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

func (tr *TRenderer) updateWithScreen(screen *Screen) {
	buffer := screen.GetBuffer()
	for iy, icolumn := range buffer {
		y := tr.bufferCoordinate.Y + iy
		lineData := make(map[int]*Char, len(icolumn))
		tr.totalBuffer[y] = lineData
		for ix, char := range icolumn {
			x := tr.bufferCoordinate.X + ix
			lineData[x] = char
		}
	}
}

func (tr *TRenderer) render(renderContext *RenderContext, abort bool, accept bool) {
	//    写入屏幕输出
	screen := tr.getNewScreen(renderContext)
	tr.updateWithScreen(screen)
	//    用户输入完毕或者放弃输入或者退出，另起一行
	if accept || abort {
		tr.bufferCoordinate.X = 0
		tr.bufferCoordinate.Y += screen.maxCursorCoordinate.Y + 1
		tr.cursorCoordinate.X = 0
		tr.cursorCoordinate.Y += screen.maxCursorCoordinate.Y + 1
		tr.showCursorCoordinate = tr.cursorCoordinate
	} else {
		//    移动光标到正确位置
		cursorCoordinate := screen.getCursorCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		tr.showCursorCoordinate.X = tr.cursorCoordinate.X + cursorCoordinate.X
		tr.showCursorCoordinate.Y = tr.cursorCoordinate.Y + cursorCoordinate.Y
	}
	tr.Show()
}

func (tr *TRenderer) renderOutput(output string) {
	if len(output) == 0 {
		return
	}
	screen := NewScreen(tr.schema, tr.getSize())
	tk := token.NewToken(token.Text, output)
	screen.WriteTokens([]token.Token{tk}, false)
	tr.updateWithScreen(screen)
	tr.bufferCoordinate.X = 0
	tr.bufferCoordinate.Y += screen.maxCursorCoordinate.Y
	tr.cursorCoordinate.X = 0
	tr.cursorCoordinate.Y += screen.maxCursorCoordinate.Y
	tr.showCursorCoordinate = tr.cursorCoordinate
	tr.Show()
}

func (tr *TRenderer) Resize() {
	tr.tscreen.Sync()
}

func (tr *TRenderer) Clear() {
	tr.WheelDown(tr.cursorCoordinate.Y)
}

// WheelUp 滚动条向上，文本向下
func (tr *TRenderer) WheelUp(n int) {
	if tr.bufferOffsetY < n {
		tr.bufferOffsetY = 0
		tr.cursorCoordinate.Y += tr.bufferOffsetY
	} else {
		tr.bufferOffsetY -= n
		tr.cursorCoordinate.Y += n
	}
}

// WheelDown 滚动条向下，文本向上
func (tr *TRenderer) WheelDown(n int) {
	if tr.cursorCoordinate.Y < n {
		tr.bufferOffsetY += tr.cursorCoordinate.Y
		tr.cursorCoordinate.Y = 0
	} else {
		tr.bufferOffsetY += n
		tr.cursorCoordinate.Y -= n
	}
}

// Show 展示到窗口画面
func (tr *TRenderer) Show() {
	tr.tscreen.Clear()
	size := tr.getSize()
	for y := 0; y < size.height; y++ {
		lineData, found := tr.totalBuffer[y+tr.bufferOffsetY]
		if found {
			for x := 0; x < size.width; x++ {
				char, found := lineData[x]
				if found {
					tstyle := tcell.StyleDefault
					if colorStyle, ok := char.style.(*terminalcolor.ColorStyle); ok {
						tstyle = terminalcolor.ToTcellStyle(colorStyle)
					}
					for i, r := range char.char {
						tr.tscreen.SetContent(x+i, y, r, nil, tstyle)
					}
					x += char.width() - 1
				}
			}
		}
	}
	tr.tscreen.ShowCursor(tr.showCursorCoordinate.X, tr.showCursorCoordinate.Y)
	tr.tscreen.Show()
}

func (tr *TRenderer) reset() {

}
