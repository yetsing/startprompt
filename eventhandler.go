package startprompt

import "github.com/yetsing/startprompt/enums/linemode"

type EventHandler interface {
	Handle(event Event)
}

type TBaseEventHandler struct {
	//    最后处理的事件
	lastEvent EventType
	//    保存每次事件的 Line 和 TCommandLine ，省去各种方法的参数传递
	//    这两个值需要从事件中获取
	line *Line
	tcli *TCommandLine
}

func newTBaseEventHandler() *TBaseEventHandler {
	return &TBaseEventHandler{}
}

func (tb *TBaseEventHandler) Handle(event Event) {
	switch ev := event.(type) {
	case *EventKey:
		tb.HandleEventKey(ev)
		tb.lastEvent = ev.Type()
	case *EventMouse:
		tb.HandleEventMouse(ev)
		tb.lastEvent = ev.Type()
	default:
		DebugLog("not support event=%s %+v", event.Type(), event)
	}
}

func (tb *TBaseEventHandler) HandleEventMouse(em *EventMouse) {
	tb.tcli = em.GetTCommandLine()
	tb.line = tb.tcli.GetLine()
	eventType := em.Type()
	switch eventType {
	case EventMouseScrollUp:
		tb.tcli.GetRenderer().ScrollUp(1)
	case EventMouseScrollDown:
		tb.tcli.GetRenderer().ScrollDown(1)
	}
}

func (tb *TBaseEventHandler) HandleEventKey(ek *EventKey) {
	eventType := ek.Type()
	tb.tcli = ek.GetTCommandLine()
	tb.line = tb.tcli.GetLine()

	if tb.needsToSave(eventType) {
		tb.line.SaveToUndoStack()
	}

	data := ek.GetData()
	switch eventType {
	case EventTypeCtrlSpace:
		tb.CtrlSpace(data)
	case EventTypeCtrlA:
		tb.CtrlA(data)
	case EventTypeCtrlB:
		tb.CtrlB(data)
	case EventTypeCtrlC:
		tb.CtrlC(data)
	case EventTypeCtrlD:
		tb.CtrlD(data)
	case EventTypeCtrlE:
		tb.CtrlE(data)
	case EventTypeCtrlF:
		tb.CtrlF(data)
	case EventTypeCtrlG:
		tb.CtrlG(data)
	case EventTypeCtrlH:
		tb.CtrlH(data)
	case EventTypeCtrlI:
		tb.CtrlI(data)
	case EventTypeCtrlJ:
		tb.CtrlJ(data)
	case EventTypeCtrlK:
		tb.CtrlK(data)
	case EventTypeCtrlL:
		tb.CtrlL(data)
	case EventTypeCtrlM:
		tb.CtrlM(data)
	case EventTypeCtrlN:
		tb.CtrlN(data)
	case EventTypeCtrlO:
		tb.CtrlO(data)
	case EventTypeCtrlP:
		tb.CtrlP(data)
	case EventTypeCtrlQ:
		tb.CtrlQ(data)
	case EventTypeCtrlR:
		tb.CtrlR(data)
	case EventTypeCtrlS:
		tb.CtrlS(data)
	case EventTypeCtrlT:
		tb.CtrlT(data)
	case EventTypeCtrlU:
		tb.CtrlU(data)
	case EventTypeCtrlV:
		tb.CtrlV(data)
	case EventTypeCtrlW:
		tb.CtrlW(data)
	case EventTypeCtrlX:
		tb.CtrlX(data)
	case EventTypeCtrlY:
		tb.CtrlY(data)
	case EventTypeCtrlZ:
		tb.CtrlZ(data)
	case EventTypeCtrlBackslash:
		tb.CtrlBackslash(data)
	case EventTypeCtrlSquareClose:
		tb.CtrlSquareClose(data)
	case EventTypeCtrlCircumflex:
		tb.CtrlCircumflex(data)
	case EventTypeCtrlUnderscore:
		tb.CtrlUnderscore(data)
	case EventTypeBackspace:
		tb.Backspace(data)
	case EventTypeArrowUp:
		tb.ArrowUp(data)
	case EventTypeArrowDown:
		tb.ArrowDown(data)
	case EventTypeArrowRight:
		tb.ArrowRight(data)
	case EventTypeArrowLeft:
		tb.ArrowLeft(data)
	case EventTypeHome:
		tb.Home(data)
	case EventTypeEnd:
		tb.End(data)
	case EventTypeDeleteAction:
		tb.DeleteAction(data)
	case EventTypePageUp:
		tb.PageUp(data)
	case EventTypePageDown:
		tb.PageDown(data)
	case EventTypeBacktab:
		tb.Backtab(data)
	case EventTypeF1:
		tb.F1(data)
	case EventTypeF2:
		tb.F2(data)
	case EventTypeF3:
		tb.F3(data)
	case EventTypeF4:
		tb.F4(data)
	case EventTypeF5:
		tb.F5(data)
	case EventTypeF6:
		tb.F6(data)
	case EventTypeF7:
		tb.F7(data)
	case EventTypeF8:
		tb.F8(data)
	case EventTypeF9:
		tb.F9(data)
	case EventTypeF10:
		tb.F10(data)
	case EventTypeF11:
		tb.F11(data)
	case EventTypeF12:
		tb.F12(data)
	case EventTypeF13:
		tb.F13(data)
	case EventTypeF14:
		tb.F14(data)
	case EventTypeF15:
		tb.F15(data)
	case EventTypeF16:
		tb.F16(data)
	case EventTypeF17:
		tb.F17(data)
	case EventTypeF18:
		tb.F18(data)
	case EventTypeF19:
		tb.F19(data)
	case EventTypeF20:
		tb.F20(data)
	case EventTypeEscape:
		tb.EscapeAction(data)
	case EventTypeInsertChar:
		tb.InsertChar(data)
	}
}

func (tb *TBaseEventHandler) CtrlSpace(_ []rune) {
}
func (tb *TBaseEventHandler) CtrlA(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorToStartOfLine(false)
}
func (tb *TBaseEventHandler) CtrlB(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorLeft()
}
func (tb *TBaseEventHandler) CtrlC(_ []rune) {
	tb.line.ToNormalMode()
	tb.tcli.SetAbort()
}
func (tb *TBaseEventHandler) CtrlD(_ []rune) {
	line := tb.line
	line.ToNormalMode()
	// 有输入文本时，表现为删除 delete ；否则是退出
	if line.HasText() {
		line.DeleteCharacterAfterCursor(1)
	} else {
		tb.tcli.SetExit()
	}
}
func (tb *TBaseEventHandler) CtrlE(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorToEndOfLine()
}
func (tb *TBaseEventHandler) CtrlF(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorRight()
}
func (tb *TBaseEventHandler) CtrlG(_ []rune) {

}
func (tb *TBaseEventHandler) CtrlH(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteCharacterBeforeCursor(1)
}
func (tb *TBaseEventHandler) CtrlI(_ []rune) {
	// ctrl_i 相当于 "\t"
	tb.tab()
}
func (tb *TBaseEventHandler) CtrlJ(_ []rune) {
	// ctrl_j 相当于按下 Enter
	tb.enter()
}
func (tb *TBaseEventHandler) CtrlK(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteUntilEndOfLine()
}
func (tb *TBaseEventHandler) CtrlL(_ []rune) {
	tb.tcli.GetRenderer().Clear()
}
func (tb *TBaseEventHandler) CtrlM(_ []rune) {
	// ctrl_m 相等于 "\r" ，我们把他当成 \n 的效果
	tb.enter()
}
func (tb *TBaseEventHandler) CtrlN(_ []rune) {
	tb.line.AutoDown()
}
func (tb *TBaseEventHandler) CtrlO(_ []rune) {

}
func (tb *TBaseEventHandler) CtrlP(_ []rune) {
	tb.line.AutoUp()
}
func (tb *TBaseEventHandler) CtrlQ(_ []rune) {}
func (tb *TBaseEventHandler) CtrlR(_ []rune) {}
func (tb *TBaseEventHandler) CtrlS(_ []rune) {}
func (tb *TBaseEventHandler) CtrlT(_ []rune) {}
func (tb *TBaseEventHandler) CtrlU(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteFromStartOfLine()
}
func (tb *TBaseEventHandler) CtrlV(_ []rune) {}
func (tb *TBaseEventHandler) CtrlW(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteWordBeforeCursor()
}
func (tb *TBaseEventHandler) CtrlX(_ []rune)           {}
func (tb *TBaseEventHandler) CtrlY(_ []rune)           {}
func (tb *TBaseEventHandler) CtrlZ(_ []rune)           {}
func (tb *TBaseEventHandler) CtrlBackslash(_ []rune)   {}
func (tb *TBaseEventHandler) CtrlSquareClose(_ []rune) {}
func (tb *TBaseEventHandler) CtrlCircumflex(_ []rune)  {}
func (tb *TBaseEventHandler) CtrlUnderscore(_ []rune)  {}
func (tb *TBaseEventHandler) Backspace(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteCharacterBeforeCursor(1)
}
func (tb *TBaseEventHandler) ArrowUp(_ []rune) {
	tb.line.AutoUp()
}
func (tb *TBaseEventHandler) ArrowDown(_ []rune) {
	tb.line.AutoDown()
}
func (tb *TBaseEventHandler) ArrowRight(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorRight()
}
func (tb *TBaseEventHandler) ArrowLeft(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.CursorLeft()
}
func (tb *TBaseEventHandler) Home(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.Home()
}
func (tb *TBaseEventHandler) End(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.End()
}
func (tb *TBaseEventHandler) DeleteAction(_ []rune) {
	tb.line.ToNormalMode()
	tb.line.DeleteCharacterAfterCursor(1)
}
func (tb *TBaseEventHandler) PageUp(_ []rune)   {}
func (tb *TBaseEventHandler) PageDown(_ []rune) {}
func (tb *TBaseEventHandler) Backtab(_ []rune)  {}
func (tb *TBaseEventHandler) F1(_ []rune)       {}
func (tb *TBaseEventHandler) F2(_ []rune)       {}
func (tb *TBaseEventHandler) F3(_ []rune)       {}
func (tb *TBaseEventHandler) F4(_ []rune)       {}
func (tb *TBaseEventHandler) F5(_ []rune)       {}
func (tb *TBaseEventHandler) F6(_ []rune)       {}
func (tb *TBaseEventHandler) F7(_ []rune)       {}
func (tb *TBaseEventHandler) F8(_ []rune)       {}
func (tb *TBaseEventHandler) F9(_ []rune)       {}
func (tb *TBaseEventHandler) F10(_ []rune)      {}
func (tb *TBaseEventHandler) F11(_ []rune)      {}
func (tb *TBaseEventHandler) F12(_ []rune)      {}
func (tb *TBaseEventHandler) F13(_ []rune)      {}
func (tb *TBaseEventHandler) F14(_ []rune)      {}
func (tb *TBaseEventHandler) F15(_ []rune)      {}
func (tb *TBaseEventHandler) F16(_ []rune)      {}
func (tb *TBaseEventHandler) F17(_ []rune)      {}
func (tb *TBaseEventHandler) F18(_ []rune)      {}
func (tb *TBaseEventHandler) F19(_ []rune)      {}
func (tb *TBaseEventHandler) F20(_ []rune)      {}
func (tb *TBaseEventHandler) EscapeAction(_ []rune) {
	tb.line.CancelComplete()
}
func (tb *TBaseEventHandler) InsertChar(data []rune) {
	tb.line.ToNormalMode()
	tb.line.InsertText(data, true)
}

func (tb *TBaseEventHandler) enter() {
	tb.line.AutoEnter()
	if tb.line.IsAccept() {
		tb.tcli.SetAccept()
	}
}

func (tb *TBaseEventHandler) needsToSave(event EventType) bool {
	// 用户输入字符时不进行保存，用户输入字符后再进行保存
	// 这样一次可以撤销用户多次输入，而不是撤销一个个字符
	return !(event == EventTypeInsertChar && tb.lastEvent == EventTypeInsertChar)
}

func (tb *TBaseEventHandler) tab() {
	line := tb.line
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
