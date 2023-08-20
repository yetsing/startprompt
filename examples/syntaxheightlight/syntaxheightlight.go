package main

/*
展示语法高亮的用法和效果
*/

import (
	"fmt"
	"github.com/yetsing/startprompt"
	"github.com/yetsing/startprompt/lexer"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
	"unicode"
)

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// 一个简单的分词器，区分整数、字符串
func getTokens(text string) []token.Token {
	var tokens []token.Token
	codeBuffer := lexer.NewCodeBuffer(text)
	for codeBuffer.HasChar() {
		r := codeBuffer.CurrentChar()
		// 解析空格
		if unicode.IsSpace(r) {
			start := codeBuffer.GetIndex()
			for unicode.IsSpace(codeBuffer.CurrentChar()) {
				codeBuffer.Advance(1)
			}
			end := codeBuffer.GetIndex()
			tk := token.NewToken(token.Whitespace, codeBuffer.Slice(start, end))
			tokens = append(tokens, tk)
			continue
		}

		// 解析整数
		if isdigit(r) {
			start := codeBuffer.GetIndex()
			for isdigit(codeBuffer.CurrentChar()) {
				codeBuffer.Advance(1)
			}
			end := codeBuffer.GetIndex()
			tk := token.NewToken(token.NumberInteger, codeBuffer.Slice(start, end))
			tokens = append(tokens, tk)
			continue
		}

		// 解析双引号字符串
		if r == '"' {
			start := codeBuffer.GetIndex()
			// 跳过开始的双引号
			codeBuffer.Advance(1)
			for codeBuffer.HasChar() && codeBuffer.CurrentChar() != '"' {
				codeBuffer.Advance(1)
			}
			// 跳过结束的双引号
			codeBuffer.Advance(1)
			end := codeBuffer.GetIndex()
			tk := token.NewToken(token.String, codeBuffer.Slice(start, end))
			tokens = append(tokens, tk)
			continue
		}

		tk := token.NewToken(token.Text, string(r))
		tokens = append(tokens, tk)
		codeBuffer.Advance(1)
	}
	return tokens
}

type HighlightCode struct {
	document *startprompt.Document
}

func newHighlightCode(document *startprompt.Document) startprompt.Code {
	return &HighlightCode{document: document}
}

func (c *HighlightCode) GetTokens() []token.Token {
	return getTokens(c.document.Text())
}

func (c *HighlightCode) Complete() string {
	return ""
}

func (c *HighlightCode) GetCompletions() []*startprompt.Completion {
	return nil
}

func (c *HighlightCode) ContinueInput() bool {
	return false
}

func main() {
	var schema = startprompt.Schema{
		token.Number: terminalcolor.NewFgColorStyleHex("#4E9A06"),
		token.String: terminalcolor.NewFgColorStyleHex("#C4A000"),
	}

	c, err := startprompt.NewCommandLine(&startprompt.CommandLineOption{
		// 如果想要自定义每个 token 的颜色，可以指定 schema
		Schema:      schema,
		NewCodeFunc: newHighlightCode,
	})
	if err != nil {
		fmt.Printf("failed to startprompt.NewCommandLine: %v\n", err)
		return
	}
	line, err := c.ReadInput()
	if err != nil {
		fmt.Printf("ReadInput error: %v\n", err)
		return
	}
	fmt.Println("echo:", line)
}
