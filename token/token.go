package token

/*
TokenType 定义参照 pygments https://github.com/pygments/pygments/blob/master/pygments/token.py

    pygments.token
    ~~~~~~~~~~~~~~

    Basic token types and the standard tokens.

    :copyright: Copyright 2006-2023 by the Pygments team, see AUTHORS.
    :license: BSD, see LICENSE for details.
*/

import (
	"fmt"
	"strings"
)

type TokenType string

//goland:noinspection GoUnusedConst
const (
	Unspecific TokenType = "Unspecific"

	Text                 TokenType = "text"
	Whitespace           TokenType = "whitespace"
	Escape               TokenType = "escape"
	Error                TokenType = "error"
	Other                TokenType = "other"
	Keyword              TokenType = "keyword"
	KeywordConstant      TokenType = Keyword + ".constant"
	KeywordDeclaration   TokenType = Keyword + ".declaration"
	KeywordNamespace     TokenType = Keyword + ".namespace"
	KeywordPseudo        TokenType = Keyword + ".pseudo"
	KeywordReserved      TokenType = Keyword + ".reserved"
	KeywordType          TokenType = Keyword + ".type"
	Name                 TokenType = "name"
	NameAttribute        TokenType = Name + ".attribute"
	NameBuiltin          TokenType = Name + ".builtin"
	NameBuiltinPseudo    TokenType = Name + ".builtin.pseudo"
	NameClass            TokenType = Name + ".class"
	NameConstant         TokenType = Name + ".constant"
	NameDecorator        TokenType = Name + ".decorator"
	NameEntity           TokenType = Name + ".entity"
	NameException        TokenType = Name + ".exception"
	NameFunction         TokenType = Name + ".function"
	NameFunctionMagic    TokenType = Name + ".function.magic"
	NameProperty         TokenType = Name + ".property"
	NameLabel            TokenType = Name + ".label"
	NameNamespace        TokenType = Name + ".namespace"
	NameOther            TokenType = Name + ".other"
	NameTag              TokenType = Name + ".tag"
	NameVariable         TokenType = Name + ".variable"
	NameVariableClass    TokenType = NameVariable + ".class"
	NameVariableGlobal   TokenType = NameVariable + ".global"
	NameVariableInstance TokenType = NameVariable + ".instance"
	NameVariableMagic    TokenType = NameVariable + ".magic"
	Literal              TokenType = "literal"
	LiteralDate          TokenType = Literal + ".date"
	String               TokenType = Literal + ".string"
	StringAffix          TokenType = String + ".affix"
	StringBacktick       TokenType = String + ".backtick"
	StringChar           TokenType = String + ".char"
	StringDelimiter      TokenType = String + ".delimiter"
	StringDoc            TokenType = String + ".doc"
	StringDouble         TokenType = String + ".double"
	StringEscape         TokenType = String + ".escape"
	StringHeredoc        TokenType = String + ".heredoc"
	StringInterpol       TokenType = String + ".interpol"
	StringOther          TokenType = String + ".other"
	StringRegex          TokenType = String + ".regex"
	StringSingle         TokenType = String + ".single"
	StringSymbol         TokenType = String + ".symbol"
	Number               TokenType = Literal + "number"
	NumberBin            TokenType = Number + ".bin"
	NumberFloat          TokenType = Number + ".float"
	NumberHex            TokenType = Number + ".hex"
	NumberInteger        TokenType = Number + ".integer"
	NumberIntegerLong    TokenType = Number + ".integer.long"
	NumberOct            TokenType = Number + ".oct"
	Operator             TokenType = "operator"
	OperatorWord         TokenType = Operator + ".word"
	// Punctuation 标点符号
	Punctuation       TokenType = "punctuation"
	PunctuationMarker TokenType = Punctuation + ".marker"
	Comment           TokenType = "comment"
	CommentHashbang   TokenType = Comment + ".hashbang"
	CommentMultiline  TokenType = Comment + ".multiline"
	// CommentPreproc 预处理类型注释
	CommentPreproc     TokenType = Comment + ".preproc"
	CommentPreprocFile TokenType = Comment + ".preprocfile"
	CommentSingle      TokenType = Comment + ".single"
	CommentSpecial     TokenType = Comment + ".special"
	Generic            TokenType = "generic"
	GenericDeleted     TokenType = Generic + ".deleted"
	// GenericEmph emph 强调
	GenericEmph       TokenType = Generic + ".emph"
	GenericError      TokenType = Generic + ".error"
	GenericHeading    TokenType = Generic + ".heading"
	GenericInserted   TokenType = Generic + ".inserted"
	GenericOutput     TokenType = Generic + ".output"
	GenericPrompt     TokenType = Generic + ".prompt"
	GenericStrong     TokenType = Generic + ".strong"
	GenericSubheading TokenType = Generic + ".subheading"
	GenericEmphStrong TokenType = Generic + ".emphstrong"
	GenericTraceback  TokenType = Generic + ".traceback"

	Prompt                 TokenType = "prompt"
	PromptSecondLinePrefix TokenType = Prompt + ".secondlineprefix"

	CompletionMenu                  TokenType = "completionmenu"
	CompletionMenuCurrentCompletion TokenType = CompletionMenu + ".currentcompletion"
	CompletionMenuCompletion        TokenType = CompletionMenu + ".completion"
	CompletionProgressButton        TokenType = CompletionMenu + ".progressbutton"
	CompletionProgressBar           TokenType = CompletionMenu + ".progressbar"

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

func NewToken(tokenType TokenType, s string) Token {
	return Token{
		Type:    tokenType,
		Literal: s,
	}
}
