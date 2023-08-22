package startprompt

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/yetsing/startprompt/terminalcode"
	"golang.org/x/term"
	"os"
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

type CommandLineOption struct {
	Schema        Schema
	History       History
	NewCodeFunc   NewCodeFunc
	NewPromptFunc NewPromptFunc

	OnExit  AbortAction
	OnAbort AbortAction

	Debug bool
}

var defaultCommandLineOption = &CommandLineOption{
	Schema:        defaultSchema,
	History:       NewMemHistory(),
	NewCodeFunc:   newBaseCode,
	NewPromptFunc: newBasePrompt,
	Debug:         false,
	OnAbort:       AbortActionRetry,
	OnExit:        AbortActionReturnError,
}

func (cp *CommandLineOption) copy() *CommandLineOption {
	return &CommandLineOption{
		Schema:        cp.Schema,
		History:       cp.History,
		NewCodeFunc:   cp.NewCodeFunc,
		NewPromptFunc: cp.NewPromptFunc,
		Debug:         cp.Debug,
		OnAbort:       cp.OnAbort,
		OnExit:        cp.OnExit,
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
	cp.Debug = other.Debug
	if other.OnExit != AbortActionUnspecific {
		cp.OnExit = other.OnExit
	}
	if other.OnAbort != AbortActionUnspecific {
		cp.OnAbort = other.OnAbort
	}
}

type CommandLine struct {
	reader *bufio.Reader
	writer *bufio.Writer

	option *CommandLineOption

	enableDebugLog bool

	onAbort AbortAction
	onExit  AbortAction
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
	line := newLine(render, c.option.NewCodeFunc, c.option.NewPromptFunc, c.option.History)
	handler := NewBaseHandler(line)
	is := NewInputStream(handler)
	render.render(line.GetRenderContext())

	var r rune
	reader := c.reader
	var inputText string
	for {
		r, _, err = reader.ReadRune()
		if err != nil {
			return "", err
		}
		DebugLog("read rune: %d", r)
		is.Feed(r)
		//c.draw(line)
		if line.exit {
			if c.onExit != AbortActionIgnore {
				render.render(line.GetRenderContext())
			}

			switch c.onExit {
			case AbortActionReturnError:
				return "", ExitError
			case AbortActionReturnNone:
				return "", nil
			case AbortActionRetry:
				line.reset()
			case AbortActionIgnore:

			}
		} else if line.abort {
			if c.onAbort != AbortActionIgnore {
				render.render(line.GetRenderContext())
			}

			switch c.onAbort {
			case AbortActionReturnError:
				return "", AbortError
			case AbortActionReturnNone:
				return "", nil
			case AbortActionRetry:
				line.reset()
			case AbortActionIgnore:

			}
		} else if line.accept {
			render.render(line.GetRenderContext())
			inputText = line.text()
			break
		}
		render.render(line.GetRenderContext())
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

func (c *CommandLine) OnAbort(action AbortAction) {
	c.onAbort = action
}

func (c *CommandLine) OnExit(action AbortAction) {
	c.onExit = action
}

func NewCommandLine(option *CommandLineOption) (*CommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}

	var finalOption *CommandLineOption
	if option != nil {
		finalOption = defaultCommandLineOption.copy()
		finalOption.update(option)
	} else {
		finalOption = defaultCommandLineOption
	}
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	if finalOption.Debug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	return &CommandLine{
		reader: reader,
		writer: writer,
		option: finalOption,

		onAbort: finalOption.OnAbort,
		onExit:  finalOption.OnExit,
	}, nil
}
