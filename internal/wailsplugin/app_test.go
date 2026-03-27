package wailsplugin

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/pluginapi"
)

type mockEditorAPI struct {
	buf      *buffer.Buffer
	row, col int
	filePath string
}

func newMockEditor(content string) *mockEditorAPI {
	b := buffer.New()
	b.InsertString(0, content)
	return &mockEditorAPI{buf: b, filePath: "test.md"}
}

func (m *mockEditorAPI) Buffer() *buffer.Buffer          { return m.buf }
func (m *mockEditorAPI) CursorRow() int                  { return m.row }
func (m *mockEditorAPI) CursorCol() int                  { return m.col }
func (m *mockEditorAPI) FilePath() string                { return m.filePath }
func (m *mockEditorAPI) Mode() string                    { return "NORMAL" }
func (m *mockEditorAPI) SetStatus(_ string, _ bool)      {}
func (m *mockEditorAPI) SetCursorRow(r int)              { m.row = r }
func (m *mockEditorAPI) SetCursorCol(c int)              { m.col = c }
func (m *mockEditorAPI) LineCount() int                  { return m.buf.LineCount() }
func (m *mockEditorAPI) Line(n int) string               { return m.buf.Line(n) }
func (m *mockEditorAPI) Content() string                 { return m.buf.Content() }
func (m *mockEditorAPI) InsertAt(row, col int, text string) {
	// Simple mock: just insert at byte position
	pos := 0
	for i := 0; i < row && i < m.buf.LineCount(); i++ {
		pos += len(m.buf.Line(i)) + 1 // +1 for newline
	}
	pos += col
	m.buf.InsertString(pos, text)
}
func (m *mockEditorAPI) OpenFile(_ string) error { return nil }
func (m *mockEditorAPI) SaveFile() error          { return nil }
func (m *mockEditorAPI) DeleteRange(_, _, _, _ int) {}
func (m *mockEditorAPI) DeleteLines(start, end int) {
	// Simple mock: delete line range
	if start >= m.buf.LineCount() {
		return
	}
	startPos := 0
	for i := 0; i < start; i++ {
		startPos += len(m.buf.Line(i)) + 1
	}
	endPos := startPos
	for i := start; i < end && i < m.buf.LineCount(); i++ {
		endPos += len(m.buf.Line(i)) + 1
	}
	if endPos > startPos {
		m.buf.DeleteRange(startPos, endPos-startPos)
	}
}

type mockBrainAPI struct {
	backlinks map[string][]string
}

func (m *mockBrainAPI) NodeCount() int                          { return 0 }
func (m *mockBrainAPI) Nodes() []pluginapi.NodeInfo             { return nil }
func (m *mockBrainAPI) GetBacklinks(name string) []string       { return m.backlinks[name] }
func (m *mockBrainAPI) GetNeighbors(name string) []string       { return nil }
func (m *mockBrainAPI) FindPath(from, to string) ([]string, bool) { return nil, false }
func (m *mockBrainAPI) GetOrphans() []string                    { return nil }

func TestAppGetContent(t *testing.T) {
	editor := newMockEditor("hello world")
	app := newTestApp(editor, nil, nil)
	content := app.GetContent()
	if content != "hello world" {
		t.Errorf("expected 'hello world', got %q", content)
	}
}

func TestAppGetCursorPos(t *testing.T) {
	editor := newMockEditor("")
	editor.row = 5
	editor.col = 10
	app := newTestApp(editor, nil, nil)
	pos := app.GetCursorPos()
	if pos.Row != 5 || pos.Col != 10 {
		t.Errorf("expected {5, 10}, got {%d, %d}", pos.Row, pos.Col)
	}
}

func TestAppGetCurrentFile(t *testing.T) {
	editor := newMockEditor("")
	editor.filePath = "notes/test.md"
	app := newTestApp(editor, nil, nil)
	if got := app.GetCurrentFile(); got != "notes/test.md" {
		t.Errorf("expected 'notes/test.md', got %q", got)
	}
}

func TestAppInsertTextEmitsOp(t *testing.T) {
	editor := newMockEditor("")
	var emitted *pluginapi.Operation
	app := newTestApp(editor, nil, nil)
	app.emit = func(op *pluginapi.Operation) { emitted = op }
	app.InsertText(0, 0, "hello")
	if emitted == nil {
		t.Fatal("expected operation to be emitted")
	}
	if emitted.Type != pluginapi.OpInsertText {
		t.Errorf("expected OpInsertText, got %s", emitted.Type)
	}
	if emitted.Data["text"] != "hello" {
		t.Errorf("expected text 'hello', got %v", emitted.Data["text"])
	}
	if emitted.Origin != "wails-gui" {
		t.Errorf("expected origin 'wails-gui', got %q", emitted.Origin)
	}
}

func TestAppDeleteRangeEmitsOp(t *testing.T) {
	editor := newMockEditor("")
	var emitted *pluginapi.Operation
	app := newTestApp(editor, nil, nil)
	app.emit = func(op *pluginapi.Operation) { emitted = op }
	app.DeleteRange(1, 0, 1, 5)
	if emitted == nil {
		t.Fatal("expected operation to be emitted")
	}
	if emitted.Type != pluginapi.OpDeleteRange {
		t.Errorf("expected OpDeleteRange, got %s", emitted.Type)
	}
}

func TestAppOpenFileEmitsOp(t *testing.T) {
	editor := newMockEditor("")
	var emitted *pluginapi.Operation
	app := newTestApp(editor, nil, nil)
	app.emit = func(op *pluginapi.Operation) { emitted = op }
	app.OpenFile("notes/readme.md")
	if emitted == nil {
		t.Fatal("expected operation to be emitted")
	}
	if emitted.Type != pluginapi.OpOpenFile {
		t.Errorf("expected OpOpenFile, got %s", emitted.Type)
	}
	if emitted.Data["path"] != "notes/readme.md" {
		t.Errorf("expected path 'notes/readme.md', got %v", emitted.Data["path"])
	}
}

func TestAppListFiles(t *testing.T) {
	app := newTestApp(newMockEditor(""), nil, nil)
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "note.md"), []byte("# Test"), 0644)
	os.WriteFile(filepath.Join(dir, "readme.rndm"), []byte("hi"), 0644)
	os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("no"), 0644)
	files := app.ListFiles(dir)
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestAppGetWordCount(t *testing.T) {
	editor := newMockEditor("hello world foo")
	app := newTestApp(editor, nil, nil)
	if got := app.GetWordCount(); got != 3 {
		t.Errorf("expected 3 words, got %d", got)
	}
}

func TestAppGetWordCountEmpty(t *testing.T) {
	editor := newMockEditor("")
	app := newTestApp(editor, nil, nil)
	if got := app.GetWordCount(); got != 0 {
		t.Errorf("expected 0 words, got %d", got)
	}
}

func TestAppGetCharCount(t *testing.T) {
	editor := newMockEditor("hello world")
	app := newTestApp(editor, nil, nil)
	if got := app.GetCharCount(); got != 11 {
		t.Errorf("expected 11 chars, got %d", got)
	}
}

func TestAppGetBacklinkCount(t *testing.T) {
	app := newTestApp(newMockEditor(""), nil, &mockBrainAPI{
		backlinks: map[string][]string{"test": {"a.md", "b.md"}},
	})
	app.editor.(*mockEditorAPI).filePath = "test.md"
	if got := app.GetBacklinkCount(); got != 2 {
		t.Errorf("expected 2 backlinks, got %d", got)
	}
}

func TestAppGetBacklinkCountNoBrain(t *testing.T) {
	app := newTestApp(newMockEditor(""), nil, nil)
	if got := app.GetBacklinkCount(); got != 0 {
		t.Errorf("expected 0 backlinks, got %d", got)
	}
}

func newTestApp(editor pluginapi.EditorAPI, nav pluginapi.NavAPI, brain pluginapi.BrainAPI) *App {
	return &App{
		editor: editor,
		nav:    nav,
		brain:  brain,
		emit:   func(op *pluginapi.Operation) {},
	}
}
