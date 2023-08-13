package startprompt

import (
	"unicode/utf8"
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

type _UndoEntry struct {
	text           string
	cursorPosition int
}

type Line struct {
	buffer []rune
	// 光标在文本 buffer 中的位置
	cursorPosition int
	undoStack      []*_UndoEntry

	render        *rRenderer
	newCodeFunc   NewCodeFunc
	newPromptFunc NewPromptFunc

	accept bool
	abort  bool
}

func newLine(render *rRenderer, newCodeFunc NewCodeFunc, newPromptFunc NewPromptFunc) *Line {
	return &Line{
		buffer:         nil,
		cursorPosition: 0,
		abort:          false,
		render:         render,
		newCodeFunc:    newCodeFunc,
		newPromptFunc:  newPromptFunc,
	}
}

func (l *Line) reset() {
	l.buffer = nil
	l.cursorPosition = 0
	l.undoStack = nil
	l.accept = false
	l.abort = false
}

func (l *Line) GetCursorPosition() int {
	return l.cursorPosition
}

func (l *Line) SetCursorPosition(v int) {
	l.cursorPosition = maxInt(0, v)
}

// SaveToUndoStack 保存当前信息（文本和光标位置），支持 undo 操作
func (l *Line) SaveToUndoStack() {
	entry := &_UndoEntry{
		text:           string(l.buffer),
		cursorPosition: l.cursorPosition,
	}
	// 如果文本与最后一个的不相同，保存当前信息
	length := len(l.undoStack)
	if length == 0 || l.undoStack[length-1].text != entry.text {
		l.undoStack = append(l.undoStack, entry)
	}
}

func (l *Line) Home() {
	l.cursorPosition = 0
}

func (l *Line) End() {
	l.cursorPosition = len(l.buffer)
}

// Abort 放弃输入（一般是用户按下 Ctrl-C）
func (l *Line) Abort() {
	l.abort = true
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

// CursorToStartOfLine 移动光标到当前行的行首
// afterWhitespace 为 true ，则是将光标移动到当前行第一个非空字符处
func (l *Line) CursorToStartOfLine(afterWhitespace bool) {
	document := l.Document()
	l.cursorPosition -= utf8.RuneCountInString(document.currentLineBeforeCursor())

	if afterWhitespace {
		ws := document.leadingWhitespaceInCurrentLine()
		l.cursorPosition += utf8.RuneCountInString(ws)
	}
}

func (l *Line) CursorToEndOfLine() {
	l.cursorPosition += utf8.RuneCountInString(l.Document().currentLineAfterCursor())
}

func (l *Line) DeleteCharacterBeforeCursor() string {
	if l.cursorPosition == 0 {
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition-1, 1)
	l.cursorPosition--
	return string(deleted)
}

// DeleteCharacterAfterCursor 删除光标后面指定数量的字符并返回删除的字符
func (l *Line) DeleteCharacterAfterCursor(count int) string {
	if l.cursorPosition >= len(l.buffer) {
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition, count)
	return string(deleted)
}

func (l *Line) Newline() {
	l.InsertText([]rune{'\n'})
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

func (l *Line) insertRunes(index int, value []rune) {
	for i, r := range value {
		l.insertRune(index+i, r)
	}
}

func (l *Line) removeRunes(index int, count int) []rune {
	buffer := l.buffer
	removeEnd := index + count
	if removeEnd > len(l.buffer) {
		removeEnd = len(l.buffer)
	}
	ret := buffer[index:removeEnd]
	buffer = append(buffer[:index], buffer[removeEnd:]...)
	l.buffer = buffer
	return ret
}

func (l *Line) ReturnInput() {
	l.accept = true
}

func (l *Line) Exit() {
	l.abort = true
}

func (l *Line) HasText() bool {
	return len(l.buffer) > 0
}

// ListCompletions 展示所有补全 todo
func (l *Line) ListCompletions() {
	results := l.CreateCodeObj().GetCompletions()
	if len(results) > 0 && l.render != nil {
		l.render.renderCompletions(results)
	}
}

// Complete 自动补全，如果有补全返回 true
func (l *Line) Complete() bool {
	result := l.CreateCodeObj().Complete()
	if len(result) > 0 {
		runes := []rune(result)
		l.insertRunes(l.cursorPosition, runes)
		l.cursorPosition += len(runes)
		return true
	} else {
		return false
	}
}

func (l *Line) IsMultiline() bool {
	res := l.CreateCodeObj().IsMultiline()
	DebugLog("IsMultiline: %v", res)
	return res
}

func (l *Line) text() string {
	return string(l.buffer)
}

func (l *Line) Document() *Document {
	s := string(l.buffer)
	return NewDocument(s, l.cursorPosition)
}

func (l *Line) CreateCodeObj() Code {
	return l.newCodeFunc(l.Document())
}

func (l *Line) GetRenderContext() *RenderContext {
	code := l.CreateCodeObj()
	prompt := l.newPromptFunc(l, code)
	return newRenderContext(
		prompt,
		code,
		l.Document(),
		l.accept, l.abort)
}
