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
	for true {
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
				is.callHandler(escape_action)
			} else {
				// 如果不是快捷键操作，那么就是正常的输入
				is.callHandler(insert_char, first)
			}
			// 剩余的字符放到缓冲中，留待下次循环的时候处理
			buffer = []rune(is.previous[utf8.RuneLen(first):])
			buffer = append(buffer, r)
			is.previous = ""
		} else {
			is.callHandler(insert_char, r)
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
	ctrl_space Event = iota
	ctrl_a
	ctrl_b
	ctrl_c
	ctrl_d
	ctrl_e
	ctrl_f
	ctrl_g
	ctrl_h
	ctrl_i
	ctrl_j
	ctrl_k
	ctrl_l
	ctrl_m
	ctrl_n
	ctrl_o
	ctrl_p
	ctrl_q
	ctrl_r
	ctrl_s
	ctrl_t
	ctrl_u
	ctrl_v
	ctrl_w
	ctrl_x
	ctrl_y
	ctrl_z
	ctrl_backslash
	ctrl_square_close
	ctrl_circumflex
	ctrl_underscore
	backspace
	arrow_up
	arrow_down
	arrow_right
	arrow_left
	home
	end
	delete_action
	page_up
	page_down
	backtab
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
	escape_action
	insert_char
)

var keyActions = map[string]Event{
	"\x00": ctrl_space,
	"\x01": ctrl_a,
	"\x02": ctrl_b,
	"\x03": ctrl_c,
	"\x04": ctrl_d,
	"\x05": ctrl_e,
	"\x06": ctrl_f,
	"\x07": ctrl_g,
	// (Identical to '\b')
	"\x08": ctrl_h,
	// (Identical to '\t')
	"\x09": ctrl_i,
	// (Identical to '\n')
	"\x0a": ctrl_j,
	"\x0b": ctrl_k,
	"\x0c": ctrl_l,
	// (Identical to '\r')
	"\x0d":     ctrl_m,
	"\x0e":     ctrl_n,
	"\x0f":     ctrl_o,
	"\x10":     ctrl_p,
	"\x11":     ctrl_q,
	"\x12":     ctrl_r,
	"\x13":     ctrl_s,
	"\x14":     ctrl_t,
	"\x15":     ctrl_u,
	"\x16":     ctrl_v,
	"\x17":     ctrl_w,
	"\x18":     ctrl_x,
	"\x19":     ctrl_y,
	"\x1a":     ctrl_z,
	"\x1c":     ctrl_backslash,
	"\x1d":     ctrl_square_close,
	"\x1e":     ctrl_circumflex,
	"\x1f":     ctrl_underscore,
	"\x7f":     backspace,
	"\x1b[A":   arrow_up,
	"\x1b[B":   arrow_down,
	"\x1b[C":   arrow_right,
	"\x1b[D":   arrow_left,
	"\x1b[H":   home,
	"\x1b[F":   end,
	"\x1b[3~":  delete_action,
	"\x1b[1~":  home,
	"\x1b[4~":  end,
	"\x1b[5~":  page_up,
	"\x1b[6~":  page_down,
	"\x1b[7~":  home,
	"\x1b[8~":  end,
	"\x1b[Z":   backtab,
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
