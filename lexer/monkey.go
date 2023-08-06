package lexer

/*
《用Go语言自制解释器》中 monkey 语言的分词器
*/

import "github.com/yetsing/startprompt/token"

type monkeyLexer struct {
	input        []rune
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func newMonkeyLexer(input string) *monkeyLexer {
	l := &monkeyLexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *monkeyLexer) nextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case ' ', '\r', '\n', '\t':
		literal := l.readWhitespace()
		tok = token.Token{
			Type:    token.WHITESPACE,
			Literal: literal,
		}
		return tok
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *monkeyLexer) readWhitespace() string {
	position := l.position
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *monkeyLexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *monkeyLexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *monkeyLexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *monkeyLexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *monkeyLexer) readString() string {
	position := l.position
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	end := l.position
	if l.ch == '"' {
		end++
	}
	return string(l.input[position:end])
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

var keywords = map[string]bool{
	"fn":     true,
	"let":    true,
	"true":   true,
	"false":  true,
	"if":     true,
	"else":   true,
	"return": true,
}

func lookupIdent(ident string) token.TokenType {
	if _, ok := keywords[ident]; ok {
		return token.KEYWORD
	}
	return token.IDENT
}

func GetMonkeyTokens(input string) []token.Token {
	l := newMonkeyLexer(input)
	var result []token.Token
	for true {
		tk := l.nextToken()
		if tk.TypeIs(token.EOF) {
			break
		}
		result = append(result, tk)
	}
	return result
}
