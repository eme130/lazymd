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
