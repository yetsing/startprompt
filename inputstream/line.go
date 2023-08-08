package inputstream

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
	buffer []rune
	// 光标在文本 buffer 中的位置
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

func (l *Line) text() string {
	return string(l.buffer)
}

func (l *Line) Document() *Document {
	s := string(l.buffer)
	return NewDocument(s, l.cursorPosition)
}

//type Document struct {
//	Text    string
//	CursorX int
//	buffer  []rune
//	// 光标在文本 buffer 中的位置
//	cursorPosition int
//}
//
//// 返回光标位置的字符，如果不存在返回 0
//func (d *Document) currentChar() rune {
//	return d.getCharRelativeToCursor(0)
//}
//
//// 返回光标前的字符，如果不存在返回 0
//func (d *Document) charBeforeCursor() rune {
//	return d.getCharRelativeToCursor(-1)
//}
//
//// 返回光标前的文本
//func (d *Document) textBeforeCursor() string {
//	return string(d.buffer[:d.cursorPosition])
//}
//
//// 返回光标后的文本
//func (d *Document) textAfterCursor() string {
//	return string(d.buffer[d.cursorPosition:])
//}
//
//// 返回从行首到光标处的文本
//func (d *Document) currentLineBeforeCursor() string {
//	text := d.textBeforeCursor()
//	splited := strings.Split(text, "\n")
//	return splited[len(splited)-1]
//}
//
//// 返回从光标到行尾的文本
//func (d *Document) currentLineAfterCursor() string {
//	text := d.textAfterCursor()
//	splited := strings.Split(text, "\n")
//	return splited[0]
//}
//
//// 返回行数组
//func (d *Document) lines() []string {
//	text := string(d.buffer)
//	return strings.Split(text, "\n")
//}
//
//// 返回当前行到最后一行的数组
//func (d *Document) linesFromCurrent() []string {
//	row := d.cursorPositionRow()
//	return d.lines()[row:]
//}
//
//// 返回行数
//func (d *Document) lineCount() int {
//	return len(d.lines())
//}
//
//// 返回光标所在行
//func (d *Document) currentLine() string {
//	return d.currentLineBeforeCursor() + d.currentLineAfterCursor()
//}
//
//// 返回当前行的行首位置的空白字符
//func (d *Document) leadingWhitespaceInCurrentLine() string {
//	currentLine := d.currentLine()
//	var i int
//	var r rune
//	for i, r = range currentLine {
//		// 找到第一个不是空白字符
//		if !unicode.IsSpace(r) {
//			break
//		}
//	}
//	return currentLine[:i]
//}
//
//// 返回相对于 cursorPosition 的字符，如果不存在返回 0
//func (d *Document) getCharRelativeToCursor(offset int) rune {
//	index := d.cursorPosition + offset
//	if index < 0 || index >= len(d.buffer) {
//		return 0
//	}
//	return d.buffer[index]
//}
//
//// 返回光标所在行号（从 0 开始计数）
//func (d *Document) cursorPositionRow() int {
//	text := d.textBeforeCursor()
//	return len(strings.Split(text, "\n")) - 1
//}
//
//// 返回光标所在列号（从 0 开始计数）
//func (d *Document) cursorPositionCol() int {
//	return len(d.textBeforeCursor())
//}
//
//// 将文本索引翻译为相应的行和列号
//func (d *Document) translateIndexToPosition(index int) (int, int) {
//	text := string(d.buffer[:index])
//	lines := strings.Split(text, "\n")
//	row := len(lines) - 1
//	col := len(lines[len(lines)-1])
//	return row, col
//}
//
//// 将行和列号翻译为文本索引
//func (d *Document) translateRowColToIndex(row int, col int) int {
//	lines := d.lines()
//	index := 0
//	for i, line := range lines {
//		if i >= row {
//			break
//		}
//		index += len(line)
//	}
//	// 有 row 个换行符，再加上列号
//	index += row + col
//	return index
//}
