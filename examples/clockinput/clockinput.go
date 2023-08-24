package main

import (
	"fmt"
	"time"

	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/token"
)

/*
动态提示符例子
在提示符中展示当前时间，可以看到提示符随时间的变化
*/

type ClockPrompt struct {
	startprompt.BasePrompt
}

func (c *ClockPrompt) GetPrompt() []token.Token {
	now := time.Now()
	return []token.Token{
		token.NewToken(token.Prompt, now.Format("15:04:05")),
		token.NewToken(token.Prompt, " Enter something: "),
	}
}

func NewClockPrompt(_ *startprompt.Line, _ startprompt.Code) startprompt.Prompt {
	return &ClockPrompt{startprompt.BasePrompt{}}
}

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		NewPromptFunc:     NewClockPrompt,
		EnableConcurrency: true,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	c.RunInExecutor(func() {
		for {
			time.Sleep(1 * time.Second)
			c.RequestRedraw()
		}
	})
	line, err := c.ReadInput()
	if err != nil {
		fmt.Printf("ReadInput error: %v\n", err)
		return
	}
	fmt.Println("You said: ", line)
}
