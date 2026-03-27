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
