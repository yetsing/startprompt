package startprompt

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.design/x/clipboard"
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
	//    命令行当前使用的 Line 和 TRenderer 对象
	line     *Line
	renderer *TRenderer
	tscreen  tcell.Screen
	//    是否按下鼠标左键
	mousePrimaryPressed bool
	lastPrimaryEvent    *tcell.EventMouse
	lastDblclickEvent   *tcell.EventMouse
	//    点击间隔，用来判断鼠标双击、三击等
	clickInterval time.Duration
	//    配置选项
	option *CommandLineOption
	//    下面几个都用用于并发的情况
	//    传递输入 channel
	inputChannel chan *inputStruct
	//    传递输出 channel
	outputChannel chan outputStruct
	//    传递重新渲染事件 channel
	redrawChannel chan struct{}
	//    传递关闭事件 channel
	closeChannel chan struct{}
	//    下面两个用于 tcell.Screen ChannelEvents
	tEventChannel chan tcell.Event
	tQuitChannel  chan struct{}
	//    缓冲读取的 EventKey
	tEventKeyChannel chan *tcell.EventKey
	//    是否正在读取用户输入
	isReadingInput bool
	running        bool
	//   下面几个对应用户的特殊操作：退出、丢弃、确定
	exitFlag   bool
	abortFlag  bool
	acceptFlag bool
	//    wg 用来等待协程结束
	wg sync.WaitGroup
}

// NewTCommandLine 新建命令行对象
func NewTCommandLine(option *CommandLineOption) (*TCommandLine, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return nil, fmt.Errorf("not in a terminal")
	}

	//     Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		return nil, err
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
	s.SetCursorStyle(tcell.CursorStyleSteadyBar)

	c := &TCommandLine{
		tscreen: s,
		option:  actualOption,

		inputChannel:  make(chan *inputStruct),
		redrawChannel: make(chan struct{}, 16),
		closeChannel:  make(chan struct{}),
		outputChannel: make(chan outputStruct, 16),
		tEventChannel: make(chan tcell.Event, 16),
		tQuitChannel:  make(chan struct{}),

		renderer: newTRenderer(s, actualOption.Schema, actualOption.PromptFactory),
	}
	c.setup()
	return c, nil
}

func (tc *TCommandLine) setup() {
	tc.reset()
	tc.clickInterval = 200 * time.Millisecond
	if tc.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
	tc.wg.Add(1)
	tc.running = true
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
}

// Close 关闭命令行，恢复终端到原先的模式
func (tc *TCommandLine) Close() {
	maybePanic := recover()
	//    取消鼠标事件的上报（鼠标移动会上报大量事件，减少后面 discardTEvent 丢弃的事件）
	tc.tscreen.DisableMouse()
	tc.running = false
	close(tc.closeChannel)
	close(tc.redrawChannel)
	close(tc.outputChannel)
	tc.wg.Wait()
	//    继续读取事件
	//    因为 tcell 内部也是用 channel 传递的事件
	//    如果有大量事件上报，又没有读取事件，tcell 内部就可能卡在 channel send 发送
	//    从而导致调用 tcell.Screen.Fini 卡住（tcell 内部开了 goroutine 读取事件发送到 channel ， Fini 会等待这些 goroutine 结束）
	//    比如说开启鼠标支持，在 Close 前调用了 time.Sleep ，这个时候就会有大量的鼠标事件堆积
	tc.discardTEvent()
	tc.tscreen.Fini()
	if maybePanic != nil {
		DebugLog("panic: %v\n%s", maybePanic, string(debug.Stack()))
		panic(maybePanic)
	}
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

// ReadRune 读取 rune ，不能与 ReadInput 同时调用
func (tc *TCommandLine) ReadRune() (rune, error) {
	for event := range tc.tEventKeyChannel {
		if event.Key() == tcell.KeyRune {
			return event.Rune(), nil
		}
	}
	return 0, nil
}

func (tc *TCommandLine) run() {
	defer func() {
		maybePanic := recover()
		tc.wg.Done()
		if maybePanic != nil {
			DebugLog("panic: %v\n%s", maybePanic, string(debug.Stack()))
			panic(maybePanic)
		}
	}()
	for tc.running {
		tc.tEventKeyChannel = make(chan *tcell.EventKey, 1024)
		tc.runOther()
		close(tc.tEventKeyChannel)
		tc.runLoop()
	}
	tc.flushOutput()
	DebugLog("run stopped")
}

// runLoop 在用户调用 ReadInput 时执行
// 对应用户输入文本阶段
func (tc *TCommandLine) runLoop() {
	if !tc.running {
		return
	}
	DebugLog("runLoop")
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
		if len(tc.tEventKeyChannel) == 0 {
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
			}
		} else {
			for eventKey := range tc.tEventKeyChannel {
				if !tc.emitEvent(eventKey) {
					continue
				}
			}
		}

		//    处理特别的输入事件结果
		if tc.exitFlag {
			//    一般是用户按了 Ctrl-D
			switch tc.option.OnExit {
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

		renderer.update()
		//    画出用户输入
		renderer.render(line.GetRenderContext(), false, false)
	}
}

// runOther 在用户调用 ReadInput 前和返回后执行
// 对应非用户输入阶段，因为我们接管了鼠标操作，所以还要继续处理鼠标事件
func (tc *TCommandLine) runOther() {
	if !tc.running {
		return
	}
	DebugLog("runOther")
	renderer := tc.renderer

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
			renderer.Show()
			continue
		case ev := <-tc.tEventChannel:
			switch tev := ev.(type) {
			case *tcell.EventKey:
				//    这里用 select 监听写入，是为了防止阻塞在这里，ref: https://stackoverflow.com/a/25657232
				//    如果 tEventKeyChannel 满了，丢弃事件
				select {
				case tc.tEventKeyChannel <- tev:
				default:
				}
			default:
				//    没有触发事件，直接进入下一次循环，避免没必要的渲染
				if !tc.emitEvent(ev) {
					continue
				}
			}
		case output := <-tc.outputChannel:
			//    用户调用了 ReadInput
			if output.flush {
				return
			}
			renderer.renderOutput(output.text)
			continue
		}

		renderer.update()
		renderer.Show()
	}
}

func (tc *TCommandLine) emitEvent(tevent tcell.Event) bool {
	var event Event
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
			event = NewEventKey(eventType, data, nil, tc)
		} else {
			DebugLog("unsupported tcell.EventKey: %+v", ev)
		}
	case *tcell.EventMouse:
		x, y := ev.Position()
		coor := Coordinate{x, y}
		switch ev.Buttons() {
		case tcell.ButtonPrimary:
			if tc.mousePrimaryPressed {
				event = NewEventMouse(EventTypeMouseMove, coor, nil, tc)
			} else if tc.IsDblClick(ev) {
				event = NewEventMouse(EventTypeMouseDblclick, coor, nil, tc)
				tc.mousePrimaryPressed = false
				tc.lastPrimaryEvent = nil
				tc.lastDblclickEvent = ev
			} else if tc.IsTripleClick(ev) {
				event = NewEventMouse(EventTypeMouseTripleClick, coor, nil, tc)
				tc.mousePrimaryPressed = false
				tc.lastPrimaryEvent = nil
				tc.lastDblclickEvent = nil
			} else {
				event = NewEventMouse(EventTypeMouseDown, coor, nil, tc)
				tc.mousePrimaryPressed = true
				tc.lastPrimaryEvent = ev
			}
		case tcell.ButtonNone:
			if tc.mousePrimaryPressed {
				event = NewEventMouse(EventTypeMouseUp, coor, nil, tc)
				tc.mousePrimaryPressed = false
			}
		case tcell.WheelUp:
			event = NewEventMouse(EventTypeMouseWheelUp, coor, nil, tc)
		case tcell.WheelDown:
			event = NewEventMouse(EventTypeMouseWheelDown, coor, nil, tc)
		}
	}
	if event == nil {
		return false
	}
	DebugLog("emit event=%s", event.Type())
	tc.option.Handler.Handle(event)
	return true
}

func (tc *TCommandLine) IsDblClick(ev *tcell.EventMouse) bool {
	//    clickInterval 内相同位置按下鼠标左键两次
	if tc.lastPrimaryEvent == nil {
		return false
	}
	lastX, lastY := tc.lastPrimaryEvent.Position()
	x, y := ev.Position()
	return lastX == x && lastY == y && tc.lastPrimaryEvent.When().Add(tc.clickInterval).After(time.Now())
}

func (tc *TCommandLine) IsTripleClick(ev *tcell.EventMouse) bool {
	//    clickInterval 内相同位置按下鼠标左键三次
	if tc.lastDblclickEvent == nil {
		return false
	}
	lastX, lastY := tc.lastDblclickEvent.Position()
	x, y := ev.Position()
	return lastX == x && lastY == y && tc.lastDblclickEvent.When().Add(tc.clickInterval).After(time.Now())
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
	//    这是使用超时，而不是等待 channel 关闭
	//    是因为我们要在 tcell.Screen.Fini 之前调用这个方法，丢弃掉 tcell 内部 channel 堆积的事件
	//    调用 tcell.Screen.Fini 之后， tEventChannel 会被关闭，这样我们就无法读取到 tcell 内部 channel 堆积的事件
	timer := time.NewTimer(100 * time.Millisecond)
	for {
		select {
		case <-tc.tEventChannel:
		case <-timer.C:
			return
		}
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
