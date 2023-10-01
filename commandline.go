package startprompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

/*
功能核心类
*/

type AbortAction string

//goland:noinspection GoUnusedConst
const (
	// AbortActionUnspecific 空值
	AbortActionUnspecific AbortAction = ""
	// AbortActionIgnore 忽略此次操作
	AbortActionIgnore AbortAction = "ignore"
	// AbortActionRetry 让用户重新输入
	AbortActionRetry AbortAction = "retry"
	// AbortActionReturnError 返回错误，一般是 AbortError 或 ExitError
	AbortActionReturnError AbortAction = "return_error"
	// AbortActionReturnNone 返回空
	AbortActionReturnNone AbortAction = "return_none"
)

// AbortError 用户中断，一般是 Ctrl-C 触发
// ExitError 用户停止，一般是 Ctrl-D 触发
var AbortError = errors.New("user abort")
var ExitError = errors.New("user exit")

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// ctrlKey 返回与 Ctrl 一起按下的键
//
//goland:noinspection GoUnusedFunction
func ctrlKey(k rune) rune {
	return k & 0x1f
}

type PollEvent string

const (
	// PollEventInput 输入事件，用户按下键盘输入
	// PollEventRedraw 重画事件，重画当前输入
	// PollEventTimeout 超时事件，一段时间内没有其他事件触发
	PollEventInput   PollEvent = "input"
	PollEventRedraw  PollEvent = "redraw"
	PollEventTimeout PollEvent = "timeout"
)

// CommandLineOption 命令行选项
type CommandLineOption struct {
	// Schema token 样式（主要是颜色、加粗等）
	Schema Schema
	// Handler 事件处理器
	Handler EventHandler
	// History 输入历史存储
	History History
	// CodeFactory Code 类工厂方法
	CodeFactory CodeFactory
	// PromptFactory Prompt 类工厂方法
	PromptFactory PromptFactory

	// OnExit 用户停止时动作（Ctrl-D）
	OnExit AbortAction
	// OnAbort 用户中断时动作（Ctrl-C）
	OnAbort AbortAction

	// 自动缩进，如果开启，新行的缩进会与上一行保持一致
	AutoIndent bool
	// 开启 debug 日志
	EnableDebug bool
}

// defaultCommandLineOption 默认命令行配置
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

// copy 复制命令行配置，返回新的配置对象
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

// update 更新配置
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
	// 标准输入和输出的缓冲读写
	reader *bufio.Reader
	writer *bufio.Writer
	//    配置选项
	option *CommandLineOption
	//    下面几个都用用于并发的情况
	//    轮询超时时间
	pollTimeout time.Duration
	//    读取错误
	readError error
	//    重画和读取 channel
	redrawChannel chan rune
	//    传输读取的 rune
	readChannel chan rune
	//    是否正在读取用户输入
	isReadingInput bool
	//    下面几个对应用户的特殊操作：退出、丢弃、确定
	//    当操作发生时，对应的 flag 会设置为 true
	exitFlag   bool
	abortFlag  bool
	acceptFlag bool
	//    命令行当前使用的 Line 和 Renderer 对象
	line     *Line
	renderer *Renderer
}

// NewCommandLine 传入配置，新建命令行对象
func NewCommandLine(option *CommandLineOption) (*CommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}

	//     组合传入配置和默认配置
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

		redrawChannel: make(chan rune, 32),
		readChannel:   make(chan rune, 1024),
	}
	c.setup()
	DebugLog("start commandline")
	return c, nil
}

// setup 命令行初始化
func (c *CommandLine) setup() {
	c.reset()
	c.pollTimeout = 100 * time.Millisecond
	//    开启 debug log
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
				//    发生错误时，停止读取（如果调用方忽略 ReadInput 返回的错误，是否应该继续读取？）
				//    这个错误会由 ReadInput 判断并返回给调用方
				c.readError = err
				c.readChannel <- 0
				break
			}
			c.readChannel <- r
		}
	}()
}

// reset 重置 flag 和错误
func (c *CommandLine) reset() {
	c.exitFlag = false
	c.abortFlag = false
	c.acceptFlag = false
	c.readError = nil
}

// Close 关闭命令行，现在这个方法啥也没做
func (c *CommandLine) Close() {
	DebugLog("closed commandline")
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

// pollEvent 轮询事件
func (c *CommandLine) pollEvent() ([]rune, PollEvent) {
	select {
	case r := <-c.readChannel:
		rbuf := []rune{r}
		//    非阻塞的读取后续事件，优化粘贴大量文本的情况，快速处理，减少多次 render 导致的停顿感
		runeReading := true
		for runeReading {
			select {
			case r = <-c.readChannel:
				rbuf = append(rbuf, r)
			default:
				runeReading = false
			}
		}
		return rbuf, PollEventInput
	case <-c.redrawChannel:
		//    将缓冲的信息都读取出来，以免循环中不断触发
		//    或许加个重绘时间限制更好，比如 1s 只能重画 30 次？
		loop := len(c.redrawChannel)
		for i := 0; i < loop; i++ {
			<-c.redrawChannel
		}
		return nil, PollEventRedraw
	case <-time.After(c.pollTimeout):
		return nil, PollEventTimeout
	}
}

// ReadInput 读取用户输入
func (c *CommandLine) ReadInput() (string, error) {
	if c.isReadingInput {
		return "", fmt.Errorf("already reading input")
	}
	c.isReadingInput = true
	c.redrawChannel = make(chan rune, 1024)
	DebugLog("reading input")

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
	//    重置各个对象状态
	resetFunc()

	//    开启 terminal raw mode
	//    这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	//    不会自动展示用户输入
	//    更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		c.redrawChannel = nil
		c.isReadingInput = false
		return "", err
	}
	//    ReadInput 调用返回后，控制流程就到了用户，我们需要恢复终端的初始状态
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		if err != nil {
			fmt.Printf("term.Restore error: %v\r\n", err)
		}
	}()

	var inputText string
	for {
		//    轮询事件
		runes, pollEvent := c.pollEvent()
		switch pollEvent {
		case PollEventInput:
			if c.readError != nil {
				c.redrawChannel = nil
				c.isReadingInput = false
				return "", c.readError
			}
			DebugLog("read rune: [%d, ...] len=%d", runes[0], len(runes))
			//    识别用户输入，触发事件
			is.FeedRunes(runes)
		case PollEventTimeout:
			//    读取用户输入超时
			if !is.FeedTimeout() {
				//    没有触发事件，进入下一次循环，减少没必要的重画
				continue
			}
		}

		//    处理特别的输入事件结果
		if c.exitFlag {
			DebugLog("handle exit flag, action: %s", c.option.OnExit)
			//    一般是用户按了 Ctrl-D ，代表退出
			switch c.option.OnExit {
			case AbortActionReturnError:
				renderer.render(line.GetRenderContext(), true, false)
				c.redrawChannel = nil
				c.isReadingInput = false
				return "", ExitError
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
				c.redrawChannel = nil
				c.isReadingInput = false
				return "", nil
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if c.abortFlag {
			DebugLog("handle abort flag, action: %s", c.option.OnAbort)
			//    一般是用户按了 Ctrl-C ，代表中断
			switch c.option.OnAbort {
			case AbortActionReturnError:
				renderer.render(line.GetRenderContext(), true, false)
				c.redrawChannel = nil
				c.isReadingInput = false
				return "", AbortError
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
				c.redrawChannel = nil
				c.isReadingInput = false
				return "", nil
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if c.acceptFlag {
			DebugLog("handle accept flag")
			//    一般是用户按了 Enter ，代表完成本次输入
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
	DebugLog("return input: <%s>, err: nil", inputText)
	return inputText, nil
}

// ReadRune 读取 rune ，不能与 ReadInput 同时调用
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

// SetOnAbort 设置用户中断时的动作
func (c *CommandLine) SetOnAbort(action AbortAction) {
	c.option.OnAbort = action
}

// SetOnExit 设置用户停止时的动作
func (c *CommandLine) SetOnExit(action AbortAction) {
	c.option.OnExit = action
}

// SetExitFlag 设置停止标志
func (c *CommandLine) SetExitFlag() {
	c.exitFlag = true
}

// SetAbortFlag 设置中断标志
func (c *CommandLine) SetAbortFlag() {
	c.abortFlag = true
}

// SetAcceptFlag 设置（本次输入）完成标志
func (c *CommandLine) SetAcceptFlag() {
	c.acceptFlag = true
}

// IsReadingInput 是否正在读取输入
func (c *CommandLine) IsReadingInput() bool {
	return c.isReadingInput
}
