package main

import (
	"fmt"

	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewTCommandLine(&startprompt.CommandLineOption{
		EnableDebug: true,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewTCommandLine: %v\n", err)
		return
	}
	defer c.Close()
	c.Println("Type some text. Press Enter confirm or Ctrl-D exit")
	line, err := c.ReadInput()
	if err != nil {
		return
	}
	c.Println("echo:", line)
}
