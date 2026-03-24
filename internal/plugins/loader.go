package plugins

import (
	"os"
	"path/filepath"
)

const pluginDir = ".config/lazymd/plugins"

// DiscoverExternal scans ~/.config/lazymd/plugins/ for external plugins
// and registers them with the engine.
func DiscoverExternal(eng *Engine) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil // silently skip if no home dir
	}

	dir := filepath.Join(home, pluginDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // no plugin dir is fine
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(dir, entry.Name(), "manifest.json")
		manifest, err := ParseManifest(manifestPath)
		if err != nil {
			continue // skip invalid plugins
		}

		entryPath := manifest.ResolveEntry(filepath.Join(dir, entry.Name()))

		// Check that the entry executable exists
		if _, err := os.Stat(entryPath); err != nil {
			continue
		}

		plugin := NewExternalPlugin(manifest, entryPath)

		switch manifest.Tier {
		case "backend":
			eng.RegisterBackend(plugin)
		case "frontend":
			// External frontend plugins: future work
			// For now, register as backend (commands still work)
			eng.RegisterBackend(plugin)
		}
	}

	return nil
}
