package startprompt

import "github.com/gdamore/tcell/v2"

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
	"<draw>",
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
	// EventTypeCtrlH Control-H (8) (Identical to '\b')
	EventTypeCtrlH
	// EventTypeCtrlI Control-I (9) (Identical to '\t')
	EventTypeCtrlI
	// EventTypeCtrlJ Control-J (10) (Identical to '\n')
	EventTypeCtrlJ
	EventTypeCtrlK
	// EventTypeCtrlL Control-L (Clear; form feed)
	EventTypeCtrlL
	// EventTypeCtrlM Control-M (13) (Identical to '\r')
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
	// EventTypeCtrlSpace Control-Space (Also for Ctrl-@)
	EventTypeCtrlSpace
	// EventTypeCtrlBackslash Both Control-\ and Ctrl-|
	EventTypeCtrlBackslash
	// EventTypeCtrlSquareClose Control-]
	EventTypeCtrlSquareClose
	// EventTypeCtrlCircumflex Control-^
	EventTypeCtrlCircumflex
	// EventTypeCtrlUnderscore Control-underscore (Also for Ctrl-hypen.)
	EventTypeCtrlUnderscore
	// EventTypeBackspace (127) Backspace
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
	EventTypeEscape
	EventTypeInsertChar

	EventTypeDraw
	EventTypeTab = EventTypeCtrlI
)

// tkeyMapping tcell key 事件映射
var tkeyMapping = map[tcell.Key]EventType{
	tcell.KeyRune:    EventTypeInsertChar,
	tcell.KeyUp:      EventTypeArrowUp,
	tcell.KeyDown:    EventTypeArrowDown,
	tcell.KeyRight:   EventTypeArrowRight,
	tcell.KeyLeft:    EventTypeArrowLeft,
	tcell.KeyPgUp:    EventTypePageUp,
	tcell.KeyPgDn:    EventTypePageDown,
	tcell.KeyHome:    EventTypeHome,
	tcell.KeyEnd:     EventTypeEnd,
	tcell.KeyBacktab: EventTypeBacktab,
	tcell.KeyEsc:     EventTypeEscape,

	tcell.KeyCtrlSpace:      EventTypeCtrlSpace,
	tcell.KeyCtrlA:          EventTypeCtrlA,
	tcell.KeyCtrlB:          EventTypeCtrlB,
	tcell.KeyCtrlC:          EventTypeCtrlC,
	tcell.KeyCtrlD:          EventTypeCtrlD,
	tcell.KeyCtrlE:          EventTypeCtrlE,
	tcell.KeyCtrlF:          EventTypeCtrlF,
	tcell.KeyCtrlG:          EventTypeCtrlG,
	tcell.KeyCtrlH:          EventTypeCtrlH,
	tcell.KeyCtrlI:          EventTypeCtrlI,
	tcell.KeyCtrlJ:          EventTypeCtrlJ,
	tcell.KeyCtrlK:          EventTypeCtrlK,
	tcell.KeyCtrlL:          EventTypeCtrlL,
	tcell.KeyCtrlM:          EventTypeCtrlM,
	tcell.KeyCtrlN:          EventTypeCtrlN,
	tcell.KeyCtrlO:          EventTypeCtrlO,
	tcell.KeyCtrlP:          EventTypeCtrlP,
	tcell.KeyCtrlQ:          EventTypeCtrlQ,
	tcell.KeyCtrlR:          EventTypeCtrlR,
	tcell.KeyCtrlS:          EventTypeCtrlS,
	tcell.KeyCtrlT:          EventTypeCtrlT,
	tcell.KeyCtrlU:          EventTypeCtrlU,
	tcell.KeyCtrlV:          EventTypeCtrlV,
	tcell.KeyCtrlW:          EventTypeCtrlW,
	tcell.KeyCtrlX:          EventTypeCtrlX,
	tcell.KeyCtrlY:          EventTypeCtrlY,
	tcell.KeyCtrlZ:          EventTypeCtrlZ,
	tcell.KeyCtrlBackslash:  EventTypeCtrlBackslash,
	tcell.KeyCtrlRightSq:    EventTypeCtrlSquareClose,
	tcell.KeyCtrlCarat:      EventTypeCtrlCircumflex,
	tcell.KeyCtrlUnderscore: EventTypeCtrlUnderscore,
	tcell.KeyF1:             EventTypeF1,
	tcell.KeyF2:             EventTypeF2,
	tcell.KeyF3:             EventTypeF3,
	tcell.KeyF4:             EventTypeF4,
	tcell.KeyF5:             EventTypeF5,
	tcell.KeyF6:             EventTypeF6,
	tcell.KeyF7:             EventTypeF7,
	tcell.KeyF8:             EventTypeF8,
	tcell.KeyF9:             EventTypeF9,
	tcell.KeyF10:            EventTypeF10,
	tcell.KeyF11:            EventTypeF11,
	tcell.KeyF12:            EventTypeF12,
	tcell.KeyF13:            EventTypeF13,
	tcell.KeyF14:            EventTypeF14,
	tcell.KeyF15:            EventTypeF15,
	tcell.KeyF16:            EventTypeF16,
	tcell.KeyF17:            EventTypeF17,
	tcell.KeyF18:            EventTypeF18,
	tcell.KeyF19:            EventTypeF19,
	tcell.KeyF20:            EventTypeF20,
}
