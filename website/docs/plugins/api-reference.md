---
title: Plugin API Reference
sidebar_position: 4
description: LazyMD plugin API reference — IR operations, events, frontend/backend interfaces, context APIs, passes, and capability declarations.
keywords: [LazyMD API, plugin API, IR operations, event types, frontend plugin, backend plugin, Go API reference]
---

# Plugin API Reference

All types live in `internal/pluginapi`.

## Operations (IR)

Operations are mutations flowing from frontends through passes to the core.

| OpType | Description |
|--------|-------------|
| `insert_text` | Insert text at a position |
| `delete_range` | Delete a range of text |
| `open_file` | Open a file |
| `save_file` | Save the current file |
| `set_cursor` | Move the cursor |
| `set_mode` | Change vim mode |
| `run_command` | Execute a `:command` |

```go
type Operation struct {
    Type     OpType
    Data     map[string]any
    Canceled bool
    Origin   string // plugin that emitted it
}
```

## Events

Events are notifications broadcast from the core to all plugins.

| EventType | Description |
|-----------|-------------|
| `file_opened` | A file was opened |
| `file_saved` | A file was saved |
| `file_closed` | A file was closed |
| `buffer_changed` | Buffer content changed |
| `mode_changed` | Vim mode changed |
| `cursor_moved` | Cursor position changed |
| `graph_updated` | Knowledge graph was refreshed |

```go
type Event struct {
    Type EventType
    Data map[string]any
}
```

## FrontendPlugin Interface

```go
type FrontendPlugin interface {
    Info() PluginInfo
    Init(ctx *FrontendContext) error
    Shutdown() error
    OnEvent(ctx *FrontendContext, event *Event)
    Render(width, height int) string
    KeyBindings() []KeyBinding
    Commands() []FrontendCommandDef
    StatusItems() []StatusItem
}
```

## BackendPlugin Interface

```go
type BackendPlugin interface {
    Info() PluginInfo
    Init(ctx *BackendContext) error
    Shutdown() error
    OnOperation(ctx *BackendContext, op *Operation)
    OnEvent(ctx *BackendContext, event *Event)
    Commands() []CommandDef
    Capabilities() []Capability
}
```

## Pass Interface

Passes transform or cancel operations before they reach the core. Lower priority runs first.

```go
type Pass interface {
    Name() string
    Priority() int // 0-100, lower runs first
    Transform(ctx *BackendContext, op *Operation) bool // false = drop
}
```

## Context APIs

### FrontendContext

```go
type FrontendContext struct {
    Emit   func(op *Operation)
    Editor EditorAPI
    Nav    NavAPI   // nil if unavailable
    Brain  BrainAPI // nil if unavailable
    Theme  ThemeAPI // nil if unavailable
    Log    func(msg string)
}
```

### BackendContext

```go
type BackendContext struct {
    Emit   func(op *Operation)
    Editor EditorAPI
    Nav    NavAPI   // nil if unavailable
    Brain  BrainAPI // nil if unavailable
    Log    func(msg string)
}
```

## EditorAPI

```go
type EditorAPI interface {
    Buffer() *buffer.Buffer
    CursorRow() int
    CursorCol() int
    FilePath() string
    Mode() string // "NORMAL", "INSERT", "COMMAND"
    SetStatus(msg string, isError bool)
    SetCursorRow(row int)
    SetCursorCol(col int)
    LineCount() int
    Line(n int) string
    Content() string
    InsertAt(row, col int, text string)
    DeleteLines(start, end int)
}
```

## Supporting Types

```go
type PluginInfo struct {
    Name        string
    Version     string
    Author      string
    Description string
}

type CommandDef struct {
    Name        string
    Description string
    Handler     func(ctx *BackendContext, args string) error
}

type FrontendCommandDef struct {
    Name        string
    Description string
    Handler     func(ctx *FrontendContext, args string) error
}

type KeyBinding struct {
    Mode    string // "normal", "insert", "*"
    Key     string // "ctrl+p", "gd"
    Handler func(ctx *FrontendContext)
}

type StatusItem struct {
    Name     string
    Align    string // "left" or "right"
    Renderer func(ctx *FrontendContext) string
}

type Capability struct {
    Name        string // "ai.complete", "export.pdf", "vcs.git"
    Description string
}
```
