package startprompt

import (
	"testing"
)

func prepareDocument() *Document {
	text := "line 1\n" + "line 2\n" + "line 3\n" + "line 4\n"
	return NewDocument(
		text,
		len("line 1\n"+"lin"),
	)
}

func testStringEqual(t *testing.T, want string, got string) {
	t.Helper()
	if got != want {
		t.Errorf("want=%q, but got=%q", want, got)
	}
}

func testIntEqual(t *testing.T, want int, got int) {
	t.Helper()
	if got != want {
		t.Errorf("want=%d, but got=%d", want, got)
	}
}

func testBoolEqual(t *testing.T, want bool, got bool) {
	t.Helper()
	if got != want {
		t.Errorf("want=%v, but got=%v", want, got)
	}
}

func TestCurrentChar(t *testing.T) {
	document := prepareDocument()
	if document.currentChar() != "e" {
		t.Fatalf("currentChar want \"e\", but got=%q", document.currentChar())
	}
	if document.charBeforeCursor() != "n" {
		t.Fatalf("charBeforeCursor want \"n\", but got=%q", document.charBeforeCursor())
	}
}

func TestTextBeforeCursor(t *testing.T) {
	document := prepareDocument()
	if document.textBeforeCursor() != "line 1\nlin" {
		t.Fatalf("got=%q", document.textBeforeCursor())
	}
}

func TestTextAfterCursor(t *testing.T) {
	document := prepareDocument()
	if document.textAfterCursor() != "e 2\n"+"line 3\n"+"line 4\n" {
		t.Fatalf("got=%q", document.textAfterCursor())
	}
}

func TestLines(t *testing.T) {
	document := prepareDocument()
	want := []string{"line 1", "line 2", "line 3", "line 4", ""}
	got := document.lines()
	for i, s := range want {
		if got[i] != s {
			t.Fatalf("%d want=%q, got=%q", i, s, got[i])
		}
	}
}

func TestLineCount(t *testing.T) {
	document := prepareDocument()
	want := 5
	got := document.lineCount()
	if got != want {
		t.Errorf("want=%d, but got=%d", want, got)
	}
}

func TestCurrentLineBeforeCursor(t *testing.T) {
	document := prepareDocument()
	want := "lin"
	got := document.currentLineBeforeCursor()
	testStringEqual(t, want, got)
}

func TestCurrentLineAfterCursor(t *testing.T) {
	document := prepareDocument()
	want := "e 2"
	got := document.currentLineAfterCursor()
	testStringEqual(t, want, got)
}

func TestCurrentLine(t *testing.T) {
	document := prepareDocument()
	want := "line 2"
	got := document.currentLine()
	testStringEqual(t, want, got)
}

func TestCursorPositionRowAndCol(t *testing.T) {
	document := prepareDocument()
	testIntEqual(t, 1, document.CursorPositionRow())
	testIntEqual(t, 3, document.CursorPositionCol())

	document = NewDocument("", 0)
	testIntEqual(t, 0, document.CursorPositionRow())
	testIntEqual(t, 0, document.CursorPositionCol())
}

func TestTranslateIndexToRowCol(t *testing.T) {
	document := prepareDocument()
	row, col := document.translateIndexToRowCol(len("line 1\nline 2\nlin"))
	testIntEqual(t, 2, row)
	testIntEqual(t, 3, col)

	row, col = document.translateIndexToRowCol(0)
	testIntEqual(t, 0, row)
	testIntEqual(t, 0, col)
}

func TestTranslateRowColToIndex(t *testing.T) {
	document := prepareDocument()
	var index int
	index = document.translateRowColToIndex(0, 0)
	testIntEqual(t, 0, index)
	index = document.translateRowColToIndex(0, 3)
	testIntEqual(t, 3, index)
	index = document.translateRowColToIndex(0, 6)
	testIntEqual(t, 6, index)
	index = document.translateRowColToIndex(0, 7)
	testIntEqual(t, 6, index)
	index = document.translateRowColToIndex(1, 7)
	testIntEqual(t, 13, index)
}

func TestIsCursorAtTheEnd(t *testing.T) {
	got := NewDocument("hello", 5).isCursorAtTheEnd()
	testBoolEqual(t, true, got)
	got = NewDocument("hello", 4).isCursorAtTheEnd()
	testBoolEqual(t, false, got)
}
