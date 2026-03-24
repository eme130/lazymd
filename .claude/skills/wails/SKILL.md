---
name: wails
description: Use when building desktop applications with Wails (Go + web frontend), scaffolding Wails projects, binding Go structs to frontend, configuring wails.json, using runtime APIs (events, dialogs, windows, menus), or troubleshooting Wails builds.
---

# Wails - Desktop Apps with Go + Web Frontend

## Overview

Wails builds cross-platform desktop apps using Go for the backend and web technologies (HTML/CSS/JS) for the frontend, packaged into a single binary. It uses the OS native webview (not Electron), resulting in small binaries and low memory usage.

**Import path:** `github.com/wailsapp/wails/v2`

## CLI Commands

```bash
# Install CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify environment
wails doctor

# Create project (templates: react-ts, vue-ts, svelte-ts, vanilla, preact, lit)
wails init -n myapp -t react-ts

# Dev mode with hot reload
wails dev

# Build production binary
wails build

# Build with optimizations
wails build -clean -upx -ldflags "-w -s"

# Cross-compile (e.g. Windows from macOS)
wails build -platform windows/amd64
```

## Project Structure

```
myapp/
  build/              # Build artifacts, appicon.png, info.plist
  frontend/           # Web frontend (React/Vue/Svelte/etc.)
    src/
    dist/             # Built frontend assets (embedded in binary)
    package.json
    wailsjs/          # Auto-generated JS bindings
      go/main/App.js  # Generated from bound Go structs
      runtime/         # Wails runtime JS API
  app.go              # Application struct + bound methods
  main.go             # Entry point with wails.Run() config
  wails.json          # Project configuration
  go.mod
```

## Core Pattern: main.go

```go
package main

import (
    "context"
    "embed"
    "log"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    app := NewApp()

    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup:  app.startup,
        OnDomReady: app.domReady,
        OnShutdown: app.shutdown,
        Bind: []interface{}{
            app,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Core Pattern: app.go (Binding Go Methods)

```go
package main

import (
    "context"
    "fmt"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx context.Context
}

func NewApp() *App {
    return &App{}
}

// Lifecycle hooks - called by Wails, NOT exposed to frontend
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx // Save context for runtime calls
}

// Public methods ARE exposed to frontend via auto-generated JS bindings
func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s!", name)
}
```

**Binding rules:**
- Add struct instances to the `Bind` slice in `options.App`
- Only **public methods** (uppercase) are exposed to frontend
- `startup`, `domReady`, `beforeClose`, `shutdown` are lifecycle hooks, not bound
- Multiple structs can be bound: `Bind: []interface{}{app, &svc1{}, &svc2{}}`
- Frontend calls: `import {Greet} from '../wailsjs/go/main/App'` (auto-generated)

## Frontend Calling Go

```javascript
// Auto-generated bindings in frontend/wailsjs/go/main/App.js
import { Greet } from '../wailsjs/go/main/App';

const result = await Greet("World"); // Returns "Hello World!"
```

## Events (Go <-> Frontend)

```go
// Go: emit event to frontend
runtime.EventsEmit(a.ctx, "eventName", data)

// Go: listen for frontend event
runtime.EventsOn(a.ctx, "fromFrontend", func(args ...interface{}) {
    message := args[0].(string)
})
```

```javascript
// JS: listen for Go event
import { EventsOn, EventsEmit } from '../wailsjs/runtime/runtime';

EventsOn("eventName", (data) => { console.log(data); });

// JS: emit event to Go
EventsEmit("fromFrontend", "hello from JS");
```

## Runtime APIs (Go)

All require the saved `context.Context` from `startup()`:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

// Dialogs
result, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
    Type:    runtime.QuestionDialog,
    Title:   "Confirm",
    Message: "Are you sure?",
})

file, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
    Title: "Select File",
    Filters: []runtime.FileFilter{
        {DisplayName: "Go Files", Pattern: "*.go"},
    },
})

// Window management
runtime.WindowSetTitle(a.ctx, "New Title")
runtime.WindowSetSize(a.ctx, 1200, 800)
runtime.WindowMinimise(a.ctx)
runtime.WindowMaximise(a.ctx)
runtime.WindowToggleMaximise(a.ctx)
runtime.WindowFullscreen(a.ctx)
runtime.WindowCenter(a.ctx)
runtime.WindowHide(a.ctx)
runtime.WindowShow(a.ctx)

// Clipboard
runtime.ClipboardSetText(a.ctx, "copied text")
text, _ := runtime.ClipboardGetText(a.ctx)

// System
runtime.BrowserOpenURL(a.ctx, "https://example.com")
env := runtime.Environment(a.ctx)
runtime.Quit(a.ctx)
```

## Application Options Reference

```go
&options.App{
    Title:             "App Name",
    Width:             1024,
    Height:            768,
    MinWidth:          400,
    MinHeight:         400,
    MaxWidth:          1920,
    MaxHeight:         1080,
    DisableResize:     false,
    Fullscreen:        false,
    Frameless:         false,        // Borderless window
    StartHidden:       false,
    HideWindowOnClose: false,        // Hide instead of quit on close
    AlwaysOnTop:       false,
    BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 255},
    WindowStartState:  options.Normal, // Normal, Maximised, Minimised, Fullscreen

    // Drag and drop
    DragAndDrop: &options.DragAndDrop{
        EnableFileDrop:  true,
        CSSDropProperty: "--wails-drop-target",
        CSSDropValue:    "drop",
    },

    // CSS-based window dragging
    CSSDragProperty: "--wails-draggable",
    CSSDragValue:    "drag",

    // macOS specific
    Mac: &mac.Options{
        TitleBar: &mac.TitleBar{
            TitlebarAppearsTransparent: true,
            HideTitle:                  true,
            FullSizeContent:            true,
            UseToolbar:                 true,
        },
        Appearance: mac.NSAppearanceNameDarkAqua,
        About: &mac.AboutInfo{
            Title:   "My App",
            Message: "v1.0.0",
        },
    },

    // Windows specific
    Windows: &windows.Options{
        WebviewIsTransparent: false,
        WindowIsTranslucent:  false,
        BackdropType:         windows.Mica,
        Theme:                windows.SystemDefault,
    },

    // Single instance lock
    SingleInstanceLock: &options.SingleInstanceLock{
        UniqueId: "unique-app-id-here",
        OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
            // Handle second instance
        },
    },
}
```

## wails.json Configuration

```json
{
  "name": "myapp",
  "outputfilename": "myapp",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "Your Name",
    "email": "you@example.com"
  }
}
```

## Menus

```go
import "github.com/wailsapp/wails/v2/pkg/menu"
import "github.com/wailsapp/wails/v2/pkg/menu/keys"

appMenu := menu.NewMenu()
fileMenu := appMenu.AddSubmenu("File")
fileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(cd *menu.CallbackData) {
    // handle open
})
fileMenu.AddSeparator()
fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(cd *menu.CallbackData) {
    runtime.Quit(a.ctx)
})

// Set in options
&options.App{
    Menu: appMenu,
}

// Update menu at runtime
runtime.MenuSetApplicationMenu(a.ctx, newMenu)
runtime.MenuUpdateApplicationMenu(a.ctx)
```

## Common Patterns

**Frameless window with CSS drag:**
```css
/* Make header draggable */
.titlebar {
    --wails-draggable: drag;
    height: 30px;
}
.titlebar button {
    --wails-draggable: no-drag; /* Buttons stay clickable */
}
```

**File drop handling:**
```go
&options.App{
    DragAndDrop: &options.DragAndDrop{
        EnableFileDrop: true,
    },
    OnDomReady: func(ctx context.Context) {
        runtime.EventsOn(ctx, "wails:file-drop", func(args ...interface{}) {
            files := args[0].([]string)
            // handle dropped files
        })
    },
}
```

## Troubleshooting

| Issue | Fix |
|---|---|
| `wails doctor` fails | Install deps: Xcode CLI tools (macOS), gcc+webkit2gtk (Linux) |
| Frontend changes not reflected | Run `wails dev` (not `go run`), check `frontend:build` in wails.json |
| Binding not generated | Run `wails generate module`, ensure methods are public |
| Build fails on Linux | `sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev` |
| Window blank on start | Check `//go:embed all:frontend/dist` path matches actual build output |
| Context nil in methods | Save `ctx` in `startup()`, use it in all runtime calls |

## When to Use context7

For the latest Wails API docs, use the context7 MCP tool with library ID `/wailsapp/wails`. Good queries:
- "Wails v3 application service API" (for v3 alpha)
- "system tray Wails"
- "Wails custom asset handler middleware"
- "Wails enum binding"
