package main

import (
	"context"
	"embed"
	"log"

	"github.com/EME130/lazymd/internal/brain"
	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/editor"
	"github.com/EME130/lazymd/internal/nav"
	"github.com/EME130/lazymd/internal/pluginadapter"
	"github.com/EME130/lazymd/internal/pluginapi"
	"github.com/EME130/lazymd/internal/plugins"
	"github.com/EME130/lazymd/internal/wailsplugin"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	buf := buffer.New()
	g := brain.NewGraph()
	ed := editor.New(buf)
	navigator := nav.NewBuiltin(buf)

	editorAdapter := &pluginadapter.EditorAdapter{Ed: ed}
	navAdapter := &pluginadapter.NavAdapter{Nav: navigator}
	brainAdapter := &pluginadapter.BrainAdapter{Graph: g}
	themeAdapter := &pluginadapter.ThemeAdapter{}

	engine := plugins.NewEngine()

	frontendCtx := &pluginapi.FrontendContext{
		Emit:   engine.Emit,
		Editor: editorAdapter,
		Nav:    navAdapter,
		Brain:  brainAdapter,
		Theme:  themeAdapter,
		Log:    func(msg string) { log.Println("[plugin]", msg) },
	}
	backendCtx := &pluginapi.BackendContext{
		Emit:   engine.Emit,
		Editor: editorAdapter,
		Nav:    navAdapter,
		Brain:  brainAdapter,
		Log:    func(msg string) { log.Println("[plugin]", msg) },
	}
	engine.SetContexts(frontendCtx, backendCtx)

	guiPlugin := wailsplugin.NewPlugin()
	engine.RegisterFrontend(guiPlugin)

	app := wailsplugin.NewApp(frontendCtx)

	err := wails.Run(&options.App{
		Title:  "LazyMD",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			guiPlugin.SetWailsRuntime(ctx, runtime.EventsEmit)
			engine.InitAll()
		},
		OnShutdown: func(_ context.Context) {
			engine.ShutdownAll()
		},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
