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
	switch action {
	// enter 按下
	case ctrl_j:
		b.line.ReturnInput()
	case ctrl_m:
		b.line.ReturnInput()
	case home:
		b.line.Home()
	case end:
		b.line.End()
	case insert_char:
		b.line.InsertText(a)
	}
}
