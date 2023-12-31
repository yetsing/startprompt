package startprompt

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

func maxInt(a ...int) int {
	m := a[0]
	for _, n := range a {
		if n > m {
			m = n
		}
	}
	return m
}

func minInt(a ...int) int {
	m := a[0]
	for _, n := range a {
		if n < m {
			m = n
		}
	}
	return m
}

// 获取终端窗口大小，参考 https://stackoverflow.com/a/67087586
func getSize(fd int) (int, int) {
	width, height, err := term.GetSize(fd)
	if err != nil {
		panic(err)
	}
	return width, height
}

func repeatByte(c byte, count int) string {
	var b bytes.Buffer
	for i := 0; i < count; i++ {
		b.WriteByte(c)
	}
	return b.String()
}

// 参考 https://stackoverflow.com/a/10030772
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 限制 n 的范围为 [0, high]
func limitInt(n int, high int) int {
	if high <= 0 {
		panic("limitInt require high > 0")
	}
	if n < 0 {
		n += high
	}
	if n < 0 {
		n = 0
	} else if n > high {
		n = high
	}
	return n
}

// Python 方式切片（索引支持负数，索引可以超出范围）
func sliceRunes(runes []rune, start int, end int) []rune {
	length := len(runes)
	return runes[limitInt(start, length):limitInt(end, length)]
}

func findRunes(runes []rune, r rune) int {
	for i, r2 := range runes {
		if r == r2 {
			return i
		}
	}
	return -1
}

func concatRunes(a ...[]rune) []rune {
	resultLength := 0
	for _, runes := range a {
		resultLength += len(runes)
	}
	result := make([]rune, resultLength)
	offset := 0
	for _, runes := range a {
		copy(result[offset:], runes)
		offset += len(runes)
	}
	return result
}

// insertRunes 将 eles 数据插入 data 的 index 位置处
// 相当于 data[:index] + eles + data[index:]
func insertRunes(data []rune, index int, eles []rune) []rune {
	if index >= len(data) {
		return append(data, eles...)
	}
	return concatRunes(data[:index], eles, data[index:])
}

func isWordDelimiter(r rune) bool {
	set := ".()[]{}"
	return unicode.IsSpace(r) || strings.ContainsRune(set, r)
}

func stringStartAt(s string, start int) string {
	c := 0
	for i := range s {
		if c == start {
			return s[i:]
		}
		c++
	}
	return ""
}

// 返回指定宽度左对齐的字符串，不够的右边补空格
func ljustWidth(s string, width int) string {
	diff := width - runewidth.StringWidth(s)
	if diff <= 0 {
		return s
	}
	return s + repeatByte(' ', diff)
}

// StringLjustWidth 返回指定宽度左对齐的字符串，不够的右边补空格
func StringLjustWidth(s string, width int) string {
	return ljustWidth(s, width)
}

// RepeatString 将 n 个字符串 s 拼接在一起
func RepeatString(s string, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(s)
	}
	return sb.String()
}

func IsSpace(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
