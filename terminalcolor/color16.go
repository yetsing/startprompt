package terminalcolor

import "strconv"

/*
原先只有 3-bit 8 种颜色，后来加了 Bright 系列，于是有了 16 种颜色
4-bit 颜色转义序列，可以表示 16 种颜色
参考：https://en.wikipedia.org/wiki/ANSI_escape_code 中 "3-bit and 4-bit" 一节

*/

// Color16 定义前景颜色的数字，背景的数字 = 前景的数字 + 10
type Color16 int

const Color16Default Color16 = -1

//goland:noinspection GoUnusedConst
const (
	Color16Black Color16 = iota + 30
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
	Color16BrightBlack Color16 = iota + 90
	Color16BrightRed
	Color16BrightGreen
	Color16BrightYellow
	Color16BrightBlue
	Color16BrightMagenta
	Color16BrightCyan
	Color16BrightWhite
)

type Color16Style struct {
	fg        Color16
	bg        Color16
	bold      bool
	underline bool
	italic    bool
}

//goland:noinspection GoUnusedExportedFunction
func NewColor16Style(fg, bg Color16, bold, underline, italic bool) *Color16Style {
	return &Color16Style{
		fg:        fg,
		bg:        bg,
		bold:      bold,
		underline: underline,
		italic:    italic,
	}
}

func (c *Color16Style) ColorEscape() string {
	var attrs []string
	if c.fg != Color16Default {
		attrs = append(attrs, strconv.Itoa(int(c.fg)))
	}
	if c.bg != Color16Default {
		// 背景的数字 = 前景的数字 + 10
		attrs = append(attrs, strconv.Itoa(int(c.bg)+10))
	}
	if c.bold {
		attrs = append(attrs, "01")
	}
	if c.underline {
		attrs = append(attrs, "04")
	}
	if c.italic {
		attrs = append(attrs, "03")
	}
	return escapeAttrs(attrs)
}

func (c *Color16Style) ResetEscape() string {
	var attrs []string
	if c.fg != Color16Default {
		attrs = append(attrs, "39")
	}
	if c.bg != Color16Default {
		attrs = append(attrs, "49")
	}
	if c.bold || c.underline || c.italic {
		attrs = append(attrs, "00")
	}
	return escapeAttrs(attrs)
}

// 常用的颜色样式
//
//goland:noinspection GoUnusedGlobalVariable
var (
	Black   = NewColor16Style(Color16Black, Color16Default, false, false, false)
	Red     = NewColor16Style(Color16Red, Color16Default, false, false, false)
	Green   = NewColor16Style(Color16Green, Color16Default, false, false, false)
	Yellow  = NewColor16Style(Color16Yellow, Color16Default, false, false, false)
	Blue    = NewColor16Style(Color16Blue, Color16Default, false, false, false)
	Magenta = NewColor16Style(Color16Magenta, Color16Default, false, false, false)
	Cyan    = NewColor16Style(Color16Cyan, Color16Default, false, false, false)
	Gray    = NewColor16Style(Color16Gray, Color16Default, false, false, false)

	BrightBlack   = NewColor16Style(Color16BrightBlack, Color16Default, false, false, false)
	BrightRed     = NewColor16Style(Color16BrightRed, Color16Default, false, false, false)
	BrightGreen   = NewColor16Style(Color16BrightGreen, Color16Default, false, false, false)
	BrightYellow  = NewColor16Style(Color16BrightYellow, Color16Default, false, false, false)
	BrightBlue    = NewColor16Style(Color16BrightBlue, Color16Default, false, false, false)
	BrightMagenta = NewColor16Style(Color16BrightMagenta, Color16Default, false, false, false)
	BrightCyan    = NewColor16Style(Color16BrightCyan, Color16Default, false, false, false)
	BrightWhite   = NewColor16Style(Color16BrightWhite, Color16Default, false, false, false)
)
