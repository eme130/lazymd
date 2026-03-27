package pluginadapter

import (
	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/editor"
)

// EditorAdapter wraps an editor.EditorModel as a pluginapi.EditorAPI.
type EditorAdapter struct {
	Ed *editor.EditorModel
}

func (a *EditorAdapter) Buffer() *buffer.Buffer          { return a.Ed.Buffer() }
func (a *EditorAdapter) CursorRow() int                  { return a.Ed.CursorRow() }
func (a *EditorAdapter) CursorCol() int                  { return a.Ed.CursorCol() }
func (a *EditorAdapter) FilePath() string                { return a.Ed.FilePath() }
func (a *EditorAdapter) Mode() string                    { return a.Ed.Mode().String() }
func (a *EditorAdapter) SetStatus(msg string, isError bool) { a.Ed.SetStatus(msg, isError) }
func (a *EditorAdapter) SetCursorRow(row int)            { a.Ed.SetCursorRow(row) }
func (a *EditorAdapter) SetCursorCol(col int)            { a.Ed.SetCursorCol(col) }

func (a *EditorAdapter) OpenFile(path string) error { return a.Ed.OpenFile(path) }
func (a *EditorAdapter) SaveFile() error             { return a.Ed.Save() }

func (a *EditorAdapter) DeleteRange(startRow, startCol, endRow, endCol int) {
	buf := a.Ed.Buffer()
	startOff := buf.PosToOffset(startRow, startCol)
	endOff := buf.PosToOffset(endRow, endCol)
	if endOff > startOff {
		buf.DeleteRange(startOff, endOff-startOff)
	}
}

func (a *EditorAdapter) LineCount() int {
	return a.Ed.Buffer().LineCount()
}

func (a *EditorAdapter) Line(n int) string {
	return a.Ed.Buffer().Line(n)
}

func (a *EditorAdapter) Content() string {
	return a.Ed.Buffer().Content()
}

func (a *EditorAdapter) InsertAt(row, col int, text string) {
	offset := a.Ed.Buffer().PosToOffset(row, col)
	a.Ed.Buffer().InsertString(offset, text)
}

func (a *EditorAdapter) DeleteLines(start, end int) {
	if start >= end || start < 0 {
		return
	}
	buf := a.Ed.Buffer()
	startOff := buf.PosToOffset(start, 0)
	var endOff int
	if end >= buf.LineCount() {
		endOff = buf.Length()
	} else {
		endOff = buf.PosToOffset(end, 0)
	}
	if endOff > startOff {
		buf.DeleteRange(startOff, endOff-startOff)
	}
}
