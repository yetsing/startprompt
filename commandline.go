package startprompt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/term"
)

type AbortAction string

//goland:noinspection GoUnusedConst
const (
	AbortActionUnspecific  AbortAction = ""
	AbortActionIgnore      AbortAction = "ignore"
	AbortActionRetry       AbortAction = "retry"
	AbortActionReturnError AbortAction = "return_error"
	AbortActionReturnNone  AbortAction = "return_none"
)

var AbortError = errors.New("user abort")
var ExitError = errors.New("user exit")

func ensureOk(err error) {
	if err != nil {
		panic(err)
	}
}

//goland:noinspection GoUnusedFunction
func ctrlKey(k rune) rune {
	return k & 0x1f
}

type PollEvent string

const (
	PollEventInput   PollEvent = "input"
	PollEventRedraw  PollEvent = "redraw"
	PollEventTimeout PollEvent = "timeout"
)

type CommandLineOption struct {
	Schema        Schema
	History       History
	CodeFactory   CodeFactory
	PromptFactory PromptFactory

	OnExit  AbortAction
	OnAbort AbortAction

	// 自动缩进，如果开启，新行的缩进会与上一行保持一致
	AutoIndent bool
	// 开启 debug 日志
	EnableDebug       bool
	EnableConcurrency bool
}

var defaultCommandLineOption = &CommandLineOption{
	Schema:        defaultSchema,
	History:       NewMemHistory(),
	CodeFactory:   newBaseCode,
	PromptFactory: newBasePrompt,
	OnAbort:       AbortActionRetry,
	OnExit:        AbortActionReturnError,
	AutoIndent:    false,
	EnableDebug:   false,
}

func (cp *CommandLineOption) copy() *CommandLineOption {
	return &CommandLineOption{
		Schema:        cp.Schema,
		History:       cp.History,
		CodeFactory:   cp.CodeFactory,
		PromptFactory: cp.PromptFactory,
		OnAbort:       cp.OnAbort,
		OnExit:        cp.OnExit,
		AutoIndent:    cp.AutoIndent,
		EnableDebug:   cp.EnableDebug,
	}
}

func (cp *CommandLineOption) update(other *CommandLineOption) {
	if other.Schema != nil {
		cp.Schema = other.Schema
	}
	if other.History != nil {
		cp.History = other.History
	}
	if other.CodeFactory != nil {
		cp.CodeFactory = other.CodeFactory
	}
	if other.PromptFactory != nil {
		cp.PromptFactory = other.PromptFactory
	}
	if other.OnExit != AbortActionUnspecific {
		cp.OnExit = other.OnExit
	}
	if other.OnAbort != AbortActionUnspecific {
		cp.OnAbort = other.OnAbort
	}
	cp.AutoIndent = other.AutoIndent
	cp.EnableDebug = other.EnableDebug
	cp.EnableConcurrency = other.EnableConcurrency
}

type CommandLine struct {
	reader *bufio.Reader
	writer *bufio.Writer

	//    配置选项
	option *CommandLineOption

	//    下面几个都用用于并发的情况
	//    等待输入超时时间
	inputTimeout time.Duration
	//    读取错误
	readError error
	//     重画和读取 channel
	redrawChannel chan rune
	readChannel   chan rune
	//    是否正在读取用户输入
	isReadingInput bool

	//   下面几个对应用户的特殊操作：退出、丢弃、确定
	exitFlag   bool
	abortFlag  bool
	returnCode Code
}

func (c *CommandLine) setup() {
	if c.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	if c.option.EnableConcurrency {
		//    新开协程读取用户输入
		go func() {
			for {
				r, _, err := c.reader.ReadRune()
				if err != nil {
					c.readError = err
					c.readChannel <- 0
					break
				}
				c.readChannel <- r
			}
		}()
	}
	c.reset()
}

func NewCommandLine(option *CommandLineOption) (*CommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}

	actualOption := defaultCommandLineOption.copy()
	if option != nil {
		actualOption.update(option)
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	c := &CommandLine{
		reader: reader,
		writer: writer,
		option: actualOption,

		redrawChannel: make(chan rune, 1024),
		readChannel:   make(chan rune),
	}
	c.setup()
	return c, nil
}

func (c *CommandLine) reset() {
	c.exitFlag = false
	c.abortFlag = false
	c.returnCode = nil
}

// RequestRedraw 请求重绘（线程安全）
func (c *CommandLine) RequestRedraw() {
	if !c.option.EnableConcurrency {
		panic("Must enable concurrency")
	}
	if c.redrawChannel != nil {
		c.redrawChannel <- 'x'
	}
}

// RunInExecutor 运行后台任务
func (c *CommandLine) RunInExecutor(callback func()) {
	if !c.option.EnableConcurrency {
		panic("Must enable concurrency")
	}

	go callback()
}

// OnInputTimeout 在等待输入超时后调用
func (c *CommandLine) OnInputTimeout(_ Code) {

}

func (c *CommandLine) OnReadInputStart() {

}

func (c *CommandLine) OnReadInputEnd() {

}

func (c *CommandLine) pollEvent() (rune, PollEvent) {
	select {
	case r := <-c.readChannel:
		return r, PollEventInput
	case <-c.redrawChannel:
		// 将缓冲的信息都读取出来，以免循环中不断触发
		loop := len(c.redrawChannel)
		for i := 0; i < loop; i++ {
			<-c.redrawChannel
		}
		return 0, PollEventRedraw
	case <-time.After(c.inputTimeout):
		return 0, PollEventTimeout
	}
}

func (c *CommandLine) ReadInput() (string, error) {
	if c.isReadingInput {
		return "", fmt.Errorf("already reading input")
	}
	c.isReadingInput = true
	c.redrawChannel = make(chan rune, 1024)

	render := newRender(c.option.Schema)
	line := newLine(
		c.option.CodeFactory,
		c.option.PromptFactory,
		c.option.History,
		newLineCallbacks(c, render),
		c.option.AutoIndent,
	)
	handler := NewBaseHandler(line)
	is := NewInputStream(handler)
	render.render(line.GetRenderContext(), false, false)

	resetFunc := func() {
		line.reset()
		c.reset()
	}

	resetFunc()
	c.OnReadInputStart()

	//    开启 terminal raw mode
	//    这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	//    不会自动展示用户输入
	//    更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
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

	var r rune
	reader := c.reader
	var inputText string
	for {
		//    读取用户输入
		if c.option.EnableConcurrency {
			//    用户有多个任务运行，不能一直阻塞在用户输入上
			var pollEvent PollEvent
			r, pollEvent = c.pollEvent()
			switch pollEvent {
			case PollEventInput:
				if c.readError != nil {
					return "", err
				}
				DebugLog("read rune: %d", r)
				//    识别用户输入，触发事件
				is.Feed(r)
			case PollEventTimeout:
				//    读取用户输入超时
				c.OnInputTimeout(line.CreateCode())
				continue
			}
		} else {
			r, _, err = reader.ReadRune()
			if err != nil {
				return "", err
			}
			DebugLog("read rune: %d", r)
			//    识别用户输入，触发事件
			is.Feed(r)
		}

		//    处理特别的输入事件结果
		if c.exitFlag {
			//    一般是用户按了 Ctrl-D
			switch c.option.OnExit {
			case AbortActionReturnError:
				render.render(line.GetRenderContext(), true, false)
				return "", ExitError
			case AbortActionReturnNone:
				render.render(line.GetRenderContext(), true, false)
				return "", nil
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if c.abortFlag {
			//    一般是用户按了 Ctrl-C
			switch c.option.OnAbort {
			case AbortActionReturnError:
				render.render(line.GetRenderContext(), true, false)
				return "", AbortError
			case AbortActionReturnNone:
				render.render(line.GetRenderContext(), true, false)
				return "", nil
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if c.returnCode != nil {
			//    一般是用户按了 Enter
			render.render(line.GetRenderContext(), false, true)
			inputText = line.text()
			break
		}

		//    画出用户输入
		render.render(line.GetRenderContext(), false, false)
	}
	c.redrawChannel = nil
	c.OnReadInputEnd()
	c.isReadingInput = false
	DebugLog("return input: <%s>", inputText)
	return inputText, nil
}

func (c *CommandLine) getCursorPosition() (int, int) {
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
		panic(fmt.Sprintf("invalid cursor position report escape: %v", buf[:i]))
	}
	sepIndex := bytes.IndexByte(buf[:], ';')
	if sepIndex == -1 {
		panic(fmt.Sprintf("invalid cursor position report separator: %v", buf[:i]))
	}
	row, err := strconv.ParseInt(string(buf[2:sepIndex]), 10, 32)
	ensureOk(err)
	col, err := strconv.ParseInt(string(buf[sepIndex+1:i]), 10, 32)
	ensureOk(err)
	return int(row), int(col)
}

// writeString 写入字符串到输出缓冲中
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

func (c *CommandLine) SetOnAbort(action AbortAction) {
	c.option.OnAbort = action
}

func (c *CommandLine) SetOnExit(action AbortAction) {
	c.option.OnExit = action
}

func (c *CommandLine) SetExit() {
	c.exitFlag = true
}

func (c *CommandLine) SetAbort() {
	c.abortFlag = true
}

func (c *CommandLine) SetReturnValue(code Code) {
	c.returnCode = code
}

type LineCallbacks struct {
	commandLine *CommandLine
	render      *rRenderer
}

func newLineCallbacks(commandLine *CommandLine, render *rRenderer) *LineCallbacks {
	return &LineCallbacks{commandLine: commandLine, render: render}
}

func (l *LineCallbacks) ClearScreen() {
	l.render.clear()
}

func (l *LineCallbacks) Exit() {
	l.commandLine.SetExit()
}

func (l *LineCallbacks) Abort() {
	l.commandLine.SetAbort()
}

func (l *LineCallbacks) ReturnInput(code Code) {
	l.commandLine.SetReturnValue(code)
}
