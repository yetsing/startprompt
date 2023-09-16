package startprompt

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
	"golang.design/x/clipboard"
)

type area struct {
	start Coordinate
	end   Coordinate
}

func (a *area) Contains(coordinate Coordinate) bool {
	start := a.getStart()
	end := a.getEnd()
	if start.Y == end.Y {
		return start.Y == coordinate.Y && start.X <= coordinate.X && coordinate.X < end.X
	}
	if coordinate.Y == start.Y {
		return start.X <= coordinate.X
	} else if coordinate.Y == end.Y {
		return coordinate.X < end.X
	}
	return coordinate.Y > start.Y && coordinate.Y < end.Y
}

// RectContains 判断点是否在开始和结束组成的矩形中
func (a *area) RectContains(coordinate Coordinate) bool {
	start := a.getStart()
	end := a.getEnd()
	return start.Y <= coordinate.Y && coordinate.Y < end.Y && start.X <= coordinate.X && coordinate.X < end.X
}

func (a *area) isEmpty() bool {
	start := a.getStart()
	end := a.getEnd()
	return start.Y > end.Y || (start.Y == end.Y && start.X >= end.X)
}

func (a *area) getStart() Coordinate {
	if a.start.gt(&a.end) {
		return a.end
	}
	return a.start
}

func (a *area) getEnd() Coordinate {
	if a.start.gt(&a.end) {
		return a.start
	}
	return a.end
}

// limitTo 坐标超出 coordinate 的设为 coordinate
func (a *area) limitTo(coordinate Coordinate) {
	if a.start.gt(&coordinate) {
		a.start = coordinate
	}
	if a.end.gt(&coordinate) {
		a.end = coordinate
	}
}

type xChar struct {
	*Char
	x int
}

type sScrollTextView struct {
	//    选中文本
	selectionText string
	//    行列二维数组
	data [][]xChar
	//    当前输入左上角在 data 第几行
	inputY int
	//    data 在 y 轴上的偏移量（滚动量）
	//    其实就是从第几行开始显示在窗口中
	//    范围在 [0, offsetLimitY]
	//        offsetY 之所以有个上界，是为了模拟终端的滚动效果，
	//        终端的滚动条默认无法移动，按下 Ctrl-L 置顶当前输入
	//        此时滚动条可以向上移动，向下只能移动到最初的位置
	offsetY int
	//    偏移量的最大值
	offsetLimitY int
	//    选中区域
	selection area
}

func newScrollTextView() *sScrollTextView {
	return &sScrollTextView{data: [][]xChar{nil}}
}

// growTo 增加数据长度， y 是从 0 开始的索引
func (st *sScrollTextView) growTo(y int) {
	for i := len(st.data) - 1; i < y; i++ {
		st.data = append(st.data, []xChar{})
	}
}

func (st *sScrollTextView) appendAt(vy int, xchar xChar) {
	st.data[vy] = append(st.data[vy], xchar)
}

func (st *sScrollTextView) readScreen(screen *Screen) {
	lastCoordinate := screen.getLastCoordinate()
	buffer := screen.GetBuffer()
	//    清空之前的输入数据
	st.data = st.data[:st.inputY]
	st.growTo(st.inputY + lastCoordinate.Y)
	for y := 0; y <= lastCoordinate.Y; y++ {
		vy := st.inputY + y
		lineBuffer, found := buffer[y]
		if found {
			//    当前行最大的 x 坐标
			endX := 0
			for x := range lineBuffer {
				if x > endX {
					endX = x
				}
			}
			x := 0
			for x <= endX {
				var char *Char
				if _, found := lineBuffer[x]; found {
					char = lineBuffer[x]
				} else {
					char = newChar(' ', nil)
				}
				st.appendAt(vy, xChar{char, x})
				x += char.width()
			}
		}
	}
}

// getLineAt 传入窗口坐标 y ，返回对应行数据
func (st *sScrollTextView) getLineAt(y int) ([]xChar, bool) {
	vy := st.offsetY + y
	return st.getLine(vy)
}

func (st *sScrollTextView) getLine(n int) ([]xChar, bool) {
	if n <= len(st.data)-1 {
		return st.data[n], true
	}
	return nil, false
}

func (st *sScrollTextView) getLastCoordinate() Coordinate {
	y := len(st.data) - 1
	lastLine := st.data[y]
	x := 0
	if len(lastLine) > 0 {
		ch := lastLine[len(lastLine)-1]
		x = ch.x + ch.width()
	}
	return Coordinate{x, y}
}

// restoreScroll 恢复原本的滚动位置
//
//	当我们滚动到之前的文本时，按下键盘，画面应该回到之前输入的位置。
//	效果参考终端
func (st *sScrollTextView) restoreScroll() {
	st.offsetY = st.offsetLimitY
}

// moveUp 文本向上移动，会增加滚动的边界
func (st *sScrollTextView) moveUp(n int) int {
	if st.offsetLimitY+n > len(st.data)-1 {
		n = len(st.data) - 1 - st.offsetLimitY
	}
	st.offsetY += n
	st.offsetLimitY += n
	return n
}

// moveDown 文本向下移动，会减少滚动的边界
func (st *sScrollTextView) moveDown(n int) int {
	if st.offsetLimitY < n {
		n = st.offsetLimitY
	}
	st.offsetY -= n
	st.offsetLimitY -= n
	return n
}

// scrollUp 文本向上滚动
func (st *sScrollTextView) scrollUp(n int) int {
	if st.offsetY+n > st.offsetLimitY {
		n = st.offsetLimitY - st.offsetY
	}
	st.offsetY += n
	return n
}

// scrollDown 文本向下滚动，返回实际滚动行数
func (st *sScrollTextView) scrollDown(n int) int {
	if n > st.offsetY {
		n = st.offsetY
	}
	st.offsetY -= n
	return n
}

func (st *sScrollTextView) inputToEnd() {
	st.inputY = len(st.data) - 1
}

func (st *sScrollTextView) acceptInput() {
	st.inputY = len(st.data)
	st.growTo(st.inputY)
}

// containLine 是否包含窗口坐标 y 处行
func (st *sScrollTextView) containLine(y int) bool {
	sy := st.offsetY + y
	return sy < len(st.data)
}

// inputContainLine 当前输入是否包含窗口坐标 y 处行
func (st *sScrollTextView) inputContainLine(y int) bool {
	sy := st.offsetY + y
	return st.inputY <= sy && sy < len(st.data)
}

// getInputStartCoordinate 返回当前输入左上角的窗口坐标
func (st *sScrollTextView) getInputStartCoordinate() Coordinate {
	return Coordinate{0, st.inputY - st.offsetY}
}

// getClosetCharCoordinate 返回最接近的字符窗口坐标，布尔值表示是否找到
func (st *sScrollTextView) getClosetCharCoordinate(coordinate Coordinate) (Coordinate, bool) {
	lineData, found := st.getLineAt(coordinate.Y)
	if !found {
		return Coordinate{-1, -1}, false
	}
	ret := Coordinate{0, coordinate.Y}
	//    找到最后一个 x 坐标小于等于的
	for _, datum := range lineData {
		if datum.x > coordinate.X {
			return ret, true
		}
		ret.X = datum.x
	}
	return Coordinate{-1, -1}, false
}

// getWordArea 返回窗口坐标处的单词区域（窗口坐标）
func (st *sScrollTextView) getWordArea(coordinate Coordinate) area {
	lineData, found := st.getLineAt(coordinate.Y)
	if !found {
		return area{}
	}
	//    找到窗口坐标所在字符索引
	index := -1
	for i, datum := range lineData {
		if datum.x > coordinate.X {
			break
		}
		index = i
	}
	if index == -1 {
		return area{}
	}

	length := len(lineData)
	//    默认是行尾
	end := Coordinate{lineData[length-1].x + lineData[length-1].width(), coordinate.Y}
	for i := index; i < length; i++ {
		xc := lineData[i]
		if IsSpace(xc.char) {
			end = Coordinate{xc.x, coordinate.Y}
			break
		}
	}
	//    默认是行首
	start := Coordinate{lineData[0].x, coordinate.Y}
	for i := index; i >= 0; i-- {
		xc := lineData[i]
		if IsSpace(xc.char) {
			start = Coordinate{xc.x + xc.width(), coordinate.Y}
			break
		}
	}
	return area{
		start: start,
		end:   end,
	}
}

// mouseDown 鼠标（左键）点击，传入窗口坐标
func (st *sScrollTextView) mouseDown(coordinate Coordinate) {
	coordinate.addY(st.offsetY)
	st.selection = area{start: coordinate, end: coordinate}
}

func (st *sScrollTextView) mouseMove(coordinate Coordinate) {
	coordinate.addY(st.offsetY)
	st.selection.end = coordinate
}

func (st *sScrollTextView) mouseUp(coordinate Coordinate) {
	coordinate.addY(st.offsetY)
	st.selection.end = coordinate

	//    如果坐标超出范围，将其设置为最后一个坐标
	last := st.getLastCoordinate()
	st.selection.limitTo(last)
}

func (st *sScrollTextView) dblclick(coordinate Coordinate) {
	st.selection = st.getWordArea(coordinate)
}

func (st *sScrollTextView) tripeClick(coordinate Coordinate) {
	coordinate.addY(st.offsetY)
	st.selection = area{
		start: Coordinate{0, coordinate.Y},
		end:   Coordinate{1 << 24, coordinate.Y},
	}
}

// inSelection 判断窗口坐标是否在选中区域内
func (st *sScrollTextView) inSelection(coordinate Coordinate) bool {
	coordinate.addY(st.offsetY)
	return st.selection.Contains(coordinate)
}

func (st *sScrollTextView) getSelectionText() string {
	var builder strings.Builder
	start := st.selection.getStart()
	end := st.selection.getEnd()
	DebugLog("getSelectionText start=%+v end=%+v", start, end)
	for y := start.Y; y <= end.Y; y++ {
		lineData, found := st.getLine(y)
		if found {
			for _, datum := range lineData {
				if st.selection.Contains(Coordinate{datum.x, y}) {
					builder.WriteString(datum.char)
				}
			}
		}
		if y != end.Y {
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

func (st *sScrollTextView) cancelSelection() {
	st.selection = area{}
}

func (st *sScrollTextView) update() {
	if st.selection.isEmpty() {
		return
	}

	DebugLog("selection %+v", st.selection)
	//   将选中文本复制到系统剪贴板
	//   文本发生变化时复制一次
	text := st.getSelectionText()
	DebugLog("selection text: %s", text)
	if text != st.selectionText {
		DebugLog("copy text")
		clipboard.Write(clipboard.FmtText, []byte(text))
		st.selectionText = text
	}

}

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
