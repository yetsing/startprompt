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
	c, err := startprompt.NewCommandLine(nil)
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
