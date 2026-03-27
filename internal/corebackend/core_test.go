package corebackend

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/pluginapi"
)

type mockEditor struct {
	buf      *buffer.Buffer
	row, col int
	filePath string
}

func newMockEditor() *mockEditor {
	return &mockEditor{buf: buffer.New(), filePath: "test.md"}
}

func (m *mockEditor) Buffer() *buffer.Buffer             { return m.buf }
func (m *mockEditor) CursorRow() int                     { return m.row }
func (m *mockEditor) CursorCol() int                     { return m.col }
func (m *mockEditor) FilePath() string                   { return m.filePath }
func (m *mockEditor) Mode() string                       { return "NORMAL" }
func (m *mockEditor) SetStatus(_ string, _ bool)         {}
func (m *mockEditor) SetCursorRow(r int)                 { m.row = r }
func (m *mockEditor) SetCursorCol(c int)                 { m.col = c }
func (m *mockEditor) LineCount() int                     { return m.buf.LineCount() }
func (m *mockEditor) Line(n int) string                  { return m.buf.Line(n) }
func (m *mockEditor) Content() string                    { return m.buf.Content() }
func (m *mockEditor) InsertAt(row, col int, text string) {
	off := m.buf.PosToOffset(row, col)
	m.buf.InsertString(off, text)
}
func (m *mockEditor) DeleteRange(startRow, startCol, endRow, endCol int) {
	s := m.buf.PosToOffset(startRow, startCol)
	e := m.buf.PosToOffset(endRow, endCol)
	if e > s {
		m.buf.DeleteRange(s, e-s)
	}
}
func (m *mockEditor) DeleteLines(_, _ int) {}
func (m *mockEditor) OpenFile(path string) error {
	if err := m.buf.LoadFile(path); err != nil {
		return err
	}
	m.filePath = path
	return nil
}
func (m *mockEditor) SaveFile() error {
	return m.buf.SaveFile(m.filePath)
}

func newTestCtx(ed *mockEditor) (*pluginapi.BackendContext, *[]*pluginapi.Event) {
	var events []*pluginapi.Event
	return &pluginapi.BackendContext{
		Editor:    ed,
		Broadcast: func(evt *pluginapi.Event) { events = append(events, evt) },
		Log:       func(msg string) {},
	}, &events
}

func TestOpenFileBroadcastsEvents(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "hello.md")
	os.WriteFile(fpath, []byte("# Hello"), 0644)

	ed := newMockEditor()
	ctx, events := newTestCtx(ed)
	p := New()

	op := pluginapi.NewOperation(pluginapi.OpOpenFile, "test")
	op.Data["path"] = fpath
	p.OnOperation(ctx, op)

	if len(*events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(*events))
	}
	if (*events)[0].Type != pluginapi.EventFileOpened {
		t.Errorf("expected EventFileOpened, got %s", (*events)[0].Type)
	}
	if (*events)[1].Type != pluginapi.EventBufferChanged {
		t.Errorf("expected EventBufferChanged, got %s", (*events)[1].Type)
	}
	content, _ := (*events)[1].Data["content"].(string)
	if content != "# Hello" {
		t.Errorf("expected content '# Hello', got %q", content)
	}
}

func TestInsertTextBroadcasts(t *testing.T) {
	ed := newMockEditor()
	ctx, events := newTestCtx(ed)
	p := New()

	op := pluginapi.NewOperation(pluginapi.OpInsertText, "gui")
	op.Data["row"] = 0
	op.Data["col"] = 0
	op.Data["text"] = "hello"
	p.OnOperation(ctx, op)

	if len(*events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(*events))
	}
	if (*events)[0].Type != pluginapi.EventBufferChanged {
		t.Errorf("expected EventBufferChanged, got %s", (*events)[0].Type)
	}
	if ed.Content() != "hello" {
		t.Errorf("expected buffer content 'hello', got %q", ed.Content())
	}
}

func TestSetCursorBroadcasts(t *testing.T) {
	ed := newMockEditor()
	ctx, events := newTestCtx(ed)
	p := New()

	op := pluginapi.NewOperation(pluginapi.OpSetCursor, "gui")
	op.Data["row"] = 5
	op.Data["col"] = 3
	p.OnOperation(ctx, op)

	if len(*events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(*events))
	}
	if (*events)[0].Type != pluginapi.EventCursorMoved {
		t.Errorf("expected EventCursorMoved, got %s", (*events)[0].Type)
	}
	if ed.row != 5 || ed.col != 3 {
		t.Errorf("expected cursor at (5,3), got (%d,%d)", ed.row, ed.col)
	}
}
