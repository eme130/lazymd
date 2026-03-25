---
title: Plugin System Overview
sidebar_position: 1
description: "LazyMD's LLVM-inspired plugin architecture — two-tier frontend/backend plugins, IR operations and events, middleware passes, and external plugin support via JSON-RPC."
keywords: [LazyMD plugins, plugin system, frontend plugin, backend plugin, IR pipeline, Go plugin, extensible editor]
---

# Plugin System Overview

lm has an LLVM-inspired two-tier plugin architecture:

- **Frontend plugins** — own UI surface: panels, keybindings, status bar widgets. They produce operations that flow into the IR pipeline.
- **Backend plugins** — extend core capabilities: AI, export, sync, indexing. They receive operations after passes and respond to events.
- **Passes** — middleware that transforms or cancels operations before they reach the core (inspired by LLVM optimization passes).

## IR Pipeline

Operations flow through an intermediate representation:

```
Frontend → Passes (transform/cancel) → Core → Backend
                                          ↓
                                    Events broadcast
                                    to all plugins
```

**Operations** are mutations (e.g., `insert_text`, `open_file`, `set_mode`). **Events** are notifications (e.g., `file_opened`, `buffer_changed`, `cursor_moved`).

## Plugin Types

| Type | Interface | Capabilities |
|------|-----------|-------------|
| Frontend | `FrontendPlugin` | Panels, keybindings, status items, `:commands` |
| Backend | `BackendPlugin` | Operations, events, `:commands`, capability declarations |
| External | `BackendPlugin` (via JSON-RPC stdio) | Same as backend, any language |

## External Plugins

External plugins run as child processes communicating over JSON-RPC 2.0 stdio. Place them in `~/.config/lazymd/plugins/<name>/` with a `manifest.json`. See [Development Guide](./development.md) for details.

## Key Packages

| Package | Role |
|---------|------|
| `internal/pluginapi` | IR types, interfaces, API contracts (imports only `buffer`) |
| `internal/plugins` | Engine, loader, registry, external plugin support |
| `internal/pluginadapter` | Adapters: nav, brain, theme, editor → pluginapi interfaces |
