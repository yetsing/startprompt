package startprompt

type Event interface {
	Type() EventType
}

// EventKey 代表键盘事件
type EventKey struct {
	cli       *CommandLine
	tcli      *TCommandLine
	data      []rune
	eventType EventType
}

//goland:noinspection GoUnusedExportedFunction
func NewEventKey(eventType EventType, data []rune, cli *CommandLine, tcli *TCommandLine) *EventKey {
	return &EventKey{
		eventType: eventType,
		data:      data,
		cli:       cli,
		tcli:      tcli,
	}
}

func (ek *EventKey) Type() EventType {
	return ek.eventType
}

func (ek *EventKey) GetData() []rune {
	return ek.data
}

func (ek *EventKey) appendData(data []rune) {
	ek.data = append(ek.data, data...)
}

func (ek *EventKey) GetCommandLine() *CommandLine {
	if ek.cli == nil {
		panic("not found CommandLine from EventKey")
	}
	return ek.cli
}

func (ek *EventKey) GetTCommandLine() *TCommandLine {
	if ek.tcli == nil {
		panic("not found TCommandLine from EventKey")
	}
	return ek.tcli
}

// EventMouse 代表鼠标事件
type EventMouse struct {
	cli        *CommandLine
	tcli       *TCommandLine
	coordinate Coordinate
	eventType  EventType
}

func NewEventMouse(
	eventType EventType,
	coordinate Coordinate,
	cli *CommandLine,
	tcli *TCommandLine,
) *EventMouse {
	return &EventMouse{
		eventType:  eventType,
		coordinate: coordinate,
		cli:        cli,
		tcli:       tcli,
	}
}

func (em *EventMouse) Type() EventType {
	return em.eventType
}

func (em *EventMouse) GetCommandLine() *CommandLine {
	if em.cli == nil {
		panic("not found CommandLine from EventMouse")
	}
	return em.cli
}

func (em *EventMouse) GetTCommandLine() *TCommandLine {
	if em.tcli == nil {
		panic("not found TCommandLine from EventMouse")
	}
	return em.tcli
}

func (em *EventMouse) GetCoordinate() Coordinate {
	return em.coordinate
}

type EventBuffer struct {
	buffer []Event
}

func newEventBuffer() *EventBuffer {
	return &EventBuffer{}
}

func (ebuf *EventBuffer) append(event Event) {
	length := len(ebuf.buffer)
	if length > 0 {
		//    尝试合并事件
		last := ebuf.buffer[length-1]
		if last.Type() == event.Type() && event.Type() == EventTypeInsertChar {
			lastk := last.(*EventKey)
			eventk := event.(*EventKey)
			lastk.appendData(eventk.GetData())
		}
	} else {
		ebuf.buffer = append(ebuf.buffer, event)
	}
}

func (ebuf *EventBuffer) getAll() []Event {
	return ebuf.buffer
}

func (ebuf *EventBuffer) reset() {
	ebuf.buffer = nil
}
