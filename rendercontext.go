package startprompt

type RenderContext struct {
	completeState   *cCompletionState
	document        *Document
	code            Code
	highlights      []section
	cancelSelection bool
}

func newRenderContext(
	code Code,
	completeState *cCompletionState,
	document *Document,
	highlights []section,
	cancelSelection bool,
) *RenderContext {
	return &RenderContext{
		code:            code,
		completeState:   completeState,
		document:        document,
		highlights:      highlights,
		cancelSelection: cancelSelection,
	}
}
