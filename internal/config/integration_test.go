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
