// Package lexer 实现各种分词器
package lexer

import "github.com/yetsing/startprompt/token"

type GetTokensFunc func(input string) []token.Token
