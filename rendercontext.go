package startprompt

type RenderContext struct {
	code Code
	// 表示用户已经输入完毕
	accept bool
	// 表示用户已经放弃输入
	abort bool
}

func newRenderContext(code Code, accept bool, abort bool) *RenderContext {
	return &RenderContext{
		code:   code,
		accept: accept,
		abort:  abort,
	}
}
