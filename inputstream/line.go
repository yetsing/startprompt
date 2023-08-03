package inputstream

import "github.com/mattn/go-runewidth"

func maxInt(a ...int) int {
	m := a[0]
	for _, n := range a {
		if n > m {
			m = n
		}
	}
	return m
}

type Line struct {
	buffer         []rune
	cursorPosition int
	finished       bool
}

func NewLine() *Line {
	return &Line{
		buffer:         nil,
		cursorPosition: 0,
	}
}

func (l *Line) reset() {
	l.buffer = []rune{}
	l.cursorPosition = 0
	l.finished = false
}

func (l *Line) SetCursorPosition(v int) {
	l.cursorPosition = maxInt(0, v)
}

func (l *Line) Home() {
	l.cursorPosition = 0
}

func (l *Line) End() {
	l.cursorPosition = len(l.buffer)
}

func (l *Line) CursorLeft() {
	if l.cursorPosition > 0 {
		l.cursorPosition--
	}
}

func (l *Line) CursorRight() {
	if l.cursorPosition < len(l.buffer) {
		l.cursorPosition++
	}
}

func (l *Line) DeleteCharacterBeforeCursor() rune {
	if l.cursorPosition == 0 {
		return 0
	}
	deleted := l.removeRune(l.cursorPosition - 1)
	l.cursorPosition--
	return deleted
}

func (l *Line) DeleteCharacterAfterCursor() rune {
	if l.cursorPosition == len(l.buffer) {
		return 0
	}
	deleted := l.removeRune(l.cursorPosition)
	return deleted
}

func (l *Line) InsertText(data []rune) {
	// 在 cursorPosition 的位置插入 data
	if len(data) == 0 {
		return
	}
	for i, r := range data {
		l.insertRune(l.cursorPosition+i, r)
	}
	l.cursorPosition += len(data)
}

func (l *Line) insertRune(index int, value rune) {
	buffer := l.buffer
	if len(buffer) == index {
		buffer = append(buffer, value)
	} else {
		buffer = append(buffer[:index+1], buffer[index:]...)
		buffer[index] = value
	}
	l.buffer = buffer
}

func (l *Line) removeRune(index int) rune {
	buffer := l.buffer
	value := buffer[index]
	buffer = append(buffer[:index], buffer[index+1:]...)
	l.buffer = buffer
	return value
}

func (l *Line) ReturnInput() {
	l.finished = true
}

func (l *Line) Finished() bool {
	return l.finished
}

func (l *Line) Text() string {
	return string(l.buffer)
}

func (l *Line) Document() *Document {
	s := string(l.buffer[:l.cursorPosition])
	return &Document{
		Text: l.Text(),
		// 光标在文件的右边，所以实际的显示要 +1
		CursorX: runewidth.StringWidth(s),
	}
}

type Document struct {
	Text    string
	CursorX int
}
