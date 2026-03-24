package plugins

import (
	"testing"

	"github.com/EME130/lazymd/internal/pluginapi"
)

func TestAllFrontendsEmpty(t *testing.T) {
	all := AllFrontends()
	if len(all) != 0 {
		t.Errorf("expected 0 frontends, got %d", len(all))
	}
}

func TestAllBackendsEmpty(t *testing.T) {
	all := AllBackends()
	if len(all) != 0 {
		t.Errorf("expected 0 backends, got %d", len(all))
	}
}

func TestEngineBasics(t *testing.T) {
	eng := NewEngine()

	if eng.FrontendCount() != 0 {
		t.Errorf("expected 0 frontends, got %d", eng.FrontendCount())
	}
	if eng.BackendCount() != 0 {
		t.Errorf("expected 0 backends, got %d", eng.BackendCount())
	}

	// ExecuteCommand should return false for unknown commands
	if eng.ExecuteCommand("nonexistent", nil, "") {
		t.Error("nonexistent command should not be found")
	}
}

func TestEngineEmit(t *testing.T) {
	eng := NewEngine()
	ctx := &pluginapi.BackendContext{}
	eng.SetContexts(&pluginapi.FrontendContext{}, ctx)

	// Emit with no plugins should not panic
	op := pluginapi.NewOperation(pluginapi.OpSaveFile, "test")
	eng.Emit(op)
}

func TestEngineBroadcastEvent(t *testing.T) {
	eng := NewEngine()
	eng.SetContexts(&pluginapi.FrontendContext{}, &pluginapi.BackendContext{})

	// Broadcast with no plugins should not panic
	eng.BroadcastEvent(pluginapi.NewEvent(pluginapi.EventFileSaved))
}

func TestEngineListPluginsEmpty(t *testing.T) {
	eng := NewEngine()
	if len(eng.ListPlugins()) != 0 {
		t.Error("expected empty plugin list")
	}
	if len(eng.ListCommands()) != 0 {
		t.Error("expected empty command list")
	}
}

func TestEnginePassCancellation(t *testing.T) {
	eng := NewEngine()
	eng.SetContexts(&pluginapi.FrontendContext{}, &pluginapi.BackendContext{})

	// Register a pass that cancels all operations
	eng.RegisterPass(&cancelPass{})

	called := false
	eng.RegisterBackend(&spyBackend{onOp: func() { called = true }})
	eng.backendCtx = &pluginapi.BackendContext{} // re-set after registering

	op := pluginapi.NewOperation(pluginapi.OpSaveFile, "test")
	eng.Emit(op)

	if !op.Canceled {
		t.Error("expected operation to be canceled")
	}
	if called {
		t.Error("backend should not have been called after pass cancellation")
	}
}

// --- test helpers ---

type cancelPass struct{}

func (p *cancelPass) Name() string                                                  { return "cancel" }
func (p *cancelPass) Priority() int                                                 { return 0 }
func (p *cancelPass) Transform(_ *pluginapi.BackendContext, op *pluginapi.Operation) bool {
	op.Cancel()
	return true
}

type spyBackend struct {
	onOp func()
}

func (s *spyBackend) Info() pluginapi.PluginInfo { return pluginapi.PluginInfo{Name: "spy"} }
func (s *spyBackend) Init(_ *pluginapi.BackendContext) error   { return nil }
func (s *spyBackend) Shutdown() error                          { return nil }
func (s *spyBackend) OnOperation(_ *pluginapi.BackendContext, _ *pluginapi.Operation) {
	if s.onOp != nil {
		s.onOp()
	}
}
func (s *spyBackend) OnEvent(_ *pluginapi.BackendContext, _ *pluginapi.Event) {}
func (s *spyBackend) Commands() []pluginapi.CommandDef                        { return nil }
func (s *spyBackend) Capabilities() []pluginapi.Capability                    { return nil }
