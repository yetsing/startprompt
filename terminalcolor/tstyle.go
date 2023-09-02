package terminalcolor

import "github.com/gdamore/tcell/v2"

var tcolorMapping = map[Color]tcell.Color{
	Color16Black:         tcell.ColorBlack,
	Color16Red:           tcell.ColorMaroon,
	Color16Green:         tcell.ColorGreen,
	Color16Yellow:        tcell.ColorOlive,
	Color16Blue:          tcell.ColorNavy,
	Color16Magenta:       tcell.ColorPurple,
	Color16Cyan:          tcell.ColorTeal,
	Color16Gray:          tcell.ColorSilver,
	Color16BrightBlack:   tcell.ColorGray,
	Color16BrightRed:     tcell.ColorRed,
	Color16BrightGreen:   tcell.ColorLime,
	Color16BrightYellow:  tcell.ColorYellow,
	Color16BrightBlue:    tcell.ColorBlue,
	Color16BrightMagenta: tcell.ColorFuchsia,
	Color16BrightCyan:    tcell.ColorAqua,
	Color16BrightWhite:   tcell.ColorWhite,

	Color256No0:   tcell.ColorBlack,
	Color256No1:   tcell.ColorMaroon,
	Color256No2:   tcell.ColorGreen,
	Color256No3:   tcell.ColorOlive,
	Color256No4:   tcell.ColorNavy,
	Color256No5:   tcell.ColorPurple,
	Color256No6:   tcell.ColorTeal,
	Color256No7:   tcell.ColorSilver,
	Color256No8:   tcell.ColorGray,
	Color256No9:   tcell.ColorRed,
	Color256No10:  tcell.ColorLime,
	Color256No11:  tcell.ColorYellow,
	Color256No12:  tcell.ColorBlue,
	Color256No13:  tcell.ColorFuchsia,
	Color256No14:  tcell.ColorAqua,
	Color256No15:  tcell.ColorWhite,
	Color256No16:  tcell.Color16,
	Color256No17:  tcell.Color17,
	Color256No18:  tcell.Color18,
	Color256No19:  tcell.Color19,
	Color256No20:  tcell.Color20,
	Color256No21:  tcell.Color21,
	Color256No22:  tcell.Color22,
	Color256No23:  tcell.Color23,
	Color256No24:  tcell.Color24,
	Color256No25:  tcell.Color25,
	Color256No26:  tcell.Color26,
	Color256No27:  tcell.Color27,
	Color256No28:  tcell.Color28,
	Color256No29:  tcell.Color29,
	Color256No30:  tcell.Color30,
	Color256No31:  tcell.Color31,
	Color256No32:  tcell.Color32,
	Color256No33:  tcell.Color33,
	Color256No34:  tcell.Color34,
	Color256No35:  tcell.Color35,
	Color256No36:  tcell.Color36,
	Color256No37:  tcell.Color37,
	Color256No38:  tcell.Color38,
	Color256No39:  tcell.Color39,
	Color256No40:  tcell.Color40,
	Color256No41:  tcell.Color41,
	Color256No42:  tcell.Color42,
	Color256No43:  tcell.Color43,
	Color256No44:  tcell.Color44,
	Color256No45:  tcell.Color45,
	Color256No46:  tcell.Color46,
	Color256No47:  tcell.Color47,
	Color256No48:  tcell.Color48,
	Color256No49:  tcell.Color49,
	Color256No50:  tcell.Color50,
	Color256No51:  tcell.Color51,
	Color256No52:  tcell.Color52,
	Color256No53:  tcell.Color53,
	Color256No54:  tcell.Color54,
	Color256No55:  tcell.Color55,
	Color256No56:  tcell.Color56,
	Color256No57:  tcell.Color57,
	Color256No58:  tcell.Color58,
	Color256No59:  tcell.Color59,
	Color256No60:  tcell.Color60,
	Color256No61:  tcell.Color61,
	Color256No62:  tcell.Color62,
	Color256No63:  tcell.Color63,
	Color256No64:  tcell.Color64,
	Color256No65:  tcell.Color65,
	Color256No66:  tcell.Color66,
	Color256No67:  tcell.Color67,
	Color256No68:  tcell.Color68,
	Color256No69:  tcell.Color69,
	Color256No70:  tcell.Color70,
	Color256No71:  tcell.Color71,
	Color256No72:  tcell.Color72,
	Color256No73:  tcell.Color73,
	Color256No74:  tcell.Color74,
	Color256No75:  tcell.Color75,
	Color256No76:  tcell.Color76,
	Color256No77:  tcell.Color77,
	Color256No78:  tcell.Color78,
	Color256No79:  tcell.Color79,
	Color256No80:  tcell.Color80,
	Color256No81:  tcell.Color81,
	Color256No82:  tcell.Color82,
	Color256No83:  tcell.Color83,
	Color256No84:  tcell.Color84,
	Color256No85:  tcell.Color85,
	Color256No86:  tcell.Color86,
	Color256No87:  tcell.Color87,
	Color256No88:  tcell.Color88,
	Color256No89:  tcell.Color89,
	Color256No90:  tcell.Color90,
	Color256No91:  tcell.Color91,
	Color256No92:  tcell.Color92,
	Color256No93:  tcell.Color93,
	Color256No94:  tcell.Color94,
	Color256No95:  tcell.Color95,
	Color256No96:  tcell.Color96,
	Color256No97:  tcell.Color97,
	Color256No98:  tcell.Color98,
	Color256No99:  tcell.Color99,
	Color256No100: tcell.Color100,
	Color256No101: tcell.Color101,
	Color256No102: tcell.Color102,
	Color256No103: tcell.Color103,
	Color256No104: tcell.Color104,
	Color256No105: tcell.Color105,
	Color256No106: tcell.Color106,
	Color256No107: tcell.Color107,
	Color256No108: tcell.Color108,
	Color256No109: tcell.Color109,
	Color256No110: tcell.Color110,
	Color256No111: tcell.Color111,
	Color256No112: tcell.Color112,
	Color256No113: tcell.Color113,
	Color256No114: tcell.Color114,
	Color256No115: tcell.Color115,
	Color256No116: tcell.Color116,
	Color256No117: tcell.Color117,
	Color256No118: tcell.Color118,
	Color256No119: tcell.Color119,
	Color256No120: tcell.Color120,
	Color256No121: tcell.Color121,
	Color256No122: tcell.Color122,
	Color256No123: tcell.Color123,
	Color256No124: tcell.Color124,
	Color256No125: tcell.Color125,
	Color256No126: tcell.Color126,
	Color256No127: tcell.Color127,
	Color256No128: tcell.Color128,
	Color256No129: tcell.Color129,
	Color256No130: tcell.Color130,
	Color256No131: tcell.Color131,
	Color256No132: tcell.Color132,
	Color256No133: tcell.Color133,
	Color256No134: tcell.Color134,
	Color256No135: tcell.Color135,
	Color256No136: tcell.Color136,
	Color256No137: tcell.Color137,
	Color256No138: tcell.Color138,
	Color256No139: tcell.Color139,
	Color256No140: tcell.Color140,
	Color256No141: tcell.Color141,
	Color256No142: tcell.Color142,
	Color256No143: tcell.Color143,
	Color256No144: tcell.Color144,
	Color256No145: tcell.Color145,
	Color256No146: tcell.Color146,
	Color256No147: tcell.Color147,
	Color256No148: tcell.Color148,
	Color256No149: tcell.Color149,
	Color256No150: tcell.Color150,
	Color256No151: tcell.Color151,
	Color256No152: tcell.Color152,
	Color256No153: tcell.Color153,
	Color256No154: tcell.Color154,
	Color256No155: tcell.Color155,
	Color256No156: tcell.Color156,
	Color256No157: tcell.Color157,
	Color256No158: tcell.Color158,
	Color256No159: tcell.Color159,
	Color256No160: tcell.Color160,
	Color256No161: tcell.Color161,
	Color256No162: tcell.Color162,
	Color256No163: tcell.Color163,
	Color256No164: tcell.Color164,
	Color256No165: tcell.Color165,
	Color256No166: tcell.Color166,
	Color256No167: tcell.Color167,
	Color256No168: tcell.Color168,
	Color256No169: tcell.Color169,
	Color256No170: tcell.Color170,
	Color256No171: tcell.Color171,
	Color256No172: tcell.Color172,
	Color256No173: tcell.Color173,
	Color256No174: tcell.Color174,
	Color256No175: tcell.Color175,
	Color256No176: tcell.Color176,
	Color256No177: tcell.Color177,
	Color256No178: tcell.Color178,
	Color256No179: tcell.Color179,
	Color256No180: tcell.Color180,
	Color256No181: tcell.Color181,
	Color256No182: tcell.Color182,
	Color256No183: tcell.Color183,
	Color256No184: tcell.Color184,
	Color256No185: tcell.Color185,
	Color256No186: tcell.Color186,
	Color256No187: tcell.Color187,
	Color256No188: tcell.Color188,
	Color256No189: tcell.Color189,
	Color256No190: tcell.Color190,
	Color256No191: tcell.Color191,
	Color256No192: tcell.Color192,
	Color256No193: tcell.Color193,
	Color256No194: tcell.Color194,
	Color256No195: tcell.Color195,
	Color256No196: tcell.Color196,
	Color256No197: tcell.Color197,
	Color256No198: tcell.Color198,
	Color256No199: tcell.Color199,
	Color256No200: tcell.Color200,
	Color256No201: tcell.Color201,
	Color256No202: tcell.Color202,
	Color256No203: tcell.Color203,
	Color256No204: tcell.Color204,
	Color256No205: tcell.Color205,
	Color256No206: tcell.Color206,
	Color256No207: tcell.Color207,
	Color256No208: tcell.Color208,
	Color256No209: tcell.Color209,
	Color256No210: tcell.Color210,
	Color256No211: tcell.Color211,
	Color256No212: tcell.Color212,
	Color256No213: tcell.Color213,
	Color256No214: tcell.Color214,
	Color256No215: tcell.Color215,
	Color256No216: tcell.Color216,
	Color256No217: tcell.Color217,
	Color256No218: tcell.Color218,
	Color256No219: tcell.Color219,
	Color256No220: tcell.Color220,
	Color256No221: tcell.Color221,
	Color256No222: tcell.Color222,
	Color256No223: tcell.Color223,
	Color256No224: tcell.Color224,
	Color256No225: tcell.Color225,
	Color256No226: tcell.Color226,
	Color256No227: tcell.Color227,
	Color256No228: tcell.Color228,
	Color256No229: tcell.Color229,
	Color256No230: tcell.Color230,
	Color256No231: tcell.Color231,
	Color256No232: tcell.Color232,
	Color256No233: tcell.Color233,
	Color256No234: tcell.Color234,
	Color256No235: tcell.Color235,
	Color256No236: tcell.Color236,
	Color256No237: tcell.Color237,
	Color256No238: tcell.Color238,
	Color256No239: tcell.Color239,
	Color256No240: tcell.Color240,
	Color256No241: tcell.Color241,
	Color256No242: tcell.Color242,
	Color256No243: tcell.Color243,
	Color256No244: tcell.Color244,
	Color256No245: tcell.Color245,
	Color256No246: tcell.Color246,
	Color256No247: tcell.Color247,
	Color256No248: tcell.Color248,
	Color256No249: tcell.Color249,
	Color256No250: tcell.Color250,
	Color256No251: tcell.Color251,
	Color256No252: tcell.Color252,
	Color256No253: tcell.Color253,
	Color256No254: tcell.Color254,
	Color256No255: tcell.Color255,
}

func ToTcellStyle(style *ColorStyle) tcell.Style {
	tstyle := tcell.Style{}
	tstyle.Foreground(tcolorMapping[style.fg])
	tstyle.Background(tcolorMapping[style.bg])
	tstyle.Bold(style.bold)
	tstyle.Underline(style.underline)
	tstyle.Italic(style.italic)
	return tstyle
}