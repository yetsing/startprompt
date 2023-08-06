package main

import (
	"fmt"
	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/lexer"
)

func main() {
	c, err := startprompt.NewCommandLine(lexer.GetMonkeyTokens, true)
	defer c.Restore()
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	for c.Running() {
		line := c.ReadInput()
		c.OutputStringf("echo: %s\r\n", line)
	}
}
