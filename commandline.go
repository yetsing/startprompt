package startprompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/yetsing/startprompt/terminalcode"
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
	NewCodeFunc   NewCodeFunc
	NewPromptFunc NewPromptFunc

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
	NewCodeFunc:   newBaseCode,
	NewPromptFunc: newBasePrompt,
	OnAbort:       AbortActionRetry,
	OnExit:        AbortActionReturnError,
	AutoIndent:    false,
	EnableDebug:   false,
}

func (cp *CommandLineOption) copy() *CommandLineOption {
	return &CommandLineOption{
		Schema:        cp.Schema,
		History:       cp.History,
		NewCodeFunc:   cp.NewCodeFunc,
		NewPromptFunc: cp.NewPromptFunc,
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
	if other.NewCodeFunc != nil {
		cp.NewCodeFunc = other.NewCodeFunc
	}
	if other.NewPromptFunc != nil {
		cp.NewPromptFunc = other.NewPromptFunc
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

	option *CommandLineOption

	inputTimeout time.Duration

	readError error

	redrawChannel chan rune
	readChannel   chan rune
}

// RequestRedraw 请求重绘（线程安全）
func (c *CommandLine) RequestRedraw() {
	if !c.option.EnableConcurrency {
		panic("not enable concurrency")
	}
	c.redrawChannel <- 'x'
}

// RunInExecutor 运行后台任务
func (c *CommandLine) RunInExecutor(callback func()) {
	if !c.option.EnableConcurrency {
		panic("not enable concurrency")
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

	render := newRender(c.option.Schema)
	line := newLine(c.option.NewCodeFunc, c.option.NewPromptFunc, c.option.History, c.option.AutoIndent)
	handler := NewBaseHandler(line)
	is := NewInputStream(handler)
	render.render(line.GetRenderContext())

	var r rune
	reader := c.reader
	var inputText string
	for {
		// 用户有多个任务运行，不能一直阻塞在用户输入上
		if c.option.EnableConcurrency {
			var pollEvent PollEvent
			r, pollEvent = c.pollEvent()
			switch pollEvent {
			case PollEventRedraw:
				render.render(line.GetRenderContext())
				continue
			case PollEventTimeout:
				c.OnInputTimeout(line.CreateCodeObj())
				continue
			}
			if c.readError != nil {
				return "", err
			}
		} else {
			r, _, err = reader.ReadRune()
			if err != nil {
				return "", err
			}
		}
		DebugLog("read rune: %d", r)
		// 识别用户输入，触发事件
		is.Feed(r)
		//c.draw(line)
		if line.exit {
			// 一般是用户按了 Ctrl-D
			if c.option.OnExit != AbortActionIgnore {
				render.render(line.GetRenderContext())
			}

			switch c.option.OnExit {
			case AbortActionReturnError:
				return "", ExitError
			case AbortActionReturnNone:
				return "", nil
			case AbortActionRetry:
				line.reset()
			case AbortActionIgnore:

			}
		} else if line.abort {
			// 一般是用户按了 Ctrl-C
			if c.option.OnAbort != AbortActionIgnore {
				render.render(line.GetRenderContext())
			}

			switch c.option.OnAbort {
			case AbortActionReturnError:
				return "", AbortError
			case AbortActionReturnNone:
				return "", nil
			case AbortActionRetry:
				line.reset()
			case AbortActionIgnore:

			}
		} else if line.accept {
			// 一般是用户按了 Enter
			render.render(line.GetRenderContext())
			inputText = line.text()
			break
		}
		switch line.renderType {
		case LineRenderClear:
			// 一般是用户按了 Ctrl-L
			render.clear()
		case LineRenderListCompletion:
			// 一般是用户按了两次 tab （实际上没有支持这个效果）
			render.renderCompletions(line.GetRenderCompletions())
		}
		// 重新画出用户输入
		render.render(line.GetRenderContext())
		line.ResetRenderType()
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

func (c *CommandLine) SetOnAbort(action AbortAction) {
	c.option.OnAbort = action
}

func (c *CommandLine) SetOnExit(action AbortAction) {
	c.option.OnExit = action
}

func (c *CommandLine) setup() {
	if c.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	if c.option.EnableConcurrency {
		// 新开协程读取用户输入
		go func() {
			for {
				r, _, err := c.reader.ReadRune()
				if err != nil {
					c.readError = err
					break
				}
				c.readChannel <- r
			}
		}()
	}
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
