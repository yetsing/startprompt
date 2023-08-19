package startprompt

import (
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type Schema map[token.TokenType]terminalcolor.Style

func (s Schema) StyleForToken(tokenType token.TokenType) terminalcolor.Style {
	if v, found := s[tokenType]; found {
		return v
	}
	// 使用父类的样式
	for t, style := range s {
		if t.HasChild(tokenType) {
			return style
		}
	}
	return nil
}

var defaultSchema = map[token.TokenType]terminalcolor.Style{
	token.Keyword:  terminalcolor.NewFgColorStyleHex("#ee00ee"),
	token.Operator: terminalcolor.NewFgColorStyleHex("#aa6666"),
	token.Number:   terminalcolor.NewFgColorStyleHex("#ff0000"),
	token.Name:     terminalcolor.NewFgColorStyleHex("#008800"),
	token.String:   terminalcolor.NewFgColorStyleHex("#440000"),

	token.Error:   terminalcolor.NewColorStyleHex("#000000", "#ff8888"),
	token.Comment: terminalcolor.NewFgColorStyleHex("#0000dd"),

	token.CompletionMenuCurrentCompletion: terminalcolor.NewColorStyleHex("#000000", "#dddddd"),
	token.CompletionMenuCompletion:        terminalcolor.NewColorStyleHex("#ffff88", "#888888"),
	token.CompletionProgressButton:        terminalcolor.NewColorStyleHex("", "#000000"),
	token.CompletionProgressBar:           terminalcolor.NewColorStyleHex("", "#aaaaaa"),
}
