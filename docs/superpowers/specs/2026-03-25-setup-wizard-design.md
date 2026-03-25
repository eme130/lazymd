# Setup Wizard Design

**Date**: 2026-03-25
**Status**: Approved

## Problem

LazyMD has no config system and no first-run detection. The app uses `cwd` as the implicit vault. All three frontends (TUI, Wails Desktop, MCP) need a consistent way to configure the vault path on first launch.

## Design

### Approach: Mandatory `RunSetup` on `FrontendPlugin`

Every frontend must implement `RunSetup`. Compile-time enforcement — you can't register a frontend without setup support.

### Config Package (`internal/config`)

Config file location: `~/.config/lazymd/config.json`

```json
{
  "vault_path": "/Users/eme/notes"
}
```

**API:**

```go
type Config struct {
    VaultPath string `json:"vault_path"`
}

func Load() (*Config, error)           // reads config file; returns zero-value Config if file missing; returns error on parse failure
func Save(cfg *Config) error           // atomic write (tmp+rename), creates dirs (0755) as needed
func NeedsSetup(cfg *Config) bool      // true if VaultPath is empty (string check only, no disk validation)
func Dir() string                      // returns ~/.config/lazymd
```

**Path rules:**
- `~` expands to `os.UserHomeDir()` (no `~user/` support)
- Relative paths resolved to absolute via `filepath.Abs()`
- Must be an existing directory (or user confirms creation)
- Symlinks followed, not resolved

### ConfigAPI in pluginapi

```go
// ConfigAPI exposes app configuration to plugins.
type ConfigAPI interface {
    VaultPath() string
}
```

Added to both `FrontendContext` and `BackendContext` as `Config ConfigAPI`.
Always available (never nil) — populated before `Init()` is called.

**Adapter:** `internal/pluginadapter/config.go` wraps `*config.Config` to implement `ConfigAPI`, consistent with existing adapter pattern.

### FrontendPlugin Contract Change

Add to `FrontendPlugin` interface:

```go
// RunSetup is called when first-run setup is needed.
// Frontends must present a vault selection UI or accept it programmatically.
// Return the chosen vault path, or error to abort.
RunSetup(cfg ConfigAPI) (vaultPath string, err error)
```

### Frontend Implementations

**TUI (`cmd/lm/`):**
- Standalone Bubble Tea program runs *before* the main app (separate `tea.Program`)
- UI flow:
  1. Show welcome message + prompt: "Enter vault path:"
  2. User types path (~ expansion applied)
  3. If path exists and is a directory → accept
  4. If path doesn't exist → prompt "Create directory? [Y/n]" → create with `os.MkdirAll(path, 0755)` or abort
  5. If path is a file or creation fails → show error, re-prompt
  6. `Esc` or `Ctrl+C` → exit app entirely
- Saves to config, then boots main app with `os.Chdir(vault)`

**Wails Desktop (`cmd/lm-desktop/`):**
- Setup screen component in Svelte frontend
- `WailsApp.SelectVaultDir() (string, error)` — bound method that opens native OS directory picker via `runtime.OpenDirectoryDialog()`
- Saves to config, then proceeds to main UI

**MCP (`--mcp-server` mode):**
- Accepts `--vault <path>` flag
- Precedence: `--vault` flag > saved config > error
- If `--vault` provided: saves to config and uses it (overrides any existing config)
- If no `--vault` and config exists: uses saved vault
- If no `--vault` and no config: exit with stderr message "run `lm` first to set up your vault, or pass `--vault <path>`"
- `RunSetup` returns the `--vault` arg (no interactive wizard)

**External plugins:** `RunSetup` returns `("", nil)` — they don't own setup, the host frontend does.

### Startup Flow

Every `cmd/` entry point follows:

```
1. Load config
2. NeedsSetup?
   → Yes → frontend.RunSetup() → Save config
   → No  → continue
3. os.Chdir(vault)  — if fails, exit with error
4. Create Engine, build contexts (ConfigAPI populated), SetContexts, InitAll
5. Normal boot (tea.Program / Wails / MCP loop)
```

Setup runs *before* the Engine exists. The entry point (`cmd/`) orchestrates setup, then creates the Engine with full context. This avoids chicken-and-egg: plugins don't need to be initialized to run setup.

**TUI**: Setup is a separate `tea.Program` that runs and exits before `NewApp()`.
**Wails**: Setup check happens in Go `main()` before `wails.Run()`. If needed, a setup-only Wails window launches first.
**MCP**: Setup is pure flag parsing in `main()`, no Engine involvement.

## New Files

- `internal/config/config.go` — Config struct, Load/Save/NeedsSetup
- `internal/config/config_test.go`
- `internal/pluginadapter/config.go` — ConfigAdapter wrapping *config.Config

## Modified Files

- `internal/pluginapi/api.go` — add `ConfigAPI` interface
- `internal/pluginapi/frontend.go` — add `RunSetup` to `FrontendPlugin`, add `Config` to `FrontendContext`
- `internal/pluginapi/backend.go` — add `Config` to `BackendContext`
- `internal/wailsplugin/plugin.go` — implement `RunSetup`
- `internal/wailsplugin/app.go` — add `RunSetup` bound method (native dir picker)
- `cmd/lm/main.go` — setup flow before app boot, `--vault` on MCP mode
- `cmd/lm-desktop/main.go` — setup flow before Wails boot

## Backward Compatibility

Adding `RunSetup` to `FrontendPlugin` is a breaking interface change. All existing implementations must add the method:
- `wailsplugin.Plugin` — gets real implementation
- External plugins (JSON-RPC) — the external plugin loader provides a default `RunSetup` that returns `("", nil)`

This is an internal interface with no published stability guarantee. All implementations are in-tree.

## Not In Scope

- Vault switching (`:vault` command) — future feature
- Multiple named vaults — future feature
- Theme selection in wizard — change anytime with `:theme`
- Windows/cross-platform path handling — macOS/Linux only for now
- Config file versioning — single-field struct, add when needed
