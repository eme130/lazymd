package corebackend

import "github.com/EME130/lazymd/internal/pluginapi"

// Plugin is the core backend that translates IR operations into actual
// editor/buffer mutations and broadcasts resulting events.
type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (p *Plugin) Info() pluginapi.PluginInfo {
	return pluginapi.PluginInfo{
		Name:        "core",
		Version:     "0.1.0",
		Description: "Core operation handler",
	}
}

func (p *Plugin) Init(_ *pluginapi.BackendContext) error { return nil }
func (p *Plugin) Shutdown() error                        { return nil }
func (p *Plugin) Commands() []pluginapi.CommandDef       { return nil }
func (p *Plugin) Capabilities() []pluginapi.Capability   { return nil }
func (p *Plugin) OnEvent(_ *pluginapi.BackendContext, _ *pluginapi.Event) {}

func (p *Plugin) OnOperation(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	switch op.Type {
	case pluginapi.OpOpenFile:
		p.handleOpenFile(ctx, op)
	case pluginapi.OpSaveFile:
		p.handleSaveFile(ctx, op)
	case pluginapi.OpInsertText:
		p.handleInsertText(ctx, op)
	case pluginapi.OpDeleteRange:
		p.handleDeleteRange(ctx, op)
	case pluginapi.OpSetCursor:
		p.handleSetCursor(ctx, op)
	}
}

func (p *Plugin) handleOpenFile(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	path, _ := op.Data["path"].(string)
	if path == "" {
		return
	}
	if err := ctx.Editor.OpenFile(path); err != nil {
		ctx.Log("open_file error: " + err.Error())
		return
	}
	evt := pluginapi.NewEvent(pluginapi.EventFileOpened)
	evt.Data["path"] = path
	ctx.Broadcast(evt)

	bufEvt := pluginapi.NewEvent(pluginapi.EventBufferChanged)
	bufEvt.Data["content"] = ctx.Editor.Content()
	bufEvt.Data["origin"] = op.Origin
	ctx.Broadcast(bufEvt)
}

func (p *Plugin) handleSaveFile(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	if err := ctx.Editor.SaveFile(); err != nil {
		ctx.Log("save_file error: " + err.Error())
		return
	}
	evt := pluginapi.NewEvent(pluginapi.EventFileSaved)
	evt.Data["path"] = ctx.Editor.FilePath()
	ctx.Broadcast(evt)
}

func (p *Plugin) handleInsertText(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	row := toInt(op.Data["row"])
	col := toInt(op.Data["col"])
	text, _ := op.Data["text"].(string)
	if text == "" {
		return
	}
	ctx.Editor.InsertAt(row, col, text)

	evt := pluginapi.NewEvent(pluginapi.EventBufferChanged)
	evt.Data["content"] = ctx.Editor.Content()
	evt.Data["origin"] = op.Origin
	ctx.Broadcast(evt)
}

func (p *Plugin) handleDeleteRange(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	startRow := toInt(op.Data["start_row"])
	startCol := toInt(op.Data["start_col"])
	endRow := toInt(op.Data["end_row"])
	endCol := toInt(op.Data["end_col"])
	ctx.Editor.DeleteRange(startRow, startCol, endRow, endCol)

	evt := pluginapi.NewEvent(pluginapi.EventBufferChanged)
	evt.Data["content"] = ctx.Editor.Content()
	evt.Data["origin"] = op.Origin
	ctx.Broadcast(evt)
}

func (p *Plugin) handleSetCursor(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {
	row := toInt(op.Data["row"])
	col := toInt(op.Data["col"])
	ctx.Editor.SetCursorRow(row)
	ctx.Editor.SetCursorCol(col)

	evt := pluginapi.NewEvent(pluginapi.EventCursorMoved)
	evt.Data["row"] = row
	evt.Data["col"] = col
	ctx.Broadcast(evt)
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	default:
		return 0
	}
}
