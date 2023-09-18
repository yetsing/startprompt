package startprompt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// 解析 VT100 输入流数据
// 参考：https://vt100.net/docs/vt100-ug/chapter3.html

func NewInputStream(handler EventHandler, cli *CommandLine) *InputStream {
	return &InputStream{handler: handler, cli: cli}
}

type InputStream struct {
	handler  EventHandler
	cli      *CommandLine
	previous string
	// 是否触发事件
	isEmitEvent bool
}

func (is *InputStream) Reset() {
	is.isEmitEvent = false
}

// FeedTimeout 超时通知，主要用来快速触发 Esc 事件
// 返回值表示是否有事件触发
//
//	因为 ANSI 转义序列都是 Esc 开头
//	导致无法区分 Esc 和其他的快捷键，只能等待后续字符，再做判断
//	因此按下 Esc 后不会有事件触发，现在通过超时来快速识别 Esc 键
func (is *InputStream) FeedTimeout() bool {
	var offset int
	is.isEmitEvent = false
	for i, r := range is.previous {
		// 触发 Esc 事件
		if r == '\x1b' {
			is.callHandler(EventTypeEscape, '\x1b')
		} else {
			offset = i
			break
		}
	}
	is.previous = is.previous[offset:]
	return is.isEmitEvent
}

func (is *InputStream) FeedData(data string) {
	for _, r := range data {
		is.Feed(r)
	}
}

func (is *InputStream) FeedRunes(runes []rune) {
	for _, r := range runes {
		is.Feed(r)
	}
}

// Feed 根据输入触发对应的事件
func (is *InputStream) Feed(r rune) {
	var buffer []rune
	is.isEmitEvent = false
	for {
		key := string(r)
		if len(is.previous) > 0 {
			key = is.previous + key
		}
		// 检查是不是快捷键操作
		action, found := keyActions[key]
		if found {
			is.callHandler(action, []rune(key)...)
			is.previous = ""
			break
		}
		// 检查是不是多字符快捷键操作
		// 因为多字符需要输入多次，所以查看有没有哪个 key 的前缀可以匹配上
		if prefixMatchKeyActions(key) {
			is.previous = key
			break
		}

		// 之前保存了 key 可能匹配快捷键，但是现在发现没有匹配的快捷键
		// 比如有个 abc 的快捷键，用户先输入 ab ，匹配到 abc 前缀，于是暂时保存 ab
		// 后面用户输入 d ，此时用户输入的 key 就变成 abd ，没有匹配，需要做其他处理
		if len(is.previous) > 0 {
			first := runeAt(is.previous, 0)
			// 按下 Esc 键就会收到 '\x1b' ，所以这里需要判断一下特殊处理
			if first == '\x1b' {
				is.callHandler(EventTypeEscape, '\x1b')
			} else {
				// 如果不是快捷键操作，那么就是正常的输入
				is.callHandler(EventTypeInsertChar, first)
			}
			// 剩余的字符放到缓冲中，留待下次循环的时候处理
			buffer = []rune(is.previous[utf8.RuneLen(first):])
			buffer = append(buffer, r)
			is.previous = ""
		} else {
			is.callHandler(EventTypeInsertChar, r)
		}
		// 如果之前有缓存字符，继续进行处理
		if len(buffer) > 0 {
			r = buffer[0]
			buffer = buffer[1:]
		} else {
			break
		}
	}
}

func (is *InputStream) callHandler(eventType EventType, a ...rune) {
	is.isEmitEvent = true
	event := NewEventKey(eventType, a, is.cli, nil)
	is.handler.Handle(event)
}

func runeAt(s string, index int) rune {
	for i, r := range s {
		if i == index {
			return r
		}
	}
	panic(fmt.Sprintf("not found rune at %d", index))
}

// 检查 prefix 是否是 keyActions 中某个 key 的前缀
func prefixMatchKeyActions(prefix string) bool {
	for key := range keyActions {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

var keyActions = map[string]EventType{
	// Control-Space (Also for Ctrl-@)
	"\x00": EventTypeCtrlSpace,
	"\x01": EventTypeCtrlA,
	"\x02": EventTypeCtrlB,
	"\x03": EventTypeCtrlC,
	"\x04": EventTypeCtrlD,
	"\x05": EventTypeCtrlE,
	"\x06": EventTypeCtrlF,
	"\x07": EventTypeCtrlG,
	// Control-H (8) (Identical to '\b')
	"\x08": EventTypeCtrlH,
	// Control-I (9) (Identical to '\t')
	"\x09": EventTypeCtrlI,
	// Control-J (10) (Identical to '\n')
	"\x0a": EventTypeCtrlJ,
	"\x0b": EventTypeCtrlK,
	// Control-L (Clear; form feed)
	"\x0c": EventTypeCtrlL,
	// Control-M (13) (Identical to '\r')
	"\x0d": EventTypeCtrlM,
	"\x0e": EventTypeCtrlN,
	"\x0f": EventTypeCtrlO,
	"\x10": EventTypeCtrlP,
	"\x11": EventTypeCtrlQ,
	"\x12": EventTypeCtrlR,
	"\x13": EventTypeCtrlS,
	"\x14": EventTypeCtrlT,
	"\x15": EventTypeCtrlU,
	"\x16": EventTypeCtrlV,
	"\x17": EventTypeCtrlW,
	"\x18": EventTypeCtrlX,
	"\x19": EventTypeCtrlY,
	"\x1a": EventTypeCtrlZ,

	// Both Control-\ and Ctrl-|
	"\x1c": EventTypeCtrlBackslash,
	// Control-]
	"\x1d": EventTypeCtrlSquareClose,
	// Control-^
	"\x1e": EventTypeCtrlCircumflex,
	// Control-underscore (Also for Ctrl-hypen.)
	"\x1f": EventTypeCtrlUnderscore,
	// (127) Backspace
	"\x7f":    EventTypeBackspace,
	"\x1b[A":  EventTypeArrowUp,
	"\x1b[B":  EventTypeArrowDown,
	"\x1b[C":  EventTypeArrowRight,
	"\x1b[D":  EventTypeArrowLeft,
	"\x1b[H":  EventTypeHome,
	"\x1bOH":  EventTypeHome,
	"\x1b[F":  EventTypeEnd,
	"\x1bOF":  EventTypeEnd,
	"\x1b[3~": EventTypeDeleteAction,
	// xterm, gnome-terminal.
	"\x1b[3;2~": EventTypeShiftDelete,
	// tmux
	"\x1b[1~": EventTypeHome,
	// tmux
	"\x1b[4~": EventTypeEnd,
	"\x1b[5~": EventTypePageUp,
	"\x1b[6~": EventTypePageDown,
	// xrvt
	"\x1b[7~": EventTypeHome,
	// xrvt
	"\x1b[8~": EventTypeEnd,
	// shift + tab
	"\x1b[Z": EventTypeBacktab,

	"\x1bOP":   EventTypeF1,
	"\x1bOQ":   EventTypeF2,
	"\x1bOR":   EventTypeF3,
	"\x1bOS":   EventTypeF4,
	"\x1b[15~": EventTypeF5,
	"\x1b[17~": EventTypeF6,
	"\x1b[18~": EventTypeF7,
	"\x1b[19~": EventTypeF8,
	"\x1b[20~": EventTypeF9,
	"\x1b[21~": EventTypeF10,
	"\x1b[23~": EventTypeF11,
	"\x1b[24~": EventTypeF12,
	"\x1b[25~": EventTypeF13,
	"\x1b[26~": EventTypeF14,
	"\x1b[28~": EventTypeF15,
	"\x1b[29~": EventTypeF16,
	"\x1b[31~": EventTypeF17,
	"\x1b[32~": EventTypeF18,
	"\x1b[33~": EventTypeF19,
	"\x1b[34~": EventTypeF20,
}
