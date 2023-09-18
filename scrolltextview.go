package startprompt

import (
	"strings"

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
func (st *sScrollTextView) getClosetCharCoordinate(windowCoordinate Coordinate) (Coordinate, bool) {
	lineData, found := st.getLineAt(windowCoordinate.Y)
	if !found {
		return Coordinate{-1, -1}, false
	}
	ret := Coordinate{0, windowCoordinate.Y}
	//    找到最后一个 x 坐标小于等于的
	for _, datum := range lineData {
		if datum.x > windowCoordinate.X {
			return ret, true
		}
		ret.X = datum.x
	}
	return Coordinate{-1, -1}, false
}

// getWordArea 返回坐标处的单词区域
func (st *sScrollTextView) getWordArea(coordinate Coordinate) area {
	lineData, found := st.getLine(coordinate.Y)
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
func (st *sScrollTextView) mouseDown(windowCoordinate Coordinate) {
	coor := st.convertWindowCoordinate(windowCoordinate)
	st.selection = area{start: coor, end: coor}
}

func (st *sScrollTextView) mouseMove(windowCoordinate Coordinate) {
	coor := st.convertWindowCoordinate(windowCoordinate)
	st.selection.end = coor
}

func (st *sScrollTextView) mouseUp(windowCoordinate Coordinate) {
	coor := st.convertWindowCoordinate(windowCoordinate)
	st.selection.end = coor

	//    如果坐标超出范围，将其设置为最后一个坐标
	last := st.getLastCoordinate()
	st.selection.limitTo(last)
}

func (st *sScrollTextView) dblclick(windowCoordinate Coordinate) {
	coor := st.convertWindowCoordinate(windowCoordinate)
	st.selection = st.getWordArea(coor)
}

func (st *sScrollTextView) tripeClick(windowCoordinate Coordinate) {
	coor := st.convertWindowCoordinate(windowCoordinate)
	st.selection = area{
		start: Coordinate{0, coor.Y},
		end:   Coordinate{1 << 24, coor.Y},
	}
}

// inSelection 判断窗口坐标是否在选中区域内
func (st *sScrollTextView) inSelection(windowCoordinate Coordinate) bool {
	coor := st.convertWindowCoordinate(windowCoordinate)
	return st.selection.Contains(coor)
}

func (st *sScrollTextView) getSelectionText() string {
	var builder strings.Builder
	start := st.selection.getStart()
	end := st.selection.getEnd()
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

func (st *sScrollTextView) convertWindowCoordinate(windowCoordinate Coordinate) Coordinate {
	coor := Coordinate{windowCoordinate.X, windowCoordinate.Y}
	coor.addY(st.offsetY)
	return coor
}

func (st *sScrollTextView) update() {
	if st.selection.isEmpty() {
		return
	}

	//   将选中文本复制到系统剪贴板
	//   文本发生变化时复制一次
	text := st.getSelectionText()
	if text != st.selectionText {
		clipboard.Write(clipboard.FmtText, []byte(text))
		st.selectionText = text
	}

}
