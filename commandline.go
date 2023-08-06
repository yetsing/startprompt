package startprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"

	"github.com/yetsing/startprompt/inputstream"
	"github.com/yetsing/startprompt/lexer"
	"github.com/yetsing/startprompt/terminalcode"
)

func panicf(format string, a ...any) {
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

type CommandLine struct {
	reader         *bufio.Reader
	writer         *bufio.Writer
	running        bool
	tokensFunc     lexer.GetTokensFunc
	enableDebugLog bool
	oldState       *term.State
}

func (c *CommandLine) ReadInput() string {
	line := inputstream.NewLine()
	handler := inputstream.NewBaseHandler(line)
	is := inputstream.NewInputStream(handler)
	var r rune
	var err error
	reader := c.reader
	var inputText string
	for true {
		r, _, err = reader.ReadRune()
		if err != nil {
			panicf("error read: %v\n", err)
		}
		DebugLog("read rune: %d", r)
		is.Feed(r)
		DebugLog("feed: %d", r)
		doc := line.Document()
		c.draw(doc)
		DebugLog("draw document")
		if r == ctrlKey('q') {
			DebugLog("exit normally")
			c.running = false
			break
		}
		if line.Finished() {
			inputText = doc.Text
			break
		}
	}
	DebugLog("return input: %s", inputText)
	if len(inputText) > 0 {
		c.OutputStringf("\r\n")
	}
	return inputText
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

	screen := NewScreen(defaultSchema)
	screen.WriteTokens(c.tokensFunc(text))
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
		panicf("invalid cursor position report escape: %v", buf[:i])
	}
	sepIndex := bytes.IndexByte(buf[:], ';')
	if sepIndex == -1 {
		panicf("invalid cursor position report separator: %v", buf[:i])
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

func (c *CommandLine) Restore() {
	err := term.Restore(int(os.Stdin.Fd()), c.oldState)
	if err != nil {
		fmt.Printf("restore error: %v\r\n", err)
	}
}

func NewCommandLine(tokensFunc lexer.GetTokensFunc, debug bool) (*CommandLine, error) {
	// 开启 terminal raw mode
	// 这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	// 不会自动展示用户输入
	// 更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	if debug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	return &CommandLine{
		reader:     reader,
		writer:     writer,
		running:    true,
		tokensFunc: tokensFunc,
		oldState:   oldState,
	}, nil
}
