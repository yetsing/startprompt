package main

import (
	"fmt"
	"github.com/yetsing/startprompt"
	"golang.org/x/term"
	"os"
)

func main() {
	// 开启 terminal raw mode
	// 这种模式下会拿到用户原始的输入，比如输入 Ctrl-c 时，不会中断当前程序，而是拿到 Ctrl-c 的表示
	// 不会自动展示用户输入
	// 更多说明解释参考：https://viewsourcecode.org/snaptoken/kilo/02.enteringRawMode.html
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("error make raw mode: %v\n", err)
		return
	}

	defer func(fd int, oldState *term.State) {
		err := term.Restore(fd, oldState)
		if err != nil {
			fmt.Printf("restore error: %v\r\n", err)
		}
	}(int(os.Stdin.Fd()), oldState)

	c := startprompt.NewCommandLine()
	for c.Running() {
		line := c.ReadInput()
		c.OutputStringf("echo: %s\r\n", line)
	}
}
