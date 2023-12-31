package startprompt

import (
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type Schema map[token.TokenType]*terminalcolor.ColorStyle

func (s Schema) StyleForToken(tokenType token.TokenType) *terminalcolor.ColorStyle {
	if v, found := s[tokenType]; found {
		return v
	}
	// 使用父类的样式
	for t, style := range s {
		if t.HasChild(tokenType) {
			return style
		}
	}
	return styleDefault
}

func (s Schema) StyleForSelection(origStyle *terminalcolor.ColorStyle) *terminalcolor.ColorStyle {
	tokenType := token.Selection
	style, found := s[tokenType]
	if !found {
		style = selectionStyleDefault
	}
	if style.FgIsColorDefault() {
		//    保留原有的文本颜色
		style = style.CopyAndFg(origStyle.Fg())
	}
	return style
}

var (
	selectionStyleDefault = terminalcolor.NewBgColorStyleHex("#40334d")
	styleDefault          = terminalcolor.NewDefaultColorStyle()
)

var defaultSchema = map[token.TokenType]*terminalcolor.ColorStyle{
	token.Keyword:  terminalcolor.NewFgColorStyleHex("#ee00ee"),
	token.Operator: terminalcolor.NewFgColorStyleHex("#aa6666"),
	token.Number:   terminalcolor.NewFgColorStyleHex("#2aacb8"),
	// token.Name:     terminalcolor.NewFgColorStyleHex("#008800"),
	token.String: terminalcolor.NewFgColorStyleHex("#6aab73"),

	token.Error:   terminalcolor.NewColorStyleHex("#000000", "#ff8888"),
	token.Comment: terminalcolor.NewFgColorStyleHex("#0000dd"),

	token.CompletionMenuCompletion:        terminalcolor.NewColorStyleHex("#ffffbb", "#888888"),
	token.CompletionMenuCompletionCurrent: terminalcolor.NewColorStyleHex("#000000", "#dddddd"),
	token.CompletionMenuMetaCurrent:       terminalcolor.NewColorStyleHex("#000000", "#bbbbbb"),
	token.CompletionMenuMeta:              terminalcolor.NewColorStyleHex("#cccccc", "#888888"),
	token.CompletionMenuProgressBar:       terminalcolor.NewColorStyleHex("", "#aaaaaa"),
	token.CompletionMenuProgressButton:    terminalcolor.NewColorStyleHex("", "#000000"),

	token.Selection: selectionStyleDefault,
}
