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
	if document.CurrentChar() != "e" {
		t.Fatalf("currentChar want \"e\", but got=%q", document.CurrentChar())
	}
	if document.CharBeforeCursor() != "n" {
		t.Fatalf("CharBeforeCursor want \"n\", but got=%q", document.CharBeforeCursor())
	}
}

func TestTextBeforeCursor(t *testing.T) {
	document := prepareDocument()
	if document.TextBeforeCursor() != "line 1\nlin" {
		t.Fatalf("got=%q", document.TextBeforeCursor())
	}
}

func TestTextAfterCursor(t *testing.T) {
	document := prepareDocument()
	if document.TextAfterCursor() != "e 2\n"+"line 3\n"+"line 4\n" {
		t.Fatalf("got=%q", document.TextAfterCursor())
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
	got := document.LineCount()
	if got != want {
		t.Errorf("want=%d, but got=%d", want, got)
	}
}

func TestCurrentLineBeforeCursor(t *testing.T) {
	document := prepareDocument()
	want := "lin"
	got := document.CurrentLineBeforeCursor()
	testStringEqual(t, want, got)
}

func TestCurrentLineAfterCursor(t *testing.T) {
	document := prepareDocument()
	want := "e 2"
	got := document.CurrentLineAfterCursor()
	testStringEqual(t, want, got)
}

func TestCurrentLine(t *testing.T) {
	document := prepareDocument()
	want := "line 2"
	got := document.CurrentLine()
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

func TestFindStartOfPreviousWord(t *testing.T) {
	var got int
	var doc *Document

	for i := 0; i < 10; i++ {
		doc = NewDocument("package  ", i)
		got = doc.findStartOfPreviousWord()
		testIntEqual(t, 0-i, got)
	}

	doc = NewDocument("hello world", 11)
	got = doc.findStartOfPreviousWord()
	testIntEqual(t, -5, got)

	doc = NewDocument("hello 中文", 8)
	got = doc.findStartOfPreviousWord()
	testIntEqual(t, -2, got)
}

func TestFindNextWordBeginning(t *testing.T) {
	var got int
	var want int
	var doc *Document

	doc = NewDocument("package", 0)
	want = 0
	got = doc.findNextWordBeginning()
	testIntEqual(t, want, got)

	doc = NewDocument(" package", 0)
	want = 1
	got = doc.findNextWordBeginning()
	testIntEqual(t, want, got)

	doc = NewDocument(" package ", 0)
	want = 1
	got = doc.findNextWordBeginning()
	testIntEqual(t, want, got)

	doc = NewDocument("e a", 0)
	got = doc.findNextWordBeginning()
	testIntEqual(t, 2, got)

	doc = NewDocument("e ab", 0)
	got = doc.findNextWordBeginning()
	testIntEqual(t, 2, got)

	doc = NewDocument("e ab      ", 0)
	got = doc.findNextWordBeginning()
	testIntEqual(t, 2, got)
}

func TestFindNextWordEnding(t *testing.T) {
	var got int
	var want int
	var doc *Document

	doc = NewDocument("package", 0)
	want = 7
	got = doc.findNextWordEnding(true)
	testIntEqual(t, want, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, want, got)

	doc = NewDocument(" package", 0)
	want = 8
	got = doc.findNextWordEnding(true)
	testIntEqual(t, want, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, want, got)

	doc = NewDocument(" package ", 0)
	want = 8
	got = doc.findNextWordEnding(true)
	testIntEqual(t, want, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, want, got)

	doc = NewDocument("e a", 0)
	got = doc.findNextWordEnding(true)
	testIntEqual(t, 1, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, 3, got)

	doc = NewDocument("e ab", 0)
	got = doc.findNextWordEnding(true)
	testIntEqual(t, 1, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, 4, got)

	doc = NewDocument("e ab      ", 0)
	got = doc.findNextWordEnding(true)
	testIntEqual(t, 1, got)
	got = doc.findNextWordEnding(false)
	testIntEqual(t, 4, got)
}
