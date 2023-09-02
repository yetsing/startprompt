package startprompt

type EventHandler interface {
	Handle(event Event)
}
