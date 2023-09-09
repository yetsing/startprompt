# startprompt

start build prompt commandline application

开始构建提示符命令行应用

~~复刻划掉~~参考 [python-prompt-toolkit](https://github.com/prompt-toolkit/python-prompt-toolkit) 库

# hello world

### CommandLine

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

### TCommandLine

```go
package main

/*
helloworld 例子，读取用户输入并将其打印出来
*/

import (
	"fmt"

	"github.com/yetsing/startprompt"
)

func main() {
	c, err := startprompt.NewTCommandLine(nil)
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
```

# 实现特性

- 支持常用快捷键操作，可看 [keybinding](./docs/keybinding.md)
- 支持输入历史（提供内存和文件两种实现）
- 支持语法高亮（通过自定义分词器实现）
- 支持鼠标 (TCommandLine 支持)

有两个实现 `CommandLine` `TCommandLine` ,
`TCommandLine` 基于 [tcell](https://github.com/gdamore/tcell) 实现，增加鼠标支持

# 更多例子

`examples` 文件夹有一些简单的例子，其中以 t 开头的是使用 `TCommandLine` 的例子。

- [startprompt-python-repl](https://github.com/yetsing/startprompt-python-repl) 一个 Python repl

- sqlite cli todo

# 开启 Debug 日志

日志内容会输出到当前目录下的 `startprompt.log` 文件中

```go
c, err := startprompt.NewTCommandLine(&startprompt.CommandLineOption{
    EnableDebug: true,
})
```

# 参考

[python-prompt-toolkit](https://github.com/prompt-toolkit/python-prompt-toolkit)

[pygments](https://github.com/pygments/pygments)

词法分析主要参考下面

[Let's Build A Simple Interpreter](https://github.com/rspivak/lsbasi)

《用Go语言自制解释器》
