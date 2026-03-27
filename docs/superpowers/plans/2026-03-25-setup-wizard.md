# Setup Wizard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add first-run setup wizard that lets users select a vault directory, enforced across all frontends via the `FrontendPlugin` interface.

**Architecture:** New `internal/config` package handles config persistence (`~/.config/lazymd/config.json`). `ConfigAPI` interface added to pluginapi. `RunSetup` method added to `FrontendPlugin` — compile-time enforcement. Each `cmd/` entry point orchestrates setup before Engine creation.

**Tech Stack:** Go, Bubble Tea v2 (TUI wizard), Wails v2 runtime dialogs (desktop wizard)

**Spec:** `docs/superpowers/specs/2026-03-25-setup-wizard-design.md`

---

### Task 1: Config package — `internal/config`

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests for Load/Save/NeedsSetup**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNeedsSetup_EmptyVault(t *testing.T) {
	cfg := &Config{}
	if !NeedsSetup(cfg) {
		t.Fatal("expected NeedsSetup=true for empty VaultPath")
	}
}

func TestNeedsSetup_WithVault(t *testing.T) {
	cfg := &Config{VaultPath: "/tmp/vault"}
	if NeedsSetup(cfg) {
		t.Fatal("expected NeedsSetup=false when VaultPath is set")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := &Config{VaultPath: "/tmp/my-vault"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.VaultPath != "/tmp/my-vault" {
		t.Fatalf("got VaultPath=%q, want %q", loaded.VaultPath, "/tmp/my-vault")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should not error on missing file: %v", err)
	}
	if cfg.VaultPath != "" {
		t.Fatalf("expected empty VaultPath, got %q", cfg.VaultPath)
	}
}

func TestLoad_CorruptJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	configDir := filepath.Join(dir, ".config", "lazymd")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{bad json"), 0644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error on corrupt JSON")
	}
}

func TestSave_AtomicWrite(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfg := &Config{VaultPath: "/tmp/v1"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Overwrite with new value
	cfg2 := &Config{VaultPath: "/tmp/v2"}
	if err := Save(cfg2); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.VaultPath != "/tmp/v2" {
		t.Fatalf("got %q, want /tmp/v2", loaded.VaultPath)
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		input string
		want  string
	}{
		{"~/notes", filepath.Join(home, "notes")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", ""}, // becomes absolute, varies by cwd
	}
	for _, tt := range tests {
		got, err := ExpandPath(tt.input)
		if err != nil {
			t.Fatalf("ExpandPath(%q): %v", tt.input, err)
		}
		if tt.want != "" && got != tt.want {
			t.Fatalf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
		if !filepath.IsAbs(got) {
			t.Fatalf("ExpandPath(%q) = %q, not absolute", tt.input, got)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/config/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Write implementation**

```go
// internal/config/config.go
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "config.json"

// Config holds application configuration.
type Config struct {
	VaultPath string `json:"vault_path"`
}

// Dir returns the config directory path (~/.config/lazymd).
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "lazymd")
}

// Load reads config from disk. Returns zero-value Config if file is missing.
// Returns error on parse failure.
func Load() (*Config, error) {
	path := filepath.Join(Dir(), configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes config to disk atomically (write to tmp, then rename).
// Creates config directory (0755) if needed.
func Save(cfg *Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp := filepath.Join(dir, configFileName+".tmp")
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, configFileName))
}

// NeedsSetup returns true if VaultPath is empty (string check only).
func NeedsSetup(cfg *Config) bool {
	return cfg.VaultPath == ""
}

// ExpandPath expands ~ to home dir and resolves relative paths to absolute.
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home
	}
	return filepath.Abs(path)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/config/ -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/
git commit -m "feat(config): add config package with Load/Save/NeedsSetup/ExpandPath"
```

---

### Task 2: ConfigAPI interface + adapter

**Files:**
- Modify: `internal/pluginapi/api.go:1-71`
- Modify: `internal/pluginapi/frontend.go:27-34`
- Modify: `internal/pluginapi/backend.go:23-30`
- Create: `internal/pluginadapter/config.go`

- [ ] **Step 1: Add ConfigAPI to pluginapi/api.go**

Add after the `ThemeAPI` interface in `internal/pluginapi/api.go`:

```go
// ConfigAPI exposes app configuration to plugins.
type ConfigAPI interface {
	VaultPath() string
}
```

- [ ] **Step 2: Add Config field to FrontendContext**

In `internal/pluginapi/frontend.go`, add `Config ConfigAPI` field to `FrontendContext`:

```go
type FrontendContext struct {
	Emit   func(op *Operation)
	Editor EditorAPI
	Nav    NavAPI
	Brain  BrainAPI
	Theme  ThemeAPI
	Config ConfigAPI
	Log    func(msg string)
}
```

- [ ] **Step 3: Add Config field to BackendContext**

In `internal/pluginapi/backend.go`, add `Config ConfigAPI` field to `BackendContext`:

```go
type BackendContext struct {
	Emit   func(op *Operation)
	Editor EditorAPI
	Nav    NavAPI
	Brain  BrainAPI
	Config ConfigAPI
	Log    func(msg string)
}
```

- [ ] **Step 4: Create ConfigAdapter**

```go
// internal/pluginadapter/config.go
package pluginadapter

import "github.com/EME130/lazymd/internal/config"

// ConfigAdapter wraps *config.Config as a pluginapi.ConfigAPI.
type ConfigAdapter struct {
	Cfg *config.Config
}

func (a *ConfigAdapter) VaultPath() string {
	if a.Cfg == nil {
		return ""
	}
	return a.Cfg.VaultPath
}
```

- [ ] **Step 5: Verify compilation**

Run: `export PATH="/usr/local/go/bin:$PATH" && go build ./internal/pluginapi/ && go build ./internal/pluginadapter/`
Expected: SUCCESS

- [ ] **Step 6: Commit**

```bash
git add internal/pluginapi/api.go internal/pluginapi/frontend.go internal/pluginapi/backend.go internal/pluginadapter/config.go
git commit -m "feat(pluginapi): add ConfigAPI interface and ConfigAdapter"
```

---

### Task 3: Add RunSetup to FrontendPlugin interface

**Files:**
- Modify: `internal/pluginapi/frontend.go:5-24`
- Modify: `internal/wailsplugin/plugin.go:51` (stub to compile)

- [ ] **Step 1: Add RunSetup to FrontendPlugin**

In `internal/pluginapi/frontend.go`, add to the `FrontendPlugin` interface:

```go
// RunSetup is called when first-run setup is needed.
// Frontends must present a vault selection UI or accept it programmatically.
// Return the chosen vault path, or error to abort.
RunSetup(cfg ConfigAPI) (vaultPath string, err error)
```

- [ ] **Step 2: Add stub to wailsplugin.Plugin so it compiles**

In `internal/wailsplugin/plugin.go`, add:

```go
func (p *Plugin) RunSetup(_ pluginapi.ConfigAPI) (string, error) {
	return "", nil // placeholder — real implementation in Task 6
}
```

- [ ] **Step 3: Verify compilation**

Run: `export PATH="/usr/local/go/bin:$PATH" && go build ./...`
Expected: SUCCESS — all frontends implement the new method

- [ ] **Step 4: Commit**

```bash
git add internal/pluginapi/frontend.go internal/wailsplugin/plugin.go
git commit -m "feat(pluginapi): add RunSetup to FrontendPlugin interface"
```

---

### Task 4: TUI setup wizard

**Files:**
- Create: `internal/ui/setup.go`
- Create: `internal/ui/setup_test.go`

- [ ] **Step 1: Write test for setup model**

```go
// internal/ui/setup_test.go
package ui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupModel_AcceptExistingDir(t *testing.T) {
	dir := t.TempDir()
	m := NewSetupModel()
	m.input = dir

	// Simulate enter
	m, result := m.submit()
	if result != dir {
		t.Fatalf("got %q, want %q", result, dir)
	}
	if m.err != "" {
		t.Fatalf("unexpected error: %s", m.err)
	}
}

func TestSetupModel_RejectFile(t *testing.T) {
	f, _ := os.CreateTemp("", "test")
	f.Close()
	defer os.Remove(f.Name())

	m := NewSetupModel()
	m.input = f.Name()

	m, result := m.submit()
	if result != "" {
		t.Fatalf("expected empty result for file path, got %q", result)
	}
	if m.err == "" {
		t.Fatal("expected error for file path")
	}
}

func TestSetupModel_TildeExpansion(t *testing.T) {
	home, _ := os.UserHomeDir()
	m := NewSetupModel()
	m.input = "~/test-vault-" + t.Name()

	expanded := m.expandInput()
	want := filepath.Join(home, "test-vault-"+t.Name())
	if expanded != want {
		t.Fatalf("got %q, want %q", expanded, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/ui/ -run TestSetup -v`
Expected: FAIL — functions don't exist

- [ ] **Step 3: Write setup model implementation**

```go
// internal/ui/setup.go
package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/EME130/lazymd/internal/config"
	"github.com/EME130/lazymd/internal/themes"
	"github.com/charmbracelet/lipgloss"
)

// SetupModel is a standalone Bubble Tea model for first-run vault selection.
type SetupModel struct {
	input       string
	cursor      int
	err         string
	confirming  bool // asking to create dir
	done        bool
	vaultPath   string // result
	aborted     bool
	width       int
	height      int
}

// NewSetupModel creates a new setup wizard model.
func NewSetupModel() SetupModel {
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "notes")
	return SetupModel{
		input:  defaultPath,
		cursor: len(defaultPath),
	}
}

// VaultPath returns the chosen vault path after the wizard completes.
func (m SetupModel) VaultPath() string { return m.vaultPath }

// Aborted returns true if the user cancelled.
func (m SetupModel) Aborted() bool { return m.aborted }

func (m SetupModel) Init() tea.Cmd { return nil }

func (m SetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		key := msg.String()

		if m.confirming {
			switch key {
			case "y", "Y", "enter":
				expanded := m.expandInput()
				if err := os.MkdirAll(expanded, 0755); err != nil {
					m.err = fmt.Sprintf("Failed to create: %v", err)
					m.confirming = false
					return m, nil
				}
				m.vaultPath = expanded
				m.done = true
				return m, tea.Quit
			case "n", "N", "esc":
				m.confirming = false
				m.err = ""
				return m, nil
			}
			return m, nil
		}

		switch key {
		case "ctrl+c", "esc":
			m.aborted = true
			m.done = true
			return m, tea.Quit

		case "enter":
			var result string
			m, result = m.submit()
			if result != "" {
				m.vaultPath = result
				m.done = true
				return m, tea.Quit
			}
			return m, nil

		case "backspace":
			if m.cursor > 0 {
				m.input = m.input[:m.cursor-1] + m.input[m.cursor:]
				m.cursor--
			}
		case "delete":
			if m.cursor < len(m.input) {
				m.input = m.input[:m.cursor] + m.input[m.cursor+1:]
			}
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right":
			if m.cursor < len(m.input) {
				m.cursor++
			}
		case "home", "ctrl+a":
			m.cursor = 0
		case "end", "ctrl+e":
			m.cursor = len(m.input)
		default:
			runes := []rune(key)
			if len(runes) == 1 {
				m.input = m.input[:m.cursor] + string(runes[0]) + m.input[m.cursor:]
				m.cursor++
			}
		}
		m.err = ""
	}
	return m, nil
}

func (m SetupModel) expandInput() string {
	expanded, err := config.ExpandPath(strings.TrimSpace(m.input))
	if err != nil {
		return strings.TrimSpace(m.input)
	}
	return expanded
}

// submit validates and returns the vault path, or sets an error.
func (m SetupModel) submit() (SetupModel, string) {
	expanded := m.expandInput()

	info, err := os.Stat(expanded)
	if err != nil {
		if os.IsNotExist(err) {
			m.confirming = true
			m.err = ""
			return m, ""
		}
		m.err = fmt.Sprintf("Cannot access: %v", err)
		return m, ""
	}

	if !info.IsDir() {
		m.err = "Path is a file, not a directory"
		return m, ""
	}

	return m, expanded
}

func (m SetupModel) View() tea.View {
	c := themes.CurrentColors()

	accent := lipgloss.NewStyle().Foreground(lipgloss.Color(c.BorderActive)).Bold(true)
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color(c.TextMuted))
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Error))

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(accent.Render("  Welcome to LazyMD"))
	b.WriteString("\n\n")
	b.WriteString(muted.Render("  Select your vault directory — where your notes live."))
	b.WriteString("\n\n")

	if m.confirming {
		b.WriteString("  Directory does not exist. Create it?\n")
		b.WriteString(muted.Render("  " + m.expandInput()))
		b.WriteString("\n\n")
		b.WriteString(accent.Render("  [Y]es") + "  " + muted.Render("[N]o"))
	} else {
		b.WriteString("  Vault path: ")
		// Render input with cursor
		before := m.input[:m.cursor]
		after := ""
		cursorCh := " "
		if m.cursor < len(m.input) {
			cursorCh = string(m.input[m.cursor])
			after = m.input[m.cursor+1:]
		}
		cursor := lipgloss.NewStyle().Reverse(true).Render(cursorCh)
		b.WriteString(before + cursor + after)

		if m.err != "" {
			b.WriteString("\n\n")
			b.WriteString("  " + errStyle.Render(m.err))
		}
	}

	b.WriteString("\n\n")
	b.WriteString(muted.Render("  enter accept • esc quit"))
	b.WriteString("\n")

	var v tea.View
	v.SetContent(b.String())
	return v
}

// RunSetupWizard runs the TUI setup wizard and returns the chosen vault path.
// Returns ("", nil) if the user aborted.
func RunSetupWizard() (string, error) {
	m := NewSetupModel()
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	final := result.(SetupModel)
	if final.Aborted() {
		return "", nil
	}
	return final.VaultPath(), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./internal/ui/ -run TestSetup -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/ui/setup.go internal/ui/setup_test.go
git commit -m "feat(ui): add TUI setup wizard for first-run vault selection"
```

---

### Task 5: Wire setup into TUI entry point (`cmd/lm/main.go`)

**Files:**
- Modify: `cmd/lm/main.go`

- [ ] **Step 1: Add --vault flag and setup flow**

Modify `cmd/lm/main.go` to:
1. Add `--vault` flag for MCP mode
2. Load config before dispatching to any mode
3. If TUI mode and `NeedsSetup` → run wizard
4. If MCP mode → use `--vault` flag or saved config
5. `os.Chdir(vault)` before normal boot

Replace the full `main()` function:

```go
package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/EME130/lazymd/internal/agent"
	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/config"
	"github.com/EME130/lazymd/internal/demo"
	"github.com/EME130/lazymd/internal/mcp"
	"github.com/EME130/lazymd/internal/ui"
	"github.com/EME130/lazymd/internal/web"
)

var version = "dev"

func main() {
	mcpServer := flag.Bool("mcp-server", false, "Start MCP server mode")
	mcpAlias := flag.Bool("mcp", false, "Start MCP server mode (alias)")
	webServer := flag.Bool("web-server", false, "Start web server mode")
	port := flag.Int("port", 8080, "Web server port")
	agentMode := flag.Bool("agent", false, "Start agent mode")
	demoMode := flag.Bool("demo", false, "Launch with demo vault")
	showVersion := flag.Bool("version", false, "Print version and exit")
	vaultFlag := flag.String("vault", "", "Vault directory path (for MCP/agent/web modes)")
	flag.Parse()

	if *showVersion {
		fmt.Println("lm " + version)
		return
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Handle --vault flag: override config
	if *vaultFlag != "" {
		expanded, err := config.ExpandPath(*vaultFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid vault path: %v\n", err)
			os.Exit(1)
		}
		cfg.VaultPath = expanded
		if err := config.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
			os.Exit(1)
		}
	}

	switch {
	case *mcpServer || *mcpAlias:
		if config.NeedsSetup(cfg) {
			fmt.Fprintln(os.Stderr, "No vault configured. Run `lm` first to set up your vault, or pass `--vault <path>`.")
			os.Exit(1)
		}
		if err := os.Chdir(cfg.VaultPath); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot access vault %q: %v\n", cfg.VaultPath, err)
			os.Exit(1)
		}

		buf := buffer.New()
		srv := mcp.New(buf)

		args := flag.Args()
		if len(args) > 0 {
			srv.SetFilePath(args[0])
			if err := srv.LoadFileIfNeeded(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load file: %v\n", err)
				os.Exit(1)
			}
		}

		if err := srv.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(1)
		}
	case *webServer:
		if config.NeedsSetup(cfg) {
			fmt.Fprintln(os.Stderr, "No vault configured. Run `lm` first to set up your vault, or pass `--vault <path>`.")
			os.Exit(1)
		}
		if err := os.Chdir(cfg.VaultPath); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot access vault %q: %v\n", cfg.VaultPath, err)
			os.Exit(1)
		}

		buf := buffer.New()
		srv := web.NewServer(buf, *port)

		args := flag.Args()
		if len(args) > 0 {
			srv.SetFilePath(args[0])
			if err := buf.LoadFile(args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load file: %v\n", err)
				os.Exit(1)
			}
		}

		if err := srv.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Web server error: %v\n", err)
			os.Exit(1)
		}
	case *agentMode:
		if config.NeedsSetup(cfg) {
			fmt.Fprintln(os.Stderr, "No vault configured. Run `lm` first to set up your vault, or pass `--vault <path>`.")
			os.Exit(1)
		}
		if err := os.Chdir(cfg.VaultPath); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot access vault %q: %v\n", cfg.VaultPath, err)
			os.Exit(1)
		}

		buf := buffer.New()
		mcpSrv := mcp.New(buf)

		args := flag.Args()
		if len(args) > 0 {
			mcpSrv.SetFilePath(args[0])
			if err := mcpSrv.LoadFileIfNeeded(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load file: %v\n", err)
				os.Exit(1)
			}
		}

		ap := agent.NewPlugin()
		ap.SetBackend(agent.NewMcpBackend())
		ap.Connect(agent.Config{Transport: agent.TransportStdio})

		fmt.Fprintln(os.Stderr, "LazyMD agent mode: MCP stdio backend active")

		if err := mcpSrv.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Agent error: %v\n", err)
			os.Exit(1)
		}
	default:
		// TUI mode
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
		} else {
			// Setup wizard if needed
			if config.NeedsSetup(cfg) {
				vaultPath, err := ui.RunSetupWizard()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
					os.Exit(1)
				}
				if vaultPath == "" {
					// User aborted
					return
				}
				cfg.VaultPath = vaultPath
				if err := config.Save(cfg); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to save config: %v\n", err)
					os.Exit(1)
				}
			}
			if err := os.Chdir(cfg.VaultPath); err != nil {
				fmt.Fprintf(os.Stderr, "Cannot access vault %q: %v\n", cfg.VaultPath, err)
				os.Exit(1)
			}
		}

		filePath := ""
		args := flag.Args()
		if len(args) > 0 {
			filePath = args[0]
		}
		if *demoMode {
			filePath = "welcome.md"
		}
		app := ui.NewApp(filePath)
		p := tea.NewProgram(app)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
```

- [ ] **Step 2: Verify compilation**

Run: `export PATH="/usr/local/go/bin:$PATH" && go build ./cmd/lm/`
Expected: SUCCESS

- [ ] **Step 3: Run all tests**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add cmd/lm/main.go
git commit -m "feat(cli): wire setup wizard into TUI, add --vault flag for MCP/agent/web"
```

---

### Task 6: Wails desktop setup — Go side

**Files:**
- Modify: `internal/wailsplugin/plugin.go:39-42` (replace stub)
- Modify: `internal/wailsplugin/app.go` (add SelectVaultDir)
- Modify: `cmd/lm-desktop/main.go`

- [ ] **Step 1: Add SelectVaultDir to WailsApp**

In `internal/wailsplugin/app.go`, add:

```go
import wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

// SelectVaultDir opens a native directory picker and returns the chosen path.
func (a *App) SelectVaultDir(ctx context.Context) (string, error) {
	return wailsRuntime.OpenDirectoryDialog(ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Vault Directory",
	})
}
```

Note: `App` needs a `wailsCtx context.Context` field set during OnStartup, or accept it as parameter. Since Wails binds methods, add the context field:

In `internal/wailsplugin/app.go`, add to `App` struct:

```go
type App struct {
	editor   pluginapi.EditorAPI
	nav      pluginapi.NavAPI
	brain    pluginapi.BrainAPI
	theme    pluginapi.ThemeAPI
	emit     func(op *pluginapi.Operation)
	wailsCtx context.Context
}
```

Add setter:
```go
func (a *App) SetWailsContext(ctx context.Context) {
	a.wailsCtx = ctx
}
```

Update `SelectVaultDir` to use stored context:
```go
func (a *App) SelectVaultDir() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.wailsCtx, wailsRuntime.OpenDialogOptions{
		Title: "Select Vault Directory",
	})
}
```

- [ ] **Step 2: Update wailsplugin.Plugin RunSetup**

Replace the stub in `internal/wailsplugin/plugin.go`:

```go
func (p *Plugin) RunSetup(_ pluginapi.ConfigAPI) (string, error) {
	// Wails setup is handled by the frontend calling SelectVaultDir on WailsApp.
	// This method is not called directly — the cmd/lm-desktop entry point orchestrates.
	return "", nil
}
```

(Keep as no-op — Wails setup is driven by the Svelte frontend calling the bound `SelectVaultDir` method, not by the plugin engine.)

- [ ] **Step 3: Wire setup into cmd/lm-desktop/main.go**

Update `cmd/lm-desktop/main.go` to load config, add ConfigAdapter to contexts, and `os.Chdir` to vault. If `NeedsSetup`, the Wails app launches with a setup flag — the Svelte frontend shows a setup screen instead of the main UI.

Add config loading before `wails.Run()`:

```go
import "github.com/EME130/lazymd/internal/config"

// In main(), before wails.Run:
cfg, err := config.Load()
if err != nil {
	log.Fatalf("Failed to load config: %v", err)
}

needsSetup := config.NeedsSetup(cfg)
```

Add `ConfigAdapter` to both contexts:

```go
configAdapter := &pluginadapter.ConfigAdapter{Cfg: cfg}

// In frontendCtx:
Config: configAdapter,

// In backendCtx:
Config: configAdapter,
```

Add `NeedsSetup() bool` method to `WailsApp` so Svelte can check:

```go
// In internal/wailsplugin/app.go:
func (a *App) NeedsSetup() bool {
	return a.vaultPath == ""
}
```

Pass needsSetup state to app. In `OnStartup`, if vault is set, chdir:

```go
OnStartup: func(ctx context.Context) {
	app.SetWailsContext(ctx)
	guiPlugin.SetWailsRuntime(ctx, runtime.EventsEmit)
	if !needsSetup {
		os.Chdir(cfg.VaultPath)
	}
	engine.InitAll()
},
```

Add `SaveVault(path string) error` to WailsApp for the Svelte setup to call after dir selection:

```go
// In internal/wailsplugin/app.go:
func (a *App) SaveVault(path string) error {
	expanded, err := config.ExpandPath(path)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}
	cfg.VaultPath = expanded
	if err := config.Save(cfg); err != nil {
		return err
	}
	return os.Chdir(expanded)
}
```

- [ ] **Step 4: Verify compilation**

Run: `export PATH="/usr/local/go/bin:$PATH" && go build ./cmd/lm-desktop/`
Expected: SUCCESS

- [ ] **Step 5: Commit**

```bash
git add internal/wailsplugin/plugin.go internal/wailsplugin/app.go cmd/lm-desktop/main.go
git commit -m "feat(desktop): add vault selection for Wails GUI setup flow"
```

---

### Task 7: Wails desktop setup — Svelte frontend

**Files:**
- Create: `cmd/lm-desktop/frontend/src/components/SetupWizard.svelte`
- Modify: `cmd/lm-desktop/frontend/src/App.svelte`

- [ ] **Step 1: Read current App.svelte to understand structure**

Read `cmd/lm-desktop/frontend/src/App.svelte` to understand the current component structure.

- [ ] **Step 2: Create SetupWizard.svelte**

```svelte
<!-- cmd/lm-desktop/frontend/src/components/SetupWizard.svelte -->
<script>
  import { SelectVaultDir, SaveVault } from '../../wailsjs/go/wailsplugin/App';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let vaultPath = '';
  let error = '';
  let saving = false;

  async function browse() {
    try {
      const path = await SelectVaultDir();
      if (path) {
        vaultPath = path;
        error = '';
      }
    } catch (e) {
      error = 'Failed to open directory picker';
    }
  }

  async function confirm() {
    if (!vaultPath) {
      error = 'Please select a directory';
      return;
    }
    saving = true;
    try {
      await SaveVault(vaultPath);
      dispatch('complete');
    } catch (e) {
      error = `Failed to save: ${e}`;
      saving = false;
    }
  }
</script>

<div class="setup">
  <h1>Welcome to LazyMD</h1>
  <p class="subtitle">Select your vault directory — where your notes live.</p>

  <div class="input-row">
    <input type="text" bind:value={vaultPath} placeholder="~/notes" />
    <button on:click={browse}>Browse</button>
  </div>

  {#if error}
    <p class="error">{error}</p>
  {/if}

  <button class="confirm" on:click={confirm} disabled={saving}>
    {saving ? 'Setting up...' : 'Get Started'}
  </button>
</div>

<style>
  .setup {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    font-family: monospace;
    background: #1a1b26;
    color: #a9b1d6;
  }
  h1 { color: #7aa2f7; margin-bottom: 0.5rem; }
  .subtitle { color: #565f89; margin-bottom: 2rem; }
  .input-row {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  input {
    background: #24283b;
    border: 1px solid #3b4261;
    color: #a9b1d6;
    padding: 0.5rem 1rem;
    font-family: monospace;
    font-size: 1rem;
    width: 400px;
    border-radius: 4px;
  }
  input:focus { border-color: #7aa2f7; outline: none; }
  button {
    background: #3b4261;
    color: #a9b1d6;
    border: none;
    padding: 0.5rem 1rem;
    font-family: monospace;
    cursor: pointer;
    border-radius: 4px;
  }
  button:hover { background: #414868; }
  .confirm {
    background: #7aa2f7;
    color: #1a1b26;
    font-weight: bold;
    padding: 0.75rem 2rem;
    margin-top: 1rem;
  }
  .confirm:hover { background: #89b4fa; }
  .confirm:disabled { opacity: 0.5; cursor: default; }
  .error { color: #f7768e; margin: 0; }
</style>
```

- [ ] **Step 3: Modify App.svelte to gate on setup**

In `App.svelte`, add setup gate:

```svelte
<script>
  import { onMount } from 'svelte';
  import { NeedsSetup } from '../wailsjs/go/wailsplugin/App';
  import SetupWizard from './components/SetupWizard.svelte';
  // ... existing imports

  let needsSetup = null; // null = loading, true/false = known

  onMount(async () => {
    needsSetup = await NeedsSetup();
  });

  function onSetupComplete() {
    needsSetup = false;
  }
</script>

{#if needsSetup === null}
  <!-- Loading -->
{:else if needsSetup}
  <SetupWizard on:complete={onSetupComplete} />
{:else}
  <!-- existing app content -->
{/if}
```

- [ ] **Step 4: Build frontend to verify**

Run: `cd cmd/lm-desktop/frontend && npm run build`
Expected: SUCCESS

- [ ] **Step 5: Commit**

```bash
git add cmd/lm-desktop/frontend/src/components/SetupWizard.svelte cmd/lm-desktop/frontend/src/App.svelte
git commit -m "feat(desktop): add Svelte setup wizard component"
```

---

### Task 8: Integration test — full setup flow

**Files:**
- Create: `internal/config/integration_test.go`

- [ ] **Step 1: Write integration test for config round-trip**

```go
// internal/config/integration_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFullSetupFlow(t *testing.T) {
	// Simulate first-run
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !NeedsSetup(cfg) {
		t.Fatal("expected NeedsSetup on fresh install")
	}

	// User selects vault
	vault := filepath.Join(home, "my-notes")
	os.MkdirAll(vault, 0755)

	expanded, err := ExpandPath(vault)
	if err != nil {
		t.Fatalf("ExpandPath: %v", err)
	}
	cfg.VaultPath = expanded
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Second launch — should not need setup
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if NeedsSetup(cfg2) {
		t.Fatal("expected NeedsSetup=false after save")
	}
	if cfg2.VaultPath != expanded {
		t.Fatalf("VaultPath=%q, want %q", cfg2.VaultPath, expanded)
	}

	// Verify config file exists at expected path
	configPath := filepath.Join(home, ".config", "lazymd", "config.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file missing: %v", err)
	}
}

func TestVaultFlagOverride(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Save initial config
	cfg := &Config{VaultPath: "/old/vault"}
	Save(cfg)

	// Override with --vault equivalent
	newVault := filepath.Join(home, "new-vault")
	os.MkdirAll(newVault, 0755)

	cfg.VaultPath = newVault
	Save(cfg)

	loaded, _ := Load()
	if loaded.VaultPath != newVault {
		t.Fatalf("got %q, want %q", loaded.VaultPath, newVault)
	}
}
```

- [ ] **Step 2: Run all tests**

Run: `export PATH="/usr/local/go/bin:$PATH" && go test ./... -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/config/integration_test.go
git commit -m "test(config): add integration tests for full setup flow"
```
