package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yetsing/startprompt/token"
)

/*
这个分词器跟在编译器或者解释器里面的有点不一样
因为用户的输入在不断变化，他需要容忍错误，尽可能地解析字符
会包括空格的 Token
*/

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func py3GenericInteger(buffer *CodeBuffer, prefix string, isLegalDigit func(rune) bool) string {
	start := buffer.GetIndex()
	if len(prefix) > 0 {
		p := buffer.PeekString(len(prefix))
		if !strings.EqualFold(p, prefix) {
			return ""
		}
		buffer.Advance(len(prefix))
	}
	for {
		if buffer.CurrentChar() == '_' {
			buffer.Advance(1)
			if isLegalDigit(buffer.CurrentChar()) {
				buffer.Advance(1)
			} else {
				buffer.Unread(1)
				break
			}
		} else if isLegalDigit(buffer.CurrentChar()) {
			buffer.Advance(1)
		} else {
			break
		}
	}
	end := buffer.GetIndex()
	s := buffer.Slice(start, end)
	if len(prefix) > 0 && strings.EqualFold(s, prefix) {
		return ""
	}
	return s
}

func py3bininteger(buffer *CodeBuffer) string {
	isbindigit := func(r rune) bool {
		return r == '0' || r == '1'
	}
	return py3GenericInteger(buffer, "0b", isbindigit)
}

func py3octinteger(buffer *CodeBuffer) string {
	isoctdigit := func(r rune) bool {
		return r >= '0' && r <= '7'
	}
	return py3GenericInteger(buffer, "0o", isoctdigit)
}

func py3hexinteger(buffer *CodeBuffer) string {
	ishexdigit := func(r rune) bool {
		return isdigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	}
	return py3GenericInteger(buffer, "0x", ishexdigit)
}

func Py3ReadNumber(buffer *CodeBuffer) string {
	prefix := buffer.PeekString(2)
	switch prefix {
	case "0b", "0B":
		return py3bininteger(buffer)
	case "0o", "0O":
		return py3octinteger(buffer)
	case "0x", "0X":
		return py3hexinteger(buffer)
	}
	//    十进制整数或者浮点数
	index := buffer.GetIndex()
	hasDot := false
	hasExp := false
	for buffer.HasChar() {
		ch := buffer.CurrentChar()
		if isdigit(ch) {
			buffer.Advance(1)
			continue
		}

		continueLoop := true
		switch ch {
		case '.':
			if hasDot {
				continueLoop = false
				break
			}
			hasDot = true
			//    存在 10. 这样的写法，所以后面的字符不用判断了
			buffer.Advance(1)
		case '_':
			if isdigit(buffer.Peek()) {
				buffer.Advance(1)
			} else {
				continueLoop = false
			}
		case 'e', 'E':
			if hasExp {
				continueLoop = false
				break
			}
			hasExp = true
			peek := buffer.Peek()
			if isdigit(peek) {
				buffer.Advance(1)
			} else if peek == '+' || peek == '-' {
				//    看下一个是不是数字
				if isdigit(buffer.PeekN(2)) {
					//    跳过 E+ E- e+ e- 等两个字符
					buffer.Advance(2)
				} else {
					continueLoop = false
				}
			} else {
				continueLoop = false
				break
			}
		default:
			continueLoop = false
		}
		if continueLoop {
			continue
		}

		break
	}
	return buffer.Slice(index, buffer.GetIndex())
}

type Py3Lexer struct {
	code   string
	buffer *CodeBuffer

	lastToken token.Token

	// 记录多行字符串的状态
	enterMultiline bool
	multilineEnd   string
	multilineType  token.TokenType

	// 记录括号层级
	parenLevel int

	// 缩进栈
	indentStack []int

	// 上一行是否以 \ 结尾，如果是下一行直接接上上一行
	continuedLine bool

	tokens []token.Token

	oneCharOps string
	twoCharOps []string
	// 字符串前面的前缀，也包括了字节字符串的前缀，因为它们两个分词都是 token.String
	stringPrefixs []string
	// 多行字符串的前缀
	multilineStringPrefixs []string
}

func NewPy3Lexer(code string) *Py3Lexer {
	return &Py3Lexer{
		code:        code,
		buffer:      NewCodeBuffer(code),
		indentStack: []int{0},
		lastToken:   token.NewToken(token.NL, ""),

		oneCharOps: ",.()+-*/=^%&|~<>[]{}:@;",
		twoCharOps: []string{
			"//", "+=", "-=", "*=", "/=", "%=",
			"&=", "|=",
			"==", ">=", "<=", "!=",
			">>", "<<",
			// type hint 返回值的标记
			"->",
			// **kwargs 前面的 **
			"**",
		},
		stringPrefixs: []string{
			"\"", "'",

			"r\"", "r'",
			"u\"", "u'",
			"R\"", "R'",
			"U\"", "U'",
			"f\"", "f'",
			"F\"", "F'",
			"fr\"", "fr'",
			"Fr\"", "Fr'",
			"fR\"", "fR'",
			"FR\"", "FR'",
			"rf\"", "rf'",
			"rF\"", "rF'",
			"Rf\"", "Rf'",
			"RF\"", "RF'",

			"b\"", "b'",
			"B\"", "B'",
			"br\"", "br'",
			"Br\"", "Br'",
			"bR\"", "bR'",
			"BR\"", "BR'",
			"rb\"", "rb'",
			"rB\"", "rB'",
			"Rb\"", "Rb'",
			"RB\"", "RB'",
		},
		multilineStringPrefixs: []string{
			"\"\"\"", "'''",

			"r\"\"\"", "r'''",
			"u\"\"\"", "u'''",
			"R\"\"\"", "R'''",
			"U\"\"\"", "U'''",
			"f\"\"\"", "f'''",
			"F\"\"\"", "F'''",
			"fr\"\"\"", "fr'''",
			"Fr\"\"\"", "Fr'''",
			"fR\"\"\"", "fR'''",
			"FR\"\"\"", "FR'''",
			"rf\"\"\"", "rf'''",
			"rF\"\"\"", "rF'''",
			"Rf\"\"\"", "Rf'''",
			"RF\"\"\"", "RF'''",

			"b\"\"\"", "b'''",
			"B\"\"\"", "B'''",
			"br\"\"\"", "br'''",
			"Br\"\"\"", "Br'''",
			"bR\"\"\"", "bR'''",
			"BR\"\"\"", "BR'''",
			"rb\"\"\"", "rb'''",
			"rB\"\"\"", "rB'''",
			"Rb\"\"\"", "Rb'''",
			"RB\"\"\"", "RB'''",
		},
	}
}

func (l *Py3Lexer) Tokens() []token.Token {
	if len(l.tokens) == 0 {
		for l.buffer.HasChar() {
			l.lineTokens()
		}

		// 多行字符串还没有输入完成
		if l.enterMultiline {
			l.buildToken(token.String)
			l.enterMultiline = false
		}

		// 跟 Python 标准库 tokenize 保持一致，如果最后没有换行符，也给他添加一个 newline token
		if !strings.HasSuffix(l.code, "\n") {
			l.buffer.Mark()
			l.buildToken(token.NewLine)
		}

		l.buffer.Mark()
		// 缩进栈有值，创建对应的 dedent token
		for l.indentStack[len(l.indentStack)-1] > 0 {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			l.buildToken(token.Dedent)
		}
	}
	return l.tokens
}

// 解析一行的 token ，这么做的原因是 Python 独有的缩进，按行可以更好地解析缩进
func (l *Py3Lexer) lineTokens() {
	buffer := l.buffer
	// 解析每行开始的缩进
	if l.parenLevel == 0 && !l.enterMultiline && !l.continuedLine {
		index := buffer.Mark()
		for buffer.CurrentChar() == ' ' {
			buffer.Advance(1)
		}
		// 代码前面的缩进才有意义
		var dents []token.Token
		if buffer.HasChar() && buffer.CurrentChar() != '\r' && buffer.CurrentChar() != '\n' &&
			buffer.CurrentChar() != '#' {
			dents = l.indentsOrDedents()
		}
		// 如果空格没有被解析为缩进，重置解析的位置
		// 还有就是 dedent 的 token 值是空，造成开头的空格被丢弃
		// 比如下面这种情况， "print(3)" 行会有一个缩进 token ，"print(4)" 行因为与上一行空格一致，所以不会有缩进 token
		// 需要解析为空格，不然输出后会发现少了一块
		// if 1:
		//   print(3)
		//   print(4)
		//   if 2:
		//     pass
		//   else:
		//     pass
		spaceCount := 0
		for _, tk := range dents {
			spaceCount += len(tk.Literal)
		}
		buffer.SetIndex(index + spaceCount)
	}

	l.continuedLine = false
	for buffer.HasChar() {
		// 处理多行字符串
		if l.enterMultiline {
			currentLine := buffer.CurrentLine()
			i := strings.Index(currentLine, l.multilineEnd)
			//     没有找到多行字符串的结束符号
			if i == -1 {
				// 跳过一整行
				buffer.Advance(utf8.RuneCountInString(currentLine))
				return
			}
			//     找到了多行字符串的结束符号
			buffer.Advance(utf8.RuneCountInString(currentLine[:i]) + len(l.multilineEnd))
			l.buildToken(l.multilineType)
			l.enterMultiline = false
			l.multilineEnd = ""
			l.multilineType = token.Unspecific
			continue
		}

		ch := buffer.CurrentChar()
		// 换行
		if ch == '\n' {
			l.newLineOrNl()
			break
		}

		// 空格
		if unicode.IsSpace(ch) {
			buffer.Mark()
			buffer.ReadSpace()
			l.buildToken(token.Whitespace)
			continue
		}

		// 检查多行字符串的开始
		match := false
		for _, prefix := range l.multilineStringPrefixs {
			peek := buffer.PeekString(len(prefix))
			if peek == prefix {
				buffer.Mark()
				l.enterMultiline = true
				// 倒数的三个
				l.multilineEnd = prefix[len(prefix)-3:]
				l.multilineType = token.String
				buffer.Advance(len(prefix))
				match = true
				break
			}
		}
		if match {
			continue
		}

		// 单行字符串
		for _, prefix := range l.stringPrefixs {
			peek := buffer.PeekString(len(prefix))
			if peek == prefix {
				l.quotedString(len(prefix), rune(prefix[len(prefix)-1]))
				match = true
				break
			}
		}
		if match {
			continue
		}

		// 三字符的操作符
		threeChar := buffer.PeekString(3)
		if stringIn(threeChar, "//=", ">>=", "<<=", "...") {
			l.operator(threeChar)
			continue
		}

		// 两字符的操作符
		twoChar := buffer.PeekString(2)
		if stringIn(twoChar, l.twoCharOps...) {
			l.operator(twoChar)
			continue
		}

		// 位于行末尾的 \ ，表示下一行接上上一行
		if twoChar == "\\\n" {
			buffer.Mark()
			buffer.Advance(2)
			l.buildToken(token.Whitespace)
			l.continuedLine = true
			return
		}

		// 注释
		if ch == '#' {
			l.comment()
			continue
		}

		// 还有类似 ".234" 这样的小数，所以有两个条件
		if isdigit(ch) || (ch == '.' && isdigit(buffer.Peek())) {
			l.number()
			continue
		}

		// 单字符的操作符
		if strings.ContainsRune(l.oneCharOps, ch) {
			// 维护括号的层级，用来判断当前是否在括号里面
			if strings.ContainsRune("([{", ch) {
				l.parenLevel++
			} else if strings.ContainsRune(")]}", ch) {
				l.parenLevel--
			}
			l.operator(string(ch))
			continue
		}

		// 变量（标志符）
		if isIdentifierStart(ch) {
			l.identifier()
			continue
		}

		// 错误字符
		buffer.Mark()
		buffer.Advance(1)
		l.buildToken(token.Error)
	}
}

func (l *Py3Lexer) identifier() token.Token {
	l.buffer.Mark()
	for isIdentifierContinue(l.buffer.CurrentChar()) {
		l.buffer.Advance(1)
	}
	return l.buildToken(token.Name)
}

func (l *Py3Lexer) number() token.Token {
	l.buffer.Mark()
	Py3ReadNumber(l.buffer)
	return l.buildToken(token.Number)
}

func (l *Py3Lexer) comment() token.Token {
	l.buffer.Mark()
	l.buffer.ReadUntil('\n')
	return l.buildToken(token.Comment)
}

func (l *Py3Lexer) quotedString(offset int, quote rune) token.Token {
	buffer := l.buffer
	buffer.Mark()
	buffer.Advance(offset)
	for buffer.HasChar() {
		c := buffer.CurrentChar()
		if c == '\\' {
			if buffer.Peek() == quote {
				buffer.Advance(2)
				continue
			} else if buffer.Peek() == '\\' {
				buffer.Advance(2)
				continue
			}
		}
		if c != quote {
			buffer.Advance(1)
		} else {
			break
		}
	}
	// 跳过末尾的引号
	buffer.Advance(1)
	return l.buildToken(token.String)
}

func (l *Py3Lexer) operator(op string) token.Token {
	l.buffer.Mark()
	l.buffer.Advance(len(op))
	return l.buildToken(token.Operator)
}

func (l *Py3Lexer) indentsOrDedents() []token.Token {
	indentLength := len(l.buffer.ReadFromMark())
	// 新的缩进层级
	if indentLength > l.indentStack[len(l.indentStack)-1] {
		l.indentStack = append(l.indentStack, indentLength)
		return []token.Token{l.buildToken(token.Indent)}
	}
	if !intSliceHas(l.indentStack, indentLength) {
		panic("unindent does not match any outer indentation level")
	}
	var dedents []token.Token
	for indentLength < l.indentStack[len(l.indentStack)-1] {
		// pop 最后一个元素
		l.indentStack = l.indentStack[:len(l.indentStack)-1]
		l.buffer.Mark()
		tk := l.buildToken(token.Dedent)
		dedents = append(dedents, tk)
	}
	return dedents
}

func (l *Py3Lexer) newLineOrNl() token.Token {
	l.buffer.Mark()
	ttype := token.NewLine
	if l.parenLevel > 0 || l.lastToken.TypeIn(token.NewLine, token.NL, token.Comment) {
		ttype = token.NL
	}
	l.buffer.Advance(1)
	return l.buildToken(ttype)
}

func (l *Py3Lexer) buildToken(tokenType token.TokenType) token.Token {
	s := l.buffer.ReadFromMark()
	tk := token.NewToken(tokenType, s)
	if tokenType != token.Comment && tokenType != token.Whitespace {
		l.lastToken = tk
	}
	l.tokens = append(l.tokens, tk)
	return tk
}

func intSliceHas(ns []int, n int) bool {
	for _, i := range ns {
		if i == n {
			return true
		}
	}
	return false
}

func stringIn(s string, es ...string) bool {
	for _, e := range es {
		if e == s {
			return true
		}
	}
	return false
}

// 参考 Python 的规则 https://docs.python.org/3/reference/lexical_analysis.html#identifiers
var idStartCategorys = map[string]uint8{
	"Lu": 1,
	"Ll": 1,
	"Lm": 1,
	"Lt": 1,
	"Lo": 1,
	"Nl": 1,
}

var idContinueCategorys = map[string]uint8{
	"Lu": 1,
	"Ll": 1,
	"Lm": 1,
	"Lt": 1,
	"Lo": 1,
	"Nl": 1,
	"Mn": 1,
	"Mc": 1,
	"Nd": 1,
	"Pc": 1,
}

func isIdentifierStart(ch rune) bool {
	switch ch {
	case '_':
		return true
	default:
		cat := UnicodeCategory(ch)
		_, ok := idStartCategorys[cat]
		return ok
	}
}

func isIdentifierContinue(ch rune) bool {
	switch ch {
	case '_':
		return true
	default:
		cat := UnicodeCategory(ch)
		_, ok := idContinueCategorys[cat]
		return ok
	}
}

// UnicodeCategory returns the Unicode Character Category of the given rune.
// code from https://stackoverflow.com/a/53507592
func UnicodeCategory(r rune) string {
	for name, table := range unicode.Categories {
		if len(name) == 2 && unicode.Is(table, r) {
			return name
		}
	}
	return "Cn"
}
