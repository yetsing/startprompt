package main

/*
展示补全的用法和效果

按下 Tab 补全当前单词
- 如果只有一个匹配的补全，补全文本直接添加在后面
- 如果有多个，展示所有补全， Ctrl-P 和 Ctrl-N 上下移动选择补全项，按 Tab 则会使用当前选中的补全
- 按 Esc 或者 Ctrl+[ 退出补全（注意：需要按两下）
*/

import (
	"fmt"
	"strings"

	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/token"
)

type AnimalCode struct {
	document *startprompt.Document
	animals  []string
}

func newCompleteCode(document *startprompt.Document) startprompt.Code {
	return &AnimalCode{
		document: document,
		animals: []string{
			"bat",
			"bear",
			"beaver",
			"bee",
			"bison",
			"butterfly",
			"cat",
			"chicken",
			"crocodile",
			"dinosaur",
			"dog",
			"dolphine",
			"dove",
			"duck",
			"eagle",
			"elephant",
			"fish",
			"goat",
			"gorilla",
			"kangoroo",
			"leopard",
			"lion",
			"mouse",
			"rabbit",
			"rat",
			"snake",
			"spider",
			"turkey",
			"turtle",
		},
	}
}

func (c *AnimalCode) GetTokens() []token.Token {
	return []token.Token{
		{
			token.Unspecific,
			c.document.Text(),
		},
	}
}

func (c *AnimalCode) Complete() string {
	completions := c.GetCompletions()
	if len(completions) == 1 {
		return completions[0].Suffix
	}
	return ""
}

func (c *AnimalCode) GetCompletions() []*startprompt.Completion {
	word := c.document.GetWordBeforeCursor()

	var completions []*startprompt.Completion
	for _, animal := range c.animals {
		if strings.HasPrefix(animal, word) {
			cp := &startprompt.Completion{
				Display: animal,
				Suffix:  animal[len(word):],
			}
			completions = append(completions, cp)
		}
	}
	return completions
}

func (c *AnimalCode) ContinueInput() bool {
	return false
}

func main() {
	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		CodeFactory: newCompleteCode,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	fmt.Println("Press tab to complete")
	line, err := c.ReadInput()
	if err != nil {
		fmt.Printf("ReadInput error: %v\n", err)
		return
	}
	fmt.Println("echo:", line)
}
