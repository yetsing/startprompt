package terminalcolor

/*
原先只有 3-bit 8 种颜色，后来加了 Bright 系列，于是有了 16 种颜色
4-bit 颜色转义序列，可以表示 16 种颜色
参考：https://en.wikipedia.org/wiki/ANSI_escape_code 中 "3-bit and 4-bit" 一节

*/

//goland:noinspection GoUnusedConst
const (
	Color16Black Color = iota + 30
	Color16Red
	Color16Green
	Color16Yellow
	Color16Blue
	Color16Magenta
	Color16Cyan
	Color16Gray
)

//goland:noinspection GoUnusedConst
const (
	Color16BrightBlack Color = iota + 90
	Color16BrightRed
	Color16BrightGreen
	Color16BrightYellow
	Color16BrightBlue
	Color16BrightMagenta
	Color16BrightCyan
	Color16BrightWhite
)

// 常用的颜色样式
//
//goland:noinspection GoUnusedGlobalVariable
var (
	Black   = NewColorStyle(Color16Black, ColorDefault)
	Red     = NewColorStyle(Color16Red, ColorDefault)
	Green   = NewColorStyle(Color16Green, ColorDefault)
	Yellow  = NewColorStyle(Color16Yellow, ColorDefault)
	Blue    = NewColorStyle(Color16Blue, ColorDefault)
	Magenta = NewColorStyle(Color16Magenta, ColorDefault)
	Cyan    = NewColorStyle(Color16Cyan, ColorDefault)
	Gray    = NewColorStyle(Color16Gray, ColorDefault)

	BrightBlack   = NewColorStyle(Color16BrightBlack, ColorDefault)
	BrightRed     = NewColorStyle(Color16BrightRed, ColorDefault)
	BrightGreen   = NewColorStyle(Color16BrightGreen, ColorDefault)
	BrightYellow  = NewColorStyle(Color16BrightYellow, ColorDefault)
	BrightBlue    = NewColorStyle(Color16BrightBlue, ColorDefault)
	BrightMagenta = NewColorStyle(Color16BrightMagenta, ColorDefault)
	BrightCyan    = NewColorStyle(Color16BrightCyan, ColorDefault)
	BrightWhite   = NewColorStyle(Color16BrightWhite, ColorDefault)
)
