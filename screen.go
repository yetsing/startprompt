package startprompt

import (
	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/terminalcode"
	"github.com/yetsing/startprompt/terminalcolor"
	"github.com/yetsing/startprompt/token"
	"strings"
)

type Grid struct {
	char  rune
	style terminalcolor.Style
}

func (g *Grid) output() string {
	if g.style != nil {
		return terminalcolor.ApplyStyle(g.style, string(g.char), true)
	} else {
		return string(g.char)
	}
}

func (g *Grid) width() int {
	n := runewidth.StringWidth(string(g.char))
	if n < 0 {
		n = 0
	}
	return n
}

func newGrid(char rune, style terminalcolor.Style) *Grid {
	return &Grid{
		char:  char,
		style: style,
	}
}

type CursorPosition struct {
	X int
	Y int
}

func NewScreen(schema Schema) *Screen {
	return &Screen{
		schema: schema,
		buffer: map[int]map[int]*Grid{},
		x:      0,
		y:      0,
	}
}

type Screen struct {
	schema Schema
	buffer map[int]map[int]*Grid
	x      int
	y      int
}

func (s *Screen) Output() (string, CursorPosition) {
	// 如果我们不手动移动光标位置，那么就需要一个个字符地输出，这样光标才会自动向右（向下）移动
	// 那么在 buffer 里面的坐标之间的空档，我们都要输出空白字符用来填充
	var result []string
	var cursorPos CursorPosition
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
				var grid *Grid
				if _, found := lineData[c]; found {
					grid = lineData[c]
				} else {
					// 32 是空格
					grid = newGrid(32, nil)
				}
				result = append(result, grid.output())
				c += grid.width()
			}
			cursorPos.X = c
		}

		// 除了最后一行的都加上换行符
		if i != rows-1 {
			result = append(result, terminalcode.CRLF)
		}
	}
	return strings.Join(result, ""), cursorPos
}

func (s *Screen) WriteTokens(tokens []token.Token) {
	for _, t := range tokens {
		if t.TypeIs(token.EOF) {
			break
		}
		s.WriteText(t.Literal, s.schema[t.Type])
	}
}

func (s *Screen) WriteText(text string, style terminalcolor.Style) {
	for _, r := range text {
		s.WriteChar(r, style)
	}
}

func (s *Screen) WriteChar(char rune, style terminalcolor.Style) {
	grid := s.writeAtPos(s.x, s.y, char, style)
	s.x += grid.width()
}

func (s *Screen) writeAtPos(x int, y int, char rune, style terminalcolor.Style) *Grid {
	grid := newGrid(char, style)
	if _, found := s.buffer[y]; !found {
		s.buffer[y] = map[int]*Grid{}
	}
	s.buffer[y][x] = grid
	return grid
}
