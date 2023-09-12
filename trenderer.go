package startprompt

import (
	"github.com/gdamore/tcell/v2"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type area struct {
	start Coordinate
	end   Coordinate
}

func (a *area) Contains(coordinate Coordinate) bool {
	if a.start.Y == a.end.Y {
		return a.start.Y == coordinate.Y && a.start.X <= coordinate.X && coordinate.X < a.end.X
	}
	if coordinate.Y == a.start.Y {
		return a.start.X <= coordinate.X
	} else if coordinate.Y == a.end.Y {
		return coordinate.X < a.end.X
	}
	return coordinate.Y > a.start.Y && coordinate.Y < a.end.Y
}

type xChar struct {
	*Char
	x int
}

type sScrollTextView struct {
	//    行列二维数组
	data [][]xChar
	//    当前输入左上角在 data 第几行
	inputY int
	//    data 在 y 轴上的偏移量（滚动量）
	//    其实就是从第几行开始显示在窗口中
	offsetY int
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
	width := screen.getWidth()
	buffer := screen.GetBuffer()
	//    如果只有一行，那么文本宽度就是最后坐标的 x
	//    否则就是 screen 宽度
	endX := lastCoordinate.X
	if lastCoordinate.Y > 0 {
		endX = width
	}
	for y := 0; y <= lastCoordinate.Y; y++ {
		vy := st.inputY + y
		st.growTo(vy)
		//    清空当前行数据
		st.data[vy] = nil
		lineBuffer, found := buffer[y]
		if found {
			x := 0
			for x < endX {
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
	if vy <= len(st.data)-1 {
		return st.data[vy], true
	}
	return nil, false
}

// scrollUp 文本向上滚动
func (st *sScrollTextView) scrollUp(n int) {
	st.offsetY += n
	length := len(st.data)
	if st.offsetY >= length {
		st.offsetY = length - 1
	}
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
	DebugLog("coordinate=%+v, offsetY=%d", coordinate, st.offsetY)
	lineData, found := st.getLineAt(coordinate.Y)
	if !found {
		DebugLog("not found lineData %+v", coordinate)
		return area{}
	}
	//    找到窗口坐标所在字符索引
	index := -1
	for i, datum := range lineData {
		DebugLog("found <%s> index=%d", datum.char, index)
		if datum.x > coordinate.X {
			break
		}
		index = i
	}
	if index == -1 {
		DebugLog("not found index %+v", coordinate)
		return area{}
	}

	DebugLog("found index=%d", index)

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
	DebugLog("word area start=%+v end=%+v", start, end)
	return area{
		start: start,
		end:   end,
	}
}

type TRenderer struct {
	tscreen        tcell.Screen
	selection      area
	scrollTextView *sScrollTextView
	schema         Schema
	promptFactory  PromptFactory

	//    xy 坐标到输入行列的映射
	inputLocationMap map[Coordinate]Location
	//    当前输入开头在窗口中的坐标（输入左上角）
	inputStartCoordinate Coordinate
	//    在窗口中显示的光标坐标
	cursorCoordinate Coordinate
}

func newTRenderer(tscreen tcell.Screen, schema Schema, promptFactory PromptFactory) *TRenderer {
	return &TRenderer{
		tscreen:        tscreen,
		selection:      area{Coordinate{0, 0}, Coordinate{0, 0}},
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
	if renderContext.completeState != nil {
		newCompletionMenu(screen, renderContext.completeState, 7).write()
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
	//    用户输入完毕或者放弃输入或者退出，另起一行
	if accept || abort {
		tr.scrollTextView.acceptInput()
		tr.cursorCoordinate = tr.scrollTextView.getInputStartCoordinate()
	} else {
		//    移动光标到正确位置
		cursorCoordinate := screen.getCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		inputStartCoordinate := tr.scrollTextView.getInputStartCoordinate()
		tr.cursorCoordinate.X = inputStartCoordinate.X + cursorCoordinate.X
		tr.cursorCoordinate.Y = inputStartCoordinate.Y + cursorCoordinate.Y
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
	tr.cursorCoordinate = tr.scrollTextView.getInputStartCoordinate()
	tr.Show()
}

func (tr *TRenderer) Resize() {
	tr.tscreen.Sync()
}

func (tr *TRenderer) Clear() {
	inputStartCoordinate := tr.scrollTextView.getInputStartCoordinate()
	tr.WheelDown(inputStartCoordinate.Y)
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
	size := tr.getSize()
	for y := 0; y < size.height; y++ {
		lineData, found := tr.scrollTextView.getLineAt(y)
		if found {
			for _, datum := range lineData {
				tstyle := tcell.StyleDefault
				if colorStyle, ok := datum.style.(*terminalcolor.ColorStyle); ok {
					tstyle = terminalcolor.ToTcellStyle(colorStyle)
				}
				if tr.selection.Contains(Coordinate{datum.x, y}) {
					tstyle = tstyle.Reverse(true)
				}
				for i, r := range datum.char {
					tr.tscreen.SetContent(datum.x+i, y, r, nil, tstyle)
				}
			}
		}
	}
	tr.tscreen.ShowCursor(tr.cursorCoordinate.X, tr.cursorCoordinate.Y)
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

// LineInInputArea InInputArea 判断坐标 y 所在行是否在当前输入区域内
func (tr *TRenderer) LineInInputArea(y int) bool {
	return tr.scrollTextView.inputContainLine(y)
}

// LineInTextArea 判断坐标 y 所在行是否在文本区域内
func (tr *TRenderer) LineInTextArea(y int) bool {
	return tr.scrollTextView.containLine(y)
}

// SelectWord 选择指定坐标处的单词（鼠标双击触发）
func (tr *TRenderer) SelectWord(coordinate Coordinate) {
	tr.selection = tr.scrollTextView.getWordArea(coordinate)
}
