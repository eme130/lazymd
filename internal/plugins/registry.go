package plugins

import "github.com/EME130/lazymd/internal/pluginapi"

// AllFrontends returns all compiled-in frontend plugins.
func AllFrontends() []pluginapi.FrontendPlugin {
	return nil
}

// AllBackends returns all compiled-in backend plugins.
func AllBackends() []pluginapi.BackendPlugin {
	return nil
}
