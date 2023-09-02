package main

import (
	"fmt"

	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewTCommandLine(&startprompt.CommandLineOption{})
	if err != nil {
		fmt.Printf("failed to startprompt.NewTCommandLine: %v\n", err)
		return
	}
	defer c.Close()
	//c.Println("Press Ctrl-D exit")
	_, err = c.ReadInput()
	if err != nil {

	}
}
