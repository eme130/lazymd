---
title: Building from Source
sidebar_position: 2
description: Build LazyMD from source — TUI binary, desktop GUI (Wails), and debug builds using go-task.
keywords: [build LazyMD, go build, compile from source, go project, Taskfile, Wails]
---

# Building from Source

lm uses [go-task](https://taskfile.dev) as its build system.

## TUI

```bash
task build            # Build TUI binary → bin/lm
task run              # Run TUI directly
task run -- myfile.md # Run with arguments
```

## Desktop (Wails)

```bash
task desktop:install  # Install frontend npm deps
task desktop:build    # Build desktop binary → bin/lm-desktop
task desktop:dev      # Wails dev mode (hot reload)
```

## Quality

```bash
task test             # Run all tests
task fmt              # Format Go code
task vet              # Run go vet
task lint             # Run all linters (vet + fmt)
```

## Debug

```bash
task debug            # Debug TUI with delve
task debug:desktop    # Debug desktop with delve
```

## MCP Server

```bash
task mcp              # Build and run as MCP server
task mcp -- myfile.md # Run MCP server with file preloaded
```

## Manual Build (without go-task)

```bash
go build -o lm ./cmd/lm                          # TUI
go build -ldflags="-s -w" -o lm ./cmd/lm         # Optimized
go run ./cmd/lm -- myfile.md                      # Run directly
```
