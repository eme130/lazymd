package pluginapi

// FrontendPlugin owns UI surface: panels, keybindings, status widgets.
// It produces operations that flow into the IR.
type FrontendPlugin interface {
	Info() PluginInfo
	Init(ctx *FrontendContext) error
	Shutdown() error

	// OnEvent is called when the core broadcasts an event.
	OnEvent(ctx *FrontendContext, event *Event)

	// Render returns content for a panel. Return "" if this plugin has no panel.
	Render(width, height int) string

	// KeyBindings returns keybindings this plugin registers.
	KeyBindings() []KeyBinding

	// Commands returns : commands this plugin registers.
	Commands() []FrontendCommandDef

	// StatusItems returns status bar widgets. Return nil if none.
	StatusItems() []StatusItem
}

// FrontendContext is the API surface available to frontend plugins.
type FrontendContext struct {
	Emit   func(op *Operation) // emit operation into the IR pipeline
	Editor EditorAPI
	Nav    NavAPI   // nil if unavailable
	Brain  BrainAPI // nil if unavailable
	Theme  ThemeAPI // nil if unavailable
	Log    func(msg string)
}

// KeyBinding maps a key combination to a handler in a given mode.
type KeyBinding struct {
	Mode    string // "normal", "insert", "*" (all modes)
	Key     string // e.g. "ctrl+p", "gd"
	Handler func(ctx *FrontendContext)
}

// StatusItem is a widget rendered in the status bar.
type StatusItem struct {
	Name     string
	Align    string // "left" or "right"
	Renderer func(ctx *FrontendContext) string
}
