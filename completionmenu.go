package startprompt

import (
	"github.com/mattn/go-runewidth"
	"github.com/yetsing/startprompt/token"
)

type cCompletionMenuInfo struct {
	area      area
	sliceFrom int
	sliceTo   int
}

// getCompleteIndex 返回坐标位置在第几个补全项上
func (c *cCompletionMenuInfo) getCompleteIndex(coordinate Coordinate) int {
	if !c.area.RectContains(coordinate) {
		return -1
	}
	startY := c.area.getStart().Y
	return c.sliceFrom + (coordinate.Y - startY)
}

// cCompletionMenu 辅助补全菜单的渲染
type cCompletionMenu struct {
	screen        *Screen
	completeState *cCompletionState
	info          *cCompletionMenuInfo
	maxHeight     int

	progressButtonToken token.Token
	progressBarToken    token.Token
}

func newCompletionMenu(screen *Screen, completeState *cCompletionState, maxHeight int) *cCompletionMenu {
	return &cCompletionMenu{
		screen:        screen,
		completeState: completeState,
		info:          &cCompletionMenuInfo{},
		maxHeight:     maxHeight,

		progressButtonToken: token.NewToken(token.CompletionMenuProgressButton, " "),
		progressBarToken:    token.NewToken(token.CompletionMenuProgressBar, " "),
	}
}

// 返回光标的位置坐标
func (c *cCompletionMenu) getOrigin() Coordinate {
	return c.screen.getCoordinate(
		c.completeState.originalDocument.CursorPositionRow(),
		c.completeState.originalDocument.CursorPositionCol())
}

// getDrawCoordinate 返回菜单渲染位置坐标（因为是从左上角开始，所以这个就是左上角的坐标）
// itemWidth 补全项的宽度
func (c *cCompletionMenu) getDrawCoordinate(itemWidth int) Coordinate {
	coordinate := c.getOrigin()
	x := coordinate.X
	y := coordinate.Y
	y++
	//    这里 x - 1 是因为前面会加个空格
	x = maxInt(0, x-1)
	if x+itemWidth > c.screen.Width() {
		x -= (x + itemWidth) - c.screen.Width() + 1
	}
	return Coordinate{
		X: x,
		Y: y,
	}
}

func (c *cCompletionMenu) showMeta() bool {
	for _, completion := range c.completeState.currentCompletions {
		if len(completion.DisplayMeta) > 0 {
			return true
		}
	}
	return false
}

// getMenuWidth 返回补全展示文本的宽度
func (c *cCompletionMenu) getMenuWidth() int {
	maxDisplay := c.screen.Width() / 2
	menuWidth := 0
	for _, completion := range c.completeState.currentCompletions {
		w := runewidth.StringWidth(completion.Display)
		if w > menuWidth {
			menuWidth = w
		}
	}
	return minInt(maxDisplay, menuWidth)
}

// getMenuMetaWidth 返回补全元信息的宽度
func (c *cCompletionMenu) getMenuMetaWidth() int {
	maxDisplayMeta := c.screen.Width() / 2
	menuMetaWidth := 0
	for _, completion := range c.completeState.currentCompletions {
		if len(completion.DisplayMeta) == 0 {
			continue
		}
		w := runewidth.StringWidth(completion.DisplayMeta)
		if w > menuMetaWidth {
			menuMetaWidth = w
		}
	}
	return minInt(maxDisplayMeta, menuMetaWidth)
}

// 将菜单写入 screen 里面
func (c *cCompletionMenu) write() {
	completions := c.completeState.currentCompletions
	index := c.completeState.completeIndex

	//    决定从哪个补全项开始展示
	sliceFrom := 0
	//    补全项多于最大高度并且当前选择项在下半部分位置，需要向上移动补全菜单
	//    尽可能地让选中的补全项位于菜单中上部分
	if len(completions) > c.maxHeight && index != -1 && index > c.maxHeight/2 {
		sliceFrom = minInt(
			index-c.maxHeight/2,          // 将选择项移到中间位置
			len(completions)-c.maxHeight, // 最后一个补全在最底部
		)
	}

	sliceTo := minInt(sliceFrom+c.maxHeight, len(completions))

	//    计算补全菜单的宽度
	menuWidth := c.getMenuWidth()
	menuMetaWidth := c.getMenuMetaWidth()
	//    获取菜单的位置坐标
	//    补全项前后总共有 5 个空格
	coordinate := c.getDrawCoordinate(menuWidth + menuMetaWidth + 5)
	showMeta := c.showMeta()
	//    写入补全到 screen
	for i, completion := range completions[sliceFrom:sliceTo] {
		//    i+sliceFrom == index 判断补全项是否已选中
		tks := []token.Token{
			token.NewToken(token.Unspecific, " "),
			c.getMenuItemToken(completion, i+sliceFrom == index, menuWidth),
		}
		if showMeta {
			tks = append(
				tks,
				c.getMenuItemMetaToken(completion, i+sliceFrom == index, menuMetaWidth),
			)
		}
		if i+sliceFrom == index {
			tks = append(tks, c.progressButtonToken)
		} else {
			tks = append(tks, c.progressBarToken)
		}
		tks = append(tks, token.NewToken(token.Unspecific, " "))
		c.screen.WriteTokensAtPos(coordinate.X, coordinate.Y+i, tks)
	}
	lastCoordinate := c.screen.getLastCoordinate()
	//    lastCoordinate 是最后一个字符的坐标，同样要包含在 area 里面
	//    而 area 是左闭右开区间，所以需要加 1
	lastCoordinate.add(1, 1)
	c.info.area = area{coordinate, lastCoordinate}
	c.info.sliceFrom = sliceFrom
	c.info.sliceTo = sliceTo
}

func (c *cCompletionMenu) getMenuItemToken(completion *Completion, isCurrentCompletion bool, width int) token.Token {
	var ttype token.TokenType
	if isCurrentCompletion {
		ttype = token.CompletionMenuCompletionCurrent
	} else {
		ttype = token.CompletionMenuCompletion
	}
	return token.NewToken(ttype, " "+ljustWidth(completion.Display, width))
}

func (c *cCompletionMenu) getMenuItemMetaToken(completion *Completion, isCurrentCompletion bool, width int) token.Token {
	var ttype token.TokenType
	if isCurrentCompletion {
		ttype = token.CompletionMenuMetaCurrent
	} else {
		ttype = token.CompletionMenuMeta
	}
	return token.NewToken(ttype, " "+ljustWidth(completion.DisplayMeta, width))
}

func (c *cCompletionMenu) getInfo() *cCompletionMenuInfo {
	return c.info
}
