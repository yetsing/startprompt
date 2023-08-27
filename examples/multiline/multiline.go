package main

/*
多行输入
需要按下两次 Enter 才能结束输入
*/

import (
	"fmt"

	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/token"
)

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
	return c.document.EmptyLineCountAtTheEnd() == 0
}

func (c *MultilineCode) CompleteAfterInsertText() bool {
	return false
}

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		CodeFactory: newMultilineCode,
		AutoIndent:  true,
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
