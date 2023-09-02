package startprompt

import "github.com/yetsing/startprompt/enums/linemode"

/*
负责处理 inputstream 触发的事件
*/

type BaseHandler struct {
	//    最后处理的事件
	lastEvent EventType
	//    用户按了两次 tab
	secondTab bool

	//    保存每次事件的 Line 和 CommandLine ，省去各种方法的参数传递
	line *Line
	cli  *CommandLine
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{line: nil, secondTab: false}
}

func (b *BaseHandler) Handle(event Event) {
	DebugLog("handle event=%s", event.Type())
	ek, ok := event.(*EventKey)
	if !ok {
		DebugLog("not support event=%s %+v", event.Type(), event)
		return
	}
	eventType := ek.Type()
	b.cli = ek.GetCommandLine()
	b.line = b.cli.GetLine()
	b.lastEvent = eventType

	if b.needsToSave(eventType) {
		b.line.SaveToUndoStack()
	}

	data := ek.GetData()
	switch eventType {
	case EventTypeCtrlSpace:
		b.CtrlSpace(data)
	case EventTypeCtrlA:
		b.CtrlA(data)
	case EventTypeCtrlB:
		b.CtrlB(data)
	case EventTypeCtrlC:
		b.CtrlC(data)
	case EventTypeCtrlD:
		b.CtrlD(data)
	case EventTypeCtrlE:
		b.CtrlE(data)
	case EventTypeCtrlF:
		b.CtrlF(data)
	case EventTypeCtrlG:
		b.CtrlG(data)
	case EventTypeCtrlH:
		b.CtrlH(data)
	case EventTypeCtrlI:
		b.CtrlI(data)
	case EventTypeCtrlJ:
		b.CtrlJ(data)
	case EventTypeCtrlK:
		b.CtrlK(data)
	case EventTypeCtrlL:
		b.CtrlL(data)
	case EventTypeCtrlM:
		b.CtrlM(data)
	case EventTypeCtrlN:
		b.CtrlN(data)
	case EventTypeCtrlO:
		b.CtrlO(data)
	case EventTypeCtrlP:
		b.CtrlP(data)
	case EventTypeCtrlQ:
		b.CtrlQ(data)
	case EventTypeCtrlR:
		b.CtrlR(data)
	case EventTypeCtrlS:
		b.CtrlS(data)
	case EventTypeCtrlT:
		b.CtrlT(data)
	case EventTypeCtrlU:
		b.CtrlU(data)
	case EventTypeCtrlV:
		b.CtrlV(data)
	case EventTypeCtrlW:
		b.CtrlW(data)
	case EventTypeCtrlX:
		b.CtrlX(data)
	case EventTypeCtrlY:
		b.CtrlY(data)
	case EventTypeCtrlZ:
		b.CtrlZ(data)
	case EventTypeCtrlBackslash:
		b.CtrlBackslash(data)
	case EventTypeCtrlSquareClose:
		b.CtrlSquareClose(data)
	case EventTypeCtrlCircumflex:
		b.CtrlCircumflex(data)
	case EventTypeCtrlUnderscore:
		b.CtrlUnderscore(data)
	case EventTypeBackspace:
		b.Backspace(data)
	case EventTypeArrowUp:
		b.ArrowUp(data)
	case EventTypeArrowDown:
		b.ArrowDown(data)
	case EventTypeArrowRight:
		b.ArrowRight(data)
	case EventTypeArrowLeft:
		b.ArrowLeft(data)
	case EventTypeHome:
		b.Home(data)
	case EventTypeEnd:
		b.End(data)
	case EventTypeDeleteAction:
		b.DeleteAction(data)
	case EventTypePageUp:
		b.PageUp(data)
	case EventTypePageDown:
		b.PageDown(data)
	case EventTypeBacktab:
		b.Backtab(data)
	case EventTypeF1:
		b.F1(data)
	case EventTypeF2:
		b.F2(data)
	case EventTypeF3:
		b.F3(data)
	case EventTypeF4:
		b.F4(data)
	case EventTypeF5:
		b.F5(data)
	case EventTypeF6:
		b.F6(data)
	case EventTypeF7:
		b.F7(data)
	case EventTypeF8:
		b.F8(data)
	case EventTypeF9:
		b.F9(data)
	case EventTypeF10:
		b.F10(data)
	case EventTypeF11:
		b.F11(data)
	case EventTypeF12:
		b.F12(data)
	case EventTypeF13:
		b.F13(data)
	case EventTypeF14:
		b.F14(data)
	case EventTypeF15:
		b.F15(data)
	case EventTypeF16:
		b.F16(data)
	case EventTypeF17:
		b.F17(data)
	case EventTypeF18:
		b.F18(data)
	case EventTypeF19:
		b.F19(data)
	case EventTypeF20:
		b.F20(data)
	case EventTypeEscapeAction:
		b.EscapeAction(data)
	case EventTypeInsertChar:
		b.InsertChar(data)
	}
}

func (b *BaseHandler) CtrlSpace([]rune) {
}
func (b *BaseHandler) CtrlA([]rune) {
	b.line.ToNormalMode()
	b.line.CursorToStartOfLine(false)
}
func (b *BaseHandler) CtrlB([]rune) {
	b.line.ToNormalMode()
	b.line.CursorLeft()
}
func (b *BaseHandler) CtrlC([]rune) {
	b.line.ToNormalMode()
	b.cli.SetAbort()
}
func (b *BaseHandler) CtrlD([]rune) {
	line := b.line
	line.ToNormalMode()
	// 有输入文本时，表现为删除 delete ；否则是退出
	if line.HasText() {
		line.DeleteCharacterAfterCursor(1)
	} else {
		b.cli.SetExit()
	}
}
func (b *BaseHandler) CtrlE([]rune) {
	b.line.ToNormalMode()
	b.line.CursorToEndOfLine()
}
func (b *BaseHandler) CtrlF([]rune) {
	b.line.ToNormalMode()
	b.line.CursorRight()
}
func (b *BaseHandler) CtrlG([]rune) {

}
func (b *BaseHandler) CtrlH([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterBeforeCursor(1)
}
func (b *BaseHandler) CtrlI([]rune) {
	// ctrl_i 相当于 "\t"
	b.tab()
}
func (b *BaseHandler) CtrlJ([]rune) {
	// ctrl_j 相当于按下 Enter
	b.enter()
}
func (b *BaseHandler) CtrlK([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteUntilEndOfLine()
}
func (b *BaseHandler) CtrlL([]rune) {
	b.cli.GetRenderer().Clear()
}
func (b *BaseHandler) CtrlM([]rune) {
	// ctrl_m 相等于 "\r" ，我们把他当成 \n 的效果
	b.enter()
}
func (b *BaseHandler) CtrlN([]rune) {
	b.line.AutoDown()
}
func (b *BaseHandler) CtrlO([]rune) {

}
func (b *BaseHandler) CtrlP([]rune) {
	b.line.AutoUp()
}
func (b *BaseHandler) CtrlQ([]rune) {}
func (b *BaseHandler) CtrlR([]rune) {}
func (b *BaseHandler) CtrlS([]rune) {}
func (b *BaseHandler) CtrlT([]rune) {}
func (b *BaseHandler) CtrlU([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteFromStartOfLine()
}
func (b *BaseHandler) CtrlV([]rune) {}
func (b *BaseHandler) CtrlW([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteWordBeforeCursor()
}
func (b *BaseHandler) CtrlX([]rune)           {}
func (b *BaseHandler) CtrlY([]rune)           {}
func (b *BaseHandler) CtrlZ([]rune)           {}
func (b *BaseHandler) CtrlBackslash([]rune)   {}
func (b *BaseHandler) CtrlSquareClose([]rune) {}
func (b *BaseHandler) CtrlCircumflex([]rune)  {}
func (b *BaseHandler) CtrlUnderscore([]rune)  {}
func (b *BaseHandler) Backspace([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterBeforeCursor(1)
}
func (b *BaseHandler) ArrowUp([]rune) {
	b.line.AutoUp()
}
func (b *BaseHandler) ArrowDown([]rune) {
	b.line.AutoDown()
}
func (b *BaseHandler) ArrowRight([]rune) {
	b.line.ToNormalMode()
	b.line.CursorRight()
}
func (b *BaseHandler) ArrowLeft([]rune) {
	b.line.ToNormalMode()
	b.line.CursorLeft()
}
func (b *BaseHandler) Home([]rune) {
	b.line.ToNormalMode()
	b.line.Home()
}
func (b *BaseHandler) End([]rune) {
	b.line.ToNormalMode()
	b.line.End()
}
func (b *BaseHandler) DeleteAction([]rune) {
	b.line.ToNormalMode()
	b.line.DeleteCharacterAfterCursor(1)
}
func (b *BaseHandler) PageUp([]rune)   {}
func (b *BaseHandler) PageDown([]rune) {}
func (b *BaseHandler) Backtab([]rune)  {}
func (b *BaseHandler) F1([]rune)       {}
func (b *BaseHandler) F2([]rune)       {}
func (b *BaseHandler) F3([]rune)       {}
func (b *BaseHandler) F4([]rune)       {}
func (b *BaseHandler) F5([]rune)       {}
func (b *BaseHandler) F6([]rune)       {}
func (b *BaseHandler) F7([]rune)       {}
func (b *BaseHandler) F8([]rune)       {}
func (b *BaseHandler) F9([]rune)       {}
func (b *BaseHandler) F10([]rune)      {}
func (b *BaseHandler) F11([]rune)      {}
func (b *BaseHandler) F12([]rune)      {}
func (b *BaseHandler) F13([]rune)      {}
func (b *BaseHandler) F14([]rune)      {}
func (b *BaseHandler) F15([]rune)      {}
func (b *BaseHandler) F16([]rune)      {}
func (b *BaseHandler) F17([]rune)      {}
func (b *BaseHandler) F18([]rune)      {}
func (b *BaseHandler) F19([]rune)      {}
func (b *BaseHandler) F20([]rune)      {}
func (b *BaseHandler) EscapeAction([]rune) {
	b.line.CancelComplete()
}
func (b *BaseHandler) InsertChar(data []rune) {
	b.line.ToNormalMode()
	b.line.InsertText(data, true)
}

func (b *BaseHandler) enter() {
	b.line.AutoEnter()
	if b.line.IsAccept() {
		b.cli.SetAccept()
	}
}

func (b *BaseHandler) needsToSave(event EventType) bool {
	// 用户输入字符时不进行保存，用户输入字符后再进行保存
	// 这样一次可以撤销用户多次输入，而不是撤销一个个字符
	return !(event == EventTypeInsertChar && b.lastEvent == EventTypeInsertChar)
}

func (b *BaseHandler) tab() {
	line := b.line
	if line.mode.Is(linemode.Complete) {
		line.AcceptComplete()
		return
	}
	if !line.Complete() {
		line.CompleteNext(1)
		// 如果没有补全，插入 4 个空格
		if line.mode.Is(linemode.Normal) {
			line.InsertText([]rune("    "), true)
		}
	}
}
