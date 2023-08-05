package terminalcolor

import "strconv"

/*
8-bit 颜色，可以表示 256 种颜色
参考：https://en.wikipedia.org/wiki/ANSI_escape_code 中 "8-bit" 一节
转义序列格式如下

ESC[38;5;⟨n⟩m Select foreground color      where n is a number from the table below
ESC[48;5;⟨n⟩m Select background color
  0-  7:  standard colors (as in ESC [ 30–37 m)
  8- 15:  high intensity colors (as in ESC [ 90–97 m)
 16-231:  6 × 6 × 6 cube (216 colors): 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
    6 个值的调色板 (0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff) ， r g b 从调色板中取，可表示的颜色就是 6 × 6 × 6
    16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
    这个公式中的 r g b 表示的是调色板索引，得到的值就是颜色表的序号
232-255:  grayscale from dark to light in 24 steps

*/

// Color256 定义颜色表的序号
type Color256 int

const Color256Default Color256 = -1

//goland:noinspection GoUnusedConst
const (
	Color256No0 Color256 = iota
	Color256No1
	Color256No2
	Color256No3
	Color256No4
	Color256No5
	Color256No6
	Color256No7
	Color256No8
	Color256No9
	Color256No10
	Color256No11
	Color256No12
	Color256No13
	Color256No14
	Color256No15
	Color256No16
	Color256No17
	Color256No18
	Color256No19
	Color256No20
	Color256No21
	Color256No22
	Color256No23
	Color256No24
	Color256No25
	Color256No26
	Color256No27
	Color256No28
	Color256No29
	Color256No30
	Color256No31
	Color256No32
	Color256No33
	Color256No34
	Color256No35
	Color256No36
	Color256No37
	Color256No38
	Color256No39
	Color256No40
	Color256No41
	Color256No42
	Color256No43
	Color256No44
	Color256No45
	Color256No46
	Color256No47
	Color256No48
	Color256No49
	Color256No50
	Color256No51
	Color256No52
	Color256No53
	Color256No54
	Color256No55
	Color256No56
	Color256No57
	Color256No58
	Color256No59
	Color256No60
	Color256No61
	Color256No62
	Color256No63
	Color256No64
	Color256No65
	Color256No66
	Color256No67
	Color256No68
	Color256No69
	Color256No70
	Color256No71
	Color256No72
	Color256No73
	Color256No74
	Color256No75
	Color256No76
	Color256No77
	Color256No78
	Color256No79
	Color256No80
	Color256No81
	Color256No82
	Color256No83
	Color256No84
	Color256No85
	Color256No86
	Color256No87
	Color256No88
	Color256No89
	Color256No90
	Color256No91
	Color256No92
	Color256No93
	Color256No94
	Color256No95
	Color256No96
	Color256No97
	Color256No98
	Color256No99
	Color256No100
	Color256No101
	Color256No102
	Color256No103
	Color256No104
	Color256No105
	Color256No106
	Color256No107
	Color256No108
	Color256No109
	Color256No110
	Color256No111
	Color256No112
	Color256No113
	Color256No114
	Color256No115
	Color256No116
	Color256No117
	Color256No118
	Color256No119
	Color256No120
	Color256No121
	Color256No122
	Color256No123
	Color256No124
	Color256No125
	Color256No126
	Color256No127
	Color256No128
	Color256No129
	Color256No130
	Color256No131
	Color256No132
	Color256No133
	Color256No134
	Color256No135
	Color256No136
	Color256No137
	Color256No138
	Color256No139
	Color256No140
	Color256No141
	Color256No142
	Color256No143
	Color256No144
	Color256No145
	Color256No146
	Color256No147
	Color256No148
	Color256No149
	Color256No150
	Color256No151
	Color256No152
	Color256No153
	Color256No154
	Color256No155
	Color256No156
	Color256No157
	Color256No158
	Color256No159
	Color256No160
	Color256No161
	Color256No162
	Color256No163
	Color256No164
	Color256No165
	Color256No166
	Color256No167
	Color256No168
	Color256No169
	Color256No170
	Color256No171
	Color256No172
	Color256No173
	Color256No174
	Color256No175
	Color256No176
	Color256No177
	Color256No178
	Color256No179
	Color256No180
	Color256No181
	Color256No182
	Color256No183
	Color256No184
	Color256No185
	Color256No186
	Color256No187
	Color256No188
	Color256No189
	Color256No190
	Color256No191
	Color256No192
	Color256No193
	Color256No194
	Color256No195
	Color256No196
	Color256No197
	Color256No198
	Color256No199
	Color256No200
	Color256No201
	Color256No202
	Color256No203
	Color256No204
	Color256No205
	Color256No206
	Color256No207
	Color256No208
	Color256No209
	Color256No210
	Color256No211
	Color256No212
	Color256No213
	Color256No214
	Color256No215
	Color256No216
	Color256No217
	Color256No218
	Color256No219
	Color256No220
	Color256No221
	Color256No222
	Color256No223
	Color256No224
	Color256No225
	Color256No226
	Color256No227
	Color256No228
	Color256No229
	Color256No230
	Color256No231
	Color256No232
	Color256No233
	Color256No234
	Color256No235
	Color256No236
	Color256No237
	Color256No238
	Color256No239
	Color256No240
	Color256No241
	Color256No242
	Color256No243
	Color256No244
	Color256No245
	Color256No246
	Color256No247
	Color256No248
	Color256No249
	Color256No250
	Color256No251
	Color256No252
	Color256No253
	Color256No254
	Color256No255
)

type Color256Style struct {
	fg        Color256
	bg        Color256
	bold      bool
	underline bool
	italic    bool
}

//goland:noinspection GoUnusedExportedFunction
func NewColor256Style(fg, bg Color256, bold, underline, italic bool) *Color256Style {
	return &Color256Style{
		fg:        fg,
		bg:        bg,
		bold:      bold,
		underline: underline,
		italic:    italic,
	}
}

func (c *Color256Style) ColorEscape() string {
	var attrs []string
	if c.fg != Color256Default {
		attrs = append(attrs, "38", "5", strconv.Itoa(int(c.fg)))
	}
	if c.bg != Color256Default {
		// 背景的数字 = 前景的数字 + 10
		attrs = append(attrs, "48", "5", strconv.Itoa(int(c.bg)+10))
	}
	if c.bold {
		attrs = append(attrs, "01")
	}
	if c.underline {
		attrs = append(attrs, "04")
	}
	if c.italic {
		attrs = append(attrs, "03")
	}
	return escapeAttrs(attrs)
}

func (c *Color256Style) ResetEscape() string {
	var attrs []string
	if c.fg != Color256Default {
		attrs = append(attrs, "39")
	}
	if c.bg != Color256Default {
		attrs = append(attrs, "49")
	}
	if c.bold || c.underline || c.italic {
		attrs = append(attrs, "00")
	}
	return escapeAttrs(attrs)
}
