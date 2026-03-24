package pluginapi

import "github.com/EME130/lazymd/internal/buffer"

// EditorAPI is the rich interface plugins use to access editor state.
type EditorAPI interface {
	Buffer() *buffer.Buffer
	CursorRow() int
	CursorCol() int
	FilePath() string
	Mode() string // "NORMAL", "INSERT", "COMMAND"
	SetStatus(msg string, isError bool)
	SetCursorRow(row int)
	SetCursorCol(col int)

	// Convenience methods
	LineCount() int
	Line(n int) string
	Content() string
	InsertAt(row, col int, text string)
	DeleteLines(start, end int) // end exclusive
}

// NavAPI exposes document navigation to plugins.
type NavAPI interface {
	ReadSection(headingPath string) (*SectionResult, error)
	ListTasks(section *string, done *bool) ([]TaskInfo, error)
	GetBreadcrumb(line int) (string, error)
}

// BrainAPI exposes the knowledge graph to plugins.
type BrainAPI interface {
	NodeCount() int
	Nodes() []NodeInfo
	GetBacklinks(name string) []string
	GetNeighbors(name string) []string
	FindPath(from, to string) ([]string, bool)
	GetOrphans() []string
}

// ThemeAPI exposes theme management to plugins.
type ThemeAPI interface {
	CurrentName() string
	SetByName(name string) bool
	ListThemes() []string
}

// SectionResult is returned by NavAPI.ReadSection.
type SectionResult struct {
	Title     string
	Level     int
	Content   string
	StartLine int
	EndLine   int
}

// TaskInfo describes a task checkbox in the document.
type TaskInfo struct {
	Line       int
	Done       bool
	Text       string
	Breadcrumb string
}

// NodeInfo describes a node in the knowledge graph.
type NodeInfo struct {
	Name      string
	Path      string
	LinkCount int
}
