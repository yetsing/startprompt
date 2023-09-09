package main

import (
	"fmt"

	"github.com/yetsing/startprompt"
)

/*
   读取用户输入并将其打印出来
*/

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		EnableDebug: true,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	line, err := c.ReadInput()
	if err != nil {
		fmt.Printf("ReadInput error: %v\n", err)
	}
	fmt.Println("echo:", line)
}
