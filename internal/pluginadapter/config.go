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
