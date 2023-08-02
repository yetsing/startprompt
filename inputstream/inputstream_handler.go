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
	//fmt.Printf("action: %s %q\r\n", action, string(a))
	switch action {
	// enter 按下
	case ctrl_j:
		b.line.ReturnInput()
	case home:
		b.line.Home()
	case end:
		b.line.End()
	case insert_char:
		b.line.InsertText(a)
	}
}
