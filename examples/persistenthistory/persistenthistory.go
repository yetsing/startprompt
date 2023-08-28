package main

/*
存储历史输入到文件中
*/

import (
	"errors"
	"fmt"

	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		History: startprompt.NewFileHistory(".example-history-file"),
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	for {
		line, err := c.ReadInput()
		if err != nil {
			if errors.Is(err, startprompt.ExitError) {
				break
			}
			fmt.Printf("ReadInput error: %v\n", err)
			break
		}
		fmt.Println("echo:", line)
	}
}
