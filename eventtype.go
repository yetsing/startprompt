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
	CtrlA EventType = iota
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

	CtrlSpace
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
	ShiftDelete
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

	Tab = CtrlI
)
