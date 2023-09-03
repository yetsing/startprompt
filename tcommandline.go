package startprompt

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type inputStruct struct {
	text string
	err  error
}

type outputStruct struct {
	text  string
	flush bool
}

var defaultTCommandLineOption = &CommandLineOption{
	Schema:        defaultSchema,
	Handler:       newTBaseEventHandler(),
	History:       NewMemHistory(),
	CodeFactory:   newBaseCode,
	PromptFactory: newBasePrompt,
	OnAbort:       AbortActionRetry,
	OnExit:        AbortActionReturnError,
	AutoIndent:    false,
	EnableDebug:   false,
}

type TCommandLine struct {
	tscreen tcell.Screen
	//    配置选项
	option *CommandLineOption
	//    下面几个都用用于并发的情况
	//    等待输入超时时间
	inputTimeout time.Duration
	//    读取错误
	readError error
	//    输入 channel
	inputChannel  chan *inputStruct
	outputChannel chan outputStruct
	redrawChannel chan struct{}
	closeChannel  chan struct{}
	//    下面两个用于 tcell.Screen ChannelEvents
	tEventChannel chan tcell.Event
	tQuitChannel  chan struct{}
	//    是否正在读取用户输入
	isReadingInput bool
	running        bool
	//   下面几个对应用户的特殊操作：退出、丢弃、确定
	exitFlag   bool
	abortFlag  bool
	acceptFlag bool
	//    命令行当前使用的 Line 和 TRenderer 对象
	line     *Line
	renderer *TRenderer
}

// NewTCommandLine 新建命令行对象
func NewTCommandLine(option *CommandLineOption) (*TCommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}

	//    update option default
	actualOption := defaultTCommandLineOption.copy()
	if option != nil {
		actualOption.update(option)
	}

	//     Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to NewScreen: %w", err)
	}
	if err := s.Init(); err != nil {
		return nil, fmt.Errorf("failed to Screen.Init: %w", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	c := &TCommandLine{
		tscreen: s,
		option:  actualOption,

		inputChannel:  make(chan *inputStruct),
		redrawChannel: make(chan struct{}, 1024),
		closeChannel:  make(chan struct{}),
		outputChannel: make(chan outputStruct, 8),
		tEventChannel: make(chan tcell.Event),
		tQuitChannel:  make(chan struct{}),

		renderer: newTRenderer(s, actualOption.Schema, actualOption.PromptFactory),
	}
	c.setup()
	return c, nil
}

func (tc *TCommandLine) setup() {
	tc.reset()
	if tc.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	//    根据 ChannelEvents 注释，需要单独开 goroutine 调用
	go func() {
		tc.tscreen.ChannelEvents(tc.tEventChannel, tc.tQuitChannel)
	}()
	//    ReadInput 返回后，控制权转移到了调用者手里，此时我们处理不了事件
	//    因为我们接管了鼠标操作，所以需要后台响应鼠标事件
	//    干脆把整个事件的处理都放到后台
	go tc.run()
}

func (tc *TCommandLine) reset() {
	tc.exitFlag = false
	tc.abortFlag = false
	tc.acceptFlag = false
	tc.readError = nil
}

// Close 关闭命令行，恢复终端到原先的模式
func (tc *TCommandLine) Close() {
	tc.running = false
	close(tc.redrawChannel)
	close(tc.outputChannel)
	close(tc.closeChannel)
	close(tc.tQuitChannel)
	tc.tscreen.Fini()
}

// RequestRedraw 请求重绘（ goroutine 安全）
func (tc *TCommandLine) RequestRedraw() {
	if tc.isReadingInput {
		tc.redrawChannel <- struct{}{}
	}
}

// RunInExecutor 运行后台任务
func (tc *TCommandLine) RunInExecutor(callback func()) {
	go callback()
}

// ReadInput 读取当前输入
func (tc *TCommandLine) ReadInput() (string, error) {
	if tc.isReadingInput {
		return "", fmt.Errorf("already reading input")
	}
	tc.isReadingInput = true

	tc.outputChannel <- outputStruct{"", true}
	in := <-tc.inputChannel

	tc.isReadingInput = false
	return in.text, in.err
}

func (tc *TCommandLine) run() {
	tc.running = true
	tc.flushOutput()
	for tc.running {
		tc.runLoop()
		tc.flushOutput()
	}
	tc.flushOutput()
	//    如果有事件没有处理，tcell.Screen.Fini 会卡在那里
	tc.discardTEvent()
}

func (tc *TCommandLine) runLoop() {
	renderer := tc.renderer
	line := newLine(
		tc.option.CodeFactory,
		tc.option.History,
		tc.option.AutoIndent,
	)
	tc.line = line

	renderer.render(line.GetRenderContext(), false, false)

	resetFunc := func() {
		line.reset()
		renderer.reset()
		tc.reset()
	}

	resetFunc()

	for {
		select {
		case <-tc.closeChannel:
			DebugLog("close")
			return
		case <-tc.redrawChannel:
			DebugLog("redraw")
			//    将缓冲的信息都读取出来，以免循环中不断触发
			loop := len(tc.redrawChannel)
			for i := 0; i < loop; i++ {
				<-tc.redrawChannel
			}
			//    渲染用户输入
			renderer.render(line.GetRenderContext(), false, false)
			continue
		case ev := <-tc.tEventChannel:
			//    没有触发事件，直接进入下一次循环，避免没必要的渲染
			if !tc.emitEvent(ev) {
				continue
			}
			DebugLog("emit event: %+v", ev)
		}

		//    处理特别的输入事件结果
		if tc.exitFlag {
			//    一般是用户按了 Ctrl-D
			switch tc.option.OnExit {
			case AbortActionReturnError:
				renderer.render(line.GetRenderContext(), true, false)
				return
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
				return
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if tc.abortFlag {
			//    一般是用户按了 Ctrl-C
			switch tc.option.OnAbort {
			case AbortActionReturnError:
				renderer.render(line.GetRenderContext(), true, false)
				tc.sendInput("", ExitError)
				return
			case AbortActionReturnNone:
				renderer.render(line.GetRenderContext(), true, false)
				tc.sendInput("", nil)
				return
			case AbortActionRetry:
				resetFunc()
			case AbortActionIgnore:

			}
		}
		if tc.acceptFlag {
			//    一般是用户按了 Enter
			//    返回用户输入的文本内容
			renderer.render(line.GetRenderContext(), false, true)
			inputText := line.text()
			DebugLog("return input: <%s>", inputText)
			tc.sendInput(inputText, nil)
			break
		}

		//    画出用户输入
		renderer.render(line.GetRenderContext(), false, false)
	}
}

func (tc *TCommandLine) emitEvent(tevent tcell.Event) bool {
	switch ev := tevent.(type) {
	case *tcell.EventResize:
		tc.renderer.Resize()
	case *tcell.EventKey:
		eventType, found := tkeyMapping[ev.Key()]
		if found {
			var data []rune
			if ev.Key() == tcell.KeyRune {
				data = []rune{ev.Rune()}
			}
			event := NewEventKey(eventType, data, nil, tc)
			tc.option.Handler.Handle(event)
			return true
		} else {
			DebugLog("unsupported tcell.EventKey: %+v", ev)
		}
	}
	return false
}

func (tc *TCommandLine) sendInput(text string, err error) {
	in := &inputStruct{text: text, err: err}
	tc.inputChannel <- in
}

// GetLine 获取当前的 Line 对象，如果为 nil ，则 panic
func (tc *TCommandLine) GetLine() *Line {
	if tc.line == nil {
		panic("not found Line from TCommandLine")
	}
	return tc.line
}

// GetRenderer 获取当前的 TRenderer 对象，如果为 nil ，则 panic
func (tc *TCommandLine) GetRenderer() *TRenderer {
	if tc.renderer == nil {
		panic("not found TRenderer from TCommandLine")
	}
	return tc.renderer
}

func (tc *TCommandLine) flushOutput() {
	ch := tc.outputChannel
	for output := range ch {
		if output.flush {
			break
		}
		tc.renderer.renderOutput(output.text)
	}
}

func (tc *TCommandLine) discardTEvent() {
	for range tc.tEventChannel {

	}
}

func (tc *TCommandLine) Write(p []byte) (int, error) {
	tc.outputChannel <- outputStruct{string(p), false}
	return len(p), nil
}

func (tc *TCommandLine) Print(a ...any) {
	_, _ = fmt.Fprint(tc, a...)
}

func (tc *TCommandLine) Println(a ...any) {
	_, _ = fmt.Fprintln(tc, a...)
}

func (tc *TCommandLine) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(tc, format, a...)
}

func (tc *TCommandLine) SetOnAbort(action AbortAction) {
	tc.option.OnAbort = action
}

func (tc *TCommandLine) SetOnExit(action AbortAction) {
	tc.option.OnExit = action
}

func (tc *TCommandLine) SetExit() {
	tc.exitFlag = true
}

func (tc *TCommandLine) SetAbort() {
	tc.abortFlag = true
}

func (tc *TCommandLine) SetAccept() {
	tc.acceptFlag = true
}

func (tc *TCommandLine) IsReadingInput() bool {
	return tc.isReadingInput
}
