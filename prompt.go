package startprompt

import (
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/token"
)

type PromptFactory func(line *Line, code Code) Prompt

type Prompt interface {
	// GetPrompt 获取输入的提示符
	GetPrompt() []token.Token
	// GetSecondLinePrefix 获取输入第二行及之后行开始的提示符
	GetSecondLinePrefix() []token.Token
}

type BasePrompt struct {
}

//goland:noinspection GoUnusedParameter
func newBasePrompt(line *Line, code Code) Prompt {
	return &BasePrompt{}
}

func (b *BasePrompt) GetPrompt() []token.Token {
	tk := token.Token{
		Type:    token.Prompt,
		Literal: "> ",
	}
	return []token.Token{tk}
}

func (b *BasePrompt) GetSecondLinePrefix() []token.Token {
	// 拿到默认提示符宽度
	var sb strings.Builder
	for _, t := range b.GetPrompt() {
		sb.WriteString(t.Literal)
	}
	promptText := sb.String()
	width := runewidth.StringWidth(strings.TrimRightFunc(promptText, unicode.IsSpace))
	spaces := runewidth.StringWidth(promptText) - width
	// 输出类似这样的 "...  " ，宽度跟默认提示符一样
	return []token.Token{
		{
			token.PromptSecondLinePrefix,
			repeatByte('.', width),
		},
		{
			token.PromptSecondLinePrefix,
			repeatByte(' ', spaces),
		},
	}
}
