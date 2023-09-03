package startprompt

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/terminalcode"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
)

var displayMapping = map[string]string{
	"\x00": "^@", // Control space
	"\x01": "^A",
	"\x02": "^B",
	"\x03": "^C",
	"\x04": "^D",
	"\x05": "^E",
	"\x06": "^F",
	"\x07": "^G",
	"\x08": "^H",
	"\x09": "^I",
	"\x0a": "^J",
	"\x0b": "^K",
	"\x0c": "^L",
	"\x0d": "^M",
	"\x0e": "^N",
	"\x0f": "^O",
	"\x10": "^P",
	"\x11": "^Q",
	"\x12": "^R",
	"\x13": "^S",
	"\x14": "^T",
	"\x15": "^U",
	"\x16": "^V",
	"\x17": "^W",
	"\x18": "^X",
	"\x19": "^Y",
	"\x1a": "^Z",
	"\x1b": "^[", // Escape
	"\x1c": "^\\",
	"\x1d": "^]",
	"\x1f": "^_",
	"\x7f": "^?", // Control backspace
}

type Char struct {
	char   string
	style  terminalcolor.Style
	cwidth int
}

func (c *Char) output() string {
	if c.style != nil {
		return terminalcolor.ApplyStyle(c.style, c.char, true)
	} else {
		return c.char
	}
}

func (c *Char) width() int {
	if c.cwidth == -1 {
		n := runewidth.StringWidth(c.char)
		if n < 0 {
			n = 0
		}
		c.cwidth = n
	}
	return c.cwidth
}

func newChar(r rune, style terminalcolor.Style) *Char {
	ch := string(r)
	if _, found := displayMapping[ch]; found {
		ch = displayMapping[ch]
	}
	return &Char{
		char:   ch,
		style:  style,
		cwidth: -1,
	}
}

type Coordinate struct {
	X int
	Y int
}

func NewScreen(schema Schema, size _Size) *Screen {
	return &Screen{
		schema:    schema,
		buffer:    map[int]map[int]*Char{},
		size:      size,
		x:         0,
		y:         0,
		inputRow:  0,
		inputCol:  0,
		cursorMap: map[Coordinate]Coordinate{},
	}
}

// Screen 以坐标维度缓冲输出字符
type Screen struct {
	schema Schema
	//    {y: {x: Char}}
	buffer map[int]map[int]*Char
	//    窗口宽度和高度
	size _Size
	//    文本中光标坐标（是一个相对于文本左上角的坐标，而不是窗口左上角）
	x int
	y int
	//    文本最后一行行尾的光标坐标（是一个相对于文本左上角的坐标，而不是窗口左上角）
	maxCursorCoordinate Coordinate
	// 文本中光标的行列
	inputRow int
	inputCol int
	// 保存光标行列到 yx 的映射
	cursorMap map[Coordinate]Coordinate

	secondLinePrefixFunc func() []token.Token
}

func (s *Screen) Width() int {
	return s.size.width
}

func (s *Screen) CurrentHeight() int {
	if len(s.buffer) == 0 {
		return 1
	} else {
		my := 0
		for y := range s.buffer {
			if y > my {
				my = y
			}
		}
		return my
	}
}

func (s *Screen) GetBuffer() map[int]map[int]*Char {
	return s.buffer
}

func (s *Screen) Output() (string, Coordinate) {
	// 如果我们不手动移动光标位置，那么就需要一个个字符地输出，这样光标才会自动向右（向下）移动
	// 那么在 buffer 里面的坐标之间的空档，我们都要输出空白字符用来填充
	var result []string
	var cursorPos Coordinate
	// 统计一下有多少行，其实就是等于最大的 y + 1
	rows := 1
	for y := range s.buffer {
		if y > rows-1 {
			rows = y + 1
		}
	}
	cursorPos.Y = rows - 1
	for i := 0; i < rows; i++ {
		lineData, found := s.buffer[i]
		if found {
			// 统计一下有多少列，其实就是等于最大的 x + 1
			cols := 1
			for x := range lineData {
				if x > cols-1 {
					cols = x + 1
				}
			}

			c := 0
			for c < cols {
				var char *Char
				if _, found := lineData[c]; found {
					char = lineData[c]
				} else {
					char = newChar(' ', nil)
				}
				result = append(result, char.output())
				c += char.width()
			}
			cursorPos.X = c
		}

		// 除了最后一行的都加上换行符
		if i != rows-1 {
			result = append(result, terminalcode.CRLF)
		}
	}
	return strings.Join(result, ""), s.maxCursorCoordinate
}

// WriteTokensAtPos 在指定位置写入 token 数组
func (s *Screen) WriteTokensAtPos(x int, y int, tokens []token.Token) {
	for _, t := range tokens {
		if t.TypeIs(token.EOF) {
			break
		}
		style := s.schema[t.Type]
		for _, r := range t.Literal {
			char := newChar(r, style)
			s.writeAtPos(x, y, char)
			x += char.width()
		}
	}
}

// WriteTokens 写入 Token 数组， saveInputPos: 是否保存输入位置
func (s *Screen) WriteTokens(tokens []token.Token, saveInputPos bool) {
	for _, t := range tokens {
		if t.TypeIs(token.EOF) {
			break
		}
		style := s.schema.StyleForToken(t.Type)
		for _, r := range t.Literal {
			s.WriteRune(r, style, saveInputPos)
		}
	}
}

func (s *Screen) WriteRune(r rune, style terminalcolor.Style, saveInputPos bool) {
	char := newChar(r, style)
	charWidth := char.width()

	//    如果宽度不够放下这个字符，另起一行
	//    如果这里用 > 的话，输入一行的最后一个字符时，光标会在字符上面，而不是正常的在字符后面
	if s.x+charWidth >= s.size.width {
		s.y++
		s.x = 0
	}

	//    记录输入位置坐标
	if saveInputPos {
		s.saveInputPos()
	}

	//    插入换行符
	if r == '\n' {
		s.y++
		s.x = 0
		if s.y > s.maxCursorCoordinate.Y {
			s.maxCursorCoordinate.Y = s.y
			s.maxCursorCoordinate.X = 0
		}

		if saveInputPos {
			s.inputRow++
			s.inputCol = 0

			if s.secondLinePrefixFunc != nil {
				s.WriteTokens(s.secondLinePrefixFunc(), false)
			}
		}
	} else {
		s.writeAtPos(s.x, s.y, char)
		if saveInputPos {
			s.inputCol++
		}
		s.x += charWidth
	}
}

func (s *Screen) writeAtPos(x int, y int, char *Char) {
	// 超出屏幕的不进行写入
	if x >= s.size.width {
		return
	}
	if _, found := s.buffer[y]; !found {
		s.buffer[y] = map[int]*Char{}
	}
	s.buffer[y][x] = char
	if y > s.maxCursorCoordinate.Y {
		s.maxCursorCoordinate.Y = y
		s.maxCursorCoordinate.X = x + char.width()
	} else if y == s.maxCursorCoordinate.Y && x+char.width() > s.maxCursorCoordinate.X {
		s.maxCursorCoordinate.X = x + char.width()
	}
}

// saveInputPos 保存行列到 xy 坐标的映射
func (s *Screen) saveInputPos() {
	s.cursorMap[Coordinate{s.inputCol, s.inputRow}] = Coordinate{s.x, s.y}
}

func (s *Screen) setSecondLinePrefix(secondLinePrefixFunc func() []token.Token) {
	s.secondLinePrefixFunc = secondLinePrefixFunc
}

// getCursorCoordinate 根据行列坐标得到 xy 坐标
func (s *Screen) getCursorCoordinate(row int, col int) Coordinate {
	return s.cursorMap[Coordinate{col, row}]
}

// 返回 cursorMap 的字符串，用于调试
func (s *Screen) cursorMapS() string {
	return fmt.Sprintf("%v", s.cursorMap)
}
