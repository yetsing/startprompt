# startprompt

start build prompt commandline application

开始构建提示符命令行应用

~~复刻划掉~~参考了 [python-prompt-toolkit](https://github.com/prompt-toolkit/python-prompt-toolkit) 库

# hello world

```go
package main

/*
读取用户输入并将其打印出来
*/

import (
	"fmt"
	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewCommandLine(nil)
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
```

 `examples` 文件夹有一些简单的例子

#### more examples

- [startprompt-python-repl](https://github.com/yetsing/startprompt-python-repl) 一个 Python repl

### keybinding

默认快捷键操作可看 [keybinding](./docs/keybinding.md)


# 参考

[python-prompt-toolkit](https://github.com/prompt-toolkit/python-prompt-toolkit)

[pygments](https://github.com/pygments/pygments)

《用Go语言自制解释器》
