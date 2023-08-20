// Package lexer 实现各种分词器
package lexer

import "github.com/yetsing/startprompt/token"

type GetTokensFunc func(input string) []token.Token

type CodeBuffer struct {
	runes  []rune
	length int
	index  int
}

func NewCodeBuffer(code string) *CodeBuffer {
	c := &CodeBuffer{
		runes:  []rune(code),
		length: 0,
		index:  0,
	}
	c.length = len(c.runes)
	return c
}

func (c *CodeBuffer) CurrentChar() rune {
	if c.index < c.length {
		return c.runes[c.index]
	}
	return 0
}

func (c *CodeBuffer) Read() rune {
	if c.index < c.length {
		r := c.runes[c.index]
		c.index++
		return r
	}
	return 0
}

func (c *CodeBuffer) Peek() rune {
	if c.index+1 < c.length {
		return c.runes[c.index+1]
	}
	return 0
}

func (c *CodeBuffer) PeekString(n int) string {
	start := c.index
	end := c.index + n
	if start > c.length {
		start = c.length
	}
	if end > c.length {
		end = c.length
	}
	return string(c.runes[start:end])
}

func (c *CodeBuffer) GetIndex() int {
	return c.index
}

func (c *CodeBuffer) SetIndex(index int) {
	c.index = index
}

func (c *CodeBuffer) Advance(step int) {
	c.index += step
}

func (c *CodeBuffer) Slice(start int, end int) string {
	if start > c.length {
		start = c.length
	}
	if end > c.length {
		end = c.length
	}
	return string(c.runes[start:end])
}

func (c *CodeBuffer) HasChar() bool {
	return c.index < c.length
}

func (c *CodeBuffer) Eof() bool {
	return c.index >= c.length
}
