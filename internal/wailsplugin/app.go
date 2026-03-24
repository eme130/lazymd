package wailsplugin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/EME130/lazymd/internal/pluginapi"
)

const pluginOrigin = "wails-gui"

type CursorPos struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type StatusInfo struct {
	File     string `json:"file"`
	Lines    int    `json:"lines"`
	Row      int    `json:"row"`
	Col      int    `json:"col"`
	Modified bool   `json:"modified"`
}

type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	LinkCount int    `json:"linkCount"`
}

type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type FileEntry struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []FileEntry `json:"children,omitempty"`
}

type App struct {
	editor pluginapi.EditorAPI
	nav    pluginapi.NavAPI
	brain  pluginapi.BrainAPI
	theme  pluginapi.ThemeAPI
	emit   func(op *pluginapi.Operation)
}

func NewApp(ctx *pluginapi.FrontendContext) *App {
	return &App{
		editor: ctx.Editor,
		nav:    ctx.Nav,
		brain:  ctx.Brain,
		theme:  ctx.Theme,
		emit:   ctx.Emit,
	}
}

// --- Editor Methods ---

func (a *App) GetContent() string {
	if a.editor == nil {
		return ""
	}
	return a.editor.Content()
}

func (a *App) GetCursorPos() CursorPos {
	if a.editor == nil {
		return CursorPos{}
	}
	return CursorPos{Row: a.editor.CursorRow(), Col: a.editor.CursorCol()}
}

func (a *App) InsertText(row, col int, text string) {
	op := pluginapi.NewOperation(pluginapi.OpInsertText, pluginOrigin)
	op.Data["row"] = row
	op.Data["col"] = col
	op.Data["text"] = text
	a.emit(op)
}

func (a *App) DeleteRange(startRow, startCol, endRow, endCol int) {
	op := pluginapi.NewOperation(pluginapi.OpDeleteRange, pluginOrigin)
	op.Data["start_row"] = startRow
	op.Data["start_col"] = startCol
	op.Data["end_row"] = endRow
	op.Data["end_col"] = endCol
	a.emit(op)
}

func (a *App) SetCursor(row, col int) {
	op := pluginapi.NewOperation(pluginapi.OpSetCursor, pluginOrigin)
	op.Data["row"] = row
	op.Data["col"] = col
	a.emit(op)
}

// --- File Tree Methods ---

func (a *App) OpenFile(path string) {
	op := pluginapi.NewOperation(pluginapi.OpOpenFile, pluginOrigin)
	op.Data["path"] = path
	a.emit(op)
}

func (a *App) SaveFile() {
	op := pluginapi.NewOperation(pluginapi.OpSaveFile, pluginOrigin)
	a.emit(op)
}

func (a *App) GetCurrentFile() string {
	if a.editor == nil {
		return ""
	}
	return a.editor.FilePath()
}

func (a *App) ListFiles(dir string) []FileEntry {
	if dir == "" {
		dir = "."
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var result []FileEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		fullPath := filepath.Join(dir, e.Name())
		if e.IsDir() {
			children := a.ListFiles(fullPath)
			if len(children) > 0 {
				result = append(result, FileEntry{
					Name: e.Name(), Path: fullPath, IsDir: true, Children: children,
				})
			}
		} else {
			ext := filepath.Ext(e.Name())
			if ext == ".md" || ext == ".rndm" {
				result = append(result, FileEntry{Name: e.Name(), Path: fullPath})
			}
		}
	}
	return result
}

// --- Brain Graph Methods ---

func (a *App) GetGraph() GraphData {
	if a.brain == nil {
		return GraphData{}
	}
	nodes := a.brain.Nodes()
	gNodes := make([]GraphNode, len(nodes))
	for i, n := range nodes {
		gNodes[i] = GraphNode{Name: n.Name, Path: n.Path, LinkCount: n.LinkCount}
	}
	var edges []GraphEdge
	seen := make(map[string]bool)
	for _, n := range nodes {
		neighbors := a.brain.GetNeighbors(n.Name)
		for _, nb := range neighbors {
			key := n.Name + "->" + nb
			if !seen[key] {
				seen[key] = true
				edges = append(edges, GraphEdge{Source: n.Name, Target: nb})
			}
		}
	}
	return GraphData{Nodes: gNodes, Edges: edges}
}

func (a *App) GetNeighbors(note string) []string {
	if a.brain == nil {
		return nil
	}
	return a.brain.GetNeighbors(note)
}

func (a *App) FindPath(from, to string) []string {
	if a.brain == nil {
		return nil
	}
	path, found := a.brain.FindPath(from, to)
	if !found {
		return nil
	}
	return path
}

// --- Navigation Methods ---

func (a *App) ListHeadings() []map[string]any {
	if a.editor == nil {
		return nil
	}
	var headings []map[string]any
	for i := 0; i < a.editor.LineCount(); i++ {
		line := a.editor.Line(i)
		level := 0
		for _, ch := range line {
			if ch == '#' {
				level++
			} else {
				break
			}
		}
		if level > 0 && level <= 6 && len(line) > level && line[level] == ' ' {
			headings = append(headings, map[string]any{
				"line": i, "level": level, "title": line[level+1:],
			})
		}
	}
	return headings
}

func (a *App) ListTasks(section, status string) []pluginapi.TaskInfo {
	if a.nav == nil {
		return nil
	}
	var secPtr *string
	if section != "" {
		secPtr = &section
	}
	var donePtr *bool
	if status == "done" {
		d := true
		donePtr = &d
	} else if status == "pending" {
		d := false
		donePtr = &d
	}
	tasks, _ := a.nav.ListTasks(secPtr, donePtr)
	return tasks
}

func (a *App) UpdateTask(line int, done bool) {
	op := pluginapi.NewOperation(pluginapi.OpRunCommand, pluginOrigin)
	op.Data["command"] = "update_task"
	op.Data["line"] = line
	op.Data["done"] = done
	a.emit(op)
}

// --- Status Methods ---

func (a *App) GetStatus() StatusInfo {
	if a.editor == nil {
		return StatusInfo{}
	}
	return StatusInfo{
		File:  a.editor.FilePath(),
		Lines: a.editor.LineCount(),
		Row:   a.editor.CursorRow(),
		Col:   a.editor.CursorCol(),
	}
}

// --- Theme Methods ---

func (a *App) GetThemeColors() map[string]string {
	colors := map[string]string{
		"--lm-bg":           "#1a1b26",
		"--lm-fg":           "#c0caf5",
		"--lm-border":       "#3b4261",
		"--lm-heading":      "#bb9af7",
		"--lm-link":         "#7aa2f7",
		"--lm-code-bg":      "#292e42",
		"--lm-highlight":    "#292e42",
		"--lm-statusbar-bg": "#16161e",
	}
	if a.theme != nil {
		colors["name"] = a.theme.CurrentName()
	} else {
		colors["name"] = "tokyo-night"
	}
	return colors
}

func (a *App) ListThemes() []string {
	if a.theme == nil {
		return nil
	}
	return a.theme.ListThemes()
}

func (a *App) SetTheme(name string) bool {
	if a.theme == nil {
		return false
	}
	return a.theme.SetByName(name)
}
