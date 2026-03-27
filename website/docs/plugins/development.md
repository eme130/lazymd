---
title: Developing Plugins
sidebar_position: 3
description: How to create plugins for LazyMD — frontend plugins (UI), backend plugins (core extensions), and external plugins (any language via JSON-RPC).
keywords: [LazyMD plugin development, Go plugin, frontend plugin, backend plugin, external plugin, plugin tutorial]
---

# Developing Plugins

See the full [Plugin Development Guide](https://github.com/EME130/lazymd/blob/main/docs/PLUGIN_GUIDE.md) for detailed instructions.

## Frontend Plugin (UI)

Frontend plugins own UI panels, keybindings, and status bar widgets.

```go
package myplugin

import "github.com/EME130/lazymd/internal/pluginapi"

type MyPanel struct{}

func (p *MyPanel) Info() pluginapi.PluginInfo {
    return pluginapi.PluginInfo{
        Name:        "my-panel",
        Version:     "0.1.0",
        Description: "A custom panel plugin",
    }
}

func (p *MyPanel) Init(ctx *pluginapi.FrontendContext) error {
    ctx.Log("my-panel loaded")
    return nil
}

func (p *MyPanel) Shutdown() error                                    { return nil }
func (p *MyPanel) OnEvent(ctx *pluginapi.FrontendContext, e *pluginapi.Event) {}
func (p *MyPanel) Render(width, height int) string                    { return "" }
func (p *MyPanel) KeyBindings() []pluginapi.KeyBinding                { return nil }
func (p *MyPanel) Commands() []pluginapi.FrontendCommandDef           { return nil }
func (p *MyPanel) StatusItems() []pluginapi.StatusItem                { return nil }
```

## Backend Plugin (Core Extension)

Backend plugins extend core capabilities (AI, export, sync).

```go
package myplugin

import "github.com/EME130/lazymd/internal/pluginapi"

type WordCounter struct{}

func (p *WordCounter) Info() pluginapi.PluginInfo {
    return pluginapi.PluginInfo{
        Name:        "word-count",
        Version:     "0.1.0",
        Description: "Word count command",
    }
}

func (p *WordCounter) Init(ctx *pluginapi.BackendContext) error { return nil }
func (p *WordCounter) Shutdown() error                         { return nil }

func (p *WordCounter) OnOperation(ctx *pluginapi.BackendContext, op *pluginapi.Operation) {}
func (p *WordCounter) OnEvent(ctx *pluginapi.BackendContext, e *pluginapi.Event)          {}

func (p *WordCounter) Commands() []pluginapi.CommandDef {
    return []pluginapi.CommandDef{{
        Name:        "wc",
        Description: "Show word count",
        Handler: func(ctx *pluginapi.BackendContext, args string) error {
            // count words using ctx.Editor.Content()
            return nil
        },
    }}
}

func (p *WordCounter) Capabilities() []pluginapi.Capability { return nil }
```

## External Plugin (Any Language)

External plugins run as child processes over JSON-RPC 2.0 stdio.

1. Create `~/.config/lazymd/plugins/my-plugin/manifest.json`:

```json
{
  "name": "my-plugin",
  "version": "0.1.0",
  "description": "My external plugin",
  "entry": "main.py",
  "runtime": "python3",
  "events": ["file_opened", "file_saved"],
  "commands": [
    {
      "name": "hello",
      "description": "Say hello"
    }
  ]
}
```

2. Implement the JSON-RPC handler in `main.py` (or any language).

## Registration

Register plugins with the engine in your application setup:

```go
engine := plugins.NewEngine()
engine.RegisterFrontend(&MyPanel{})
engine.RegisterBackend(&WordCounter{})
engine.RegisterPass(&MyPass{})
```
