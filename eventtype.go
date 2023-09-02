package startprompt

type EventType int

var eventTypeStr = []string{
	"<ctrl_a>",
	"<ctrl_b>",
	"<ctrl_c>",
	"<ctrl_d>",
	"<ctrl_e>",
	"<ctrl_f>",
	"<ctrl_g>",
	"<ctrl_h>",
	"<ctrl_i>",
	"<ctrl_j>",
	"<ctrl_k>",
	"<ctrl_l>",
	"<ctrl_m>",
	"<ctrl_n>",
	"<ctrl_o>",
	"<ctrl_p>",
	"<ctrl_q>",
	"<ctrl_r>",
	"<ctrl_s>",
	"<ctrl_t>",
	"<ctrl_u>",
	"<ctrl_v>",
	"<ctrl_w>",
	"<ctrl_x>",
	"<ctrl_y>",
	"<ctrl_z>",

	"<ctrl_space>",
	"<ctrl_backslash>",
	"<ctrl_square_close>",
	"<ctrl_circumflex>",
	"<ctrl_underscore>",
	"<backspace>",

	"<arrow_up>",
	"<arrow_down>",
	"<arrow_right>",
	"<arrow_left>",
	"<home>",
	"<end>",
	"<delete_action>",
	"<ShiftDelete>",
	"<page_up>",
	"<page_down>",
	"<backtab>",

	"<F1>",
	"<F2>",
	"<F3>",
	"<F4>",
	"<F5>",
	"<F6>",
	"<F7>",
	"<F8>",
	"<F9>",
	"<F10>",
	"<F11>",
	"<F12>",
	"<F13>",
	"<F14>",
	"<F15>",
	"<F16>",
	"<F17>",
	"<F18>",
	"<F19>",
	"<F20>",
	"<escape>",
	"<insert_char>",
}

func (a EventType) String() string {
	return eventTypeStr[a]
}

//goland:noinspection GoUnusedConst
const (
	EventTypeCtrlA EventType = iota
	EventTypeCtrlB
	EventTypeCtrlC
	EventTypeCtrlD
	EventTypeCtrlE
	EventTypeCtrlF
	EventTypeCtrlG
	EventTypeCtrlH
	EventTypeCtrlI
	EventTypeCtrlJ
	EventTypeCtrlK
	EventTypeCtrlL
	EventTypeCtrlM
	EventTypeCtrlN
	EventTypeCtrlO
	EventTypeCtrlP
	EventTypeCtrlQ
	EventTypeCtrlR
	EventTypeCtrlS
	EventTypeCtrlT
	EventTypeCtrlU
	EventTypeCtrlV
	EventTypeCtrlW
	EventTypeCtrlX
	EventTypeCtrlY
	EventTypeCtrlZ

	EventTypeCtrlSpace
	EventTypeCtrlBackslash
	EventTypeCtrlSquareClose
	EventTypeCtrlCircumflex
	EventTypeCtrlUnderscore
	EventTypeBackspace
	EventTypeArrowUp
	EventTypeArrowDown
	EventTypeArrowRight
	EventTypeArrowLeft
	EventTypeHome
	EventTypeEnd
	EventTypeDeleteAction
	EventTypeShiftDelete
	EventTypePageUp
	EventTypePageDown
	EventTypeBacktab
	EventTypeF1
	EventTypeF2
	EventTypeF3
	EventTypeF4
	EventTypeF5
	EventTypeF6
	EventTypeF7
	EventTypeF8
	EventTypeF9
	EventTypeF10
	EventTypeF11
	EventTypeF12
	EventTypeF13
	EventTypeF14
	EventTypeF15
	EventTypeF16
	EventTypeF17
	EventTypeF18
	EventTypeF19
	EventTypeF20
	EventTypeEscapeAction
	EventTypeInsertChar

	EventTypeTab = EventTypeCtrlI
)
