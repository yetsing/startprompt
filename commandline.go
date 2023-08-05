package startprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/yetsing/startprompt/lexer"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
	"os"
	"strconv"

	"github.com/yetsing/startprompt/inputstream"
	"github.com/yetsing/startprompt/terminalcode"
)

func mpanic(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
}

func ensureOk(err error) {
	if err != nil {
		panic(err)
	}
}

func ctrlKey(k rune) rune {
	return k & 0x1f
}

//goland:noinspection GoUnusedFunction
func iscntrl(k rune) bool {
	return k <= 0x1f || k == 127
}

type LexerFactory func(text string) LexerAbc

type LexerAbc interface {
	GetTokens() []token.Token
}

type Schema map[token.TokenType]terminalcolor.Style

var defaultSchema = map[token.TokenType]terminalcolor.Style{
	token.INT:     terminalcolor.BrightCyan,
	token.STRING:  terminalcolor.BrightGreen,
	token.ILLEGAL: terminalcolor.BrightRed,
	token.KEYWORD: terminalcolor.BrightMagenta,
}

type CommandLine struct {
	reader  *bufio.Reader
	writer  *bufio.Writer
	running bool
}

func (c *CommandLine) ReadInput() string {
	line := inputstream.NewLine()
	handler := inputstream.NewBaseHandler(line)
	is := inputstream.NewInputStream(handler)
	var r rune
	var err error
	reader := c.reader
	for true {
		r, _, err = reader.ReadRune()
		if err != nil {
			mpanic("error read: %v\n", err)
		}
		is.Feed(r)
		c.draw(line.Document())
		if r == ctrlKey('q') {
			c.running = false
			break
		}
		if line.Finished() {
			break
		}
	}
	text := line.Text()
	if len(text) > 0 {
		c.OutputStringf("\r\n")
	}
	return text
}

func (c *CommandLine) draw(doc *inputstream.Document) {
	text := doc.Text
	cursorX := doc.CursorX
	// 隐藏光标
	c.writeString(terminalcode.HideCursor)
	// 移动光标到行首
	c.writeString(terminalcode.CarriageReturn)
	// 删除当行到屏幕下方
	c.writeString(terminalcode.EraseDown)

	l := lexer.New(text)
	screen := NewScreen(defaultSchema)
	screen.WriteTokens(l.GetTokens())
	result, lastPosition := screen.Output()
	lastX := lastPosition.X
	c.writeString(result)

	// 移动光标
	if lastX > cursorX {
		c.writeString(terminalcode.CursorBackward(lastX - cursorX))
	} else if lastX < cursorX {
		c.writeString(terminalcode.CursorForward(cursorX - lastX))
	}
	// 显示光标
	c.writeString("\x1b[?25h")
	c.flush()
}

// 写入字符串到输出缓冲中
func (c *CommandLine) writeString(s string) {
	_, err := c.writer.WriteString(s)
	ensureOk(err)
}

func (c *CommandLine) flush() {
	err := c.writer.Flush()
	ensureOk(err)
}

// OutputString 直接输出字符串（无缓冲）
func (c *CommandLine) OutputString(s string) {
	c.writeString(s)
	c.flush()
}

func (c *CommandLine) OutputStringf(format string, a ...any) {
	c.writeString(fmt.Sprintf(format, a...))
	c.flush()
}

type Position struct {
	row int
	col int
}

func (c *CommandLine) getCursorPosition() Position {
	c.OutputString("\x1b[6n")
	var buf [32]byte
	var i int
	// 回复格式为 \x1b[A;BR
	// A 和 B 就是光标的行和列
	// todo 这里读取光标位置实际使用发现有一个问题
	// 如果用户的输入没有及时处理（比如一直按着 A），下面的循环就会读到剩余的用户输入
	// 而不是简单的 \x1b[A;BR 转义序列
	for i = 0; i < 32; i++ {
		c, err := c.reader.ReadByte()
		ensureOk(err)
		if c == 'R' {
			break
		} else {
			buf[i] = c
		}
	}
	if buf[0] != '\x1b' || buf[1] != '[' {
		mpanic("invalid cursor position report escape: %v", buf[:i])
	}
	sepIndex := bytes.IndexByte(buf[:], ';')
	if sepIndex == -1 {
		mpanic("invalid cursor position report separator: %v", buf[:i])
	}
	row, err := strconv.ParseInt(string(buf[2:sepIndex]), 10, 32)
	ensureOk(err)
	col, err := strconv.ParseInt(string(buf[sepIndex+1:i]), 10, 32)
	ensureOk(err)
	return Position{
		row: int(row),
		col: int(col),
	}
}

func (c *CommandLine) Running() bool {
	return c.running
}

func NewCommandLine() *CommandLine {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	return &CommandLine{reader: reader, writer: writer, running: true}
}
