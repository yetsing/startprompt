package startprompt

import (
	"testing"

	"github.com/yetsing/startprompt/keys"
)

type tKey struct {
	event keys.Event
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

func (h *testHandler) Handle(event keys.Event, a ...rune) {
	k := tKey{
		event: event,
		data:  string(a),
	}
	h.keys = append(h.keys, k)
}

func TestInputStreamControlKeys(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x01\x02\x10")

	testIntEqual(t, 3, len(handler.keys))
	testKeyEventEqual(t, keys.CtrlA, handler.keys[0].event)
	testKeyEventEqual(t, keys.CtrlB, handler.keys[1].event)
	testKeyEventEqual(t, keys.CtrlP, handler.keys[2].event)
	testStringEqual(t, "\x01", handler.keys[0].data)
	testStringEqual(t, "\x02", handler.keys[1].data)
	testStringEqual(t, "\x10", handler.keys[2].data)
}

func TestInputStreamArrows(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x1b[A\x1b[B\x1b[C\x1b[D")

	testIntEqual(t, 4, len(handler.keys))
	testKeyEventEqual(t, keys.ArrowUp, handler.keys[0].event)
	testKeyEventEqual(t, keys.ArrowDown, handler.keys[1].event)
	testKeyEventEqual(t, keys.ArrowRight, handler.keys[2].event)
	testKeyEventEqual(t, keys.ArrowLeft, handler.keys[3].event)
	testStringEqual(t, "\x1b[A", handler.keys[0].data)
	testStringEqual(t, "\x1b[B", handler.keys[1].data)
	testStringEqual(t, "\x1b[C", handler.keys[2].data)
	testStringEqual(t, "\x1b[D", handler.keys[3].data)
}

func TestInputStreamEscape(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x1bhello")

	testIntEqual(t, 1+len("hello"), len(handler.keys))
	testKeyEventEqual(t, keys.EscapeAction, handler.keys[0].event)
	testKeyEventEqual(t, keys.InsertChar, handler.keys[1].event)
	testStringEqual(t, "\x1b", handler.keys[0].data)
	testStringEqual(t, "h", handler.keys[1].data)
}

func TestInputStreamMetaArrows(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x1b\x1b[D")

	testIntEqual(t, 2, len(handler.keys))
	testKeyEventEqual(t, keys.EscapeAction, handler.keys[0].event)
	testKeyEventEqual(t, keys.ArrowLeft, handler.keys[1].event)
}

func TestInputStreamControlSquareClose(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x1dC")

	testIntEqual(t, 2, len(handler.keys))
	testKeyEventEqual(t, keys.CtrlSquareClose, handler.keys[0].event)
	testKeyEventEqual(t, keys.InsertChar, handler.keys[1].event)
	testStringEqual(t, "C", handler.keys[1].data)
}

func TestInputStreamInvalid(t *testing.T) {
	handler := newTestHandler()
	stream := NewInputStream(handler)
	stream.FeedData("\x1b[*")

	testIntEqual(t, 3, len(handler.keys))
	testKeyEventEqual(t, keys.EscapeAction, handler.keys[0].event)
	testKeyEventEqual(t, keys.InsertChar, handler.keys[1].event)
	testKeyEventEqual(t, keys.InsertChar, handler.keys[2].event)
	testStringEqual(t, "[", handler.keys[1].data)
	testStringEqual(t, "*", handler.keys[2].data)

}
