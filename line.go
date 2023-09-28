package startprompt

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yetsing/startprompt/enums/linemode"
)

type cCompletionState struct {
	// 补全开始时的 document
	originalDocument *Document
	// 当前的补全列表
	currentCompletions []*Completion
	// 当前补全位置
	completeIndex int
}

func newCompletionState(
	originalDocument *Document,
	currentCompletions []*Completion,
) *cCompletionState {
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

// _LineArea 输入中的选中区域
type _LineArea struct {
	start int
	end   int
}

type _UndoEntry struct {
	text           string
	cursorPosition int
}

type Line struct {
	//    string 类型的下标切片是按字节来的，而不是 Unicode
	//    为了方便增删改，选择使用 []rune
	buffer []rune
	//    光标在文本 buffer 中的位置
	cursorPosition int
	undoStack      []*_UndoEntry
	mode           linemode.LineMode
	completeState  *cCompletionState

	codeFactory   CodeFactory
	promptFactory PromptFactory

	history History

	workingLines []string
	workingIndex int

	//    自动缩进，如果开启，新行的缩进会与上一行保持一致
	autoIndent bool
	//    用户是否确定本次输入
	accept bool

	//    选中区域
	selection _LineArea
	//    取消选中
	//    当有文本被选中时，输入会替换掉选中的文本，此时选中已被取消
	//    需要通知 renderer 不要再渲染选中
	cancelSelection bool
}

func newLine(
	codeFactory CodeFactory,
	history History,
	autoIndent bool,
) *Line {
	line := &Line{
		codeFactory:    codeFactory,
		history:        history,
		cursorPosition: 0,

		autoIndent: autoIndent,
	}
	line.reset()
	return line
}

func (l *Line) reset() {
	l.mode = linemode.Normal
	l.buffer = nil
	l.cursorPosition = 0

	l.completeState = nil

	l.undoStack = nil

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
	// 如果文本与最后一个的相同，只更新光标位置
	length := len(l.undoStack)
	if length > 0 && l.undoStack[length-1].text == l.text() {
		l.undoStack[length-1].cursorPosition = l.cursorPosition
	} else {
		entry := &_UndoEntry{
			text:           l.text(),
			cursorPosition: l.cursorPosition,
		}
		l.undoStack = append(l.undoStack, entry)
	}
}

// TransformLines 转换指定行文本，索引支持负数，会忽略超出范围的部分
//
//	例如想让某几行转成大写：
//		TransformLines(5, 10, strings.ToUpper)
func (l *Line) TransformLines(start int, end int, transformCallback func(string) string) {
	lines := strings.Split(l.text(), "\n")
	length := len(lines)
	start = limitInt(start, length)
	end = limitInt(end, length)
	for i := start; i < end; i++ {
		lines[i] = transformCallback(lines[i])
	}
	result := strings.Join(lines, "\n")
	l.setText([]rune(result))
}

// TransformRegion 转换指定索引区域文本，索引支持负数，会忽略超出范围的部分
func (l *Line) TransformRegion(from int, to int, transformCallback func(string) string) {
	if from >= to {
		panic(fmt.Sprintf("TransformRegion from=%d not less than to=%d", from, to))
	}
	from = limitInt(from, len(l.buffer))
	to = limitInt(to, len(l.buffer))
	transformed := transformCallback(string(l.buffer[from:to]))
	result := concatRunes(
		l.buffer[:from],
		[]rune(transformed),
		l.buffer[to:],
	)
	l.setText(result)
}

func (l *Line) Document() *Document {
	return NewDocument(l.text(), l.cursorPosition)
}

func (l *Line) Home() {
	l.SetCursorPosition(0)
}

func (l *Line) End() {
	l.SetCursorPosition(len(l.buffer))
}

func (l *Line) CursorLeft() {
	if l.cursorPosition > 0 {
		l.SetCursorPosition(l.cursorPosition - 1)
	}
}

func (l *Line) CursorRight() {
	if l.cursorPosition < len(l.buffer) {
		l.SetCursorPosition(l.cursorPosition + 1)
	}
}

func (l *Line) CursorUp() {
	newPos := l.Document().CursorUpPosition()
	if newPos != -1 {
		l.SetCursorPosition(newPos)
	}
}

func (l *Line) CursorDown() {
	newPos := l.Document().CursorDownPosition()
	if newPos != -1 {
		l.SetCursorPosition(newPos)
	}
}

// AutoUp 根据不同的情况，触发不同的效果，如下
// 如果在补全状态下，会移动到上一个补全
// 如果光标不在第一行，移动光标到上一行
// 否则切换到上一个历史输入
func (l *Line) AutoUp() {
	if l.mode.Is(linemode.Complete) {
		l.CompletePrevious(1)
	} else if l.Document().CursorPositionRow() > 0 {
		l.CursorUp()
	} else {
		l.HistoryBackward()
	}
}

// AutoDown 根据不同的情况，触发不同的效果，如下
// 如果在补全状态下，会移动到下一个补全
// 如果光标不在第一行，移动光标到下一行
// 否则切换到下一个历史输入
func (l *Line) AutoDown() {
	if l.mode.In(linemode.Complete) {
		l.CompleteNext(1)
	} else if !l.Document().OnLastLine() {
		l.CursorDown()
	} else {
		l.HistoryForward()
	}
}

// CursorWordBack 移动光标到前一个单词的开头
func (l *Line) CursorWordBack() {
	l.SetCursorPosition(l.cursorPosition + l.Document().findStartOfPreviousWord())
}

// CursorWordForward 移动光标到下一个单词的开头
func (l *Line) CursorWordForward() {
	l.SetCursorPosition(l.cursorPosition + l.Document().findNextWordBeginning())
}

// CursorToEndOfWord 向右移动光标到单词的最后一个字符的左边
func (l *Line) CursorToEndOfWord() {
	pos := l.Document().findNextWordEnding(false)
	if pos > 1 {
		// 因为要移动的最后一个字符的左边，所以减一
		l.SetCursorPosition(l.cursorPosition + pos - 1)
	}
}

// CursorToEndOfLine 移动光标到当前行的行尾
func (l *Line) CursorToEndOfLine() {
	count := utf8.RuneCountInString(l.Document().CurrentLineAfterCursor())
	l.SetCursorPosition(l.cursorPosition + count)
}

// CursorToStartOfLine 移动光标到当前行的行首
// afterWhitespace 为 true ，则是将光标移动到当前行第一个非空字符处
func (l *Line) CursorToStartOfLine(afterWhitespace bool) {
	document := l.Document()
	l.SetCursorPosition(
		l.cursorPosition - utf8.RuneCountInString(document.CurrentLineBeforeCursor()),
	)

	if afterWhitespace {
		ws := document.LeadingWhitespaceInCurrentLine()
		l.SetCursorPosition(l.cursorPosition + utf8.RuneCountInString(ws))
	}
}

func (l *Line) DeleteCharacterBeforeCursor(count int) string {
	if l.cursorPosition == 0 {
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition-count, count)
	l.SetCursorPosition(l.cursorPosition - len(deleted))
	return string(deleted)
}

// DeleteCharacterAfterCursor 删除光标后面指定数量的字符并返回删除的字符
func (l *Line) DeleteCharacterAfterCursor(count int) string {
	if l.cursorPosition >= len(l.buffer) {
		// 光标后没有字符可删
		return ""
	}
	deleted := l.removeRunes(l.cursorPosition, count)
	return string(deleted)
}

// DeleteWord 删除光标后的单词
func (l *Line) DeleteWord() string {
	toDelete := l.Document().findNextWordBeginning()
	return l.DeleteCharacterAfterCursor(toDelete)
}

// DeleteWordBeforeCursor 删除光标前的单词，返回删除的单词
func (l *Line) DeleteWordBeforeCursor() string {
	toDelete := -l.Document().findStartOfPreviousWord()
	if toDelete > 0 {
		return l.DeleteCharacterBeforeCursor(toDelete)
	}
	return ""
}

// DeleteUntilEndOfLine 删除从光标到行尾处的字符，返回删除的文本
func (l *Line) DeleteUntilEndOfLine() string {
	after := l.Document().CurrentLineAfterCursor()
	l.DeleteCharacterAfterCursor(utf8.RuneCountInString(after))
	return after
}

// DeleteFromStartOfLine 删除从行首到光标处的字符，返回删除的文本
func (l *Line) DeleteFromStartOfLine() string {
	before := l.Document().CurrentLineBeforeCursor()
	l.DeleteCharacterBeforeCursor(utf8.RuneCountInString(before))
	return before
}

// DeleteCurrentLine 删除当前行，返回删除的文本
func (l *Line) DeleteCurrentLine() string {
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
	l.CursorToEndOfLine()
	l.DeleteCharacterAfterCursor(1)
}

// SwapCharactersBeforeCursor 交换光标前两个字符
func (l *Line) SwapCharactersBeforeCursor() {
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
	pos := l.getMatchingBracket()
	if pos >= 0 {
		l.SetCursorPosition(pos)
	}
}

func (l *Line) getMatchingBracket() int {
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
					// l.SetCursorPosition(l.cursorPosition + step + 1)
					// break
					return l.cursorPosition + step + 1
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
					// l.SetCursorPosition(l.cursorPosition - step - 1)
					// break
					return l.cursorPosition - step - 1
				}
				step++
			}
		}
	}
	return -1
}

func (l *Line) CreateCode() Code {
	return l.codeFactory(l.Document())
}

// ListCompletions 列出所有补全
func (l *Line) ListCompletions() {
}

// Complete 自动补全，如果有补全返回 true
func (l *Line) Complete() bool {
	result := l.CreateCode().Complete()
	if len(result) > 0 {
		runes := []rune(result)
		l.insertText(runes, true)
		return true
	} else {
		return false
	}
}

// CompleteNext 选择下面第 count 个补全
func (l *Line) CompleteNext(count int) {
	if !l.mode.Is(linemode.Complete) {
		l.StartComplete(true)
	} else {
		completionsCount := len(l.completeState.currentCompletions)

		var index int
		if l.completeState.completeIndex == -1 {
			index = 0
		} else {
			index = l.completeState.completeIndex + count
			if index >= completionsCount {
				index -= completionsCount
			}
		}
		l.gotoCompletion(index)
	}
}

// CompletePrevious 选择上面第 count 个补全
func (l *Line) CompletePrevious(count int) {
	if !l.mode.Is(linemode.Complete) {
		l.StartComplete(false)
	}

	if l.completeState != nil {
		var index int
		if l.completeState.completeIndex == -1 {
			index = len(l.completeState.currentCompletions) - 1
		} else {
			index = l.completeState.completeIndex - count
			if index < 0 {
				index += len(l.completeState.currentCompletions)
			}
		}
		l.gotoCompletion(index)
	}
}

// StartComplete 开始补全
func (l *Line) StartComplete(gotoFirst bool) {
	currentCompletions := l.CreateCode().GetCompletions()

	if len(currentCompletions) > 0 {
		l.completeState = newCompletionState(l.Document(), currentCompletions)
		l.mode = linemode.Complete
		if gotoFirst {
			l.gotoCompletion(0)
		}
	} else {
		l.mode = linemode.Normal
		l.completeState = nil
	}
}

// AcceptComplete 接受当前选中的补全
func (l *Line) AcceptComplete() {
	l.mode = linemode.Normal
	l.completeState = nil
}

// CancelComplete 取消补全
func (l *Line) CancelComplete() {
	if l.mode.Is(linemode.Complete) {
		l.gotoCompletion(-1)
		l.mode = linemode.Normal
		l.completeState = nil
	}
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
	l.insertText([]rune(l.completeState.currentCompletionText()), true)

	l.mode = linemode.Complete
}

// GetRenderContext 返回渲染上下文信息
func (l *Line) GetRenderContext() *RenderContext {
	code := l.CreateCode()
	var completeState *cCompletionState
	if l.mode.Is(linemode.Complete) {
		completeState = l.completeState
	} else {
		completeState = nil
	}

	var highlights []section
	document := l.Document()
	matchIndex := l.getMatchingBracket()
	if matchIndex >= 0 {
		r, c := document.translateIndexToRowCol(matchIndex)
		start := Location{r, c}
		end := Location{document.CursorPositionRow(), document.CursorPositionCol()}
		//    交换开始和结束位置
		if matchIndex > l.cursorPosition {
			start, end = end, start
		}
		end = Location{end.Row, end.Col + 1}
		highlights = append(highlights, section{start, end})
	}

	renderCtx := newRenderContext(
		code,
		completeState,
		document,
		highlights,
		l.cancelSelection,
	)
	l.cancelSelection = false
	return renderCtx
}

// HistoryForward 选择下一个历史输入
func (l *Line) HistoryForward() {
	if l.workingIndex < len(l.workingLines)-1 {
		l.workingIndex++
		l.buffer = []rune(l.workingLines[l.workingIndex])
		l.SetCursorPosition(len(l.buffer))
	}
}

// HistoryBackward 选择上一个历史输入
func (l *Line) HistoryBackward() {
	if l.workingIndex > 0 {
		l.workingIndex--
		l.buffer = []rune(l.workingLines[l.workingIndex])
		l.SetCursorPosition(len(l.buffer))
	}
}

func (l *Line) Newline() {
	spaces := l.Document().LeadingWhitespaceInCurrentLine()
	l.insertText([]rune{'\n'}, true)
	if l.autoIndent {
		l.insertText([]rune(spaces), true)
	}
}

// AutoEnter 自动处理 Enter
func (l *Line) AutoEnter() {
	if l.mode.Is(linemode.Complete) {
		l.AcceptComplete()
		return
	}
	if l.IsMultiline() {
		l.Newline()
	} else {
		l.AcceptInput()
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
	l.insertText([]rune(insert), true)
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
	l.insertText([]rune(insert), true)
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
	//    输入中有选中文本，直接删除
	if l.selection.start != -1 && l.selection.end > l.selection.start {
		l.DeleteCharacterBeforeCursor(l.selection.end - l.selection.start)
		l.selection = _LineArea{-1, -1}
		l.cancelSelection = true
	}
	l.insertText(data, moveCursor)

	if l.CompleteAfterInsertText() {
		l.StartComplete(false)
	}
}

func (l *Line) insertText(data []rune, moveCursor bool) {
	result := insertRunes(l.buffer, l.cursorPosition, data)
	// result := concatRunes(l.buffer[:l.cursorPosition], data, l.buffer[l.cursorPosition:])
	l.setText(result)
	if moveCursor {
		l.SetCursorPosition(l.cursorPosition + len(data))
	}
}

// CompleteAfterInsertText 返回 true 表示每次插入文本我们都获取一次补全
func (l *Line) CompleteAfterInsertText() bool {
	return l.CreateCode().CompleteAfterInsertText()
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
	if len(l.undoStack) > 0 {
		top := l.undoStack[len(l.undoStack)-1]
		l.undoStack = l.undoStack[:len(l.undoStack)-1]
		l.setText([]rune(top.text))
		l.SetCursorPosition(top.cursorPosition)
	}
}

// AcceptInput 确定用户输入（一般是用户按下 Enter）
func (l *Line) AcceptInput() {
	text := l.text()

	// 文本与最后一个不相同时，保存到历史中
	if l.history.Length() == 0 || l.history.GetAt(-1) != text {
		if len(text) > 0 {
			l.history.Append(text)
		}
	}

	l.accept = true
}

// IsAccept 返回 true 表示返回用户输入，否则不返回
func (l *Line) IsAccept() bool {
	return l.accept
}

func (l *Line) HasText() bool {
	return len(l.buffer) > 0
}

func (l *Line) IsMultiline() bool {
	res := l.CreateCode().ContinueInput()
	return res
}

func (l *Line) ToNormalMode() {
	l.ToMode(linemode.Normal)
}

func (l *Line) ToMode(modes ...linemode.LineMode) {
	if l.mode.In(modes...) {
		return
	}
	// todo something
	if l.mode.Is(linemode.IncrementalSearch) {
	} else if l.mode.Is(linemode.Complete) {
		l.AcceptComplete()
	}
}

func (l *Line) MouseDown(info *MouseInfoOfInput) {
	location := info.location
	if location.Row == -1 || location.Col == -1 {
		l.selection = _LineArea{-1, -1}
	} else {
		pos := l.Document().translateRowColToIndex(location.Row, location.Col)
		l.SetCursorPosition(pos)
		l.selection = _LineArea{start: pos, end: pos}
	}
	if info.completeIndex != -1 {
		l.gotoCompletion(info.completeIndex)
	}
}

func (l *Line) MouseMove(location Location) {
	if location.Row == -1 || location.Col == -1 {
		return
	}
	pos := l.Document().translateRowColToIndex(location.Row, location.Col)
	l.SetCursorPosition(pos)
	if l.selection.start != -1 {
		l.selection.end = pos
	}
}

func (l *Line) MouseUp(location Location) {
	if location.Row == -1 || location.Col == -1 {
		return
	}
	pos := l.Document().translateRowColToIndex(location.Row, location.Col)
	l.SetCursorPosition(pos)
	if l.selection.start != -1 {
		l.selection.end = pos
	}
}

func (l *Line) Dblclick(location Location) {
	document := l.Document()
	pos := document.translateRowColToIndex(location.Row, location.Col)

	end := len(l.buffer)
	for i := pos + 1; i < len(l.buffer); i++ {
		if unicode.IsSpace(l.buffer[i]) {
			end = i - 1
			break
		}
	}

	start := 0
	for i := pos - 1; i >= 0; i-- {
		if unicode.IsSpace(l.buffer[i]) {
			start = i + 1
			break
		}
	}
	if start >= end {
		return
	}

	l.selection = _LineArea{start: start, end: end}
	l.SetCursorPosition(end)
}
