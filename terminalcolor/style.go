package terminalcolor

import (
	"strconv"
	"strings"
)

func ApplyStyle(style *ColorStyle, text string, reset bool) string {
	m := make([]string, 3)
	m[0] = style.ColorEscape()
	m[1] = text
	if reset {
		m[2] = style.ResetEscape()
	}
	return strings.Join(m, "")
}

func escapeAttrs(attrs []string) string {
	if len(attrs) > 0 {
		return "\x1b[" + strings.Join(attrs, ";") + "m"
	}
	return ""
}

type Color int

const ColorDefault Color = -1

type ColorStyle struct {
	fg        Color
	bg        Color
	bold      bool
	underline bool
	italic    bool
	reverse   bool
}

func NewDefaultColorStyle() *ColorStyle {
	return NewColorStyle(ColorDefault, ColorDefault)
}

func NewFgColorStyleHex(fg string) *ColorStyle {
	return NewColorStyleGeneric(Color256IndexFromHexRGB(fg), ColorDefault, false, false, false)
}

func NewBgColorStyleHex(bg string) *ColorStyle {
	return NewColorStyleGeneric(ColorDefault, Color256IndexFromHexRGB(bg), false, false, false)
}

func NewColorStyle(fg Color, bg Color) *ColorStyle {
	return NewColorStyleGeneric(fg, bg, false, false, false)
}

func NewColorStyleHex(fg string, bg string) *ColorStyle {
	return NewColorStyleGeneric(
		Color256IndexFromHexRGB(fg),
		Color256IndexFromHexRGB(bg),
		false, false, false)
}

//goland:noinspection GoUnusedExportedFunction
func NewColorStyleGeneric(fg, bg Color, bold, underline, italic bool) *ColorStyle {
	return &ColorStyle{
		fg:        fg,
		bg:        bg,
		bold:      bold,
		underline: underline,
		italic:    italic,
	}
}

func (c *ColorStyle) Fg() Color {
	return c.fg
}

func (c *ColorStyle) FgIsColorDefault() bool {
	return c.fg == ColorDefault
}

func (c *ColorStyle) CopyAndFg(fg Color) *ColorStyle {
	return &ColorStyle{
		fg:        fg,
		bg:        c.bg,
		bold:      c.bold,
		underline: c.underline,
		italic:    c.italic,
		reverse:   c.reverse,
	}
}

func (c *ColorStyle) CopyAndReverse(on bool) *ColorStyle {
	return &ColorStyle{
		fg:        c.fg,
		bg:        c.bg,
		bold:      c.bold,
		underline: c.underline,
		italic:    c.italic,
		reverse:   on,
	}
}

func (c *ColorStyle) ColorEscape() string {
	var attrs []string
	if c.reverse {
		attrs = append(attrs, "7")
	}
	if c.fg != ColorDefault {
		if c.fg >= Color256Start {
			attrs = append(attrs, "38", "5", strconv.Itoa(int(c.fg-Color256Start)))
		} else {
			attrs = append(attrs, strconv.Itoa(int(c.fg)))
		}
	}
	if c.bg != ColorDefault {
		if c.bg >= Color256Start {
			attrs = append(attrs, "48", "5", strconv.Itoa(int(c.bg-Color256Start)))
		} else {
			// 背景的数字 = 前景的数字 + 10
			attrs = append(attrs, strconv.Itoa(int(c.bg)+10))
		}
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

func (c *ColorStyle) ResetEscape() string {
	var attrs []string
	if c.fg != ColorDefault {
		attrs = append(attrs, "39")
	}
	if c.bg != ColorDefault {
		attrs = append(attrs, "49")
	}
	if c.bold || c.underline || c.italic || c.reverse {
		attrs = append(attrs, "00")
	}
	return escapeAttrs(attrs)
}
