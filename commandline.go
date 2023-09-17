package startprompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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

func panicIfError(err error) {
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
	Handler       EventHandler
	History       History
	CodeFactory   CodeFactory
	PromptFactory PromptFactory

	OnExit  AbortAction
	OnAbort AbortAction

	// 自动缩进，如果开启，新行的缩进会与上一行保持一致
	AutoIndent bool
	// 开启 debug 日志
	EnableDebug bool
}

var defaultCommandLineOption = &CommandLineOption{
	Schema:        defaultSchema,
	Handler:       newBaseHandler(),
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
		Handler:       cp.Handler,
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
	if other.Handler != nil {
		cp.Handler = other.Handler
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
	acceptFlag bool
	//    命令行当前使用的 Line 和 Render 对象
	line     *Line
	renderer *Renderer
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
		readChannel:   make(chan rune, 16),
	}
	c.setup()
	return c, nil
}

func (c *CommandLine) setup() {
	c.reset()
	c.inputTimeout = 100 * time.Millisecond
	if c.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
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

func (c *CommandLine) reset() {
	c.exitFlag = false
	c.abortFlag = false
	c.acceptFlag = false
	c.readError = nil
}

func (c *CommandLine) Close() {

}

// RequestRedraw 请求重绘（ goroutine 安全）
func (c *CommandLine) RequestRedraw() {
	if c.redrawChannel != nil {
		c.redrawChannel <- 'x'
	}
}

// RunInExecutor 运行后台任务
func (c *CommandLine) RunInExecutor(callback func()) {
	go callback()
}

func (c *CommandLine) pollEvent() (rune, PollEvent) {
	select {
	case r := <-c.readChannel:
		return r, PollEventInput
	case <-c.redrawChannel:
		//    将缓冲的信息都读取出来，以免循环中不断触发
		//    或许加个重绘时间限制更好，比如 1s 只能重画 30 次？
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

	renderer := newRenderer(c.option.Schema, c.option.PromptFactory)
	c.renderer = renderer
	line := newLine(
		c.option.CodeFactory,
		c.option.History,
		c.option.AutoIndent,
	)
	c.line = line
	handler := c.option.Handler
	is := NewInputStream(handler, c)
	renderer.render(line.GetRenderContext(), false, false)

	resetFunc := func() {
		is.Reset()
		line.reset()
		renderer.reset()
		c.reset()
	}

	resetFunc()

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
	var inputText string
	for {
		//    读取用户输入
		var pollEvent PollEvent
		r, pollEvent = c.pollEvent()
		switch pollEvent {
		case PollEventInput:
			if c.readError != nil {
				return "", c.readError
			}
			DebugLog("read rune: %d", r)
			//    识别用户输入，触发事件
			is.Feed(r)
		case PollEventTimeout:
			//    读取用户输入超时
			if !is.FeedTimeout() {
				//    没有触发事件，进入下一次循环，减少没必要的重画
				continue
			}
		}

		//    处理特别的输入事件结果
		if c.exitFlag {
			//    一般是用户按了 Ctrl-D
			switch c.option.OnExit {
			case AbortActionReturnError:
				renderer.render(line.GetRenderContext(), true, false)
				return "", ExitError
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
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
				renderer.render(line.GetRenderContext(), true, false)
				return "", AbortError
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
				return "", nil
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if c.acceptFlag {
			//    一般是用户按了 Enter
			renderer.render(line.GetRenderContext(), false, true)
			inputText = line.text()
			break
		}

		//    画出用户输入
		renderer.render(line.GetRenderContext(), false, false)
	}
	//    返回用户输入的文本内容
	c.redrawChannel = nil
	c.isReadingInput = false
	DebugLog("return input: <%s>", inputText)
	return inputText, nil
}

func (c *CommandLine) ReadRune() (rune, error) {
	if c.readError != nil {
		return 0, c.readError
	}
	return <-c.readChannel, nil
}

// GetLine 获取当前的 Line 对象，如果为 nil ，则 panic
func (c *CommandLine) GetLine() *Line {
	if c.line == nil {
		panic("not found Line from CommandLine")
	}
	return c.line
}

// GetRenderer 获取当前的 Renderer 对象，如果为 nil ，则 panic
func (c *CommandLine) GetRenderer() *Renderer {
	if c.renderer == nil {
		panic("not found Renderer from CommandLine")
	}
	return c.renderer
}

// Print 输出字符串，类似 fmt.Print
func (c *CommandLine) Print(a ...any) {
	_, err := fmt.Fprint(c.writer, a...)
	panicIfError(err)
	c.flush()
}

// Printf 输出格式化字符串，类似 fmt.Printf
func (c *CommandLine) Printf(format string, a ...any) {
	_, err := fmt.Fprintf(c.writer, format, a...)
	panicIfError(err)
	c.flush()
}

// Println 输出一行字符串，类似 fmt.Println
func (c *CommandLine) Println(a ...any) {
	_, err := fmt.Fprintln(c.writer, a...)
	panicIfError(err)
	c.flush()
}

// flush 写入缓冲数据
func (c *CommandLine) flush() {
	err := c.writer.Flush()
	panicIfError(err)
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

func (c *CommandLine) SetAccept() {
	c.acceptFlag = true
}

func (c *CommandLine) IsReadingInput() bool {
	return c.isReadingInput
}
