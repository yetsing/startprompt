package startprompt

import (
	"github.com/yetsing/startprompt/token"
)

/*
分词和补全相关的接口和结构体
*/

type Completion struct {
	// Display 展示给用户看的
	Display string
	// Suffix 加到用户输入后面的
	Suffix string
	// DisplayMeta 补全的元信息，比如补全是变量或者方法之类的
	DisplayMeta string
}

type CodeFactory func(document *Document) Code

type Code interface {
	// GetTokens 返回分词后的 Token 列表
	GetTokens() []token.Token
	// Complete 返回补全文本，可以直接添加在用户输入后，例如按一次 tab 便出现的补全
	// 返回空字符串表示没有可直接添加的补全
	Complete() string
	// GetCompletions 返回当前可选的补全列表，供用户选择，例如连按两次 tab 出现的补全列表
	GetCompletions() []*Completion
	// ContinueInput 用户按下 Enter 键时调用，
	// 返回 true 时，会插入换行符
	// 返回 false 时，表示用户本次输入完成， CommandLine.ReadInput(TCommandLine.ReadInput) 则会返回用户输入
	ContinueInput() bool
	// CompleteAfterInsertText 返回 true 表示每次插入文本我们都获取一次补全
	CompleteAfterInsertText() bool
}

// _BaseCode Code 的默认实现
type _BaseCode struct {
	document *Document
}

func newBaseCode(document *Document) Code {
	return &_BaseCode{document: document}
}

func (c *_BaseCode) GetTokens() []token.Token {
	return []token.Token{
		{
			token.Unspecific,
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

func (c *_BaseCode) ContinueInput() bool {
	return false
}

func (c *_BaseCode) CompleteAfterInsertText() bool {
	return false
}
