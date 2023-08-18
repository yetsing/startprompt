package token

import (
	"fmt"
	"strings"
)

type TokenType string

//goland:noinspection GoUnusedConst
const (
	Unspecific             TokenType = "unspecific"
	Prompt                 TokenType = "prompt"
	PromptSecondLinePrefix TokenType = "Prompt.SecondLinePrefix"

	CompletionMenu                  TokenType = "CompletionMenu"
	CompletionMenuCurrentCompletion TokenType = "CompletionMenu.CurrentCompletion"
	CompletionMenuCompletion        TokenType = "CompletionMenu.Completion"
	CompletionProgressButton        TokenType = "CompletionMenu.ProgressButton"
	CompletionProgressBar           TokenType = "CompletionMenu.ProgressBar"

	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 1343456
	STRING TokenType = "STRING" // "foobar"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	DOT      TokenType = "."

	LT TokenType = "<"
	GT TokenType = ">"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	KEYWORD TokenType = "keyword"

	WHITESPACE TokenType = "whitespace"
)

func (t TokenType) HasChild(child TokenType) bool {
	return len(child) > len(t) && strings.HasPrefix(string(child), fmt.Sprintf("%s.", t))
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) TypeIs(ttype TokenType) bool {
	return t.Type == ttype
}

func (t *Token) TypeIn(ttypes ...TokenType) bool {
	for _, ttype := range ttypes {
		if t.TypeIs(ttype) {
			return true
		}
	}
	return false
}
