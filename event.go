package startprompt

type Event interface {
	Type() EventType
}

// EventKey 代表键盘事件
type EventKey struct {
	eventType EventType
	data      []rune
	cli       *CommandLine
	tcli      *TCommandLine
}

//goland:noinspection GoUnusedExportedFunction
func NewEventKey(eventType EventType, data []rune) *EventKey {
	return &EventKey{
		eventType: eventType,
		data:      data,
	}
}

func (ek *EventKey) Type() EventType {
	return ek.eventType
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
	eventType EventType
	cli       *CommandLine
	tcli      *TCommandLine
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