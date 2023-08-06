package lexer

import "github.com/yetsing/startprompt/token"

type GetTokensFunc func(input string) []token.Token
