package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Manifest describes an external plugin loaded from manifest.json.
type Manifest struct {
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	Author       string           `json:"author,omitempty"`
	Description  string           `json:"description,omitempty"`
	Tier         string           `json:"tier"` // "frontend" or "backend"
	Entry        string           `json:"entry"`
	Commands     []ManifestCmd    `json:"commands,omitempty"`
	Events       []string         `json:"events,omitempty"`
	Capabilities []string         `json:"capabilities,omitempty"`
	Operations   []string         `json:"operations,omitempty"`
}

// ManifestCmd describes a command in the manifest.
type ManifestCmd struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ParseManifest reads and parses a manifest.json file.
func ParseManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("manifest %s: missing name", path)
	}
	if m.Entry == "" {
		return nil, fmt.Errorf("manifest %s: missing entry", path)
	}
	if m.Tier == "" {
		m.Tier = "backend" // default to backend
	}
	if m.Tier != "frontend" && m.Tier != "backend" {
		return nil, fmt.Errorf("manifest %s: tier must be 'frontend' or 'backend'", path)
	}
	return &m, nil
}

// ResolveEntry resolves the entry path relative to the manifest directory.
func (m *Manifest) ResolveEntry(manifestDir string) string {
	if filepath.IsAbs(m.Entry) {
		return m.Entry
	}
	return filepath.Join(manifestDir, m.Entry)
}
