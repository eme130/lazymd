package plugins

import (
	"testing"

	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/editor"
)

// mockEditor implements editor.PluginEditor for testing.
type mockEditor struct {
	buf    *buffer.Buffer
	row    int
	col    int
	file   string
	status string
	isErr  bool
}

func newMockEditor(content string) *mockEditor {
	buf := buffer.New()
	if content != "" {
		buf.InsertString(0, content)
	}
	return &mockEditor{buf: buf, file: "test.md"}
}

func (m *mockEditor) Buffer() *buffer.Buffer       { return m.buf }
func (m *mockEditor) CursorRow() int               { return m.row }
func (m *mockEditor) CursorCol() int               { return m.col }
func (m *mockEditor) FilePath() string             { return m.file }
func (m *mockEditor) EditorMode() editor.Mode      { return editor.ModeNormal }
func (m *mockEditor) SetStatus(msg string, e bool) { m.status = msg; m.isErr = e }
func (m *mockEditor) SetCursorRow(row int)         { m.row = row }
func (m *mockEditor) SetCursorCol(col int)         { m.col = col }

func TestAllPluginsEmpty(t *testing.T) {
	all := AllPlugins()
	if len(all) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(all))
	}
}

func TestPluginManagerBasics(t *testing.T) {
	pm := NewManager()

	if pm.PluginCount() != 0 {
		t.Errorf("expected 0 plugins, got %d", pm.PluginCount())
	}

	ed := newMockEditor("")
	if pm.ExecuteCommand("nonexistent", ed, "") {
		t.Error("nonexistent command should not be found")
	}
}

func TestRegisterAll(t *testing.T) {
	pm := NewManager()
	ed := newMockEditor("")
	RegisterAll(pm, ed)

	if pm.PluginCount() != 0 {
		t.Errorf("expected 0 plugins after RegisterAll, got %d", pm.PluginCount())
	}
}
