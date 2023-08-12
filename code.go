package startprompt

import (
	"github.com/yetsing/startprompt/token"
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
	// Enter 用户按下 Enter 键时调用，
	// 返回 true 表示用户输入完成， CommandLine.ReadInput 则会返回用户输入
	// 返回 false 表示用户可以继续输入，文本会另起一行
	Enter() bool
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
	return nil
}

func (c *_BaseCode) Enter() bool {
	return false
}
