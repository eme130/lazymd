package pluginadapter

import "github.com/EME130/lazymd/internal/themes"

// ThemeAdapter wraps the themes package as a pluginapi.ThemeAPI.
type ThemeAdapter struct{}

func (a *ThemeAdapter) CurrentName() string {
	return themes.Current().Name
}

func (a *ThemeAdapter) SetByName(name string) bool {
	return themes.SetByName(name)
}

func (a *ThemeAdapter) ListThemes() []string {
	out := make([]string, len(themes.BuiltinThemes))
	for i, t := range themes.BuiltinThemes {
		out[i] = t.Name
	}
	return out
}
