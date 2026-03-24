package wailsplugin

import (
	"context"
	"testing"

	"github.com/EME130/lazymd/internal/pluginapi"
)

func TestPluginInfo(t *testing.T) {
	p := NewPlugin()
	info := p.Info()
	if info.Name != "wails-gui" {
		t.Errorf("expected name 'wails-gui', got %q", info.Name)
	}
}

func TestPluginInit(t *testing.T) {
	p := NewPlugin()
	ctx := &pluginapi.FrontendContext{
		Emit: func(op *pluginapi.Operation) {},
		Log:  func(msg string) {},
	}
	if err := p.Init(ctx); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestPluginRenderEmpty(t *testing.T) {
	p := NewPlugin()
	if got := p.Render(80, 24); got != "" {
		t.Errorf("expected empty render, got %q", got)
	}
}

func TestEventMapping(t *testing.T) {
	tests := []struct {
		in  pluginapi.EventType
		out string
	}{
		{pluginapi.EventBufferChanged, "buffer:changed"},
		{pluginapi.EventFileOpened, "file:opened"},
		{pluginapi.EventFileSaved, "file:saved"},
		{pluginapi.EventCursorMoved, "cursor:moved"},
		{pluginapi.EventGraphUpdated, "graph:updated"},
		{pluginapi.EventModeChanged, "mode:changed"},
	}
	for _, tt := range tests {
		got := WailsEventName(tt.in)
		if got != tt.out {
			t.Errorf("WailsEventName(%q) = %q, want %q", tt.in, got, tt.out)
		}
	}
}

func TestPluginOnEventForwards(t *testing.T) {
	p := NewPlugin()
	ctx := &pluginapi.FrontendContext{
		Emit: func(op *pluginapi.Operation) {},
		Log:  func(msg string) {},
	}
	p.Init(ctx)

	var emittedName string
	var emittedData map[string]any
	p.SetWailsRuntime(context.Background(), func(_ context.Context, name string, data ...interface{}) {
		emittedName = name
		if len(data) > 0 {
			emittedData, _ = data[0].(map[string]any)
		}
	})

	event := pluginapi.NewEvent(pluginapi.EventBufferChanged)
	event.Data["content"] = "hello"
	p.OnEvent(ctx, event)

	if emittedName != "buffer:changed" {
		t.Errorf("expected event name 'buffer:changed', got %q", emittedName)
	}
	if emittedData["content"] != "hello" {
		t.Errorf("expected data content 'hello', got %v", emittedData["content"])
	}
}

func TestPluginOnEventNoEmitter(t *testing.T) {
	p := NewPlugin()
	ctx := &pluginapi.FrontendContext{
		Emit: func(op *pluginapi.Operation) {},
		Log:  func(msg string) {},
	}
	p.Init(ctx)

	// Should not panic when emitter is nil
	event := pluginapi.NewEvent(pluginapi.EventBufferChanged)
	p.OnEvent(ctx, event)
}
