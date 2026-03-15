package editor

import (
	"testing"

	"github.com/EME130/lazymd/internal/buffer"
)

func newTestEditor(content string) *EditorModel {
	buf := buffer.New()
	if content != "" {
		buf.InsertString(0, content)
	}
	return New(buf)
}

func TestEditorInitAndInsert(t *testing.T) {
	e := newTestEditor("")
	if e.Mode() != ModeNormal {
		t.Errorf("expected normal mode, got %v", e.Mode())
	}

	// 'i' enters insert mode
	e.HandleKey(CharKey('i'))
	if e.Mode() != ModeInsert {
		t.Errorf("expected insert mode, got %v", e.Mode())
	}

	// Type 'H'
	e.HandleKey(CharKey('H'))
	if e.Buf.Length() != 1 {
		t.Errorf("expected length 1, got %d", e.Buf.Length())
	}

	// Escape returns to normal
	e.HandleKey(SpecialKey(KeyEscape))
	if e.Mode() != ModeNormal {
		t.Errorf("expected normal mode after escape, got %v", e.Mode())
	}
}

func TestCursorMovement(t *testing.T) {
	e := newTestEditor("Hello\nWorld\nTest")

	if e.Row != 0 {
		t.Errorf("expected row 0, got %d", e.Row)
	}

	// j moves down
	e.HandleKey(CharKey('j'))
	if e.Row != 1 {
		t.Errorf("expected row 1, got %d", e.Row)
	}

	// l moves right
	e.HandleKey(CharKey('l'))
	if e.Col != 1 {
		t.Errorf("expected col 1, got %d", e.Col)
	}

	// k moves up
	e.HandleKey(CharKey('k'))
	if e.Row != 0 {
		t.Errorf("expected row 0, got %d", e.Row)
	}
}

func TestCountPrefix(t *testing.T) {
	e := newTestEditor("Line1\nLine2\nLine3\nLine4\nLine5")

	// 3j moves down 3 lines
	e.HandleKey(CharKey('3'))
	e.HandleKey(CharKey('j'))
	if e.Row != 3 {
		t.Errorf("expected row 3 after 3j, got %d", e.Row)
	}
}

func TestDeleteLine(t *testing.T) {
	e := newTestEditor("Line1\nLine2\nLine3\n")

	// dd deletes current line
	e.HandleKey(CharKey('d'))
	e.HandleKey(CharKey('d'))
	content := e.Buf.Content()
	if content != "Line2\nLine3\n" {
		t.Errorf("expected Line2\\nLine3, got %q", content)
	}
}

func TestInsertMode(t *testing.T) {
	e := newTestEditor("AB")

	// Move to position 1, enter insert
	e.HandleKey(CharKey('l'))
	e.HandleKey(CharKey('i'))

	// Type 'X'
	e.HandleKey(CharKey('X'))
	content := e.Buf.Content()
	if content != "AXB" {
		t.Errorf("expected AXB, got %q", content)
	}
}

func TestCommandMode(t *testing.T) {
	e := newTestEditor("Hello")

	// Enter command mode
	e.HandleKey(CharKey(':'))
	if e.Mode() != ModeCommand {
		t.Errorf("expected command mode")
	}

	// Type "q!" and enter
	for _, c := range "q!" {
		e.HandleKey(CharKey(c))
	}
	e.HandleKey(SpecialKey(KeyEnter))

	if !e.ShouldQuit {
		t.Error("expected should_quit after :q!")
	}
}

func TestWordMotions(t *testing.T) {
	e := newTestEditor("hello world foo")

	// w moves to next word
	e.HandleKey(CharKey('w'))
	if e.Col != 6 {
		t.Errorf("expected col 6 after w, got %d", e.Col)
	}

	// b moves back
	e.HandleKey(CharKey('b'))
	if e.Col != 0 {
		t.Errorf("expected col 0 after b, got %d", e.Col)
	}
}

func TestOpenLineBelow(t *testing.T) {
	e := newTestEditor("Line1\nLine2")

	// 'o' opens line below
	e.HandleKey(CharKey('o'))
	if e.Mode() != ModeInsert {
		t.Error("expected insert mode after o")
	}
	if e.Row != 1 {
		t.Errorf("expected row 1, got %d", e.Row)
	}
}

func TestPluginEditorInterface(t *testing.T) {
	e := newTestEditor("Test content")
	e.File = "test.md"

	// Verify PluginEditor interface
	var pe PluginEditor = e
	if pe.CursorRow() != 0 {
		t.Error("expected row 0")
	}
	if pe.FilePath() != "test.md" {
		t.Error("expected test.md")
	}
	pe.SetStatus("hello", false)
	if e.Status.Text != "hello" {
		t.Error("expected status 'hello'")
	}
}

func TestModeString(t *testing.T) {
	if ModeNormal.String() != "NORMAL" {
		t.Error("expected NORMAL")
	}
	if ModeInsert.String() != "INSERT" {
		t.Error("expected INSERT")
	}
	if ModeCommand.String() != "COMMAND" {
		t.Error("expected COMMAND")
	}
}

func TestLoadContent(t *testing.T) {
	buf := buffer.New()
	ed := New(buf)

	ed.LoadContent("[help:test]", "# Help\n\nSome content here")

	if ed.File != "[help:test]" {
		t.Errorf("expected file '[help:test]', got %q", ed.File)
	}
	if ed.Buf.LineCount() != 3 {
		t.Fatalf("expected 3 lines, got %d", ed.Buf.LineCount())
	}
	if ed.Buf.Line(0) != "# Help" {
		t.Errorf("expected first line '# Help', got %q", ed.Buf.Line(0))
	}
	if ed.Row != 0 || ed.Col != 0 {
		t.Errorf("expected cursor at 0,0 — got %d,%d", ed.Row, ed.Col)
	}
	if ed.ScrollRow != 0 {
		t.Errorf("expected scroll at 0, got %d", ed.ScrollRow)
	}
	if ed.Buf.IsDirty() {
		t.Fatal("expected buffer to not be dirty after LoadContent")
	}
}

func TestHelpCommandLoadsContent(t *testing.T) {
	buf := buffer.New()
	ed := New(buf)
	ed.CmdBuf = "help brain"
	ed.executeCommand()

	if ed.File != "[help:brain]" {
		t.Errorf("expected file '[help:brain]', got %q", ed.File)
	}
	if ed.Buf.LineCount() < 2 {
		t.Error("expected help content to be loaded")
	}
}

func TestHelpCommandDefaultsToOverview(t *testing.T) {
	buf := buffer.New()
	ed := New(buf)
	ed.CmdBuf = "help"
	ed.executeCommand()

	if ed.File != "[help:overview]" {
		t.Errorf("expected file '[help:overview]', got %q", ed.File)
	}
}

func TestHelpCommandUnknownTopic(t *testing.T) {
	buf := buffer.New()
	ed := New(buf)
	ed.CmdBuf = "help nonexistent"
	ed.executeCommand()

	if !ed.Status.IsError {
		t.Error("expected error status for unknown topic")
	}
}

func TestHelpCommandBlockedByDirtyBuffer(t *testing.T) {
	buf := buffer.New()
	buf.InsertString(0, "dirty content")
	ed := New(buf)
	ed.CmdBuf = "help"
	ed.executeCommand()

	if !ed.Status.IsError {
		t.Error("expected error status for dirty buffer")
	}
	if ed.Buf.Line(0) != "dirty content" {
		t.Error("expected buffer content to be unchanged")
	}
}
