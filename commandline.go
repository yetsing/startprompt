package startprompt

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"

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
}

func (c *CommandLine) ReadInput() (string, error) {
	// 开启 terminal raw mode
	// 这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	// 不会自动展示用户输入
	// 更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			fmt.Printf("term.Restore error: %v\r\n", err)
		}
	}()

	render := newRender(defaultSchema)
	line := newLine(render, newBaseCode, newBasePrompt, NewMemHistory())
	handler := NewBaseHandler(line)
	is := NewInputStream(handler)
	render.render(line.GetRenderContext())

	var r rune
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
		render.render(line.GetRenderContext())
		DebugLog("terminal width: %d", render.getWidth())
		//c.draw(line)
		if line.abort || line.accept {
			inputText = line.text()
			break
		}
		DebugLog("draw document")
		if r == ctrlKey('q') {
			DebugLog("exit normally")
			c.running = false
			break
		}
	}
	DebugLog("return input: <%s>", inputText)
	return inputText, nil
}

func (c *CommandLine) draw(line *Line) {
	renderCtx := line.GetRenderContext()
	// 为了防止在重画屏幕的过程中，光标出现闪烁，我们先隐藏光标，最后在显示光标
	// 参考：https://viewsourcecode.org/snaptoken/kilo/03.rawInputAndOutput.html#hide-the-cursor-when-repainting
	// 隐藏光标
	c.writeString(terminalcode.HideCursor)
	// 移动光标到行首
	c.writeString(terminalcode.CarriageReturn)
	// 删除当行到屏幕下方
	c.writeString(terminalcode.EraseDown)

	screen := NewScreen(defaultSchema, 0)
	screen.WriteTokens(renderCtx.code.GetTokens(), true)
	screen.saveInputPos()
	result, lastCoordinate := screen.Output()
	c.writeString(result)

	// 用户输入完毕或者放弃输入
	if renderCtx.accept || renderCtx.abort {
		// 另起一行
		c.writeString(terminalcode.CRLF)
	} else {
		doc := line.Document()
		// 移动光标
		cursorCoordinate := screen.getCursorCoordinate(doc.CursorPositionRow(), doc.CursorPositionCol())
		lastX := lastCoordinate.X
		if lastCoordinate.Y > cursorCoordinate.Y {
			c.writeString(terminalcode.CursorUp(lastCoordinate.Y - cursorCoordinate.Y))
		}
		if lastX > cursorCoordinate.X {
			c.writeString(terminalcode.CursorBackward(lastX - cursorCoordinate.X))
		} else if lastX < cursorCoordinate.X {
			c.writeString(terminalcode.CursorForward(cursorCoordinate.X - lastX))
		}
	}
	// 显示光标
	c.writeString(terminalcode.DisplayCursor)
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

// Print 输出字符串
func (c *CommandLine) Print(s string) {
	c.writeString(s)
	c.flush()
}

// Printf 输出格式化字符串
func (c *CommandLine) Printf(format string, a ...any) {
	c.writeString(fmt.Sprintf(format, a...))
	c.flush()
}

type Position struct {
	row int
	col int
}

func (c *CommandLine) getCursorPosition() Position {
	c.Print("\x1b[6n")
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

func NewCommandLine(tokensFunc lexer.GetTokensFunc, debug bool) (*CommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
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
	}, nil
}
