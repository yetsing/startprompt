package startprompt

import (
	"fmt"
	"strings"
	"unicode/utf8"
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
				is.callHandler(EscapeAction)
			} else {
				// 如果不是快捷键操作，那么就是正常的输入
				is.callHandler(InsertChar, first)
			}
			// 剩余的字符放到缓冲中，留待下次循环的时候处理
			buffer = []rune(is.previous[utf8.RuneLen(first):])
			buffer = append(buffer, r)
			is.previous = ""
		} else {
			is.callHandler(InsertChar, r)
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

func (is *InputStream) callHandler(event Event, a ...rune) {
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

type Event int

var eventStr = []string{
	"ctrl_space", "ctrl_a", "ctrl_b", "ctrl_c", "ctrl_d", "ctrl_e", "ctrl_f", "ctrl_g", "ctrl_h", "ctrl_i", "ctrl_j",
	"ctrl_k", "ctrl_l", "ctrl_m", "ctrl_n", "ctrl_o", "ctrl_p", "ctrl_q", "ctrl_r", "ctrl_s", "ctrl_t", "ctrl_u",
	"ctrl_v", "ctrl_w", "ctrl_x", "ctrl_y", "ctrl_z",
	"ctrl_backslash", "ctrl_square_close", "ctrl_circumflex", "ctrl_underscore",
	"backspace",
	"arrow_up", "arrow_down", "arrow_right", "arrow_left",
	"home", "end", "delete_action",
	"page_up", "page_down", "backtab",
	"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10",
	"F11", "F12", "F13", "F14", "F15", "F16", "F17", "F18", "F19", "F20",
	"escape",
	"insert_char",
}

func (a Event) String() string {
	return eventStr[a]
}

const (
	CtrlSpace Event = iota
	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	CtrlI
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	CtrlBackslash
	CtrlSquareClose
	CtrlCircumflex
	CtrlUnderscore
	Backspace
	ArrowUp
	ArrowDown
	ArrowRight
	ArrowLeft
	Home
	End
	DeleteAction
	PageUp
	PageDown
	Backtab
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
	F16
	F17
	F18
	F19
	F20
	EscapeAction
	InsertChar
)

var keyActions = map[string]Event{
	// Control-Space (Also for Ctrl-@)
	"\x00": CtrlSpace,
	"\x01": CtrlA,
	"\x02": CtrlB,
	"\x03": CtrlC,
	"\x04": CtrlD,
	"\x05": CtrlE,
	"\x06": CtrlF,
	"\x07": CtrlG,
	// Control-H (8) (Identical to '\b')
	"\x08": CtrlH,
	// Control-I (9) (Identical to '\t')
	"\x09": CtrlI,
	// Control-J (10) (Identical to '\n')
	"\x0a": CtrlJ,
	"\x0b": CtrlK,
	// Control-L (clear; form feed)
	"\x0c": CtrlL,
	// Control-M (13) (Identical to '\r')
	"\x0d": CtrlM,
	"\x0e": CtrlN,
	"\x0f": CtrlO,
	"\x10": CtrlP,
	"\x11": CtrlQ,
	"\x12": CtrlR,
	"\x13": CtrlS,
	"\x14": CtrlT,
	"\x15": CtrlU,
	"\x16": CtrlV,
	"\x17": CtrlW,
	"\x18": CtrlX,
	"\x19": CtrlY,
	"\x1a": CtrlZ,
	// Both Control-\ and Ctrl-|
	"\x1c": CtrlBackslash,
	// Control-]
	"\x1d": CtrlSquareClose,
	// Control-^
	"\x1e": CtrlCircumflex,
	// Control-underscore (Also for Ctrl-hypen.)
	"\x1f": CtrlUnderscore,
	// (127) Backspace
	"\x7f":    Backspace,
	"\x1b[A":  ArrowUp,
	"\x1b[B":  ArrowDown,
	"\x1b[C":  ArrowRight,
	"\x1b[D":  ArrowLeft,
	"\x1b[H":  Home,
	"\x1b[F":  End,
	"\x1b[3~": DeleteAction,
	// tmux
	"\x1b[1~": Home,
	// tmux
	"\x1b[4~": End,
	"\x1b[5~": PageUp,
	"\x1b[6~": PageDown,
	// xrvt
	"\x1b[7~": Home,
	// xrvt
	"\x1b[8~": End,
	// shift + tab
	"\x1b[Z": Backtab,

	"\x1bOP":   F1,
	"\x1bOQ":   F2,
	"\x1bOR":   F3,
	"\x1bOS":   F4,
	"\x1b[15~": F5,
	"\x1b[17~": F6,
	"\x1b[18~": F7,
	"\x1b[19~": F8,
	"\x1b[20~": F9,
	"\x1b[21~": F10,
	"\x1b[23~": F11,
	"\x1b[24~": F12,
	"\x1b[25~": F13,
	"\x1b[26~": F14,
	"\x1b[28~": F15,
	"\x1b[29~": F16,
	"\x1b[31~": F17,
	"\x1b[32~": F18,
	"\x1b[33~": F19,
	"\x1b[34~": F20,
}
