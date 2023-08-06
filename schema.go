package startprompt

import (
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

type Schema map[token.TokenType]terminalcolor.Style

var defaultSchema = map[token.TokenType]terminalcolor.Style{
	token.INT:     terminalcolor.BrightCyan,
	token.STRING:  terminalcolor.BrightGreen,
	token.ILLEGAL: terminalcolor.BrightRed,
	token.KEYWORD: terminalcolor.BrightMagenta,
}
