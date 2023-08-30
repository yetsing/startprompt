package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

/*
自定义提示符
*/

var inputCount = 1
var schema = map[token.TokenType]terminalcolor.Style{
	token.Prompt: terminalcolor.NewFgColorStyleHex("#004400"),
}

type Prompt struct {
	line *startprompt.Line
	code startprompt.Code
}

func NewPrompt(line *startprompt.Line, code startprompt.Code) startprompt.Prompt {
	return &Prompt{line: line, code: code}
}

func (p *Prompt) GetPrompt() []token.Token {
	tk := token.NewToken(token.Prompt, fmt.Sprintf("\nIn [%d]: ", inputCount))
	return []token.Token{tk}
}

func (p *Prompt) GetSecondLinePrefix() []token.Token {
	// 拿到默认提示符宽度
	var sb strings.Builder
	for _, t := range p.GetPrompt() {
		sb.WriteString(t.Literal)
	}
	promptText := sb.String()
	spaces := runewidth.StringWidth(promptText) - 5
	// 输出类似这样的 "...: " ，宽度跟默认提示符一样
	return []token.Token{
		{
			token.PromptSecondLinePrefix,
			startprompt.RepeatString(" ", spaces),
		},
		{
			token.PromptSecondLinePrefix,
			startprompt.RepeatString(".", 3) + ": ",
		},
	}
}

type MultilineCode struct {
	document *startprompt.Document
}

func newMultilineCode(document *startprompt.Document) startprompt.Code {
	return &MultilineCode{document: document}
}

func (c *MultilineCode) GetTokens() []token.Token {
	return []token.Token{
		{
			token.Unspecific,
			c.document.Text(),
		},
	}
}

func (c *MultilineCode) Complete() string {
	return ""
}

func (c *MultilineCode) GetCompletions() []*startprompt.Completion {
	return nil
}

func (c *MultilineCode) ContinueInput() bool {
	// 用于需要连续按下两次 Enter 才结束当前输入
	text := c.document.Text()
	return len(text) > 0 && !strings.HasSuffix(c.document.Text(), "\n")
}

func (c *MultilineCode) CompleteAfterInsertText() bool {
	return false
}

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		CodeFactory:   newMultilineCode,
		PromptFactory: NewPrompt,
		Schema:        schema,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	defer c.Close()
	fmt.Println("Press Ctrl-D exit")
	for {
		line, err := c.ReadInput()
		if err != nil {
			if errors.Is(err, startprompt.ExitError) {
				break
			}
			fmt.Printf("ReadInput error: %v\n", err)
			break
		}
		if len(line) == 0 {
			continue
		}
		inputCount++
		fmt.Println(line)
	}
}
