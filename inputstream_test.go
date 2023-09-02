package startprompt

import (
	"testing"
)

type tKey struct {
	event EventType
	data  string
}

type testHandler struct {
	BaseHandler
	keys []tKey
}

func newTestHandler() *testHandler {
	return &testHandler{
		BaseHandler: BaseHandler{},
		keys:        nil,
	}
}

func (h *testHandler) Handle(event EventType, a ...rune) {
	k := tKey{
		event: event,
		data:  string(a),
	}
	h.keys = append(h.keys, k)
}

func TestInputStreamControlKeys(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x01\x02\x10")

	testIntEqual(t, 3, len(handler.keys))
	testKeyEventEqual(t, EventTypeCtrlA, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeCtrlB, handler.keys[1].event)
	testKeyEventEqual(t, EventTypeCtrlP, handler.keys[2].event)
	testStringEqual(t, "\x01", handler.keys[0].data)
	testStringEqual(t, "\x02", handler.keys[1].data)
	testStringEqual(t, "\x10", handler.keys[2].data)
}

func TestInputStreamArrows(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x1b[A\x1b[B\x1b[C\x1b[D")

	testIntEqual(t, 4, len(handler.keys))
	testKeyEventEqual(t, EventTypeArrowUp, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeArrowDown, handler.keys[1].event)
	testKeyEventEqual(t, EventTypeArrowRight, handler.keys[2].event)
	testKeyEventEqual(t, EventTypeArrowLeft, handler.keys[3].event)
	testStringEqual(t, "\x1b[A", handler.keys[0].data)
	testStringEqual(t, "\x1b[B", handler.keys[1].data)
	testStringEqual(t, "\x1b[C", handler.keys[2].data)
	testStringEqual(t, "\x1b[D", handler.keys[3].data)
}

func TestInputStreamEscape(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x1bhello")

	testIntEqual(t, 1+len("hello"), len(handler.keys))
	testKeyEventEqual(t, EventTypeEscape, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeInsertChar, handler.keys[1].event)
	testStringEqual(t, "\x1b", handler.keys[0].data)
	testStringEqual(t, "h", handler.keys[1].data)
}

func TestInputStreamMetaArrows(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x1b\x1b[D")

	testIntEqual(t, 2, len(handler.keys))
	testKeyEventEqual(t, EventTypeEscape, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeArrowLeft, handler.keys[1].event)
}

func TestInputStreamControlSquareClose(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x1dC")

	testIntEqual(t, 2, len(handler.keys))
	testKeyEventEqual(t, EventTypeCtrlSquareClose, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeInsertChar, handler.keys[1].event)
	testStringEqual(t, "C", handler.keys[1].data)
}

func TestInputStreamInvalid(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler, nil)
	stream.FeedData("\x1b[*")

	testIntEqual(t, 3, len(handler.keys))
	testKeyEventEqual(t, EventTypeEscape, handler.keys[0].event)
	testKeyEventEqual(t, EventTypeInsertChar, handler.keys[1].event)
	testKeyEventEqual(t, EventTypeInsertChar, handler.keys[2].event)
	testStringEqual(t, "[", handler.keys[1].data)
	testStringEqual(t, "*", handler.keys[2].data)

}
