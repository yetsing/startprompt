package startprompt

import (
	"github.com/gdamore/tcell/v2"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

// MouseInfoOfInput 在当前输入中的一些鼠标信息
type MouseInfoOfInput struct {
	//    鼠标位置在输入的行列
	location Location
	//    鼠标位置在哪个补全项上，用于点击时切换补全
	completeIndex int
}

type TRenderer struct {
	tscreen        tcell.Screen
	scrollTextView *sScrollTextView
	//    补全菜单信息
	completionMenuInfo *cCompletionMenuInfo

	schema        Schema
	promptFactory PromptFactory

	//    xy 坐标到输入行列的映射
	inputLocationMap map[Coordinate]Location
	//    光标相对于输入左上角的（相对）坐标
	cursorRelativeCoordinate Coordinate

	//    键盘事件是否触发
	triggerEventKey bool
	//    鼠标事件是否触发
	triggerEventMouse bool
}

func newTRenderer(tscreen tcell.Screen, schema Schema, promptFactory PromptFactory) *TRenderer {
	return &TRenderer{
		tscreen:        tscreen,
		scrollTextView: newScrollTextView(),
		schema:         schema,
		promptFactory:  promptFactory,
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
	tr.completionMenuInfo = nil
	if renderContext.completeState != nil {
		menu := newCompletionMenu(screen, renderContext.completeState, 7)
		menu.write()
		tr.completionMenuInfo = menu.getInfo()
		//    转换补全的坐标为窗口坐标
		inputStartCoordinate := tr.scrollTextView.getInputStartCoordinate()
		tr.completionMenuInfo.area.start.addY(inputStartCoordinate.Y)
		tr.completionMenuInfo.area.end.addY(inputStartCoordinate.Y)
	}

	return screen
}

func (tr *TRenderer) updateWithScreen(screen *Screen) {
	tr.scrollTextView.readScreen(screen)
	locationMap := screen.getLocationMap()
	tr.inputLocationMap = make(map[Coordinate]Location, len(locationMap))
	inputStartCoordinate := tr.scrollTextView.getInputStartCoordinate()
	for coordinate, location := range locationMap {
		newCoordinate := Coordinate{
			X: coordinate.X + inputStartCoordinate.X,
			Y: coordinate.Y + inputStartCoordinate.Y,
		}
		tr.inputLocationMap[newCoordinate] = location
	}
}

func (tr *TRenderer) render(renderContext *RenderContext, abort bool, accept bool) {
	//    写入屏幕输出
	screen := tr.getNewScreen(renderContext)
	tr.updateWithScreen(screen)

	if renderContext.cancelSelection {
		tr.scrollTextView.cancelSelection()
	}

	//    用户输入完毕或者放弃输入或者退出，另起一行
	if accept || abort {
		tr.scrollTextView.acceptInput()
		tr.cursorRelativeCoordinate = Coordinate{}
	} else {
		//    移动光标到正确位置
		relativeCoordinate := screen.getCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		tr.cursorRelativeCoordinate = relativeCoordinate
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
	tr.scrollTextView.inputToEnd()
	tr.cursorRelativeCoordinate = Coordinate{}
	tr.Show()
}

func (tr *TRenderer) update() {
	tr.scrollTextView.update()
}

func (tr *TRenderer) getCursorCoordinate() Coordinate {
	inputStartCoordinate := tr.scrollTextView.getInputStartCoordinate()
	return Coordinate{
		X: tr.cursorRelativeCoordinate.X + inputStartCoordinate.X,
		Y: tr.cursorRelativeCoordinate.Y + inputStartCoordinate.Y,
	}
}

func (tr *TRenderer) Resize() {
	tr.tscreen.Sync()
}

// Clear 按下 Ctrl-L 触发，置顶光标所在行
func (tr *TRenderer) Clear() {
	tr.scrollTextView.moveUp(tr.getCursorCoordinate().Y)
}

// WheelUp 滚动条向上，文本向下
func (tr *TRenderer) WheelUp(n int) {
	tr.scrollTextView.scrollDown(n)
}

// WheelDown 滚动条向下，文本向上
func (tr *TRenderer) WheelDown(n int) {
	tr.scrollTextView.scrollUp(n)
}

// Show 展示到窗口画面
func (tr *TRenderer) Show() {
	tr.tscreen.HideCursor()
	tr.tscreen.Clear()

	//    只有键盘导致的光标移动，才将其移动到窗口内
	if tr.triggerEventKey {
		//    检查光标是否在窗口内
		cursorCoordinate := tr.getCursorCoordinate()
		//    光标在窗口的上面
		if cursorCoordinate.Y < 0 {
			tr.scrollTextView.moveDown(-cursorCoordinate.Y)
		}
		_, height := tr.tscreen.Size()
		//    光标在窗口的下面
		if cursorCoordinate.Y > height-1 {
			tr.scrollTextView.moveUp(cursorCoordinate.Y - (height - 1))
		}
	}

	size := tr.getSize()
	for y := 0; y < size.height; y++ {
		lineData, found := tr.scrollTextView.getLineAt(y)
		if found {
			for _, datum := range lineData {
				tstyle := tcell.StyleDefault
				if colorStyle, ok := datum.style.(*terminalcolor.ColorStyle); ok {
					tstyle = terminalcolor.ToTcellStyle(colorStyle)
				}
				if tr.scrollTextView.inSelection(Coordinate{datum.x, y}) {
					tstyle = tstyle.Reverse(true)
				}
				for i, r := range datum.char {
					tr.tscreen.SetContent(datum.x+i, y, r, nil, tstyle)
				}
			}
		}
	}
	cursorCoordinate := tr.getCursorCoordinate()
	tr.tscreen.ShowCursor(cursorCoordinate.X, cursorCoordinate.Y)
	tr.tscreen.Show()
}

func (tr *TRenderer) reset() {

}

// GetClosetLocation 返回跟坐标最接近的行列，返回的布尔值表示是否找到
func (tr *TRenderer) GetClosetLocation(coordinate Coordinate) (Location, bool) {
	// 在 (0, y) ~ (x, y) 的范围内寻找行列
	for x := coordinate.X; x >= 0; x-- {
		loc, found := tr.inputLocationMap[Coordinate{x, coordinate.Y}]
		if found {
			return loc, found
		}
	}
	return Location{-1, -1}, false
}

func (tr *TRenderer) GetMouseInfoOfInput(coordinate Coordinate) *MouseInfoOfInput {
	loc, _ := tr.GetClosetLocation(coordinate)
	completeIndex := -1
	if tr.completionMenuInfo != nil {
		completeIndex = tr.completionMenuInfo.getCompleteIndex(coordinate)
	}
	return &MouseInfoOfInput{loc, completeIndex}
}

// LineInInputArea InInputArea 判断坐标 y 所在行是否在当前输入区域内
func (tr *TRenderer) LineInInputArea(y int) bool {
	return tr.scrollTextView.inputContainLine(y)
}

// LineInTextArea 判断坐标 y 所在行是否在文本区域内
func (tr *TRenderer) LineInTextArea(y int) bool {
	return tr.scrollTextView.containLine(y)
}

// MouseDown 鼠标（左键）按下
func (tr *TRenderer) MouseDown(coordinate Coordinate) {
	tr.scrollTextView.mouseDown(coordinate)
}

func (tr *TRenderer) MouseMove(coordinate Coordinate) {
	tr.scrollTextView.mouseMove(coordinate)
}

func (tr *TRenderer) MouseUp(coordinate Coordinate) {
	tr.scrollTextView.mouseUp(coordinate)
}

// Dblclick 鼠标双击
func (tr *TRenderer) Dblclick(coordinate Coordinate) {
	tr.scrollTextView.dblclick(coordinate)
}

// TripeClick 鼠标三击
func (tr *TRenderer) TripeClick(coordinate Coordinate) {
	tr.scrollTextView.tripeClick(coordinate)
}

// TriggerEventKey 键盘事件触发
func (tr *TRenderer) TriggerEventKey() {
	tr.scrollTextView.restoreScroll()
	tr.triggerEventKey = true
	tr.triggerEventMouse = false
}

// TriggerEventMouse 鼠标事件触发
func (tr *TRenderer) TriggerEventMouse() {
	tr.triggerEventKey = false
	tr.triggerEventMouse = true
}
