package startprompt

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yetsing/startprompt/enums/linemode"
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

// 切换 line 到正确的模式（mode）
func toMode(line *Line, modes ...linemode.LineMode) {
	if line.mode.In(modes...) {
		return
	}
	// todo something
	if line.mode.Is(linemode.IncrementalSearch) {

	} else if line.mode.Is(linemode.Complete) {
		line.mode = linemode.Normal
	}

}

type cCompletionState struct {
	// 补全开始时的 document
	originalDocument *Document
	// 当前的补全列表
	currentCompletions []*Completion
	// 当前补全位置
	completeIndex int
}

func newCompletionState(originalDocument *Document, currentCompletions []*Completion) *cCompletionState {
	return &cCompletionState{
		originalDocument:   originalDocument,
		currentCompletions: currentCompletions,
		completeIndex:      -1,
	}
}

func (c *cCompletionState) originalCursorPosition() int {
	return c.originalDocument.CursorPosition()
}

func (c *cCompletionState) currentCompletionText() string {
	if c.completeIndex == -1 {
		return ""
	} else {
		return c.currentCompletions[c.completeIndex].Suffix
	}
}

type _UndoEntry struct {
	text           string
	cursorPosition int
}

type Line struct {
	// string 类型的下标切片是按字节来的，而不是 Unicode
	// 为了方便增删改，选择使用 []rune
	buffer []rune
	// 光标在文本 buffer 中的位置
	cursorPosition int
	undoStack      []*_UndoEntry
	mode           linemode.LineMode
	completeState  *cCompletionState

	renderer      *rRenderer
	newCodeFunc   NewCodeFunc
	newPromptFunc NewPromptFunc

	history History

	// accept 表示用户输入完成（一般是按下 Enter）
	accept bool
	// abort 表示用户丢弃当前输入（一般是按下 Ctrl-C）
	abort bool
	// exit 表示用户退出（一般是按下 Ctrl-D）
	exit bool

	workingLines []string
	workingIndex int

	// 自动缩进，如果开启，新行的缩进会与上一行保持一致
	autoIndent bool
}

func newLine(
	render *rRenderer,
	newCodeFunc NewCodeFunc,
	newPromptFunc NewPromptFunc,
	history History,
	autoIndent bool,
) *Line {
	line := &Line{
		renderer:      render,
		newCodeFunc:   newCodeFunc,
		newPromptFunc: newPromptFunc,
		history:       history,
		autoIndent:    autoIndent,
	}
	line.reset()
	return line
}

func (l *Line) reset() {
	l.mode = linemode.Normal
	l.buffer = nil
	l.cursorPosition = 0

	l.undoStack = nil

	l.completeState = nil

	l.accept = false
	l.abort = false
	l.exit = false

	lines := l.history.GetAll()
	// +1 是因为当前输入也要占个位置
	l.workingLines = make([]string, len(lines)+1)
	copy(l.workingLines, lines)
	l.workingIndex = len(l.workingLines) - 1
}

func (l *Line) text() string {
	return l.workingLines[l.workingIndex]
}

func (l *Line) setText(buffer []rune) {
	l.buffer = buffer
	l.workingLines[l.workingIndex] = string(buffer)
	l.textChanged()
}

func (l *Line) GetCursorPosition() int {
	return l.cursorPosition
}

func (l *Line) SetCursorPosition(v int) {
	l.cursorPosition = maxInt(0, v)
}

func (l *Line) getWorkingIndex() int {
	return l.workingIndex
}

func (l *Line) setWorkingIndex(value int) {
	l.workingIndex = value
	l.textChanged()
}

func (l *Line) textChanged() {
}

// SaveToUndoStack 保存当前信息（文本和光标位置），支持 undo 操作
func (l *Line) SaveToUndoStack() {
	entry := &_UndoEntry{
		text:           l.text(),
		cursorPosition: l.cursorPosition,
	}
	// 如果文本与最后一个的不相同，保存当前信息
	length := len(l.undoStack)
	if length == 0 || l.undoStack[length-1].text != entry.text {
		l.undoStack = append(l.undoStack, entry)
	}
}

// 替换当前行
func (l *Line) setCurrentLine(value []rune) {
	// 移动光标到行首
	l.CursorToStartOfLine(false)

	// 替换文本
	l.DeleteUntilEndOfLine()
	l.InsertText(value, false)
}

// 转换指定索引范围内的行文本，索引支持负数
func (l *Line) transformLines(start int, end int, transformCallback func(s string) string) {
	lines := strings.Split(l.text(), "\n")
	length := len(lines)
	start = limitInt(start, length)
	end = limitInt(end, length)
	for i := start; i < end; i++ {
		if i >= length {
			break
		}
		lines[i] = transformCallback(lines[i])
	}
	result := strings.Join(lines, "\n")
	l.setText([]rune(result))
}

func (l *Line) Document() *Document {
	return NewDocument(l.text(), l.cursorPosition)
}

func (l *Line) Home() {
	toMode(l, linemode.Normal)
	l.SetCursorPosition(0)
}

func (l *Line) End() {
	toMode(l, linemode.Normal)
	l.SetCursorPosition(len(l.buffer))
}

func (l *Line) CursorLeft() {
	toMode(l, linemode.Normal)
	if l.cursorPosition > 0 {
		l.SetCursorPosition(l.cursorPosition - 1)
	}
}

func (l *Line) CursorRight() {
	toMode(l, linemode.Normal)
	if l.cursorPosition < len(l.buffer) {
		l.SetCursorPosition(l.cursorPosition + 1)
	}
}

func (l *Line) CursorUp() {
	toMode(l, linemode.Normal)
	newPos := l.Document().CursorUpPosition()
	if newPos != -1 {
		l.SetCursorPosition(newPos)
	}
}

func (l *Line) CursorDown() {
	toMode(l, linemode.Normal)
	newPos := l.Document().CursorDownPosition()
	if newPos != -1 {
		l.SetCursorPosition(newPos)
	}
}

// AutoUp 根据不同的情况，向上触发不同的效果，如下
// 如果在补全状态下，会移动到上一个补全
// 如果光标不在第一行，移动光标到上一行
// 否则切换到上一个历史输入
func (l *Line) AutoUp() {
	toMode(l, linemode.Normal, linemode.Complete)
	if l.mode.Is(linemode.Complete) {
		l.CompletePrevious(1)
	} else if l.Document().CursorPositionRow() > 0 {
		l.CursorUp()
	} else {
		l.HistoryBackward()
	}
}

// AutoDown 根据不同的情况，向上触发不同的效果，如下
// 如果在补全状态下，会移动到下一个补全
// 如果光标不在第一行，移动光标到下一行
// 否则切换到下一个历史输入
func (l *Line) AutoDown() {
	if l.mode.In(linemode.Complete) {
		l.CompleteNext(1)
	} else if l.Document().CursorPositionRow() > 0 {
		l.CursorDown()
	} else {
		l.HistoryForward()
	}
}

// CursorWordBack 移动光标到前一个单词的开头
func (l *Line) CursorWordBack() {
	toMode(l, linemode.Normal)
	l.SetCursorPosition(l.cursorPosition + l.Document().findStartOfPreviousWord())
}

// CursorWordForward 移动光标到下一个单词的开头
func (l *Line) CursorWordForward() {
	toMode(l, linemode.Normal)
	l.SetCursorPosition(l.cursorPosition + l.Document().findNextWordBeginning())
}

// CursorToEndOfWord 向右移动光标到单词的最后一个字符的左边
func (l *Line) CursorToEndOfWord() {
	toMode(l, linemode.Normal)
	pos := l.Document().findNextWordEnding(false)
	if pos > 1 {
		// 因为要移动的最后一个字符的左边，所以减一
		l.SetCursorPosition(l.cursorPosition + pos - 1)
	}
}

// CursorToEndOfLine 移动光标到当前行的行尾
func (l *Line) CursorToEndOfLine() {
	toMode(l, linemode.Normal)
	count := utf8.RuneCountInString(l.Document().CurrentLineAfterCursor())
	l.SetCursorPosition(l.cursorPosition + count)
}

// CursorToStartOfLine 移动光标到当前行的行首
// afterWhitespace 为 true ，则是将光标移动到当前行第一个非空字符处
func (l *Line) CursorToStartOfLine(afterWhitespace bool) {
	toMode(l, linemode.Normal)
	document := l.Document()
	l.SetCursorPosition(l.cursorPosition - utf8.RuneCountInString(document.CurrentLineBeforeCursor()))

	if afterWhitespace {
		ws := document.LeadingWhitespaceInCurrentLine()
		l.SetCursorPosition(l.cursorPosition + utf8.RuneCountInString(ws))
	}
}

func (l *Line) DeleteCharacterBeforeCursor(count int) string {
	toMode(l, linemode.Normal, linemode.IncrementalSearch)
	if l.cursorPosition == 0 {
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition-count, count)
	l.SetCursorPosition(l.cursorPosition - len(deleted))
	return string(deleted)
}

// DeleteCharacterAfterCursor 删除光标后面指定数量的字符并返回删除的字符
func (l *Line) DeleteCharacterAfterCursor(count int) string {
	toMode(l, linemode.Normal)
	if l.cursorPosition >= len(l.buffer) {
		// 光标后没有字符可删
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition, count)
	return string(deleted)
}

// DeleteWord 删除光标后的单词
func (l *Line) DeleteWord() string {
	toMode(l, linemode.Normal)
	toDelete := l.Document().findNextWordBeginning()
	return l.DeleteCharacterAfterCursor(toDelete)
}

// DeleteWordBeforeCursor 删除光标前的单词，返回删除的单词
func (l *Line) DeleteWordBeforeCursor() string {
	toMode(l, linemode.Normal)
	toDelete := -l.Document().findStartOfPreviousWord()
	if toDelete > 0 {
		return l.DeleteCharacterBeforeCursor(toDelete)
	}
	return ""
}

// DeleteUntilEndOfLine 删除从光标到行尾处的字符，返回删除的文本
func (l *Line) DeleteUntilEndOfLine() string {
	toMode(l, linemode.Normal)
	after := l.Document().CurrentLineAfterCursor()
	l.DeleteCharacterAfterCursor(utf8.RuneCountInString(after))
	return after
}

// DeleteFromStartOfLine 删除从行首到光标处的字符，返回删除的文本
func (l *Line) DeleteFromStartOfLine() string {
	toMode(l, linemode.Normal)
	before := l.Document().CurrentLineBeforeCursor()
	l.DeleteCharacterBeforeCursor(utf8.RuneCountInString(before))
	return before
}

// DeleteCurrentLine 删除当前行，返回删除的文本
func (l *Line) DeleteCurrentLine() string {
	toMode(l, linemode.Normal)
	document := l.Document()

	deleted := document.CurrentLine()

	// 删除对应行
	lines := document.lines()
	row := document.CursorPositionRow()
	newLines := make([]string, len(lines)-1)
	copy(newLines, lines[:row])
	copy(newLines[row:], lines[row+1:])
	l.setText([]rune(strings.Join(newLines, "\n")))

	// 移动光标到新行文本的第一个字符位置
	beforeCursor := document.CurrentLineBeforeCursor()
	l.SetCursorPosition(l.cursorPosition - utf8.RuneCountInString(beforeCursor))
	l.CursorToStartOfLine(true)
	return deleted
}

// JoinNextLine 将当前行和下一行拼接为一行
func (l *Line) JoinNextLine() {
	toMode(l, linemode.Normal)
	l.CursorToEndOfLine()
	l.DeleteCharacterAfterCursor(1)
}

// SwapCharactersBeforeCursor 交换光标前两个字符
func (l *Line) SwapCharactersBeforeCursor() {
	toMode(l, linemode.Normal)
	pos := l.cursorPosition
	if pos >= 2 {
		a := l.buffer[pos-2]
		b := l.buffer[pos-1]

		result := concatRunes(
			l.buffer[:pos-2], []rune{b, a}, l.buffer[pos:])
		l.setText(result)
	}
}

// GotoMatchingBracket 跳转到匹配 [ ( { < 的括号
func (l *Line) GotoMatchingBracket() {
	toMode(l, linemode.Normal)
	brackets := []struct {
		left  string
		right string
	}{
		{"(", ")"},
		{"[", "]"},
		{"{", "}"},
		{"<", ">"},
	}
	document := l.Document()
	stack := 1
	for _, bracket := range brackets {
		if document.CurrentChar() == bracket.left {
			// 寻找匹配的右括号
			text := document.TextAfterCursor()
			step := 0
			for _, r := range stringStartAt(text, 1) {
				if string(r) == bracket.left {
					stack++
				} else if string(r) == bracket.right {
					stack--
				}
				if stack == 0 {
					// 是从 1 开始遍历的，所以这里需要加 1
					l.SetCursorPosition(l.cursorPosition + step + 1)
					break
				}
				step++
			}
		} else if document.CurrentChar() == bracket.right {
			// 寻找匹配的左括号
			text := document.TextBeforeCursor()
			text = reverseString(text)
			step := 0
			for _, r := range text {
				if string(r) == bracket.right {
					stack++
				} else if string(r) == bracket.left {
					stack--
				}
				if stack == 0 {
					// 比如这种情况 () 光标在括号中间
					// stack == 0 的时候， step = 0 ，需要向左移动一格，所以还需要减一
					l.SetCursorPosition(l.cursorPosition - step - 1)
					break
				}
				step++
			}
		}
	}
}

func (l *Line) CreateCodeObj() Code {
	return l.newCodeFunc(l.Document())
}

// ListCompletions 列出所有补全
func (l *Line) ListCompletions() {
	toMode(l, linemode.Normal)
	results := l.CreateCodeObj().GetCompletions()
	if len(results) > 0 && l.renderer != nil {
		l.renderer.renderCompletions(results)
	}
}

// Complete 自动补全，如果有补全返回 true
func (l *Line) Complete() bool {
	toMode(l, linemode.Normal)
	result := l.CreateCodeObj().Complete()
	if len(result) > 0 {
		runes := []rune(result)
		l.InsertText(runes, true)
		return true
	} else {
		return false
	}
}

// CompleteNext 选择下面第 count 个补全
func (l *Line) CompleteNext(count int) {
	toMode(l, linemode.Normal, linemode.Complete)
	if !l.mode.Is(linemode.Complete) {
		l.StartComplete()
	} else {
		completionsCount := len(l.completeState.currentCompletions)

		var index int
		if l.completeState.completeIndex == -1 {
			index = 0
		} else if l.completeState.completeIndex == completionsCount-1 {
			index = -1
		} else {
			index = minInt(completionsCount-1, l.completeState.completeIndex+count)
		}
		l.gotoCompletion(index)
	}
}

// CompletePrevious 选择上面第 count 个补全
func (l *Line) CompletePrevious(count int) {
	toMode(l, linemode.Normal, linemode.Complete)

	if !l.mode.Is(linemode.Complete) {
		l.StartComplete()
	}

	if l.completeState != nil {
		var index int
		if l.completeState.completeIndex == 0 {
			index = -1
		} else if l.completeState.completeIndex == -1 {
			index = len(l.completeState.currentCompletions) - 1
		} else {
			index = maxInt(0, l.completeState.completeIndex-count)
		}
		l.gotoCompletion(index)
	}
}

// StartComplete 开始补全
func (l *Line) StartComplete() {
	currentCompletions := l.CreateCodeObj().GetCompletions()

	if len(currentCompletions) > 0 {
		l.completeState = newCompletionState(l.Document(), currentCompletions)
		text := l.completeState.currentCompletionText()
		l.InsertText([]rune(text), true)
		l.mode = linemode.Complete
	} else {
		l.mode = linemode.Normal
		l.completeState = nil
	}
}

// ExitComplete 退出补全
func (l *Line) ExitComplete() {
	l.mode = linemode.Normal
	l.completeState = nil
}

// 选择指定位置的补全
func (l *Line) gotoCompletion(index int) {
	if !l.mode.Is(linemode.Complete) {
		panic(fmt.Sprintf("line mode want=Complete, but got=%s", l.mode))
	}

	// 撤销之前的补全
	count := utf8.RuneCountInString(l.completeState.currentCompletionText())
	if count > 0 {
		l.DeleteCharacterBeforeCursor(count)
	}

	// 设置新的补全
	l.completeState.completeIndex = index
	l.InsertText([]rune(l.completeState.currentCompletionText()), true)

	l.mode = linemode.Complete
}

// GetRenderContext 返回渲染上下文信息
func (l *Line) GetRenderContext() *RenderContext {
	code := l.CreateCodeObj()
	prompt := l.newPromptFunc(l, code)
	var completeState *cCompletionState
	if l.mode.Is(linemode.Complete) && !l.abort && !l.accept {
		completeState = l.completeState
	} else {
		completeState = nil
	}

	return newRenderContext(
		prompt,
		code,
		completeState,
		l.Document(),
		l.accept,
		l.abort,
		l.exit,
	)
}

// HistoryForward 选择下一个历史输入
func (l *Line) HistoryForward() {
	toMode(l, linemode.Normal)
	if l.workingIndex < len(l.workingLines)-1 {
		l.workingIndex++
		l.buffer = []rune(l.workingLines[l.workingIndex])
		l.SetCursorPosition(len(l.buffer))
	}
}

// HistoryBackward 选择上一个历史输入
func (l *Line) HistoryBackward() {
	toMode(l, linemode.Normal)
	if l.workingIndex > 0 {
		l.workingIndex--
		l.buffer = []rune(l.workingLines[l.workingIndex])
		l.SetCursorPosition(len(l.buffer))
	}
}

func (l *Line) Newline() {
	toMode(l, linemode.Normal)
	spaces := l.Document().LeadingWhitespaceInCurrentLine()
	l.InsertText([]rune{'\n'}, true)
	if l.autoIndent {
		l.InsertText([]rune(spaces), true)
	}
}

// InsertLineAbove 在当前行的上方插入一个空行，
// copyMargin 表示新行前面是否要保持同样的空格
func (l *Line) InsertLineAbove(copyMargin bool) {
	var insert string
	if copyMargin {
		insert = l.Document().LeadingWhitespaceInCurrentLine() + "\n"
	} else {
		insert = "\n"
	}

	l.CursorToStartOfLine(false)
	l.InsertText([]rune(insert), true)
	l.SetCursorPosition(l.cursorPosition - 1)
}

// InsertLineBelow 在当前行的下方插入一个空行，
// copyMargin 表示新行前面是否要保持同样的空格
func (l *Line) InsertLineBelow(copyMargin bool) {
	var insert string
	if copyMargin {
		insert = "\n" + l.Document().LeadingWhitespaceInCurrentLine()
	} else {
		insert = "\n"
	}

	l.CursorToEndOfLine()
	l.InsertText([]rune(insert), true)
}

// OverwriteText 覆盖光标位置到行尾最多 len(data) 长度数据，
// moveCursor 表示插入后是否移动光标
func (l *Line) OverwriteText(data []rune, moveCursor bool) {
	overwrittenRunes := sliceRunes(l.buffer, l.cursorPosition, l.cursorPosition+len(data))
	nlIndex := findRunes(overwrittenRunes, '\n')
	// 最多覆盖到行尾
	if nlIndex > -1 {
		overwrittenRunes = sliceRunes(l.buffer, l.cursorPosition, nlIndex)
	}
	result := concatRunes(
		l.buffer[:l.cursorPosition],
		data,
		l.buffer[l.cursorPosition+len(overwrittenRunes):],
	)
	l.setText(result)
	if moveCursor {
		l.SetCursorPosition(l.cursorPosition + len(data))
	}
}

// InsertText 在 cursorPosition 的位置插入 data
// moveCursor 表示插入后是否移动光标
func (l *Line) InsertText(data []rune, moveCursor bool) {
	result := concatRunes(l.buffer[:l.cursorPosition], data, l.buffer[l.cursorPosition:])
	l.setText(result)
	if moveCursor {
		l.SetCursorPosition(l.cursorPosition + len(data))
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
	l.setText(buffer)
	return ret
}

func (l *Line) Undo() {
	toMode(l, linemode.Normal)
	if len(l.undoStack) > 0 {
		top := l.undoStack[len(l.undoStack)-1]
		l.undoStack = l.undoStack[:len(l.undoStack)-1]
		l.setText([]rune(top.text))
		l.SetCursorPosition(top.cursorPosition)
	}
}

// Abort 丢弃输入（一般是用户按下 Ctrl-C）
func (l *Line) Abort() {
	toMode(l, linemode.Normal)
	l.abort = true
}

// Exit 停止输入（一般是用户按下 Ctrl-D）
func (l *Line) Exit() {
	toMode(l, linemode.Normal)
	l.exit = true
}

func (l *Line) ReturnInput() {
	toMode(l, linemode.Normal)
	text := l.text()

	// 文本与最后一个不相同时，保存到历史中
	if l.history.Length() == 0 || l.history.GetAt(-1) != text {
		if len(text) > 0 {
			l.history.Append(text)
		}
	}

	l.accept = true
}

func (l *Line) HasText() bool {
	return len(l.buffer) > 0
}

func (l *Line) IsMultiline() bool {
	res := l.CreateCodeObj().ContinueInput()
	return res
}

func (l *Line) Clear() {
	toMode(l, linemode.Normal)
	if l.renderer != nil {
		l.renderer.clear()
	}
}
