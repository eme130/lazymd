package pluginadapter

import (
	"path/filepath"

	"github.com/EME130/lazymd/internal/buffer"
)

// EditorAdapter wraps a buffer.Buffer as a pluginapi.EditorAPI.
type EditorAdapter struct {
	Buf      *buffer.Buffer
	Row      int
	Col      int
	File     string
	FileType string
	Status   StatusMsg
}

type StatusMsg struct {
	Text    string
	IsError bool
}

func (a *EditorAdapter) Buffer() *buffer.Buffer { return a.Buf }
func (a *EditorAdapter) CursorRow() int         { return a.Row }
func (a *EditorAdapter) CursorCol() int         { return a.Col }
func (a *EditorAdapter) FilePath() string       { return a.File }
func (a *EditorAdapter) Mode() string           { return "NORMAL" }

func (a *EditorAdapter) SetStatus(msg string, isError bool) {
	a.Status = StatusMsg{Text: msg, IsError: isError}
}

func (a *EditorAdapter) SetCursorRow(row int) {
	if row >= 0 && row < a.Buf.LineCount() {
		a.Row = row
	}
}

func (a *EditorAdapter) SetCursorCol(col int) {
	if col >= 0 {
		a.Col = col
	}
}

func (a *EditorAdapter) OpenFile(path string) error {
	if err := a.Buf.LoadFile(path); err != nil {
		return err
	}
	a.File = path
	if filepath.Ext(path) == ".tex" {
		a.FileType = "latex"
	} else {
		a.FileType = "markdown"
	}
	a.Row = 0
	a.Col = 0
	return nil
}

func (a *EditorAdapter) SaveFile() error {
	if a.File == "" {
		return nil
	}
	return a.Buf.SaveFile(a.File)
}

func (a *EditorAdapter) DeleteRange(startRow, startCol, endRow, endCol int) {
	startOff := a.Buf.PosToOffset(startRow, startCol)
	endOff := a.Buf.PosToOffset(endRow, endCol)
	if endOff > startOff {
		a.Buf.DeleteRange(startOff, endOff-startOff)
	}
}

func (a *EditorAdapter) LineCount() int {
	return a.Buf.LineCount()
}

func (a *EditorAdapter) Line(n int) string {
	return a.Buf.Line(n)
}

func (a *EditorAdapter) Content() string {
	return a.Buf.Content()
}

func (a *EditorAdapter) InsertAt(row, col int, text string) {
	offset := a.Buf.PosToOffset(row, col)
	a.Buf.InsertString(offset, text)
}

func (a *EditorAdapter) DeleteLines(start, end int) {
	if start >= end || start < 0 {
		return
	}
	startOff := a.Buf.PosToOffset(start, 0)
	var endOff int
	if end >= a.Buf.LineCount() {
		endOff = a.Buf.Length()
	} else {
		endOff = a.Buf.PosToOffset(end, 0)
	}
	if endOff > startOff {
		a.Buf.DeleteRange(startOff, endOff-startOff)
	}
}
