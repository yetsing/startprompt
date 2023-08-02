package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"

	"startprompt/inputstream"
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

type CommandLine struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c *CommandLine) ReadInput() string {
	line := inputstream.NewLine()
	handler := inputstream.NewBaseInputStreamHandler(line)
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
		if r == ctrlKey('q') {
			break
		}
		if line.Finished() {
			break
		}
		if reader.Buffered() == 0 {
			c.draw(line.Document())
		}
	}
	return line.Text()
}

func (c *CommandLine) draw(doc *inputstream.Document) {
	fmt.Printf("buffer: %d\n", c.reader.Buffered())
	text := doc.Text
	//pos := c.getCursorPosition()
	//pos.col = doc.CursorPosition + 1
	// 隐藏光标
	c.writeString("\x1b[?25l")
	// 删除整行
	c.writeString("\x1b[2K")
	//c.writeString(fmt.Sprintf("\x1b[%d;0H", pos.row))
	c.writeString(text)
	// 定位光标
	//c.writeString(fmt.Sprintf("\x1b[%d;%dH", pos.row, pos.col))
	// 显示光标
	c.writeString("\x1b[?25h")
	c.flush()
}

func (c *CommandLine) writeString(s string) {
	_, err := c.writer.WriteString(s)
	ensureOk(err)
}

func (c *CommandLine) flushString(s string) {
	c.writeString(s)
	c.flush()
}

func (c *CommandLine) flush() {
	err := c.writer.Flush()
	ensureOk(err)
}

type Position struct {
	row int
	col int
}

func (c *CommandLine) getCursorPosition() Position {
	c.flushString("\x1b[6n")
	var buf [32]byte
	var i int
	// 回复格式为 \x1b[A;BR
	// A 和 B 就是光标的行和列
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

func NewCommandLine() *CommandLine {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	return &CommandLine{reader: reader, writer: writer}
}

func main() {
	// 开启 terminal raw mode
	// 这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	// 不会自动展示用户输入
	// 更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("error make raw mode: %v\n", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error: \n", r)
		}
	}()

	defer func(fd int, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {
			fmt.Printf("error restore: %v\n", err)
		}
	}(int(os.Stdin.Fd()), oldState)

	c := NewCommandLine()
	c.ReadInput()
}
