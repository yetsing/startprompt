package startprompt

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type TCommandLine struct {
	screen tcell.Screen
	//    配置选项
	option *CommandLineOption
	//    下面几个都用用于并发的情况
	//    等待输入超时时间
	inputTimeout time.Duration
	//    读取错误
	readError error
	//     重画 channel
	redrawChannel chan rune
	//    是否正在读取用户输入
	isReadingInput bool

	line   *Line
	render *TRenderer

	//   下面几个对应用户的特殊操作：退出、丢弃、确定
	exitFlag   bool
	abortFlag  bool
	returnCode Code
}

// Close 关闭命令行，恢复终端到原先的模式
func (c *TCommandLine) Close() {
	c.screen.Fini()
}

// GetLine 获取当前的 Line 对象，如果为 nil ，则 panic
func (c *TCommandLine) GetLine() *Line {
	if c.line == nil {
		panic("not found Line from TCommandLine")
	}
	return c.line
}

// GetRender 获取当前的 TRenderer 对象，如果为 nil ，则 panic
func (c *TCommandLine) GetRender() *TRenderer {
	if c.render == nil {
		panic("not found TRenderer from TCommandLine")
	}
	return c.render
}

// ReadInput 读取当前输入
func (c *TCommandLine) ReadInput() (string, error) {
	return "", nil
}

func (c *TCommandLine) Print(a ...any) {
	fmt.Print(a...)
}

func (c *TCommandLine) Println(a ...any) {
	fmt.Println(a...)
}

func (c *TCommandLine) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
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
	actualOption := defaultCommandLineOption.copy()
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
		screen: s,
		option: actualOption,

		redrawChannel: make(chan rune, 1024),
	}
	c.setup()
	return c, nil
}

func (c *TCommandLine) reset() {
	c.exitFlag = false
	c.abortFlag = false
	c.returnCode = nil
	c.readError = nil
}

func (c *TCommandLine) setup() {
	c.reset()
	if c.option.EnableDebug {
		enableDebugLog()
	} else {
		disableDebugLog()
	}
}
