package startprompt

type RenderContext struct {
	code          Code
	completeState *cCompletionState
	document      *Document
	selection     _LineArea
}

func newRenderContext(
	code Code,
	completeState *cCompletionState,
	document *Document,
	selection _LineArea,
) *RenderContext {
	return &RenderContext{
		code:          code,
		completeState: completeState,
		document:      document,
		selection:     selection,
	}
}
