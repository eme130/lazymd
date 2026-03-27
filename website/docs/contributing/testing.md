---
title: Running Tests
sidebar_position: 3
description: How to run the LazyMD test suite with go test. Unit tests cover buffer, markdown, nav, brain, highlight, mcp, editor, pluginapi, plugins, pluginadapter, themes, ui, web, agent, and wailsplugin packages.
keywords: [LazyMD tests, go test, unit tests, test suite, Go testing]
---

# Running Tests

```bash
task test
# or directly:
go test ./...
```

This runs all unit tests across all packages (buffer, markdown, nav, brain, highlight, mcp, editor, pluginapi, plugins, pluginadapter, themes, ui, web, agent, and wailsplugin).
