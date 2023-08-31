package startprompt

import (
	"github.com/gdamore/tcell"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type TScreen struct {
	screen tcell.Screen
	schema Schema
	// 保存全部文本内容
	// {y: {x: Char}}
	buffer map[int]map[int]*Char
	// y 的偏移量，表示屏幕应该显示 buffer 中 >= offsetY 的文本数据
	offsetY int
	// 当前输入左上角的光标坐标（窗口内的绝对坐标）
	inputX int
	inputY int
	// 窗口中光标坐标（是一个相对于文本左上角的坐标，而不是窗口左上角）
	tx int
	ty int
	// 文本中光标的行列（相对于文本的）
	inputRow int
	inputCol int
	// 保存光标行列到 ty tx 的映射
	cursorMap map[Coordinate]Coordinate

	secondLinePrefixFunc func() []token.Token
}

func (ts *TScreen) WriteTokensAtPos(tx int, ty int, tokens []token.Token) {

}

func (ts *TScreen) WriteTokens(tokens []token.Token, saveInputPos bool) {

}

func (ts *TScreen) WriteRune(r rune, style terminalcolor.Style, saveInputPos bool) {

}

func (ts *TScreen) writeAtPos(tx int, ty int, char *Char) {

}
