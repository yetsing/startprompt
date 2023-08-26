package startprompt

type RenderContext struct {
	prompt        Prompt
	code          Code
	completeState *cCompletionState
	document      *Document
}

func newRenderContext(
	prompt Prompt, code Code,
	completeState *cCompletionState,
	document *Document,
) *RenderContext {
	return &RenderContext{
		prompt:        prompt,
		code:          code,
		completeState: completeState,
		document:      document,
	}
}
