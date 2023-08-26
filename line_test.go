package startprompt

import "testing"

func newTestLine() *Line {
	return newLine(
		newBaseCode, newBasePrompt, NewMemHistory(), NewNoopCallbacks(), false)
}

func TestLineInitial(t *testing.T) {
	cli := newTestLine()
	testStringEqual(t, cli.text(), "")
	testIntEqual(t, cli.GetCursorPosition(), 0)
}

func TestLine_InsertText(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("some_text"), true)
	testStringEqual(t, cli.text(), "some_text")
	testIntEqual(t, cli.GetCursorPosition(), len("some_text"))
}

func TestLine_Movement(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("some_text"), true)
	cli.CursorLeft()
	cli.CursorLeft()
	cli.CursorLeft()
	cli.CursorRight()
	cli.InsertText([]rune{'A'}, true)

	testStringEqual(t, cli.text(), "some_teAxt")
	testIntEqual(t, cli.GetCursorPosition(), len("some_teA"))
}

func TestLine_Backspace(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("some_text"), true)
	cli.CursorLeft()
	cli.CursorLeft()
	cli.DeleteCharacterBeforeCursor(1)

	testStringEqual(t, cli.text(), "some_txt")
	testIntEqual(t, cli.GetCursorPosition(), len("some_t"))
}

func TestLine_CursorUp(t *testing.T) {
	//    向上移动到更长的行
	cli := newTestLine()
	cli.InsertText([]rune("long line1\nline2"), true)
	cli.CursorUp()

	testIntEqual(t, 5, cli.GetCursorPosition())
	testIntEqual(t, 5, cli.Document().CursorPosition())

	cli.CursorUp()
	testIntEqual(t, 5, cli.GetCursorPosition())
	testIntEqual(t, 5, cli.Document().CursorPosition())

	//    向上移动到更短的行
	cli.reset()
	cli.InsertText([]rune("line1\nlong line2"), true)

	cli.CursorUp()
	testIntEqual(t, 5, cli.GetCursorPosition())
	testIntEqual(t, 5, cli.Document().CursorPosition())
}

func TestLine_CursorDown(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("line1\nline2"), true)
	cli.SetCursorPosition(3)

	cli.CursorDown()
	testIntEqual(t, len("line1\nlin"), cli.Document().CursorPosition())

	cli.reset()
	cli.InsertText([]rune("long line1\na\nb"), true)
	cli.SetCursorPosition(3)

	cli.CursorDown()
	testIntEqual(t, len("long line1\na"), cli.Document().CursorPosition())
}

func TestLine_JoinNextLine(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("line1\nline2\nline3"), true)
	cli.CursorUp()
	cli.JoinNextLine()

	testStringEqual(t, "line1\nline2line3", cli.text())

	cli.reset()
	cli.InsertText([]rune("line1"), true)
	cli.SetCursorPosition(0)
	cli.JoinNextLine()
	testStringEqual(t, "line1", cli.text())
}

func TestLine_Newline(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("hello world"), true)
	cli.Newline()
	testStringEqual(t, "hello world\n", cli.text())
}

func TestLine_SwapCharactersBeforeCursor(t *testing.T) {
	cli := newTestLine()
	cli.InsertText([]rune("hello world"), true)
	cli.CursorLeft()
	cli.CursorLeft()
	cli.SwapCharactersBeforeCursor()
	testStringEqual(t, "hello wrold", cli.text())
}
