# TUI Removal & Wails-Only Pivot

**Date:** 2026-03-28
**Status:** Approved

## Summary

LazyMD pivots from dual TUI+GUI to Wails GUI only. All TUI code, CLI modes (MCP server, web server, agent), and the vim modal editor are removed. The desktop binary (`cmd/lm-desktop`) becomes the sole entry point.

## Packages Deleted

| Package | Reason |
|---------|--------|
| `cmd/lm/` | TUI + CLI mode dispatcher |
| `internal/ui/` | Bubble Tea TUI (13 files) |
| `internal/demo/` | Demo vault (TUI only) |
| `internal/help/` | Embedded help topics (TUI `:help` command) |
| `internal/editor/` | Vim modal editor (6 files) |

## Type Migration: `editor` → `pluginapi`

Types previously in `editor` (to break circular imports) move to `pluginapi`:

- `PluginSummary` → `pluginapi.PluginSummary`
- `CommandSummary` → `pluginapi.CommandSummary`
- `CommandExecutor` → `pluginapi.CommandExecutor` (simplified: drop `PluginEditor` param)

## Adapter Rewrite: `pluginadapter/editor.go`

`EditorAdapter` currently wraps `editor.EditorModel`. Rewritten to wrap `buffer.Buffer` directly with simple fields for cursor, file path, status, and mode. Implements `pluginapi.EditorAPI` unchanged.

## Engine Update: `plugins/engine.go`

Replace all `editor.*` references with `pluginapi.*` equivalents. The `CommandExecutor` implementation drops the unused `editor.PluginEditor` parameter.

## Agent Update: `agent/agent.go`

Replace `editor.PluginEditor` with `pluginapi.EditorAPI`. Mode is a plain string instead of `editor.Mode`.

## Desktop Entry Point: `cmd/lm-desktop/main.go`

Remove `editor.New(buf)`. Create `EditorAdapter` directly with buffer.

## Dependencies Removed

- `charm.land/bubbletea/v2`
- `github.com/charmbracelet/glamour`
- `github.com/charmbracelet/lipgloss`
- `github.com/charmbracelet/harmonica`
