package startprompt

import (
	"bytes"
	"golang.org/x/term"
)

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
