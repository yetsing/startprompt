package startprompt

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type Area struct {
	start Coordinate
	end   Coordinate
}

func (a *Area) Contains(coordinate Coordinate) bool {
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

type XChar struct {
	char *Char
	x    int
}

type TRenderer struct {
	tscreen       tcell.Screen
	schema        Schema
	promptFactory PromptFactory
	selection     *Area

	//    xy 坐标到输入行列的映射
	inputLocationMap map[Coordinate]Location
	//    保存至今为止全部的输出
	//    {y: {x: Char}}
	totalBuffer map[int]map[int]*Char
	//    当前输入在 totalBuffer 的坐标（输入左上角）
	bufferCoordinate Coordinate
	//    渲染 totalBuffer 中 >= bufferOffsetY 的内容
	bufferOffsetY int
	//    当前输入在窗口中的坐标（输入左上角）
	inputCoordinate Coordinate
	//    在窗口中显示的光标坐标
	showCursorCoordinate Coordinate
}

func newTRenderer(tscreen tcell.Screen, schema Schema, promptFactory PromptFactory) *TRenderer {
	return &TRenderer{
		totalBuffer:   map[int]map[int]*Char{},
		tscreen:       tscreen,
		schema:        schema,
		promptFactory: promptFactory,
		selection:     &Area{Coordinate{0, 0}, Coordinate{0, 0}},
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
	locationMap := screen.getLocationMap()
	tr.inputLocationMap = make(map[Coordinate]Location, len(locationMap))
	for coordinate, location := range locationMap {
		newCoordinate := Coordinate{
			X: coordinate.X + tr.inputCoordinate.X,
			Y: coordinate.Y + tr.inputCoordinate.Y,
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
		tr.bufferCoordinate.X = 0
		tr.bufferCoordinate.Y += screen.maxCursorCoordinate.Y + 1
		tr.inputCoordinate.X = 0
		tr.inputCoordinate.Y += screen.maxCursorCoordinate.Y + 1
		tr.showCursorCoordinate = tr.inputCoordinate
	} else {
		//    移动光标到正确位置
		cursorCoordinate := screen.getCoordinate(
			renderContext.document.CursorPositionRow(),
			renderContext.document.CursorPositionCol())
		tr.showCursorCoordinate.X = tr.inputCoordinate.X + cursorCoordinate.X
		tr.showCursorCoordinate.Y = tr.inputCoordinate.Y + cursorCoordinate.Y
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
	tr.inputCoordinate.X = 0
	tr.inputCoordinate.Y += screen.maxCursorCoordinate.Y
	tr.showCursorCoordinate = tr.inputCoordinate
	tr.Show()
}

func (tr *TRenderer) Resize() {
	tr.tscreen.Sync()
}

func (tr *TRenderer) Clear() {
	tr.WheelDown(tr.inputCoordinate.Y)
}

// WheelUp 滚动条向上，文本向下
func (tr *TRenderer) WheelUp(n int) {
	if tr.bufferOffsetY < n {
		tr.bufferOffsetY = 0
		tr.inputCoordinate.Y += tr.bufferOffsetY
	} else {
		tr.bufferOffsetY -= n
		tr.inputCoordinate.Y += n
	}
}

// WheelDown 滚动条向下，文本向上
func (tr *TRenderer) WheelDown(n int) {
	if tr.inputCoordinate.Y < n {
		tr.bufferOffsetY += tr.inputCoordinate.Y
		tr.inputCoordinate.Y = 0
	} else {
		tr.bufferOffsetY += n
		tr.inputCoordinate.Y -= n
	}
}

// Show 展示到窗口画面
func (tr *TRenderer) Show() {
	tr.tscreen.HideCursor()
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
					if tr.selection.Contains(Coordinate{x, y}) {
						tstyle = tstyle.Reverse(true)
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

// GetClosetLocation 返回跟坐标最接近的行列，返回的布尔值是否可以找到
func (tr *TRenderer) GetClosetLocation(coordinate Coordinate) (Location, bool) {
	// 在 (x-4, y) ~ (x+4, y) 的范围内寻找行列
	end := maxInt(0, coordinate.X-4)
	for x := coordinate.X; x >= end; x-- {
		loc, found := tr.inputLocationMap[Coordinate{x, coordinate.Y}]
		if found {
			return loc, found
		}
	}
	end = coordinate.X + 4
	for x := coordinate.X; x <= end; x++ {
		loc, found := tr.inputLocationMap[Coordinate{x, coordinate.Y}]
		if found {
			return loc, found
		}
	}
	return Location{-1, -1}, false
}

// InInputArea 判断坐标是否在当前输入区域内（以行为准）
func (tr *TRenderer) InInputArea(coordinate Coordinate) bool {
	width, _ := tr.tscreen.Size()
	for x := 0; x < width; x++ {
		if _, found := tr.inputLocationMap[Coordinate{x, coordinate.Y}]; found {
			return true
		}
	}
	return false
}

// InTextArea 判断坐标是否在文本区域内（以行为准）
func (tr *TRenderer) InTextArea(coordinate Coordinate) bool {
	by := coordinate.Y + tr.bufferOffsetY
	_, found := tr.totalBuffer[by]
	return found
}

// 返回指定坐标处的字符开始坐标，调用者要保证坐标在文本区域内
func (tr *TRenderer) getCharCoordinate(coordinate Coordinate) Coordinate {
	by := coordinate.Y + tr.bufferOffsetY
	lineData, found := tr.totalBuffer[by]
	if !found {
		panic(fmt.Errorf("invalid coordinate: %+v", coordinate))
	}
	for x := coordinate.X; x >= 0; x-- {
		_, found := lineData[x]
		if found {
			return Coordinate{x, coordinate.Y}
		}
	}
	return Coordinate{0, coordinate.Y}
}

// SelectWord 选择指定坐标处的单词（鼠标双击触发）
func (tr *TRenderer) SelectWord(coordinate Coordinate) {
	by := coordinate.Y + tr.bufferOffsetY
	lineData, found := tr.totalBuffer[by]
	if !found {
		//    点击处没有文本
		return
	}
	//    获取单词的开始和结束
	width, _ := tr.tscreen.Size()
	DebugLog("select word coordinate: %+v", coordinate)
	coordinate = tr.getCharCoordinate(coordinate)
	DebugLog("select word adjust coordinate: %+v", coordinate)
	var end Coordinate
	for x := coordinate.X; x < width; x++ {
		char, found := lineData[x]
		if !found {
			end = Coordinate{x, coordinate.Y}
			break
		}
		if IsSpace(char.char) {
			end = Coordinate{x, coordinate.Y}
			break
		}
		x += char.width() - 1
	}

	var start Coordinate
	for x := coordinate.X; x >= 0; x-- {
		char, found := lineData[x]
		if !found {
			start = Coordinate{x, coordinate.Y}
			break
		}
		if IsSpace(char.char) {
			start = Coordinate{x + char.width(), coordinate.Y}
			break
		}
		x -= char.width() - 1
	}
	if start.equal(&end) {
		return
	}
	DebugLog("select word start: %+v, end: %+v", start, end)
	tr.selection = &Area{
		start: start,
		end:   end,
	}
}
