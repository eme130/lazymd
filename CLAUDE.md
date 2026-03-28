# LazyMD

Desktop markdown editor written in Go + Wails. Inspired by lazygit/lazydocker.

## Tech Stack

- Language: Go (1.24.2+)
- Desktop: Wails v2 + Svelte 4 + Vite
- File format: `.rndm` (100% backward compatible with `.md`)

## Build Commands

```bash
task test             # Run all tests
task desktop:build    # Build desktop → bin/lm-desktop
task desktop:dev      # Wails dev mode (hot reload)
```

## Project Structure

```
cmd/
  lm-desktop/
    main.go              # Entry point (Wails desktop app)
    frontend/            # Svelte 4 + Vite frontend
internal/
  buffer/
    buffer.go            # Gap buffer with undo/redo, line tracking, file I/O
  markdown/
    parser.go            # Markdown tokenizer (28 token types)
    syntax.go            # Token colors and theme mapping
  nav/
    navigator.go         # Navigation interface (heading paths, tasks, breadcrumbs)
  brain/
    graph.go             # Knowledge graph (nodes, edges, backlinks, BFS)
    scanner.go           # Recursive vault scanner for [[wiki-links]]
  highlight/
    highlighter.go       # Highlighter interface (switchable backend)
    builtin.go           # Keyword-based tokenizer (default backend)
    languages.go         # 16 language definitions
  plugins/
    engine.go            # Plugin engine (registry, lifecycle, IR pipeline)
  pluginapi/
    api.go               # EditorAPI, NavAPI, BrainAPI, ThemeAPI, ConfigAPI interfaces
    ir.go                # Operations, Events (IR types)
    plugin.go            # PluginInfo, CommandExecutor, PluginSummary
    frontend.go          # FrontendPlugin interface, FrontendContext
    backend.go           # BackendPlugin interface, BackendContext, Pass
  pluginadapter/
    editor.go            # EditorAdapter (buffer-backed EditorAPI)
    nav.go               # NavAdapter
    brain.go             # BrainAdapter
    theme.go             # ThemeAdapter
    config.go            # ConfigAdapter
  corebackend/
    core.go              # Core backend plugin (executes IR operations)
  wailsplugin/
    app.go               # WailsApp bound struct (19 methods)
    plugin.go            # Wails GUI frontend plugin
    events.go            # Go → Svelte event forwarding
  themes/
    themes.go            # Color themes (Tokyo Night, Solarized, etc.)
  mcp/
    server.go            # MCP server (JSON-RPC 2.0 over stdio)
    tools.go             # 22 tool handlers
  web/
    server.go            # HTTP + WebSocket server
    websocket.go         # RFC 6455 WebSocket implementation
  agent/
    agent.go             # Agent types, Backend interface, Plugin
    mcp_backend.go       # MCP stdio backend
    websocket_backend.go # WebSocket backend
  config/
    config.go            # App configuration (vault path, etc.)
go.mod
go.sum
```

## Planned Features

- Built-in version control support
- Extensible plugin system
- ACP agent mode (host LazyMD as agent in Zed/JetBrains)

## Slash Commands

Go dev: `/release`, `/changelog`, `/issue-triage`, `/contrib-guide`, `/go-test`, `/go-debug`

## Hooks

- `PostToolUse` on `Edit|Write`: auto-runs `gofmt` on `.go` files via `.claude/hooks/go-fmt.sh`

## Reference Projects

- [lazygit](https://github.com/jesseduffield/lazygit)
- [lazydocker](https://github.com/jesseduffield/lazydocker)
