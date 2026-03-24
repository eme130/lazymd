package plugins

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	content := `{
		"name": "test-plugin",
		"version": "1.0.0",
		"tier": "backend",
		"entry": "./test-plugin",
		"commands": [{"name": "test", "description": "test command"}],
		"events": ["file_saved"],
		"capabilities": ["test.cap"]
	}`
	os.WriteFile(path, []byte(content), 0644)

	m, err := ParseManifest(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "test-plugin" {
		t.Errorf("expected name 'test-plugin', got %q", m.Name)
	}
	if m.Tier != "backend" {
		t.Errorf("expected tier 'backend', got %q", m.Tier)
	}
	if len(m.Commands) != 1 {
		t.Errorf("expected 1 command, got %d", len(m.Commands))
	}
	if len(m.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(m.Events))
	}
}

func TestParseManifestDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	content := `{"name": "minimal", "entry": "./minimal"}`
	os.WriteFile(path, []byte(content), 0644)

	m, err := ParseManifest(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Tier != "backend" {
		t.Errorf("expected default tier 'backend', got %q", m.Tier)
	}
}

func TestParseManifestErrors(t *testing.T) {
	dir := t.TempDir()

	// Missing name
	path := filepath.Join(dir, "no-name.json")
	os.WriteFile(path, []byte(`{"entry": "./x"}`), 0644)
	if _, err := ParseManifest(path); err == nil {
		t.Error("expected error for missing name")
	}

	// Missing entry
	path = filepath.Join(dir, "no-entry.json")
	os.WriteFile(path, []byte(`{"name": "x"}`), 0644)
	if _, err := ParseManifest(path); err == nil {
		t.Error("expected error for missing entry")
	}

	// Invalid tier
	path = filepath.Join(dir, "bad-tier.json")
	os.WriteFile(path, []byte(`{"name": "x", "entry": "./x", "tier": "invalid"}`), 0644)
	if _, err := ParseManifest(path); err == nil {
		t.Error("expected error for invalid tier")
	}

	// Missing file
	if _, err := ParseManifest(filepath.Join(dir, "nonexistent.json")); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestResolveEntry(t *testing.T) {
	m := &Manifest{Entry: "./my-plugin"}
	resolved := m.ResolveEntry("/home/user/.config/lazymd/plugins/my-plugin")
	expected := "/home/user/.config/lazymd/plugins/my-plugin/my-plugin"
	if resolved != expected {
		t.Errorf("expected %q, got %q", expected, resolved)
	}

	m2 := &Manifest{Entry: "/usr/local/bin/plugin"}
	resolved2 := m2.ResolveEntry("/some/dir")
	if resolved2 != "/usr/local/bin/plugin" {
		t.Errorf("absolute path should stay absolute, got %q", resolved2)
	}
}
