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
	errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(c.ErrColor))

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
