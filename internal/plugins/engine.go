package plugins

import (
	"sort"

	"github.com/EME130/lazymd/internal/editor"
	"github.com/EME130/lazymd/internal/pluginapi"
)

// Engine is the central plugin orchestrator. It manages frontend and backend
// plugins, runs passes, and routes operations through the IR pipeline.
//
// It also implements editor.CommandExecutor for backward compatibility.
type Engine struct {
	frontends  []pluginapi.FrontendPlugin
	backends   []pluginapi.BackendPlugin
	passes     []pluginapi.Pass
	commandMap map[string]cmdEntry

	frontendCtx *pluginapi.FrontendContext
	backendCtx  *pluginapi.BackendContext
}

type cmdEntry struct {
	pluginName string
	handler    func(args string) error
}

// NewEngine creates a new plugin Engine.
func NewEngine() *Engine {
	return &Engine{
		commandMap: make(map[string]cmdEntry),
	}
}

// SetContexts sets the frontend and backend contexts used by plugins.
// Must be called before InitAll.
func (e *Engine) SetContexts(fc *pluginapi.FrontendContext, bc *pluginapi.BackendContext) {
	e.frontendCtx = fc
	e.backendCtx = bc
}

// RegisterFrontend adds a frontend plugin.
func (e *Engine) RegisterFrontend(p pluginapi.FrontendPlugin) {
	e.frontends = append(e.frontends, p)
	info := p.Info()
	for _, cmd := range p.Commands() {
		name := cmd.Name
		handler := cmd.Handler
		ctx := e.frontendCtx
		e.commandMap[name] = cmdEntry{
			pluginName: info.Name,
			handler:    func(args string) error { return handler(ctx, args) },
		}
	}
}

// RegisterBackend adds a backend plugin.
func (e *Engine) RegisterBackend(p pluginapi.BackendPlugin) {
	e.backends = append(e.backends, p)
	info := p.Info()
	for _, cmd := range p.Commands() {
		name := cmd.Name
		handler := cmd.Handler
		ctx := e.backendCtx
		e.commandMap[name] = cmdEntry{
			pluginName: info.Name,
			handler:    func(args string) error { return handler(ctx, args) },
		}
	}
}

// RegisterPass adds a pass to the pipeline, maintaining priority order.
func (e *Engine) RegisterPass(p pluginapi.Pass) {
	e.passes = append(e.passes, p)
	sort.Slice(e.passes, func(i, j int) bool {
		return e.passes[i].Priority() < e.passes[j].Priority()
	})
}

// InitAll initializes all registered plugins.
func (e *Engine) InitAll() error {
	for _, f := range e.frontends {
		if err := f.Init(e.frontendCtx); err != nil {
			return err
		}
	}
	for _, b := range e.backends {
		if err := b.Init(e.backendCtx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAll gracefully stops all plugins.
func (e *Engine) ShutdownAll() {
	for _, f := range e.frontends {
		f.Shutdown()
	}
	for _, b := range e.backends {
		b.Shutdown()
	}
}

// Emit sends an operation through the IR pipeline:
// passes → core → backends.
func (e *Engine) Emit(op *pluginapi.Operation) {
	// 1. Run passes
	for _, pass := range e.passes {
		if !pass.Transform(e.backendCtx, op) {
			return
		}
		if op.Canceled {
			return
		}
	}

	// 2. Notify backends
	for _, b := range e.backends {
		b.OnOperation(e.backendCtx, op)
	}
}

// BroadcastEvent sends an event to all frontend and backend plugins.
func (e *Engine) BroadcastEvent(event *pluginapi.Event) {
	for _, b := range e.backends {
		b.OnEvent(e.backendCtx, event)
	}
}

// --- editor.CommandExecutor implementation (backward compat) ---

func (e *Engine) ExecuteCommand(name string, _ editor.PluginEditor, args string) bool {
	entry, ok := e.commandMap[name]
	if !ok {
		return false
	}
	entry.handler(args)
	return true
}

func (e *Engine) Broadcast(eventType string, _ editor.PluginEditor) {
	e.BroadcastEvent(&pluginapi.Event{
		Type: pluginapi.EventType(eventType),
		Data: make(map[string]any),
	})
}

func (e *Engine) ListPlugins() []editor.PluginSummary {
	var out []editor.PluginSummary
	for _, f := range e.frontends {
		info := f.Info()
		out = append(out, editor.PluginSummary{
			Name:        info.Name,
			Version:     info.Version,
			Description: info.Description,
		})
	}
	for _, b := range e.backends {
		info := b.Info()
		out = append(out, editor.PluginSummary{
			Name:        info.Name,
			Version:     info.Version,
			Description: info.Description,
		})
	}
	return out
}

func (e *Engine) ListCommands() []editor.CommandSummary {
	var out []editor.CommandSummary
	for name, entry := range e.commandMap {
		out = append(out, editor.CommandSummary{
			Name:       name,
			PluginName: entry.pluginName,
		})
	}
	return out
}

// --- Query methods ---

// Frontends returns all registered frontend plugins.
func (e *Engine) Frontends() []pluginapi.FrontendPlugin {
	return e.frontends
}

// Backends returns all registered backend plugins.
func (e *Engine) Backends() []pluginapi.BackendPlugin {
	return e.backends
}

// Passes returns all registered passes.
func (e *Engine) Passes() []pluginapi.Pass {
	return e.passes
}

// FrontendCount returns the number of frontend plugins.
func (e *Engine) FrontendCount() int {
	return len(e.frontends)
}

// BackendCount returns the number of backend plugins.
func (e *Engine) BackendCount() int {
	return len(e.backends)
}

// Capabilities returns all capabilities from backend plugins.
func (e *Engine) Capabilities() []pluginapi.Capability {
	var caps []pluginapi.Capability
	for _, b := range e.backends {
		caps = append(caps, b.Capabilities()...)
	}
	return caps
}
