package main

/*
读取用户输入并将其打印出来
*/

import (
	"fmt"
	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/lexer"
)

func main() {
	c, err := startprompt.NewCommandLine(lexer.GetMonkeyTokens, true)
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	for c.Running() {
		line, err := c.ReadInput()
		if err != nil {
			fmt.Printf("ReadInput error: %v\n", err)
			break
		}
		fmt.Println("echo:", line)
	}
}
