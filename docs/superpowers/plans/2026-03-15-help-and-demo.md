# Help Command & Demo Mode Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a scoped `:help` command with topic subcommands and a `--demo` flag that launches with a sample vault showcasing LazyMD features.

**Architecture:** Two new packages (`internal/help`, `internal/demo`) using Go's `//go:embed` to bundle markdown content. The help system loads static topic files into the editor buffer; the demo vault extracts embedded files to a temp directory. Both integrate through minimal changes to `buffer.go`, `editor.go`, `command.go`, `plugin.go`, and `main.go`.

**Tech Stack:** Go 1.24+, `embed` stdlib, existing buffer/editor/plugin infrastructure.

**Spec:** `docs/superpowers/specs/2026-03-15-help-and-demo-design.md`

---

## Chunk 1: Foundation (Buffer, Editor, Plugin DTOs)

### Task 1: Add `Buffer.SetContent()` method

**Files:**
- Modify: `internal/buffer/buffer.go:257-274` (after `LoadFile`)
- Test: `internal/buffer/buffer_test.go`

- [ ] **Step 1: Write the failing test**

Add to `internal/buffer/buffer_test.go`:

```go
func TestSetContent(t *testing.T) {
	buf := New()
	buf.InsertString(0, "old content")

	buf.SetContent("Hello\nNew World")

	if buf.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", buf.LineCount())
	}
	if buf.Line(0) != "Hello" {
		t.Fatalf("expected line 0 = %q, got %q", "Hello", buf.Line(0))
	}
	if buf.Line(1) != "New World" {
		t.Fatalf("expected line 1 = %q, got %q", "New World", buf.Line(1))
	}
	if buf.IsDirty() {
		t.Fatal("expected buffer to not be dirty after SetContent")
	}
	if buf.Length() != 15 {
		t.Fatalf("expected length 15, got %d", buf.Length())
	}
}

func TestSetContentClearsUndo(t *testing.T) {
	buf := New()
	buf.InsertString(0, "first")
	buf.InsertString(5, " second")

	buf.SetContent("replaced")

	buf.Undo()
	// Undo should be a no-op — stack was cleared
	if buf.Content() != "replaced" {
		t.Fatalf("expected undo to be no-op after SetContent, got %q", buf.Content())
	}
}

func TestSetContentEmpty(t *testing.T) {
	buf := New()
	buf.InsertString(0, "something")

	buf.SetContent("")

	if buf.Length() != 0 {
		t.Fatalf("expected length 0, got %d", buf.Length())
	}
	if buf.LineCount() != 1 {
		t.Fatalf("expected 1 line, got %d", buf.LineCount())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/buffer/ -run TestSetContent -v`
Expected: FAIL — `buf.SetContent` undefined

- [ ] **Step 3: Write the implementation**

Add to `internal/buffer/buffer.go` after the `LoadFile` method (after line 274):

```go
// SetContent replaces the buffer content from a string (no disk I/O).
// Resets undo/redo and marks the buffer as clean.
func (b *Buffer) SetContent(text string) {
	data := []byte(text)
	newCap := len(data) + initialGap
	b.data = make([]byte, newCap)
	copy(b.data, data)
	b.gapStart = len(data)
	b.gapEnd = newCap
	b.dirty = false
	b.undoStack = nil
	b.redoStack = nil
	b.rebuildLineStarts()
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/buffer/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/buffer/buffer.go internal/buffer/buffer_test.go
git commit -m "feat(buffer): add SetContent method for loading text without file I/O"
```

---

### Task 2: Add DTO types and extend `CommandExecutor` interface

**Files:**
- Modify: `internal/editor/editor.go:30-35` (CommandExecutor interface)

- [ ] **Step 1: Add DTO types and extend interface**

Add the DTO types before `CommandExecutor` (after line 28 in `editor.go`), and extend the interface:

```go
// PluginSummary describes a plugin (avoids circular import with plugins package).
type PluginSummary struct {
	Name        string
	Version     string
	Description string
}

// CommandSummary describes a plugin command (avoids circular import with plugins package).
type CommandSummary struct {
	Name        string
	Description string
	PluginName  string
}

// CommandExecutor is an interface for plugin command dispatch.
// This avoids circular imports between editor and plugins packages.
type CommandExecutor interface {
	ExecuteCommand(name string, ed PluginEditor, args string) bool
	Broadcast(eventType string, ed PluginEditor)
	ListPlugins() []PluginSummary
	ListCommands() []CommandSummary
}
```

- [ ] **Step 2: Verify the editor package itself compiles**

Run: `go build ./internal/editor/`
Expected: PASS (the editor package has no dependency on plugins, so the extended interface compiles fine on its own). Note: `go build ./...` will fail until Task 3 adds the new methods to `PluginManager`. Do NOT commit yet — Tasks 2 and 3 share a commit to avoid a broken intermediate state.

---

### Task 3: Implement `ListPlugins()` and `ListCommands()` on `PluginManager`

**Files:**
- Modify: `internal/plugins/plugin.go:108-125` (after existing `AllCommands` method)
- Test: `internal/plugins/plugins_test.go`

- [ ] **Step 1: Write the failing test**

Add to `internal/plugins/plugins_test.go`:

```go
func TestListPlugins(t *testing.T) {
	pm := NewManager()
	ed := newMockEditor("")
	pm.Register(&WordCountPlugin{}, ed)
	pm.Register(&KanbanPlugin{}, ed)

	summaries := pm.ListPlugins()
	if len(summaries) != 2 {
		t.Fatalf("expected 2 plugin summaries, got %d", len(summaries))
	}
	if summaries[0].Name != "word-count" {
		t.Errorf("expected first plugin name 'word-count', got %q", summaries[0].Name)
	}
}

func TestListCommands(t *testing.T) {
	pm := NewManager()
	ed := newMockEditor("")
	pm.Register(&KanbanPlugin{}, ed)

	cmds := pm.ListCommands()
	if len(cmds) != 3 {
		t.Fatalf("expected 3 command summaries, got %d", len(cmds))
	}
	// Verify plugin name is populated
	for _, cmd := range cmds {
		if cmd.PluginName != "kanban" {
			t.Errorf("expected PluginName 'kanban', got %q", cmd.PluginName)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/plugins/ -run "TestListPlugins|TestListCommands" -v`
Expected: FAIL — `pm.ListPlugins` undefined

- [ ] **Step 3: Write the implementation**

Add to `internal/plugins/plugin.go` after the `AllCommands` method (after line 125):

```go
// ListPlugins returns plugin summaries as editor DTOs.
func (pm *PluginManager) ListPlugins() []editor.PluginSummary {
	summaries := make([]editor.PluginSummary, len(pm.plugins))
	for i, p := range pm.plugins {
		info := p.Info()
		summaries[i] = editor.PluginSummary{
			Name:        info.Name,
			Version:     info.Version,
			Description: info.Description,
		}
	}
	return summaries
}

// ListCommands returns command summaries as editor DTOs.
func (pm *PluginManager) ListCommands() []editor.CommandSummary {
	var cmds []editor.CommandSummary
	for _, p := range pm.plugins {
		info := p.Info()
		for _, cmd := range p.Commands() {
			cmds = append(cmds, editor.CommandSummary{
				Name:        cmd.Name,
				Description: cmd.Description,
				PluginName:  info.Name,
			})
		}
	}
	return cmds
}
```

- [ ] **Step 4: Run all tests to verify everything compiles and passes**

Run: `go test ./internal/... -v`
Expected: All PASS (entire project compiles with the extended interface)

- [ ] **Step 5: Commit (Tasks 2 + 3 together)**

Both the interface extension and its implementation are committed together so the build never breaks:

```bash
git add internal/editor/editor.go internal/plugins/plugin.go internal/plugins/plugins_test.go
git commit -m "feat: add PluginSummary/CommandSummary DTOs and ListPlugins/ListCommands methods"
```

---

### Task 4: Add `EditorModel.LoadContent()` method

**Files:**
- Modify: `internal/editor/editor.go:137-151` (after `OpenFile` method)
- Test: `internal/editor/editor_test.go`

- [ ] **Step 1: Write the failing test**

Add to `internal/editor/editor_test.go`. First check what's there:

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/editor/ -run TestLoadContent -v`
Expected: FAIL — `ed.LoadContent` undefined

- [ ] **Step 3: Write the implementation**

Add to `internal/editor/editor.go` after the `OpenFile` method (after line 151):

```go
// LoadContent loads text directly into the buffer without file I/O.
// Used for help pages and other generated content.
func (e *EditorModel) LoadContent(name string, content string) {
	e.Buf.SetContent(content)
	e.File = name
	e.Row = 0
	e.Col = 0
	e.ScrollRow = 0
	e.ScrollCol = 0
	e.SetStatus("", false)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/editor/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/editor/editor.go internal/editor/editor_test.go
git commit -m "feat(editor): add LoadContent method for loading text without file I/O"
```

---

## Chunk 2: Help System

### Task 5: Create help topic registry with embedded content

**Files:**
- Create: `internal/help/help.go`
- Create: `internal/help/content/overview.md`
- Create: `internal/help/content/keys.md`
- Create: `internal/help/content/brain.md`
- Create: `internal/help/content/commands.md`
- Create: `internal/help/content/mcp.md`
- Create: `internal/help/content/panels.md`
- Test: `internal/help/help_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/help/help_test.go`:

```go
package help

import "testing"

func TestGetTopicOverview(t *testing.T) {
	content, ok := GetTopic("overview")
	if !ok {
		t.Fatal("expected overview topic to exist")
	}
	if len(content) == 0 {
		t.Fatal("expected overview content to be non-empty")
	}
}

func TestGetTopicAllStatic(t *testing.T) {
	topics := []string{"overview", "keys", "brain", "commands", "mcp", "panels"}
	for _, name := range topics {
		content, ok := GetTopic(name)
		if !ok {
			t.Errorf("topic %q not found", name)
		}
		if len(content) == 0 {
			t.Errorf("topic %q is empty", name)
		}
	}
}

func TestGetTopicUnknown(t *testing.T) {
	_, ok := GetTopic("nonexistent")
	if ok {
		t.Error("expected unknown topic to return false")
	}
}

func TestTopics(t *testing.T) {
	topics := Topics()
	if len(topics) < 6 {
		t.Errorf("expected at least 6 topics, got %d", len(topics))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/help/ -v`
Expected: FAIL — package doesn't exist

- [ ] **Step 3: Create the content markdown files**

Create `internal/help/content/overview.md`:

```markdown
# LazyMD Help

Terminal-native markdown editor with vim keybindings, knowledge graph, and 60+ plugins.

## Available Topics

| Topic | Command | Description |
|-------|---------|-------------|
| Keybindings | `:help keys` | Vim motions, panel toggles, mouse |
| Brain | `:help brain` | Knowledge graph and wiki-links |
| Plugins | `:help plugins` | All plugins and their commands |
| Commands | `:help commands` | Built-in editor commands |
| MCP Server | `:help mcp` | AI agent integration |
| Panels | `:help panels` | File tree, preview, brain panel |

## Quick Start

- Press `i` to enter insert mode, `Esc` to return to normal mode
- Type `:` to enter command mode
- Press `Alt+1` / `Alt+2` / `Alt+3` to toggle panels
- Use `[[note-name]]` syntax to create wiki-links between files
```

Create `internal/help/content/keys.md`:

```markdown
# Keybindings

## Normal Mode

| Key | Action |
|-----|--------|
| `h` `j` `k` `l` | Move left, down, up, right |
| `w` `b` `e` | Word forward, word back, end of word |
| `0` `$` | Start of line, end of line |
| `gg` `G` | Top of file, bottom of file |
| `d` + motion | Delete (dd = line, dw = word) |
| `c` + motion | Change (cc = line, cw = word) |
| `y` + motion | Yank (yy = line) |
| `p` | Paste |
| `u` | Undo |
| `Ctrl+r` | Redo |
| `i` `a` `o` `O` | Enter insert mode |
| `:` | Enter command mode |

## Insert Mode

| Key | Action |
|-----|--------|
| `Esc` | Return to normal mode |
| Type normally | Insert text |
| `Backspace` | Delete character behind cursor |

## Panel Toggles

| Key | Action |
|-----|--------|
| `Alt+1` | Toggle file tree |
| `Alt+2` | Toggle preview |
| `Alt+3` | Toggle brain graph |
| `Tab` | Cycle focus between panels |

## Panel Navigation

**File Tree:** `j`/`k` to navigate, `Enter` to open, `r` to refresh

**Preview:** `j`/`k` scroll, `d`/`u` page down/up

**Brain:** `j`/`k` select node, `f` toggle local mode, `+`/`-` zoom, `h`/`l` pan

## Mouse

- Click to focus panel or place cursor
- Scroll wheel to scroll within panels
```

Create `internal/help/content/brain.md`:

```markdown
# Brain — Knowledge Graph

The Brain panel visualizes connections between your notes using wiki-links.

## Wiki-Links

Use double brackets to link notes:

    [[note-name]]

Links are resolved by filename, case-insensitive. Path prefixes and display aliases are supported:

    [[folder/note]]       resolved by last segment
    [[note|display text]] resolved by part before |

## Brain Panel

Toggle with `Alt+3`. The panel shows a force-directed ASCII graph of your vault.

| Key | Action |
|-----|--------|
| `j` / `k` | Select next/previous node |
| `f` | Toggle local mode (show only neighbors) |
| `+` / `-` | Zoom in / out |
| `h` / `l` | Pan left / right |
| `Esc` | Return focus to editor |

## Graph Features

- **Backlinks**: See which notes link TO a given note
- **Neighbors**: Find connected notes within N hops
- **Path finding**: Shortest path between any two notes
- **Orphan detection**: Notes with zero connections
- **Hub notes**: Most connected notes ranked by link count

## MCP Brain Tools

When running as an MCP server (`lm --mcp-server`), AI agents can query the graph:

- `list_links` — outgoing wiki-links from the current document
- `get_backlinks` — files that link TO a given note
- `get_graph` — full graph as JSON (nodes, edges, stats)
- `get_neighbors` — directly connected notes with link counts
- `find_path` — shortest path between two notes
- `get_orphans` — notes with zero links
- `get_hub_notes` — most connected notes ranked by link count
```

Create `internal/help/content/commands.md`:

```markdown
# Built-in Commands

Enter command mode by pressing `:` in normal mode.

## File Commands

| Command | Description |
|---------|-------------|
| `:w` | Save current file |
| `:w <path>` | Save to a specific path |
| `:q` | Quit (fails if unsaved changes) |
| `:q!` | Force quit without saving |
| `:wq` or `:x` | Save and quit |
| `:e <path>` | Open a file |

## Theme Commands

| Command | Description |
|---------|-------------|
| `:theme` | Show current theme name |
| `:theme.cycle` | Cycle to next theme |
| `:theme.list` | List all available themes |
| `:theme <name>` | Switch to a specific theme |

## Help Commands

| Command | Description |
|---------|-------------|
| `:help` | Show this help overview |
| `:help keys` | Keybindings reference |
| `:help brain` | Knowledge graph guide |
| `:help plugins` | List all plugins and commands |
| `:help commands` | Built-in commands (this page) |
| `:help mcp` | MCP server setup |
| `:help panels` | Panel usage guide |

## Plugin Commands

Plugins register additional commands. Use `:help plugins` to see all available plugin commands.
```

Create `internal/help/content/mcp.md`:

```markdown
# MCP Server

LazyMD is an MCP (Model Context Protocol) server. AI agents connect via stdio using JSON-RPC 2.0.

## Starting the Server

    lm --mcp-server              # Start MCP server
    lm --mcp-server myfile.md    # Start with file preloaded

## Tools (22 total)

### Document Tools
`open_file`, `read_document`, `write_document`, `list_headings`, `edit_section`, `insert_text`, `delete_lines`, `search_content`, `get_structure`

### Navigation Tools
`read_section`, `list_tasks`, `update_task`, `get_breadcrumb`, `move_section`, `read_section_range`

### Brain Tools
`list_links`, `get_backlinks`, `get_graph`, `get_neighbors`, `find_path`, `get_orphans`, `get_hub_notes`

## Claude Code Setup

    claude mcp add LazyMD -- /path/to/lm --mcp-server

## Gemini CLI Setup

Add to `~/.gemini/settings.json`:

    {
      "mcpServers": {
        "LazyMD": {
          "command": "/path/to/lm",
          "args": ["--mcp-server"]
        }
      }
    }
```

Create `internal/help/content/panels.md`:

```markdown
# Panels

LazyMD has four panels: File Tree, Editor, Preview, and Brain.

## File Tree (Alt+1)

Browse files in the current directory.

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Open selected file |
| `r` | Refresh file listing |
| `Esc` | Return focus to editor |

Click a file entry to select it; click again to open.

## Editor (always visible)

Vim-style modal editor with syntax-aware rendering.

- **Normal mode**: Navigation and commands
- **Insert mode**: Text input (press `i` to enter)
- **Command mode**: Ex commands (press `:` to enter)

Line numbers in the gutter. Active line is highlighted.

## Preview (Alt+2)

Rendered markdown preview using Glamour.

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down/up |
| `d` / `u` | Page down/up |
| `Esc` | Return focus to editor |

Preview updates automatically with a 250ms debounce after edits.

## Brain (Alt+3)

Force-directed ASCII graph of wiki-link connections.

| Key | Action |
|-----|--------|
| `j` / `k` | Select next/previous node |
| `f` | Toggle local mode |
| `+` / `-` | Zoom in/out |
| `h` / `l` | Pan left/right |
| `Esc` | Return focus to editor |

Local mode shows only direct neighbors of the selected node.
```

- [ ] **Step 4: Create the help package**

Create `internal/help/help.go`:

```go
package help

import (
	"embed"
	"strings"
)

//go:embed content/*.md
var content embed.FS

var topicFiles = map[string]string{
	"overview": "content/overview.md",
	"keys":     "content/keys.md",
	"brain":    "content/brain.md",
	"commands": "content/commands.md",
	"mcp":      "content/mcp.md",
	"panels":   "content/panels.md",
}

// GetTopic returns the markdown content for a help topic.
func GetTopic(name string) (string, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "help" {
		name = "overview"
	}
	file, ok := topicFiles[name]
	if !ok {
		return "", false
	}
	data, err := content.ReadFile(file)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// Topics returns all available static topic names.
func Topics() []string {
	names := make([]string, 0, len(topicFiles))
	for name := range topicFiles {
		names = append(names, name)
	}
	return names
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/help/ -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add internal/help/
git commit -m "feat(help): add help package with embedded topic markdown files"
```

---

### Task 6: Wire `:help` command into the editor

**Files:**
- Modify: `internal/editor/command.go:90-107` (before `default` case)
- Test: `internal/editor/editor_test.go`

- [ ] **Step 1: Write failing tests for help command routing**

Add to `internal/editor/editor_test.go`:

```go
func TestHelpCommandLoadsContent(t *testing.T) {
	buf := buffer.New()
	ed := New(buf)
	// No dirty buffer, no CmdExec needed for static topics
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
	// Buffer content should be unchanged
	if ed.Buf.Line(0) != "dirty content" {
		t.Error("expected buffer content to be unchanged")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/editor/ -run "TestHelpCommand" -v`
Expected: FAIL — `executeCommand` is unexported but tests are in the same package, so should be accessible. The tests fail because the help case doesn't exist yet.

Note: `executeCommand` is a method on `*EditorModel` (see `command.go:31`). Since the test file is in `package editor`, it can call unexported methods directly.

- [ ] **Step 3: Add help command routing**

In `internal/editor/command.go`, add an import for the help package and a new case before the `default` block. The file currently has this import section:

```go
import (
	"fmt"
	"strings"

	"github.com/EME130/lazymd/internal/themes"
)
```

Add the help import:

```go
import (
	"fmt"
	"strings"

	"github.com/EME130/lazymd/internal/help"
	"github.com/EME130/lazymd/internal/themes"
)
```

Then add the `help` case before the `default:` block (before line 90 `default:`):

```go
	case cmd == "help" || strings.HasPrefix(cmd, "help "):
		if e.Buf.IsDirty() {
			e.SetStatus("Unsaved changes! Save with :w first", true)
			return
		}
		topic := "overview"
		if strings.HasPrefix(cmd, "help ") {
			topic = strings.TrimSpace(cmd[5:])
		}
		if topic == "plugins" {
			e.loadPluginsHelp()
			return
		}
		content, ok := help.GetTopic(topic)
		if !ok {
			e.SetStatus("Unknown help topic. Try :help for available topics", true)
			return
		}
		e.LoadContent("[help:"+topic+"]", content)
```

- [ ] **Step 4: Add the `loadPluginsHelp` method**

Add to `internal/editor/command.go` (at the end of the file):

```go
func (e *EditorModel) loadPluginsHelp() {
	if e.CmdExec == nil {
		e.SetStatus("Plugin system not available", true)
		return
	}

	var sb strings.Builder
	sb.WriteString("# Plugins\n\n")
	sb.WriteString(fmt.Sprintf("%d plugins registered.\n\n", len(e.CmdExec.ListPlugins())))

	pluginCmds := make(map[string][]CommandSummary)
	for _, cmd := range e.CmdExec.ListCommands() {
		pluginCmds[cmd.PluginName] = append(pluginCmds[cmd.PluginName], cmd)
	}

	for _, p := range e.CmdExec.ListPlugins() {
		sb.WriteString(fmt.Sprintf("## %s\n\n", p.Name))
		if p.Description != "" {
			sb.WriteString(p.Description + "\n\n")
		}
		cmds := pluginCmds[p.Name]
		if len(cmds) > 0 {
			sb.WriteString("| Command | Description |\n")
			sb.WriteString("|---------|-------------|\n")
			for _, cmd := range cmds {
				sb.WriteString(fmt.Sprintf("| `:%s` | %s |\n", cmd.Name, cmd.Description))
			}
			sb.WriteString("\n")
		}
	}

	e.LoadContent("[help:plugins]", sb.String())
}
```

- [ ] **Step 5: Run all tests**

Run: `go test ./... -v`
Expected: All PASS (including the new `TestHelpCommand*` tests)

- [ ] **Step 6: Commit**

```bash
git add internal/editor/command.go internal/editor/editor_test.go
git commit -m "feat(editor): wire :help command with topic routing and dynamic plugins help"
```

---

## Chunk 3: Demo Mode

### Task 7: Create demo vault embedded files

**Files:**
- Create: `internal/demo/vault/welcome.md`
- Create: `internal/demo/vault/project-plan.md`
- Create: `internal/demo/vault/knowledge/golang.md`
- Create: `internal/demo/vault/knowledge/markdown.md`
- Create: `internal/demo/vault/knowledge/wiki-links.md`
- Create: `internal/demo/vault/knowledge/zettelkasten.md`
- Create: `internal/demo/vault/daily-notes/2026-03-15.md`
- Create: `internal/demo/vault/kanban.md`

- [ ] **Step 1: Create all vault files**

Create each file with the content from the spec. The exact content is in `docs/superpowers/specs/2026-03-15-help-and-demo-design.md` lines 202-439.

`internal/demo/vault/welcome.md`:
```markdown
# Welcome to LazyMD

Your terminal-native markdown editor. This demo vault showcases the core features.

## Quick Start

- Press `Alt+3` to open the **Brain** panel and see the knowledge graph
- Press `Alt+2` to toggle **Preview** for rendered markdown
- Press `Alt+1` to toggle the **File Tree**
- Type `:help` for the built-in help system

## Explore This Vault

- [[project-plan]] — A sample project roadmap with tasks
- [[golang]] — Knowledge base: Go programming
- [[markdown]] — Knowledge base: Markdown syntax
- [[wiki-links]] — Knowledge base: How wiki-links work
- [[zettelkasten]] — Knowledge base: The Zettelkasten method
- [[kanban]] — A kanban board with task columns

## Try These Commands

- `:kanban` — View kanban board summary
- `:help plugins` — See all available plugins
- `:theme.cycle` — Cycle through color themes
```

`internal/demo/vault/project-plan.md`:
```markdown
# LazyMD Roadmap

A sample project plan demonstrating headings, tasks, and navigation.

## Phase 1: Core Editor

- [x] Gap buffer implementation
- [x] Vim modal editing
- [x] Syntax-aware line rendering
- [ ] Multi-cursor support

## Phase 2: Knowledge Graph

- [x] [[wiki-links]] parser
- [x] Backlink indexing
- [ ] Graph clustering
- [ ] See [[zettelkasten]] for methodology

## Phase 3: Integrations

- [x] MCP server (22 tools)
- [ ] Plugin marketplace
- [ ] Web publishing
- [ ] See [[golang]] for implementation notes
```

`internal/demo/vault/knowledge/golang.md`:
````markdown
# Go Programming

Notes on Go as used in LazyMD.

## Why Go?

- Fast compilation, single binary output
- Excellent concurrency primitives
- Strong standard library

## Key Packages

- `bubbletea` — Terminal UI framework
- `lipgloss` — Style definitions
- `glamour` — Markdown rendering

## Code Example

```go
func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Related

- [[markdown]] — The file format we edit
- [[wiki-links]] — How notes connect
````

`internal/demo/vault/knowledge/markdown.md`:
```markdown
# Markdown Syntax

LazyMD is a markdown-native editor. Here's a quick syntax reference.

## Text Formatting

- **Bold** with `**text**`
- *Italic* with `*text*`
- ~~Strikethrough~~ with `~~text~~`
- `Code` with backticks

## Block Elements

> Blockquotes use `>` prefix

- Unordered lists with `-`
1. Ordered lists with numbers

## Task Lists

- [x] Completed task
- [ ] Pending task

## Links

- Standard: [text](url)
- Wiki-link: [[note-name]]

## Related

- [[wiki-links]] — Deep dive on wiki-link syntax
- [[golang]] — The language LazyMD is built with
```

`internal/demo/vault/knowledge/wiki-links.md`:
```markdown
# Wiki-Links

Wiki-links are the backbone of LazyMD's knowledge graph.

## Syntax

Use double brackets to link notes:

    [[note-name]]

LazyMD resolves links by filename (case-insensitive, ignores path prefixes).

## How It Works

1. The scanner finds all `[[...]]` patterns in your vault
2. Each unique note becomes a node in the graph
3. Each wiki-link becomes a directed edge
4. The Brain panel (Alt+3) visualizes the result

## Features

- **Backlinks**: See what links TO a note
- **Neighbors**: Find connected notes within N hops
- **Path finding**: Shortest path between any two notes
- **Orphan detection**: Notes with zero connections
- **Hub notes**: Most connected notes ranked by link count

## Related

- [[zettelkasten]] — The method behind linked notes
- [[markdown]] — Base syntax for everything
- [[welcome]] — Back to the start
```

`internal/demo/vault/knowledge/zettelkasten.md`:
```markdown
# Zettelkasten Method

A note-taking method based on atomic, interconnected notes.

## Principles

1. **One idea per note** — Keep notes focused and atomic
2. **Link liberally** — Connect related ideas with [[wiki-links]]
3. **Use your own words** — Restate ideas to deepen understanding
4. **Let structure emerge** — Don't force hierarchy, let links create it

## In LazyMD

LazyMD's Brain panel visualizes your Zettelkasten as a force-directed graph.
Try pressing `Alt+3` to see how this demo vault's notes connect.

The MCP tools let AI agents traverse your graph:
- `get_neighbors` — explore nearby notes
- `find_path` — discover connections between distant ideas
- `get_hub_notes` — find your most central ideas

## Related

- [[wiki-links]] — The linking mechanism
- [[golang]] — Implementation language
```

`internal/demo/vault/daily-notes/2026-03-15.md`:
```markdown
# Daily Note — 2026-03-15

## Tasks

- [x] Review the [[project-plan]] roadmap
- [ ] Explore the [[wiki-links]] knowledge base
- [ ] Try the `:kanban` command
- [ ] Toggle Brain panel with Alt+3

## Notes

Today I'm exploring LazyMD's demo vault. The knowledge graph
shows how notes connect through wiki-links.

## Ideas

- Could use [[zettelkasten]] method for research notes
- Check [[markdown]] reference for table syntax
```

`internal/demo/vault/kanban.md`:
```markdown
# Kanban Board

Try `:kanban` to see a summary of this board.

## Backlog

- [ ] Add table of contents plugin
- [ ] Implement split-pane editing
- [ ] Add [[markdown]] table support

## In Progress

- [ ] Knowledge graph clustering
- [ ] Review [[project-plan]] milestones

## Done

- [x] Core editor with vim bindings
- [x] Brain panel with [[wiki-links]]
- [x] MCP server with 22 tools
- [x] Plugin system with 50+ plugins
```

- [ ] **Step 2: Commit vault files**

```bash
git add internal/demo/vault/
git commit -m "feat(demo): add embedded demo vault markdown files"
```

---

### Task 8: Create demo package with CreateVault/CleanupVault

**Files:**
- Create: `internal/demo/demo.go`
- Test: `internal/demo/demo_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/demo/demo_test.go`:

```go
package demo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateVault(t *testing.T) {
	path, err := CreateVault()
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}
	defer CleanupVault(path)

	// Check directory exists
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("vault dir not found: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected vault path to be a directory")
	}

	// Check expected files exist
	expectedFiles := []string{
		"welcome.md",
		"project-plan.md",
		"kanban.md",
		"knowledge/golang.md",
		"knowledge/markdown.md",
		"knowledge/wiki-links.md",
		"knowledge/zettelkasten.md",
		"daily-notes/2026-03-15.md",
	}
	for _, f := range expectedFiles {
		p := filepath.Join(path, f)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s not found: %v", f, err)
		}
	}
}

func TestCreateVaultContent(t *testing.T) {
	path, err := CreateVault()
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}
	defer CleanupVault(path)

	data, err := os.ReadFile(filepath.Join(path, "welcome.md"))
	if err != nil {
		t.Fatalf("failed to read welcome.md: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("welcome.md is empty")
	}
}

func TestCleanupVault(t *testing.T) {
	path, err := CreateVault()
	if err != nil {
		t.Fatalf("CreateVault failed: %v", err)
	}

	CleanupVault(path)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected vault dir to be removed, but it still exists")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/demo/ -v`
Expected: FAIL — package doesn't exist

- [ ] **Step 3: Write the implementation**

Create `internal/demo/demo.go`:

```go
package demo

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed vault
var vaultFS embed.FS

// CreateVault extracts the embedded demo vault to a temp directory.
func CreateVault() (string, error) {
	tmpDir, err := os.MkdirTemp("", "lm-demo-*")
	if err != nil {
		return "", err
	}

	err = fs.WalkDir(vaultFS, "vault", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Strip the "vault/" prefix to get the relative path
		rel, err := filepath.Rel("vault", path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(tmpDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := vaultFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})

	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

// CleanupVault removes a previously created demo vault.
func CleanupVault(path string) {
	os.RemoveAll(path)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/demo/ -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/demo/demo.go internal/demo/demo_test.go
git commit -m "feat(demo): add CreateVault/CleanupVault with embedded vault extraction"
```

---

### Task 9: Wire `--demo` flag into main.go

**Files:**
- Modify: `cmd/lm/main.go:17-22` (flag definitions) and `cmd/lm/main.go:85-97` (default TUI case)

- [ ] **Step 1: Add import and flag**

In `cmd/lm/main.go`, add the demo import:

```go
import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/EME130/lazymd/internal/agent"
	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/demo"
	"github.com/EME130/lazymd/internal/mcp"
	"github.com/EME130/lazymd/internal/ui"
	"github.com/EME130/lazymd/internal/web"
)
```

Add the flag after the existing flags (after line 21, before `flag.Parse()`):

```go
	demoMode := flag.Bool("demo", false, "Launch with demo vault")
```

- [ ] **Step 2: Add demo vault setup in the default TUI case**

Replace the `default:` case (lines 85-97) with:

```go
	default:
		filePath := ""
		args := flag.Args()
		if len(args) > 0 {
			filePath = args[0]
		}
		if *demoMode {
			vaultPath, err := demo.CreateVault()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create demo vault: %v\n", err)
				os.Exit(1)
			}
			defer demo.CleanupVault(vaultPath)
			if err := os.Chdir(vaultPath); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to enter demo vault: %v\n", err)
				os.Exit(1)
			}
			filePath = "welcome.md"
		}
		app := ui.NewApp(filePath)
		p := tea.NewProgram(app)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./cmd/lm/`
Expected: PASS

- [ ] **Step 4: Run all tests**

Run: `go test ./... -v`
Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/lm/main.go
git commit -m "feat(cli): add --demo flag to launch with embedded sample vault"
```

---

## Chunk 4: Final Verification

### Task 10: Full test suite and manual verification

- [ ] **Step 1: Run full test suite**

Run: `go test ./... -v`
Expected: All PASS

- [ ] **Step 2: Build the binary**

Run: `go build -o lm ./cmd/lm/`
Expected: Binary created successfully

- [ ] **Step 3: Manual test — demo mode**

Run: `./lm --demo`
Expected:
- App launches with `welcome.md` loaded
- File tree shows all 8 demo files
- Press `Alt+3` — Brain panel shows graph with connected nodes
- Press `Alt+2` — Preview shows rendered markdown

- [ ] **Step 4: Manual test — help command**

While in the TUI, type `:help` and press Enter.
Expected: Editor loads overview help with topic listing.

Type `:help brain` and press Enter.
Expected: Editor loads brain help content.

Type `:help plugins` and press Enter.
Expected: Editor loads dynamically generated plugin list with all 66 plugins.

Type `:help nonexistent` and press Enter.
Expected: Error message "Unknown help topic. Try :help for available topics"

- [ ] **Step 5: Manual test — dirty buffer guard**

Open a file, type `i` to enter insert mode, add some text, press `Esc`, then type `:help`.
Expected: Error message "Unsaved changes! Save with :w first"

- [ ] **Step 6: Final commit**

```bash
git add -A
git commit -m "feat: add :help command and --demo mode

Adds scoped :help system with 7 topics (keys, brain, plugins, commands,
mcp, panels, overview). :help plugins is dynamically generated from the
plugin registry.

Adds --demo flag that launches with an embedded 8-file sample vault
showcasing wiki-links, brain graph, kanban, tasks, and navigation."
```
