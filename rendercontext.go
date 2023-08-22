package startprompt

type RenderContext struct {
	prompt        Prompt
	code          Code
	completeState *cCompletionState
	document      *Document
	// 表示用户已经输入完毕
	accept bool
	// 表示用户已经放弃输入
	abort bool
	// 表示用户退出
	exit bool
}

func newRenderContext(
	prompt Prompt, code Code,
	completeState *cCompletionState,
	document *Document,
	accept bool,
	abort bool,
	exit bool,
) *RenderContext {
	return &RenderContext{
		prompt:        prompt,
		code:          code,
		completeState: completeState,
		document:      document,
		accept:        accept,
		abort:         abort,
		exit:          exit,
	}
}
