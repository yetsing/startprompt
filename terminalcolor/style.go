package terminalcolor

import (
	"strings"
)

type Style interface {
	ColorEscape() string
	ResetEscape() string
}

func ApplyStyle(style Style, text string, reset bool) string {
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
