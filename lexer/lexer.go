// Package lexer 实现各种分词器
package lexer

import (
	"unicode"

	"github.com/yetsing/startprompt/token"
)

type GetTokensFunc func(input string) []token.Token

type CodeBuffer struct {
	runes     []rune
	length    int
	index     int
	markIndex int
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
	return c.PeekN(1)
}

func (c *CodeBuffer) PeekN(n int) rune {
	if c.index+n < c.length {
		return c.runes[c.index+n]
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
	if c.index >= len(c.runes) {
		c.index = len(c.runes)
	}
}

func (c *CodeBuffer) Unread(step int) {
	c.index -= step
	if c.index < 0 {
		c.index = 0
	}
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

func (c *CodeBuffer) ReadSpace() string {
	start := c.index
	for c.CurrentChar() != '\n' && unicode.IsSpace(c.CurrentChar()) {
		c.Advance(1)
	}
	return c.Slice(start, c.index)
}

// Mark 标记当前位置
func (c *CodeBuffer) Mark() int {
	c.markIndex = c.index
	return c.index
}

// ReadFromMark 读取从标记位置到当前位置的字符串
func (c *CodeBuffer) ReadFromMark() string {
	return string(c.runes[c.markIndex:c.index])
}

// CurrentLine 返回从 index 位置开始到行尾的字符串（包括换行符）
func (c *CodeBuffer) CurrentLine() string {
	start := c.index
	end := len(c.runes)
	for i := c.index; i < len(c.runes); i++ {
		if c.runes[i] == '\n' {
			end = i + 1
			break
		}
	}
	return string(c.runes[start:end])
}

func (c *CodeBuffer) ReadUntil(delimiter rune) string {
	index := c.index
	for c.HasChar() && c.CurrentChar() != delimiter {
		c.Advance(1)
	}
	return string(c.runes[index:c.GetIndex()])
}
