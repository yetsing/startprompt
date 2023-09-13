package main

/*
多行输入
需要按下两次 Enter 才能结束输入
*/

import (
	"fmt"
	"time"

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
	//    需要连续按下两次 Enter 才结束当前输入
	return c.document.EmptyLineCountAtTheEnd() == 0
}

func (c *MultilineCode) CompleteAfterInsertText() bool {
	return false
}

func main() {
	c, err := startprompt.NewTCommandLine(&startprompt.CommandLineOption{
		CodeFactory: newMultilineCode,
		AutoIndent:  true,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	defer c.Close()
	c.Println("Type multiline text. Press twice Enter confirm or Ctrl-D exit")
	line, err := c.ReadInput()
	if err != nil {
		c.Printf("ReadInput error: %v\n", err)
		return
	}
	c.Println("echo:", line)
	time.Sleep(1 * time.Second)
}
