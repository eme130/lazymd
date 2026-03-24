package wailsplugin

import (
	"context"

	"github.com/EME130/lazymd/internal/pluginapi"
)

type Plugin struct {
	ctx     *pluginapi.FrontendContext
	wailsRT context.Context
	emitter EventEmitter
}

type EventEmitter func(ctx context.Context, eventName string, data ...interface{})

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) SetWailsRuntime(ctx context.Context, emitter EventEmitter) {
	p.wailsRT = ctx
	p.emitter = emitter
}

func (p *Plugin) Info() pluginapi.PluginInfo {
	return pluginapi.PluginInfo{
		Name:        "wails-gui",
		Version:     "0.1.0",
		Description: "Wails desktop GUI surface",
	}
}

func (p *Plugin) Init(ctx *pluginapi.FrontendContext) error {
	p.ctx = ctx
	return nil
}

func (p *Plugin) Shutdown() error {
	return nil
}

func (p *Plugin) OnEvent(_ *pluginapi.FrontendContext, event *pluginapi.Event) {
	if p.emitter == nil || p.wailsRT == nil {
		return
	}
	name := WailsEventName(event.Type)
	p.emitter(p.wailsRT, name, event.Data)
}

func (p *Plugin) Render(_, _ int) string                    { return "" }
func (p *Plugin) KeyBindings() []pluginapi.KeyBinding       { return nil }
func (p *Plugin) Commands() []pluginapi.FrontendCommandDef  { return nil }
func (p *Plugin) StatusItems() []pluginapi.StatusItem       { return nil }
