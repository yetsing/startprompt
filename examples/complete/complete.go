package main

/*
展示补全的用法和效果
*/

import (
	"fmt"
	"strings"

	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/token"
)

type CompleteCode struct {
	document *startprompt.Document
}

func newCompleteCode(document *startprompt.Document) startprompt.Code {
	return &CompleteCode{document: document}
}

func (c *CompleteCode) GetTokens() []token.Token {
	return []token.Token{
		{
			token.Unspecific,
			c.document.Text(),
		},
	}
}

func (c *CompleteCode) Complete() string {
	completions := c.GetCompletions()
	r := c.document.CharBeforeCursor()
	if len(r) == 0 {
		return ""
	}
	for _, completion := range completions {
		if strings.HasPrefix(completion.Display, r) {
			return completion.Suffix[len(r):]
		}
	}
	return ""
}

func (c *CompleteCode) GetCompletions() []*startprompt.Completion {
	// 仅做展示，所以返回固定值
	return []*startprompt.Completion{
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

func (c *CompleteCode) ContinueInput() bool {
	return false
}

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		NewCodeFunc: newCompleteCode,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	line, err := c.ReadInput()
	if err != nil {
		fmt.Printf("ReadInput error: %v\n", err)
		return
	}
	fmt.Println("echo:", line)
}
