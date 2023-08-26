package startprompt

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yetsing/startprompt/keys"
)

// 解析 VT100 输入流数据
// 参考：https://vt100.net/docs/vt100-ug/chapter3.html

func NewInputStream(handler InputStreamHandler) *InputStream {
	return &InputStream{handler: handler}
}

type InputStream struct {
	handler  InputStreamHandler
	previous string
}

// Feed 根据输入触发对应的事件
func (is *InputStream) Feed(r rune) {
	var buffer []rune
	for {
		key := string(r)
		if len(is.previous) > 0 {
			key = is.previous + key
		}
		// 检查是不是快捷键操作
		action, found := keyActions[key]
		if found {
			is.callHandler(action)
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
				is.callHandler(keys.EscapeAction)
			} else {
				// 如果不是快捷键操作，那么就是正常的输入
				is.callHandler(keys.InsertChar, first)
			}
			// 剩余的字符放到缓冲中，留待下次循环的时候处理
			buffer = []rune(is.previous[utf8.RuneLen(first):])
			buffer = append(buffer, r)
			is.previous = ""
		} else {
			is.callHandler(keys.InsertChar, r)
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

func (is *InputStream) callHandler(event keys.Event, a ...rune) {
	is.handler.Handle(event, a...)
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

var keyActions = map[string]keys.Event{
	// Control-Space (Also for Ctrl-@)
	"\x00": keys.CtrlSpace,
	"\x01": keys.CtrlA,
	"\x02": keys.CtrlB,
	"\x03": keys.CtrlC,
	"\x04": keys.CtrlD,
	"\x05": keys.CtrlE,
	"\x06": keys.CtrlF,
	"\x07": keys.CtrlG,
	// Control-H (8) (Identical to '\b')
	"\x08": keys.CtrlH,
	// Control-I (9) (Identical to '\t')
	"\x09": keys.CtrlI,
	// Control-J (10) (Identical to '\n')
	"\x0a": keys.CtrlJ,
	"\x0b": keys.CtrlK,
	// Control-L (clear; form feed)
	"\x0c": keys.CtrlL,
	// Control-M (13) (Identical to '\r')
	"\x0d": keys.CtrlM,
	"\x0e": keys.CtrlN,
	"\x0f": keys.CtrlO,
	"\x10": keys.CtrlP,
	"\x11": keys.CtrlQ,
	"\x12": keys.CtrlR,
	"\x13": keys.CtrlS,
	"\x14": keys.CtrlT,
	"\x15": keys.CtrlU,
	"\x16": keys.CtrlV,
	"\x17": keys.CtrlW,
	"\x18": keys.CtrlX,
	"\x19": keys.CtrlY,
	"\x1a": keys.CtrlZ,

	// Both Control-\ and Ctrl-|
	"\x1c": keys.CtrlBackslash,
	// Control-]
	"\x1d": keys.CtrlSquareClose,
	// Control-^
	"\x1e": keys.CtrlCircumflex,
	// Control-underscore (Also for Ctrl-hypen.)
	"\x1f": keys.CtrlUnderscore,
	// (127) Backspace
	"\x7f":    keys.Backspace,
	"\x1b[A":  keys.ArrowUp,
	"\x1b[B":  keys.ArrowDown,
	"\x1b[C":  keys.ArrowRight,
	"\x1b[D":  keys.ArrowLeft,
	"\x1b[H":  keys.Home,
	"\x1bOH":  keys.Home,
	"\x1b[F":  keys.End,
	"\x1bOF":  keys.End,
	"\x1b[3~": keys.DeleteAction,
	// xterm, gnome-terminal.
	"\x1b[3;2~": keys.ShiftDelete,
	// tmux
	"\x1b[1~": keys.Home,
	// tmux
	"\x1b[4~": keys.End,
	"\x1b[5~": keys.PageUp,
	"\x1b[6~": keys.PageDown,
	// xrvt
	"\x1b[7~": keys.Home,
	// xrvt
	"\x1b[8~": keys.End,
	// shift + tab
	"\x1b[Z": keys.Backtab,

	"\x1bOP":   keys.F1,
	"\x1bOQ":   keys.F2,
	"\x1bOR":   keys.F3,
	"\x1bOS":   keys.F4,
	"\x1b[15~": keys.F5,
	"\x1b[17~": keys.F6,
	"\x1b[18~": keys.F7,
	"\x1b[19~": keys.F8,
	"\x1b[20~": keys.F9,
	"\x1b[21~": keys.F10,
	"\x1b[23~": keys.F11,
	"\x1b[24~": keys.F12,
	"\x1b[25~": keys.F13,
	"\x1b[26~": keys.F14,
	"\x1b[28~": keys.F15,
	"\x1b[29~": keys.F16,
	"\x1b[31~": keys.F17,
	"\x1b[32~": keys.F18,
	"\x1b[33~": keys.F19,
	"\x1b[34~": keys.F20,
}
