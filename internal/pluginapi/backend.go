package pluginapi

// BackendPlugin extends core capabilities: AI, export, sync, indexing, VCS.
// It receives operations and events, and can register new capabilities.
type BackendPlugin interface {
	Info() PluginInfo
	Init(ctx *BackendContext) error
	Shutdown() error

	// OnOperation is called after passes, when an operation reaches backends.
	OnOperation(ctx *BackendContext, op *Operation)

	// OnEvent is called when the core broadcasts an event.
	OnEvent(ctx *BackendContext, event *Event)

	// Commands returns : commands this plugin registers.
	Commands() []CommandDef

	// Capabilities declares what this plugin extends.
	Capabilities() []Capability
}

// BackendContext is the API surface available to backend plugins.
type BackendContext struct {
	Emit      func(op *Operation)    // emit new operations into the IR
	Broadcast func(event *Event)     // broadcast events to all plugins
	Editor    EditorAPI
	Nav       NavAPI  // nil if unavailable
	Brain     BrainAPI // nil if unavailable
	Config    ConfigAPI
	Log       func(msg string)
}

// Capability declares a named capability a backend plugin provides.
type Capability struct {
	Name        string // e.g. "ai.complete", "export.pdf", "vcs.git"
	Description string
}

// Pass transforms operations before they reach the core.
// Inspired by LLVM optimization passes.
type Pass interface {
	Name() string
	Priority() int // lower runs first (0-100)

	// Transform processes an operation. Return false to drop it entirely.
	Transform(ctx *BackendContext, op *Operation) bool
}
