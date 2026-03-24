# LazyMD Wails GUI Design

## Overview

Add a native desktop GUI surface to LazyMD using Wails (Go + Svelte). The GUI is a first-class surface alongside TUI, MCP, and web — not a replacement. It ships as a separate `lm-desktop` binary.

## Architecture: Embedded Core

The Wails GUI embeds the Go core directly. The `lm-desktop` binary is a single process: Wails binds Go structs to the Svelte frontend. The GUI registers as a `FrontendPlugin` in the plugin engine, emitting operations into the IR pipeline and receiving events back.

```
Svelte UI (webview)
    |
    | Wails Go bindings (auto-generated)
    v
WailsApp (bound struct)
    |
    | ctx.Emit(operation)
    v
WailsFrontendPlugin (FrontendPlugin)
    |
    v
Plugin Engine (passes -> core -> backends -> event broadcast)
    |
    | OnEvent()
    v
WailsFrontendPlugin -> Wails EventsEmit -> Svelte EventsOn
```

No MCP, no WebSocket, no subprocess. Direct in-process integration.

## Entry Point

**Binary:** `cmd/lm-desktop/main.go`

Initialization sequence:
1. Create core packages: `buffer.New()`, `editor.New()`, `nav.New()`, `brain.New()`
2. Create plugin engine: `plugins.NewEngine()`
3. Create `WailsFrontendPlugin`, register with engine
4. Create `WailsApp` (wraps plugin context for Go-Svelte bindings)
5. Call `wails.Run()` with Svelte frontend embedded via `//go:embed all:frontend/dist`

## Prerequisites: Extending the Plugin Interfaces

The current `FrontendPlugin` interface lacks two things the Wails GUI needs:

1. **`OnEvent()` method** — Currently only `BackendPlugin` has `OnEvent()`, and `Engine.BroadcastEvent()` only iterates backends. The Wails GUI needs to receive events (buffer changes, file opens, graph updates) to keep the Svelte UI in sync. **Fix:** Add `OnEvent(Event)` to the `FrontendPlugin` interface and update `Engine.BroadcastEvent()` to iterate both frontends and backends.

2. **`BrainAPI` on `FrontendContext`** — Currently `FrontendContext` exposes `EditorAPI`, `NavAPI`, and `ThemeAPI` but not `BrainAPI` (only `BackendContext` has it). The Wails GUI needs `BrainAPI` for the brain graph panel. **Fix:** Add `Brain BrainAPI` field to `FrontendContext`.

These are small, backward-compatible changes: existing frontend plugins that don't implement `OnEvent` get a no-op default, and the new `Brain` field is simply available for those that need it.

## WailsFrontendPlugin

**Package:** `internal/wailsplugin/plugin.go`

Implements `pluginapi.FrontendPlugin`:
- `Info() PluginInfo` — returns plugin metadata (name: "wails-gui", version, etc.)
- `Init(ctx FrontendContext)` — stores context, saves Wails runtime context for event emission
- `Shutdown()` — cleanup
- `OnEvent(event Event)` — forwards plugin events to Svelte via Wails `runtime.EventsEmit()`
- `Render(w, h int) string` — returns empty (Svelte handles rendering)
- `KeyBindings() []KeyBinding` — returns empty (browser handles input)
- `Commands() []Command` — GUI-specific commands if needed
- `StatusItems() []StatusItem` — returns empty (Svelte renders status bar)

**Event forwarding (`OnEvent`):**
When the plugin engine broadcasts events, `OnEvent()` forwards them to Svelte via `runtime.EventsEmit()`:

| Plugin Event | Wails Event | Svelte Action |
|---|---|---|
| `EventBufferChanged` | `buffer:changed` | Update CodeMirror content |
| `EventFileOpened` | `file:opened` | Update title, refresh file tree |
| `EventFileSaved` | `file:saved` | Clear modified indicator |
| `EventCursorMoved` | `cursor:moved` | Update status bar |
| `EventGraphUpdated` | `graph:updated` | Re-render D3 graph |
| `EventModeChanged` | `mode:changed` | (not used in GUI, no vim modes) |

## WailsApp Bindings

**Package:** `internal/wailsplugin/app.go`

Bound struct whose public methods are callable from Svelte.

### Editor Methods
- `GetContent() string` — full buffer content
- `GetCursorPos() CursorPos` — `{Row, Col int}`
- `InsertText(row, col int, text string) error` — emits `OpInsertText`
- `DeleteRange(startRow, startCol, endRow, endCol int) error` — emits `OpDeleteRange`
- `SetCursor(row, col int)` — emits `OpSetCursor`

### File Tree Methods
- `ListFiles(dir string) []FileEntry` — recursive listing, `.md`/`.rndm` files
- `OpenFile(path string) error` — emits `OpOpenFile`
- `SaveFile() error` — emits `OpSaveFile`
- `GetCurrentFile() string` — current file path

### Brain Graph Methods
- `GetGraph() GraphData` — nodes, edges, stats via `FrontendContext.Brain` (`BrainAPI`)
- `GetNeighbors(note string) []Neighbor` — direct connections
- `FindPath(from, to string) []string` — shortest path

### Navigation Methods
- `ListHeadings() []Heading` — document structure
- `ListTasks(section, status string) []Task` — task checkboxes
- `UpdateTask(line int, done bool) error` — toggle task

### Status Methods
- `GetStatus() StatusInfo` — file, line count, cursor, modified flag

### Window Methods
- `GetTheme() ThemeColors` — current theme as CSS-compatible values
- `OpenFileDialog() (string, error)` — native file picker via `runtime.OpenFileDialog`

## Svelte Frontend

**Location:** `cmd/lm-desktop/frontend/`

### Layout

CSS Grid with four resizable panels:
- Left: File Tree
- Center: Editor
- Right: Preview OR Brain Graph (mutually exclusive, toggled)
- Bottom: Status Bar + Command Bar

### Components

**`App.svelte`** — Root layout. Manages panel visibility, keyboard shortcuts (Cmd+P command palette, Cmd+B toggle sidebar, Cmd+Shift+P toggle preview/brain).

**`FileTree.svelte`** — Tree view with expand/collapse. Click to open file. File icons by extension. Calls `ListFiles()` on mount, refreshes on `file:opened`/`file:saved` events. Filesystem watcher on Go side emits events for external changes.

**`Editor.svelte`** — CodeMirror 6 with `@codemirror/lang-markdown`. Conventional editing: click-to-place cursor, mouse selection, Cmd+Z/Y undo/redo, Cmd+S save. On CodeMirror `onChange`, computes diff and calls `InsertText`/`DeleteRange`. Listens for `buffer:changed` events to apply external updates via CodeMirror `dispatch()`.

**`Preview.svelte`** — Renders markdown using remark + rehype with custom CSS matching LazyMD theme. Re-renders on `buffer:changed` events (debounced).

**`BrainGraph.svelte`** — D3.js force-directed graph on `<canvas>`. Calls `GetGraph()` on mount and on `graph:updated` events. Nodes are clickable (opens note). Supports zoom and pan. Node size proportional to link count.

**`StatusBar.svelte`** — Bottom bar: file path, cursor position (Ln/Col), modified indicator, encoding.

**`CommandBar.svelte`** — Command palette triggered by Cmd+P. Fuzzy search over: files, headings, commands. Calls appropriate bindings on selection.

### Theming

On startup, calls `GetTheme()` and injects values as CSS custom properties on `:root`. Theme switching re-queries and re-applies. All components use `var(--lm-bg)`, `var(--lm-fg)`, etc.

## State Synchronization

**Source of truth:** Go `buffer.Buffer` is authoritative. Svelte holds a local copy for rendering.

**User action flow:**
1. User types in CodeMirror
2. CodeMirror `onChange` fires
3. Svelte computes diff, calls bound Go method (e.g. `InsertText`)
4. Go method emits operation through IR pipeline
5. Plugin engine processes: passes -> core -> backends
6. Engine broadcasts `EventBufferChanged`
7. `WailsFrontendPlugin.OnEvent()` emits `buffer:changed` via Wails runtime
8. Svelte receives event, updates CodeMirror via `dispatch()`

**Self-echo prevention:** When Svelte receives a `buffer:changed` event that originated from its own action (checked via operation `Origin` field), it skips the redundant CodeMirror update.

**File tree refresh:** Triggered by file events plus `fsnotify` watcher on the Go side for external filesystem changes.

## Project Structure

```
cmd/
  lm-desktop/
    main.go                       # Wails entry point
    frontend/
      package.json
      svelte.config.js
      vite.config.js
      src/
        App.svelte
        components/
          FileTree.svelte
          Editor.svelte
          Preview.svelte
          BrainGraph.svelte
          StatusBar.svelte
          CommandBar.svelte
        lib/
          bindings.ts             # Wails auto-generated
          events.ts               # Event listener setup
          theme.ts                # CSS variable injection
internal/
  wailsplugin/
    plugin.go                     # WailsFrontendPlugin
    app.go                        # WailsApp bound struct
    events.go                     # Event forwarding logic
```

## Dependencies

**Go (go.mod):**
- `github.com/wailsapp/wails/v2`
- `github.com/fsnotify/fsnotify` (filesystem watcher for file tree refresh)

**JS (frontend/package.json):**
- `codemirror`, `@codemirror/lang-markdown`, `@codemirror/theme-one-dark`
- `d3`, `d3-force`
- `remark`, `remark-html`, `rehype-sanitize`

## Build & Run

```bash
# Development (hot reload)
cd cmd/lm-desktop && wails dev

# Production build
cd cmd/lm-desktop && wails build

# The binary embeds the Svelte frontend — single file distribution
```

## What Does NOT Change

- `internal/ui/` — TUI remains untouched
- `internal/mcp/` — MCP server unchanged
- `internal/web/` — Web server unchanged
- `internal/plugins/` — Plugin engine: minor change to `BroadcastEvent()` to include frontend plugins
- `internal/editor/` — Editor model unchanged (GUI uses conventional editing via CodeMirror, not vim keybindings)
- `cmd/lm/` — TUI binary unchanged

## Design Decisions

1. **Embedded core over MCP/WebSocket** — Single process, no IPC overhead, natural Wails pattern.
2. **Frontend plugin over direct package access** — Consistent with LLVM-inspired architecture, operations flow through IR pipeline.
3. **CodeMirror over custom editor** — Proven, extensible, handles conventional editing natively. No need to reimplement text editing in the browser.
4. **D3.js canvas over SVG** — Better performance for large graphs, supports smooth zoom/pan.
5. **remark over Glamour-to-HTML** — Browser-native rendering is higher fidelity in a webview context.
6. **Separate binary over flag** — Clean separation, independent build/release, no bloating the TUI binary with webview dependencies.
