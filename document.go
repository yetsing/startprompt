package startprompt

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type _DocumentCache struct {
	lines       []string
	lineLengths []int
	lineIndexes []int
}

type Document struct {
	text string
	// 光标在文本中的索引
	cursorPosition int
	runes          []rune

	// 行数组和索引缓存，只有在需要的时候才会计算并缓存
	cache *_DocumentCache
}

func NewDocument(text string, cursorPosition int) *Document {
	return &Document{
		text:           text,
		cursorPosition: cursorPosition,
		runes:          []rune(text),
		cache:          &_DocumentCache{},
	}
}

func (d *Document) String() string {
	return fmt.Sprintf("Document(%q, %d)", d.text, d.cursorPosition)
}

func (d *Document) Equal(other *Document) bool {
	return d.text == other.text && d.cursorPosition == other.cursorPosition
}

func (d *Document) Text() string {
	return d.text
}

func (d *Document) CursorPosition() int {
	return d.cursorPosition
}

func (d *Document) CurrentChar() string {
	return d.getCharRelativeToCursor(0)
}

func (d *Document) CharBeforeCursor() string {
	return d.getCharRelativeToCursor(-1)
}

func (d *Document) TextBeforeCursor() string {
	return string(d.runes[:d.cursorPosition])
}

func (d *Document) TextAfterCursor() string {
	return string(d.runes[d.cursorPosition:])
}

// CurrentLineBeforeCursor 返回行首到光标处的文本
func (d *Document) CurrentLineBeforeCursor() string {
	text := d.TextBeforeCursor()
	// 返回最后一行
	index := strings.LastIndexByte(text, '\n')
	if index == -1 {
		return text
	}
	return text[index+1:]
}

// CurrentLineAfterCursor 返回光标到行尾的文本（不包括换行符）
func (d *Document) CurrentLineAfterCursor() string {
	text := d.TextAfterCursor()
	// 返回第一行
	index := strings.IndexByte(text, '\n')
	if index == -1 {
		return text
	}
	return text[:index]
}

// 返回当前行到最后一行
func (d *Document) linesFromCurrent() []string {
	return d.lines()[d.CursorPositionRow():]
}

func (d *Document) LineCount() int {
	return len(d.lines())
}

// CurrentLine 返回光标所在行文本（不包括换行符）
func (d *Document) CurrentLine() string {
	return d.CurrentLineBeforeCursor() + d.CurrentLineAfterCursor()
}

// LeadingWhitespaceInCurrentLine 返回当前行开始处的空白字符
func (d *Document) LeadingWhitespaceInCurrentLine() string {
	currentLine := d.CurrentLine()
	var i int
	var r rune
	for i, r = range currentLine {
		if !unicode.IsSpace(r) {
			break
		}
	}
	return currentLine[:i]
}

// 获取相对光标位置的字符
func (d *Document) getCharRelativeToCursor(offset int) string {
	index := d.cursorPosition + offset
	if index < 0 || index >= len(d.runes) {
		return ""
	}
	return string(d.runes[index])
}

// OnFirstLine 光标是否在第一行
func (d *Document) OnFirstLine() bool {
	return d.CursorPositionRow() == 0
}

// OnLastLine 光标是否在最后一行
func (d *Document) OnLastLine() bool {
	return d.CursorPositionRow() == d.LineCount()-1
}

// CursorPositionRow 返回光标所在行号（从 0 开始计数）
func (d *Document) CursorPositionRow() int {
	row, _ := d.findLineStartIndex(d.cursorPosition)
	return row
}

// CursorPositionCol 返回光标所在的列号（从 0 开始计数）
func (d *Document) CursorPositionCol() int {
	_, lineStartIndex := d.findLineStartIndex(d.cursorPosition)
	return d.cursorPosition - lineStartIndex
}

// index 表示文本中的索引，返回所在行号和行首索引
func (d *Document) findLineStartIndex(index int) (int, int) {
	indexes := d.lineStartIndexes()
	lineno := sort.Search(len(indexes), func(i int) bool {
		return indexes[i] > index
	})
	// 得到 lineno 实际上是在 index 所在的下一行，减一便能得到 index 所在行号
	lineno--
	return lineno, indexes[lineno]
}

// 将文本索引 index 转化成行号和列号
func (d *Document) translateIndexToRowCol(index int) (int, int) {
	row, lineStartIndex := d.findLineStartIndex(index)
	col := index - lineStartIndex
	return row, col
}

// 将行号和列号转成文本的索引
func (d *Document) translateRowColToIndex(row int, col int) int {
	lineCount := d.LineCount()
	if row < 0 {
		row = 0
	} else if row >= lineCount {
		row = lineCount - 1
	}
	result := d.lineStartIndexes()[row]
	lineLength := d.lineLengths()[row]

	if col < 0 {
		col = 0
	} else if col > lineLength {
		col = lineLength
	}
	result += col
	if result > len(d.runes) {
		result = len(d.runes)
	}
	return result
}

// CursorUpPosition 返回光标往上移动一行后位置。如果光标位于第一行，返回 -1 。
func (d *Document) CursorUpPosition() int {
	text := d.TextBeforeCursor()
	if strings.ContainsRune(text, '\n') {
		lines := strings.Split(text, "\n")
		length := len(lines)
		// 光标所在行
		currentLine := lines[length-1]
		// 光标的上一行
		previousLine := lines[length-2]

		currentLineLength := utf8.RuneCountInString(currentLine)
		previousLineLength := utf8.RuneCountInString(previousLine)

		// 如果光标前文本比上一行长，光标会移动到上一行行尾
		if currentLineLength > previousLineLength {
			return d.cursorPosition - currentLineLength - 1
		} else {
			// 否则找到上一行对应的位置
			return d.cursorPosition - previousLineLength - 1
		}
	}
	return -1
}

// CursorDownPosition 返回光标往下移动一行后位置。如果光标位于最后一行，返回 -1 。
func (d *Document) CursorDownPosition() int {
	text := d.TextAfterCursor()
	if strings.ContainsRune(text, '\n') {
		pos := utf8.RuneCountInString(d.CurrentLineBeforeCursor())
		lines := strings.Split(text, "\n")
		// 光标所在行
		currentLine := lines[0]
		// 光标的下一行
		nextLine := lines[1]

		currentLineLength := utf8.RuneCountInString(currentLine)
		nextLineLength := utf8.RuneCountInString(nextLine)

		// 如果光标前文本比下一行长，光标会移动到下一行行尾
		if pos > nextLineLength {
			return d.cursorPosition + currentLineLength + nextLineLength + 1
		} else {
			// 否则找到下一行对应的位置
			return d.cursorPosition + currentLineLength + pos + 1
		}
	}
	return -1
}

// 光标是否在文本的最后面（最后一行的行尾）
func (d *Document) isCursorAtTheEnd() bool {
	return d.cursorPosition == len(d.runes)
}

// 光标是否在行尾
func (d *Document) isCursorAtTheEndOfLine() bool {
	return d.CurrentChar() == "\n" || d.CurrentChar() == ""
}

// 当光标位于字符串 sub 开头时返回 true
func (d *Document) hasMatchAtCurrentPosition(sub string) bool {
	return strings.HasPrefix(d.TextAfterCursor(), sub)
}

// 找到光标前第一个单词开头的位置记为 S ，返回 S 与光标的相对位置
// 找不到返回 0
func (d *Document) findStartOfPreviousWord() int {
	text := d.TextBeforeCursor()
	if len(text) == 0 {
		return 0
	}
	text = reverseString(text)

	step := 0
	inWord := false
	for _, r := range text {
		if isWordDelimiter(r) {
			if inWord {
				break
			}
			// 忽略开头的空格
		} else {
			inWord = true
		}
		step++
	}
	// 文本是反过来搜索的，所以要返回负数
	return -step
}

// 找到光标后第一个单词开头的位置记为 S ，返回 S 与光标的相对位置
// 找不到返回 0
func (d *Document) findNextWordBeginning() int {
	text := d.TextAfterCursor()
	if len(text) == 0 {
		return 0
	}
	// 下面这段代码有一点需要注意的是，需要忽略光标所在处的单词
	step := 0
	inSpace := false
	for _, r := range text {
		if isWordDelimiter(r) {
			inSpace = true
		} else {
			if inSpace {
				// 从空格到单词，说明是单词的开始
				return step
			}
		}
		step++
	}
	// 没有找到下一个单词
	return 0
}

// 找到光标后第一个单词结尾的位置记为 S ，返回 S 与光标的相对位置
// 找不到返回 0
// includeCurrentPosition 是否包括光标处字符，之所以有这个选项，说明如下
// 对于 vim 来说，按 e 可以移动到单词末尾，实际上光标是在单词最后一个字符的左边
// 这个时候如果再按 e 跳到下一个单词末尾，那么在判断的时候需要忽略这个字符
func (d *Document) findNextWordEnding(includeCurrentPosition bool) int {
	text := d.TextAfterCursor()
	if !includeCurrentPosition {
		text = stringStartAt(text, 1)
	}
	step := 0
	inWord := false
	for _, r := range text {
		if isWordDelimiter(r) {
			if inWord {
				break
			}
			// 忽略开头的空格
		} else {
			inWord = true
		}
		step++
	}
	if includeCurrentPosition {
		return step
	} else {
		// 因为跳过了光标处的字符，所以实际位置需要加 1
		return step + 1
	}
}

// ================
// 下面的都是缓存相关的方法
// ================
func (d *Document) lines() []string {
	if d.cache.lines == nil {
		d.cache.lines = strings.Split(d.text, "\n")
	}
	return d.cache.lines
}

func (d *Document) lineLengths() []int {
	if d.cache.lineLengths == nil {
		lines := d.lines()
		lineLengths := make([]int, len(lines))
		for i, line := range lines {
			lineLengths[i] = utf8.RuneCountInString(line)
		}
		d.cache.lineLengths = lineLengths
	}
	return d.cache.lineLengths
}

// 返回每行行首在文本中的索引
func (d *Document) lineStartIndexes() []int {
	if d.cache.lineIndexes == nil {
		lineLengths := d.lineLengths()
		count := len(lineLengths)
		indexes := make([]int, len(lineLengths))
		indexes[0] = 0
		pos := 0
		for i, lineLength := range lineLengths {
			// 最后一个不需要
			if i == count-1 {
				break
			}
			// +1 是因为要算上换行符
			pos += lineLength + 1
			indexes[i+1] = pos
		}
		d.cache.lineIndexes = indexes
	}
	return d.cache.lineIndexes
}
