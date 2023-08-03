package inputstream

type InputStreamHandler interface {
	Handle(action Action, a ...rune)
}

type BaseInputStreamHandler struct {
	line *Line
}

func NewBaseInputStreamHandler(line *Line) *BaseInputStreamHandler {
	return &BaseInputStreamHandler{line: line}
}

func (b *BaseInputStreamHandler) Handle(action Action, a ...rune) {
	line := b.line
	switch action {
	case ctrl_space:

	case ctrl_a:
		line.Home()
	case ctrl_b:
		line.CursorLeft()
	case ctrl_c:
	case ctrl_d:
	case ctrl_e:
		line.End()
	case ctrl_f:
		line.CursorRight()
	case ctrl_g:
	case ctrl_h:
	case ctrl_i:
	// enter 按下
	case ctrl_j:
		line.ReturnInput()
	case ctrl_k:
	case ctrl_l:
	case ctrl_m:
		line.ReturnInput()
	case ctrl_n:
	case ctrl_o:
	case ctrl_p:
	case ctrl_q:
	case ctrl_r:
	case ctrl_s:
	case ctrl_t:
	case ctrl_u:
	case ctrl_v:
	case ctrl_w:
	case ctrl_x:
	case ctrl_y:
	case ctrl_z:
	case ctrl_backslash:
	case ctrl_square_close:
	case ctrl_circumflex:
	case ctrl_underscore:
	case backspace:
		line.DeleteCharacterBeforeCursor()
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
		line.DeleteCharacterAfterCursor()
	case page_up:
	case page_down:
	case backtab:
	case escape_action:

	case insert_char:
		line.InsertText(a)
	}
}
