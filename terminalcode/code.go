package terminalcode

import "fmt"

//goland:noinspection GoUnusedConst
const (
	// EraseScreen
	// https://vt100.net/docs/vt100-ug/chapter3.html#ED
	// 擦除整个屏幕并且移动光标到左上角
	EraseScreen = "\x1b[2J"

	// EraseEndOfLine https://vt100.net/docs/vt100-ug/chapter3.html#EL
	// 擦除从光标位置到当前行尾的范围（包括光标位置）
	EraseEndOfLine = "\x1b[K"

	// EraseDown https://vt100.net/docs/vt100-ug/chapter3.html#ED
	// 擦除从当前行到屏幕底部的范围
	EraseDown = "\x1b[J"

	// CarriageReturn 移动光标到行首
	CarriageReturn = "\r"
	NEWLINE        = "\n"
	CRLF           = "\r\n"

	// HideCursor 隐藏光标
	HideCursor = "\x1b[?25l"
	// DisplayCursor 显示光标
	DisplayCursor = "\x1b[?25h"
)

// CursorGoto 移动光标到指定位置
//
//goland:noinspection GoUnusedExportedFunction
func CursorGoto(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", x, y)
}

// CursorUp https://vt100.net/docs/vt100-ug/chapter3.html#CUU
// 向上移动光标
//
//goland:noinspection GoUnusedExportedFunction
func CursorUp(amount int) string {
	return fmt.Sprintf("\x1b[%dA", amount)
}

// CursorDown 向下移动光标
//
//goland:noinspection GoUnusedExportedFunction
func CursorDown(amount int) string {
	return fmt.Sprintf("\x1b[%dB", amount)
}

// CursorForward 向右移动光标
func CursorForward(amount int) string {
	return fmt.Sprintf("\x1b[%dC", amount)
}

// CursorBackward 向左移动光标
func CursorBackward(amount int) string {
	return fmt.Sprintf("\x1b[%dD", amount)
}
