package pluginapi

// PluginInfo describes a plugin.
type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
}

// CommandDef defines a plugin command available via `:command_name`.
type CommandDef struct {
	Name        string
	Description string
	Handler     func(ctx *BackendContext, args string) error
}

// FrontendCommandDef defines a frontend plugin command.
type FrontendCommandDef struct {
	Name        string
	Description string
	Handler     func(ctx *FrontendContext, args string) error
}
