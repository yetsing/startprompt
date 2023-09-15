package startprompt

type RenderContext struct {
	code          Code
	completeState *cCompletionState
	document      *Document
}

func newRenderContext(
	code Code,
	completeState *cCompletionState,
	document *Document,
) *RenderContext {
	return &RenderContext{
		code:          code,
		completeState: completeState,
		document:      document,
	}
}
