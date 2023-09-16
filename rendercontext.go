package startprompt

type RenderContext struct {
	completeState   *cCompletionState
	document        *Document
	code            Code
	cancelSelection bool
}

func newRenderContext(
	code Code,
	completeState *cCompletionState,
	document *Document,
	cancelSelection bool,
) *RenderContext {
	return &RenderContext{
		code:            code,
		completeState:   completeState,
		document:        document,
		cancelSelection: cancelSelection,
	}
}
