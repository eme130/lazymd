---
title: System Design
sidebar_position: 1
description: Architecture overview of LazyMD — gap buffer, Bubble Tea TUI, LLVM-inspired plugin engine, Wails desktop GUI, and package separation.
keywords: [LazyMD architecture, gap buffer, Bubble Tea, Go terminal editor, plugin engine, system design, Wails desktop]
---

# System Design

lm is built with clean package separation:

```
cmd/
  lm/main.go              Entry point — TUI, MCP, Web, Agent mode dispatch
  lm-desktop/              Wails desktop GUI binary (Svelte + Go)
internal/
  buffer/                  Gap buffer, undo/redo, file I/O
  markdown/                Tokenizer (28 token types) + theme
  editor/                  Vim modes, cursor, keybindings
  nav/                     Navigation interface (headings, tasks, breadcrumbs)
  brain/                   Knowledge graph + vault scanner
  highlight/
    highlighter.go         Highlighter interface
    builtin.go             Keyword-based tokenizer
    languages.go           Language definitions (LaTeX, Python)
  pluginapi/               IR types, interfaces, API contracts
  plugins/                 Plugin engine, registry, loader, external support
  pluginadapter/           Adapters: nav, brain, theme, editor → pluginapi
  themes/                  Color themes
  ui/                      Bubble Tea TUI (app, layout, panels, styles)
  mcp/                     MCP server (JSON-RPC 2.0 over stdio)
  web/                     HTTP + WebSocket server
  agent/                   Agent backend interface + implementations
  wailsplugin/             Wails GUI frontend plugin + WailsApp bindings
```

## Key design decisions

- **Gap buffer** — Efficient for insert/delete at cursor position
- **Bubble Tea v2** — Elm Architecture TUI framework for composable models
- **Lip Gloss** — Declarative terminal styling
- **Glamour** — Markdown rendering for preview panel
- **LLVM-inspired plugin engine** — Two-tier (frontend/backend) with IR operations, events, and middleware passes
- **Interface-based highlighter** — Switchable backend for syntax highlighting
- **Wails desktop** — Svelte 4 + Vite frontend, Go backend. The GUI is a `FrontendPlugin` that emits operations into the same IR pipeline as the TUI
