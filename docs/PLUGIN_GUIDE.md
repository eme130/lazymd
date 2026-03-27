# LazyMD Plugin Development Guide

Build plugins to extend LazyMD with custom commands, UI panels, and editor integrations.

## Overview

LazyMD's plugin system is LLVM-inspired with a two-tier architecture:

- **Frontend plugins** own the UI surface (panels, keybindings, status widgets). They emit operations into the IR pipeline.
- **Backend plugins** extend core capabilities (AI, export, sync, indexing). They receive operations and events.

Operations flow: **Frontend → Passes → Core → Backend**
Events flow: **Core → All plugins** (both frontend and backend)

## Quick Start: Backend Plugin (In-Process)

### 1. Implement `BackendPlugin`

```go
package myplugin

import "github.com/EME130/lazymd/internal/pluginapi"

type WordCount struct {
    count int
}

func (w *WordCount) Info() pluginapi.PluginInfo {
    return pluginapi.PluginInfo{
        Name:        "word-count",
        Version:     "0.1.0",
        Author:      "Your Name",
        Description: "Displays word count in the status bar",
    }
}

func (w *WordCount) Init(ctx *pluginapi.BackendContext) error {
    ctx.Log("word-count plugin loaded")
    return nil
}

func (w *WordCount) Shutdown() error { return nil }

func (w *WordCount) OnOperation(_ *pluginapi.BackendContext, _ *pluginapi.Operation) {}

func (w *WordCount) OnEvent(ctx *pluginapi.BackendContext, event *pluginapi.Event) {
    if event.Type != pluginapi.EventBufferChanged {
        return
    }
    content := ctx.Editor.Content()
    w.count = len(strings.Fields(content))
    ctx.Editor.SetStatus(fmt.Sprintf("Words: %d", w.count), false)
}

func (w *WordCount) Commands() []pluginapi.CommandDef { return nil }

func (w *WordCount) Capabilities() []pluginapi.Capability { return nil }
```

### 2. Register with the Engine

```go
eng := plugins.NewEngine(buf)
eng.RegisterBackend(&myplugin.WordCount{})
```

### 3. Build and Test

```bash
task build
task run
```

## Quick Start: External Plugin (Child Process)

External plugins run as separate executables and communicate via JSON-RPC 2.0 over stdio.

### 1. Create a Manifest

Place in `~/.config/lazymd/plugins/my-plugin/manifest.json`:

```json
{
    "name": "my-plugin",
    "version": "0.1.0",
    "author": "Your Name",
    "description": "An external LazyMD plugin",
    "tier": "backend",
    "entry": "./my-plugin",
    "events": ["buffer_changed", "file_opened"],
    "commands": [
        {
            "name": "greet",
            "description": "Display a greeting"
        }
    ],
    "capabilities": ["custom.greet"]
}
```

### 2. Implement the JSON-RPC Handler

Your executable must handle these JSON-RPC methods on stdin/stdout:

- `initialize` — called on startup with host version and capabilities
- `event` — notifications for subscribed events (no response expected)
- `command` — invoked when the user runs `:my-plugin.greet`
- `shutdown` — called on exit (graceful teardown)

### 3. Install

```bash
chmod +x ~/.config/lazymd/plugins/my-plugin/my-plugin
```

LazyMD discovers plugins automatically on startup.

## Plugin Interfaces

### BackendPlugin

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

### FrontendPlugin

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
    RunSetup(cfg ConfigAPI) (vaultPath string, err error)
}
```

## IR: Operations and Events

### Operations (Mutations)

Operations flow from frontends through passes into the core:

```go
op := pluginapi.NewOperation(pluginapi.OpInsertText, "my-plugin")
op.Data["row"] = 5
op.Data["col"] = 0
op.Data["text"] = "Hello, world!\n"
ctx.Emit(op)
```

Available operation types: `insert_text`, `delete_range`, `open_file`, `save_file`, `set_cursor`, `set_mode`, `run_command`.

### Events (Notifications)

Events are broadcast from the core to all plugins:

```go
func (p *MyPlugin) OnEvent(ctx *pluginapi.BackendContext, event *pluginapi.Event) {
    switch event.Type {
    case pluginapi.EventFileOpened:
        path := event.Data["path"].(string)
        ctx.Log("Opened: " + path)
    case pluginapi.EventBufferChanged:
        // React to content changes
    case pluginapi.EventCursorMoved:
        // React to cursor movement
    }
}
```

Available event types: `file_opened`, `file_saved`, `file_closed`, `buffer_changed`, `mode_changed`, `cursor_moved`, `graph_updated`.

### Passes (Operation Middleware)

Passes transform operations before they reach the core — inspired by LLVM optimization passes:

```go
type AutoSavePass struct{}

func (p *AutoSavePass) Name() string     { return "auto-save" }
func (p *AutoSavePass) Priority() int    { return 50 }

func (p *AutoSavePass) Transform(ctx *pluginapi.BackendContext, op *pluginapi.Operation) bool {
    if op.Type == pluginapi.OpInsertText {
        // Could trigger auto-save logic here
    }
    return true // return false to drop the operation
}
```

## API Surface

Plugins access editor state through context objects:

### EditorAPI

```go
ctx.Editor.Buffer()              // *buffer.Buffer
ctx.Editor.Content()             // full document text
ctx.Editor.Line(n)               // nth line
ctx.Editor.LineCount()           // total lines
ctx.Editor.CursorRow()           // cursor row
ctx.Editor.CursorCol()           // cursor col
ctx.Editor.FilePath()            // current file path
ctx.Editor.Mode()                // "NORMAL", "INSERT", "COMMAND"
ctx.Editor.SetStatus(msg, isErr) // set status bar message
ctx.Editor.InsertAt(row, col, text)
ctx.Editor.DeleteRange(startRow, startCol, endRow, endCol)
ctx.Editor.DeleteLines(start, end)
ctx.Editor.OpenFile(path)
ctx.Editor.SaveFile()
```

### NavAPI (Document Navigation)

```go
ctx.Nav.ReadSection("Plan/Step 1")           // read section by heading path
ctx.Nav.ListTasks(nil, nil)                  // list all tasks
ctx.Nav.GetBreadcrumb(42)                    // heading hierarchy for line 42
```

### BrainAPI (Knowledge Graph)

```go
ctx.Brain.NodeCount()
ctx.Brain.Nodes()
ctx.Brain.GetBacklinks("my-note")
ctx.Brain.GetNeighbors("my-note")
ctx.Brain.FindPath("note-a", "note-b")
ctx.Brain.GetOrphans()
```

### ThemeAPI (Frontend Only)

```go
ctx.Theme.CurrentName()
ctx.Theme.SetByName("tokyo-night")
ctx.Theme.ListThemes()
```

## Plugin Directory Structure

```
~/.config/lazymd/plugins/
  my-plugin/
    manifest.json    # Plugin metadata and configuration
    my-plugin        # Executable (any language)
```

## Best Practices

1. **Keep plugins focused** — one plugin, one responsibility
2. **Handle errors gracefully** — never crash the editor; return errors from `Init()`
3. **Be event-efficient** — return quickly from `OnEvent()`, skip events you don't need
4. **Use the IR** — emit operations instead of mutating state directly
5. **Document commands** — use clear descriptions in `CommandDef`
6. **Clean up resources** — release allocations in `Shutdown()`

## API Reference

See `internal/pluginapi/` for the complete type definitions:
- `plugin.go` — `PluginInfo`, `CommandDef`, `FrontendCommandDef`
- `ir.go` — `Operation`, `Event`, `OpType`, `EventType`
- `frontend.go` — `FrontendPlugin`, `FrontendContext`, `KeyBinding`, `StatusItem`
- `backend.go` — `BackendPlugin`, `BackendContext`, `Capability`, `Pass`
- `api.go` — `EditorAPI`, `NavAPI`, `BrainAPI`, `ThemeAPI`, `ConfigAPI`
