package startprompt

type InputStreamHandler interface {
	Handle(action Event, a ...rune)
}

type BaseHandler struct {
	line *Line
	// 最后处理的事件
	lastEvent Event
	// 用户按了两次 tab
	secondTab bool
}

func NewBaseHandler(line *Line) *BaseHandler {
	return &BaseHandler{line: line, secondTab: false}
}

func (b *BaseHandler) Handle(event Event, a ...rune) {
	line := b.line
	b.lastEvent = event
	if b.needsToSave(event) {
		line.SaveToUndoStack()
	}
	switch event {
	case ctrl_space:

	case ctrl_a:
		line.CursorToStartOfLine(false)
	case ctrl_b:
		line.CursorLeft()
	case ctrl_c:
		line.Abort()
	case ctrl_d:
		// 有输入文本时，表现为删除 delete ；否则是退出
		if line.HasText() {
			line.DeleteCharacterAfterCursor(1)
		} else {
			line.Exit()
		}
	case ctrl_e:
		line.CursorToEndOfLine()
	case ctrl_f:
		line.CursorRight()
	case ctrl_g:
		// todo Abort an incremental search and restore the original line
	case ctrl_h:
		line.DeleteCharacterBeforeCursor(1)
	case ctrl_i:
		// ctrl_i 相当于 "\t"
		b.tab()
	// enter 按下
	case ctrl_j:
		b.enter()
	case ctrl_k:
		line.DeleteUntilEndOfLine()
	case ctrl_l:
		line.Clear()
	case ctrl_m:
		// ctrl_m 相等于 "\r" ，我们把他当成 \n 的效果
		b.enter()
	case ctrl_n:
	case ctrl_o:
	case ctrl_p:
	case ctrl_q:
	case ctrl_r:
	case ctrl_s:
	case ctrl_t:
	case ctrl_u:
		line.DeleteFromStartOfLine()
	case ctrl_v:
	case ctrl_w:
		line.DeleteWordBeforeCursor()
	case ctrl_x:
	case ctrl_y:
	case ctrl_z:
	case ctrl_backslash:
	case ctrl_square_close:
	case ctrl_circumflex:
	case ctrl_underscore:
	case backspace:
		line.DeleteCharacterBeforeCursor(1)
	case arrow_up:
	case arrow_down:
	case arrow_right:
		line.CursorRight()
	case arrow_left:
		line.CursorLeft()
	case home:
		line.Home()
	case end:
		line.End()
	case delete_action:
		line.DeleteCharacterAfterCursor(1)
	case page_up:
	case page_down:
	case backtab:
	case escape_action:

	case insert_char:
		line.InsertText(a)
	}
}

func (b *BaseHandler) needsToSave(event Event) bool {
	// 用户输入字符时不进行保存，用户输入字符后再进行保存
	// 这样一次可以撤销用户多次输入，而不是撤销一个个字符
	return !(event == insert_char && b.lastEvent == insert_char)
}

func (b *BaseHandler) tab() {
	if b.secondTab {
		// 两次 tab 会展示补全列表
		// 效果类似于 bash 里面按 tab 两次
		b.line.ListCompletions()
		b.secondTab = false
	} else {
		// 有补全的话，就不需要触发两次 tab 的效果了
		b.secondTab = !b.line.Complete()
	}
}

func (b *BaseHandler) enter() {
	line := b.line
	if line.IsMultiline() {
		line.Newline()
	} else {
		line.ReturnInput()
	}
}
