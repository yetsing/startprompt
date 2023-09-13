package main

/*
读取用户输入并将其打印出来
*/

import (
	"errors"
	"fmt"

	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewTCommandLine(&startprompt.CommandLineOption{
		//EnableDebug: true,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewTCommandLine: %v\n", err)
		return
	}
	defer c.Close()
	c.Println("Type some text. Press Enter confirm or Ctrl-D exit")
	for {
		line, err := c.ReadInput()
		if err != nil {
			if errors.Is(err, startprompt.ExitError) {
				break
			}
			c.Printf("ReadInput error: %v\n", err)
			break
		}
		c.Println("echo:", line)
	}
}
