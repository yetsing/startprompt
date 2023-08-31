package startprompt

import (
	"github.com/yetsing/startprompt/enums/linemode"
)

type InputStreamHandler interface {
	Handle(action EventType, a ...rune)
	// 下面的方法跟事件一一对应

	CtrlSpace(...rune)
	CtrlA(...rune)
	CtrlB(...rune)
	CtrlC(...rune)
	CtrlD(...rune)
	CtrlE(...rune)
	CtrlF(...rune)
	CtrlG(...rune)
	CtrlH(...rune)
	CtrlI(...rune)
	CtrlJ(...rune)
	CtrlK(...rune)
	CtrlL(...rune)
	CtrlM(...rune)
	CtrlN(...rune)
	CtrlO(...rune)
	CtrlP(...rune)
	CtrlQ(...rune)
	CtrlR(...rune)
	CtrlS(...rune)
	CtrlT(...rune)
	CtrlU(...rune)
	CtrlV(...rune)
	CtrlW(...rune)
	CtrlX(...rune)
	CtrlY(...rune)
	CtrlZ(...rune)
	CtrlBackslash(...rune)
	CtrlSquareClose(...rune)
	CtrlCircumflex(...rune)
	CtrlUnderscore(...rune)
	Backspace(...rune)
	ArrowUp(...rune)
	ArrowDown(...rune)
	ArrowRight(...rune)
	ArrowLeft(...rune)
	Home(...rune)
	End(...rune)
	DeleteAction(...rune)
	PageUp(...rune)
	PageDown(...rune)
	Backtab(...rune)
	F1(...rune)
	F2(...rune)
	F3(...rune)
	F4(...rune)
	F5(...rune)
	F6(...rune)
	F7(...rune)
	F8(...rune)
	F9(...rune)
	F10(...rune)
	F11(...rune)
	F12(...rune)
	F13(...rune)
	F14(...rune)
	F15(...rune)
	F16(...rune)
	F17(...rune)
	F18(...rune)
	F19(...rune)
	F20(...rune)
	EscapeAction(...rune)
	InsertChar(...rune)
}

type BaseHandler struct {
	line *Line
	// 最后处理的事件
	lastEvent EventType
	// 用户按了两次 tab
	secondTab bool
}

func NewBaseHandler(line *Line) *BaseHandler {
	return &BaseHandler{line: line, secondTab: false}
}

func (b *BaseHandler) Handle(event EventType, a ...rune) {
	b.lastEvent = event
	if b.needsToSave(event) {
		b.line.SaveToUndoStack()
	}

	switch event {
	case event.CtrlSpace:
		b.CtrlSpace(a...)
	case event.CtrlA:
		b.CtrlA(a...)
	case event.CtrlB:
		b.CtrlB(a...)
	case event.CtrlC:
		b.CtrlC(a...)
	case event.CtrlD:
		b.CtrlD(a...)
	case event.CtrlE:
		b.CtrlE(a...)
	case event.CtrlF:
		b.CtrlF(a...)
	case event.CtrlG:
		b.CtrlG(a...)
	case event.CtrlH:
		b.CtrlH(a...)
	case event.CtrlI:
		b.CtrlI(a...)
	case event.CtrlJ:
		b.CtrlJ(a...)
	case event.CtrlK:
		b.CtrlK(a...)
	case event.CtrlL:
		b.CtrlL(a...)
	case event.CtrlM:
		b.CtrlM(a...)
	case event.CtrlN:
		b.CtrlN(a...)
	case event.CtrlO:
		b.CtrlO(a...)
	case event.CtrlP:
		b.CtrlP(a...)
	case event.CtrlQ:
		b.CtrlQ(a...)
	case event.CtrlR:
		b.CtrlR(a...)
	case event.CtrlS:
		b.CtrlS(a...)
	case event.CtrlT:
		b.CtrlT(a...)
	case event.CtrlU:
		b.CtrlU(a...)
	case event.CtrlV:
		b.CtrlV(a...)
	case event.CtrlW:
		b.CtrlW(a...)
	case event.CtrlX:
		b.CtrlX(a...)
	case event.CtrlY:
		b.CtrlY(a...)
	case event.CtrlZ:
		b.CtrlZ(a...)
	case event.CtrlBackslash:
		b.CtrlBackslash(a...)
	case event.CtrlSquareClose:
		b.CtrlSquareClose(a...)
	case event.CtrlCircumflex:
		b.CtrlCircumflex(a...)
	case event.CtrlUnderscore:
		b.CtrlUnderscore(a...)
	case event.Backspace:
		b.Backspace(a...)
	case event.ArrowUp:
		b.ArrowUp(a...)
	case event.ArrowDown:
		b.ArrowDown(a...)
	case event.ArrowRight:
		b.ArrowRight(a...)
	case event.ArrowLeft:
		b.ArrowLeft(a...)
	case event.Home:
		b.Home(a...)
	case event.End:
		b.End(a...)
	case event.DeleteAction:
		b.DeleteAction(a...)
	case event.PageUp:
		b.PageUp(a...)
	case event.PageDown:
		b.PageDown(a...)
	case event.Backtab:
		b.Backtab(a...)
	case event.F1:
		b.F1(a...)
	case event.F2:
		b.F2(a...)
	case event.F3:
		b.F3(a...)
	case event.F4:
		b.F4(a...)
	case event.F5:
		b.F5(a...)
	case event.F6:
		b.F6(a...)
	case event.F7:
		b.F7(a...)
	case event.F8:
		b.F8(a...)
	case event.F9:
		b.F9(a...)
	case event.F10:
		b.F10(a...)
	case event.F11:
		b.F11(a...)
	case event.F12:
		b.F12(a...)
	case event.F13:
		b.F13(a...)
	case event.F14:
		b.F14(a...)
	case event.F15:
		b.F15(a...)
	case event.F16:
		b.F16(a...)
	case event.F17:
		b.F17(a...)
	case event.F18:
		b.F18(a...)
	case event.F19:
		b.F19(a...)
	case event.F20:
		b.F20(a...)
	case event.EscapeAction:
		b.EscapeAction(a...)
	case event.InsertChar:
		b.InsertChar(a...)
	}
}

func (b *BaseHandler) CtrlSpace(...rune) {
}
func (b *BaseHandler) CtrlA(...rune) {
	b.line.ToNormalMode()
	b.line.CursorToStartOfLine(false)
}
func (b *BaseHandler) CtrlB(...rune) {
	b.line.ToNormalMode()
	b.line.CursorLeft()
}
func (b *BaseHandler) CtrlC(...rune) {
	b.line.ToNormalMode()
	b.line.Abort()
}
func (b *BaseHandler) CtrlD(...rune) {
	b.line.ToNormalMode()
	line := b.line
	// 有输入文本时，表现为删除 delete ；否则是退出
	if line.HasText() {
		line.DeleteCharacterAfterCursor(1)
	} else {
		line.Exit()
	}
}
func (b *BaseHandler) CtrlE(...rune) {
	b.line.ToNormalMode()
	b.line.CursorToEndOfLine()
}
func (b *BaseHandler) CtrlF(...rune) {
	b.line.ToNormalMode()
	b.line.CursorRight()
}
func (b *BaseHandler) CtrlG(...rune) {

}
func (b *BaseHandler) CtrlH(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterBeforeCursor(1)
}
func (b *BaseHandler) CtrlI(...rune) {
	// ctrl_i 相当于 "\t"
	b.tab()
}
func (b *BaseHandler) CtrlJ(...rune) {
	// ctrl_j 相当于按下 Enter
	b.line.AutoEnter()
}
func (b *BaseHandler) CtrlK(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteUntilEndOfLine()
}
func (b *BaseHandler) CtrlL(...rune) {
	b.line.Clear()
}
func (b *BaseHandler) CtrlM(...rune) {
	// ctrl_m 相等于 "\r" ，我们把他当成 \n 的效果
	b.line.AutoEnter()
}
func (b *BaseHandler) CtrlN(...rune) {
	b.line.AutoDown()
}
func (b *BaseHandler) CtrlO(...rune) {

}
func (b *BaseHandler) CtrlP(...rune) {
	b.line.AutoUp()
}
func (b *BaseHandler) CtrlQ(...rune) {}
func (b *BaseHandler) CtrlR(...rune) {}
func (b *BaseHandler) CtrlS(...rune) {}
func (b *BaseHandler) CtrlT(...rune) {}
func (b *BaseHandler) CtrlU(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteFromStartOfLine()
}
func (b *BaseHandler) CtrlV(...rune) {}
func (b *BaseHandler) CtrlW(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteWordBeforeCursor()
}
func (b *BaseHandler) CtrlX(...rune)           {}
func (b *BaseHandler) CtrlY(...rune)           {}
func (b *BaseHandler) CtrlZ(...rune)           {}
func (b *BaseHandler) CtrlBackslash(...rune)   {}
func (b *BaseHandler) CtrlSquareClose(...rune) {}
func (b *BaseHandler) CtrlCircumflex(...rune)  {}
func (b *BaseHandler) CtrlUnderscore(...rune)  {}
func (b *BaseHandler) Backspace(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterBeforeCursor(1)
}
func (b *BaseHandler) ArrowUp(...rune) {
	b.line.AutoUp()
}
func (b *BaseHandler) ArrowDown(...rune) {
	b.line.AutoDown()
}
func (b *BaseHandler) ArrowRight(...rune) {
	b.line.ToNormalMode()
	b.line.CursorRight()
}
func (b *BaseHandler) ArrowLeft(...rune) {
	b.line.ToNormalMode()
	b.line.CursorLeft()
}
func (b *BaseHandler) Home(...rune) {
	b.line.ToNormalMode()
	b.line.Home()
}
func (b *BaseHandler) End(...rune) {
	b.line.ToNormalMode()
	b.line.End()
}
func (b *BaseHandler) DeleteAction(...rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterAfterCursor(1)
}
func (b *BaseHandler) PageUp(...rune)   {}
func (b *BaseHandler) PageDown(...rune) {}
func (b *BaseHandler) Backtab(...rune)  {}
func (b *BaseHandler) F1(...rune)       {}
func (b *BaseHandler) F2(...rune)       {}
func (b *BaseHandler) F3(...rune)       {}
func (b *BaseHandler) F4(...rune)       {}
func (b *BaseHandler) F5(...rune)       {}
func (b *BaseHandler) F6(...rune)       {}
func (b *BaseHandler) F7(...rune)       {}
func (b *BaseHandler) F8(...rune)       {}
func (b *BaseHandler) F9(...rune)       {}
func (b *BaseHandler) F10(...rune)      {}
func (b *BaseHandler) F11(...rune)      {}
func (b *BaseHandler) F12(...rune)      {}
func (b *BaseHandler) F13(...rune)      {}
func (b *BaseHandler) F14(...rune)      {}
func (b *BaseHandler) F15(...rune)      {}
func (b *BaseHandler) F16(...rune)      {}
func (b *BaseHandler) F17(...rune)      {}
func (b *BaseHandler) F18(...rune)      {}
func (b *BaseHandler) F19(...rune)      {}
func (b *BaseHandler) F20(...rune)      {}
func (b *BaseHandler) EscapeAction(...rune) {
	b.line.CancelComplete()
}
func (b *BaseHandler) InsertChar(a ...rune) {
	b.line.ToNormalMode()
	b.line.InsertText(a, true)
}

func (b *BaseHandler) needsToSave(event EventType) bool {
	// 用户输入字符时不进行保存，用户输入字符后再进行保存
	// 这样一次可以撤销用户多次输入，而不是撤销一个个字符
	return !(event == event.InsertChar && b.lastEvent == event.InsertChar)
}

func (b *BaseHandler) tab() {
	if b.line.mode.Is(linemode.Complete) {
		b.line.AcceptComplete()
		return
	}
	if !b.line.Complete() {
		b.line.CompleteNext(1)
		// 如果没有补全，插入 4 个空格
		if b.line.mode.Is(linemode.Normal) {
			b.line.InsertText([]rune("    "), true)
		}
	}
}
