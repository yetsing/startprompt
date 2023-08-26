package startprompt

import (
	"github.com/yetsing/startprompt/enums/linemode"
	"github.com/yetsing/startprompt/keys"
)

type InputStreamHandler interface {
	Handle(action keys.Event, a ...rune)
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
	lastEvent keys.Event
	// 用户按了两次 tab
	secondTab bool
}

func NewBaseHandler(line *Line) *BaseHandler {
	return &BaseHandler{line: line, secondTab: false}
}

func (b *BaseHandler) Handle(event keys.Event, a ...rune) {
	b.lastEvent = event
	if b.needsToSave(event) {
		b.line.SaveToUndoStack()
	}

	switch event {
	case keys.CtrlSpace:
		b.CtrlSpace(a...)
	case keys.CtrlA:
		b.CtrlA(a...)
	case keys.CtrlB:
		b.CtrlB(a...)
	case keys.CtrlC:
		b.CtrlC(a...)
	case keys.CtrlD:
		b.CtrlD(a...)
	case keys.CtrlE:
		b.CtrlE(a...)
	case keys.CtrlF:
		b.CtrlF(a...)
	case keys.CtrlG:
		b.CtrlG(a...)
	case keys.CtrlH:
		b.CtrlH(a...)
	case keys.CtrlI:
		b.CtrlI(a...)
	case keys.CtrlJ:
		b.CtrlJ(a...)
	case keys.CtrlK:
		b.CtrlK(a...)
	case keys.CtrlL:
		b.CtrlL(a...)
	case keys.CtrlM:
		b.CtrlM(a...)
	case keys.CtrlN:
		b.CtrlN(a...)
	case keys.CtrlO:
		b.CtrlO(a...)
	case keys.CtrlP:
		b.CtrlP(a...)
	case keys.CtrlQ:
		b.CtrlQ(a...)
	case keys.CtrlR:
		b.CtrlR(a...)
	case keys.CtrlS:
		b.CtrlS(a...)
	case keys.CtrlT:
		b.CtrlT(a...)
	case keys.CtrlU:
		b.CtrlU(a...)
	case keys.CtrlV:
		b.CtrlV(a...)
	case keys.CtrlW:
		b.CtrlW(a...)
	case keys.CtrlX:
		b.CtrlX(a...)
	case keys.CtrlY:
		b.CtrlY(a...)
	case keys.CtrlZ:
		b.CtrlZ(a...)
	case keys.CtrlBackslash:
		b.CtrlBackslash(a...)
	case keys.CtrlSquareClose:
		b.CtrlSquareClose(a...)
	case keys.CtrlCircumflex:
		b.CtrlCircumflex(a...)
	case keys.CtrlUnderscore:
		b.CtrlUnderscore(a...)
	case keys.Backspace:
		b.Backspace(a...)
	case keys.ArrowUp:
		b.ArrowUp(a...)
	case keys.ArrowDown:
		b.ArrowDown(a...)
	case keys.ArrowRight:
		b.ArrowRight(a...)
	case keys.ArrowLeft:
		b.ArrowLeft(a...)
	case keys.Home:
		b.Home(a...)
	case keys.End:
		b.End(a...)
	case keys.DeleteAction:
		b.DeleteAction(a...)
	case keys.PageUp:
		b.PageUp(a...)
	case keys.PageDown:
		b.PageDown(a...)
	case keys.Backtab:
		b.Backtab(a...)
	case keys.F1:
		b.F1(a...)
	case keys.F2:
		b.F2(a...)
	case keys.F3:
		b.F3(a...)
	case keys.F4:
		b.F4(a...)
	case keys.F5:
		b.F5(a...)
	case keys.F6:
		b.F6(a...)
	case keys.F7:
		b.F7(a...)
	case keys.F8:
		b.F8(a...)
	case keys.F9:
		b.F9(a...)
	case keys.F10:
		b.F10(a...)
	case keys.F11:
		b.F11(a...)
	case keys.F12:
		b.F12(a...)
	case keys.F13:
		b.F13(a...)
	case keys.F14:
		b.F14(a...)
	case keys.F15:
		b.F15(a...)
	case keys.F16:
		b.F16(a...)
	case keys.F17:
		b.F17(a...)
	case keys.F18:
		b.F18(a...)
	case keys.F19:
		b.F19(a...)
	case keys.F20:
		b.F20(a...)
	case keys.EscapeAction:
		b.EscapeAction(a...)
	case keys.InsertChar:
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
	b.line.ToNormalMode()
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
	b.line.ToNormalMode()
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
	b.line.ToNormalMode()
}
func (b *BaseHandler) InsertChar(a ...rune) {
	b.line.ToNormalMode()
	b.line.InsertText(a, true)
}

func (b *BaseHandler) needsToSave(event keys.Event) bool {
	// 用户输入字符时不进行保存，用户输入字符后再进行保存
	// 这样一次可以撤销用户多次输入，而不是撤销一个个字符
	return !(event == keys.InsertChar && b.lastEvent == keys.InsertChar)
}

func (b *BaseHandler) tab() {
	if b.line.mode.Is(linemode.Complete) {
		b.line.CompleteNext(0)
		b.line.ExitComplete()
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
