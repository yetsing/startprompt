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

func (d *Document) currentChar() string {
	return d.getCharRelativeToCursor(0)
}

func (d *Document) charBeforeCursor() string {
	return d.getCharRelativeToCursor(-1)
}

func (d *Document) textBeforeCursor() string {
	return string(d.runes[:d.cursorPosition])
}

func (d *Document) textAfterCursor() string {
	return string(d.runes[d.cursorPosition:])
}

// 返回行首到光标处的文本
func (d *Document) currentLineBeforeCursor() string {
	text := d.textBeforeCursor()
	// 返回最后一行
	index := strings.LastIndexByte(text, '\n')
	if index == -1 {
		return text
	}
	return text[index+1:]
}

// 返回光标到行尾的文本（不包括换行符）
func (d *Document) currentLineAfterCursor() string {
	text := d.textAfterCursor()
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

func (d *Document) lineCount() int {
	return len(d.lines())
}

// 返回光标所在行文本（不包括换行符）
func (d *Document) currentLine() string {
	return d.currentLineBeforeCursor() + d.currentLineAfterCursor()
}

// 返回当前行开始处的空白字符
func (d *Document) leadingWhitespaceInCurrentLine() string {
	currentLine := d.currentLine()
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

// 光标是否在第一行
func (d *Document) onFirstLine() bool {
	return d.CursorPositionRow() == 0
}

// 光标是否在最后一行
func (d *Document) onLastLine() bool {
	return d.CursorPositionRow() == d.lineCount()-1
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
	lineCount := d.lineCount()
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

// 光标是否在文本的最后面（最后一行的行尾）
func (d *Document) isCursorAtTheEnd() bool {
	return d.cursorPosition == len(d.runes)
}

// 光标是否在行尾
func (d *Document) isCursorAtTheEndOfLine() bool {
	return d.currentChar() == "\n" || d.currentChar() == ""
}

// 当光标位于字符串 sub 开头时返回 true
func (d *Document) hasMatchAtCurrentPosition(sub string) bool {
	return strings.HasPrefix(d.textAfterCursor(), sub)
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
