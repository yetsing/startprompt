package startprompt

import (
	"github.com/yetsing/startprompt/token"
	"strings"
)

type Completion struct {
	// 展示给用户看的
	Display string
	// 加到用户输入后面的
	Suffix string
}

type NewCodeFunc func(document *Document) Code

type Code interface {
	// GetTokens 返回分词后的 Token 列表
	GetTokens() []token.Token
	// Complete 返回补全文本，可以直接添加在用户输入后，例如按一次 tab 便出现的补全
	// 返回空字符串表示没有可直接添加的补全
	Complete() string
	// GetCompletions 返回当前可选的补全列表，供用户选择，例如连按两次 tab 出现的补全列表
	GetCompletions() []*Completion
	// IsMultiline 用户按下 Enter 键时调用，
	// 返回 true 表示需要多行输入
	IsMultiline() bool
}

type _BaseCode struct {
	document *Document
}

func newBaseCode(document *Document) Code {
	return &_BaseCode{document: document}
}

func (c *_BaseCode) GetTokens() []token.Token {
	return []token.Token{
		{
			token.UNSPECIFIC,
			c.document.Text(),
		},
	}
}

func (c *_BaseCode) Complete() string {
	return ""
}

func (c *_BaseCode) GetCompletions() []*Completion {
	// for test
	return []*Completion{
		{
			Display: "hello",
			Suffix:  "hello",
		},
		{
			Display: "world",
			Suffix:  "world",
		},
		{
			Display: "中文",
			Suffix:  "中文",
		},
	}
}

func (c *_BaseCode) IsMultiline() bool {
	// for test
	text := c.document.Text()
	return !strings.HasSuffix(text, "\n")
}
