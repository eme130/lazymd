---
title: Plugin Catalog
sidebar_position: 2
description: LazyMD plugin catalog — external plugins, community contributions, and planned built-in plugins.
keywords: [LazyMD plugins, plugin catalog, external plugins, community plugins]
---

# Plugin Catalog

lm's plugin engine is ready for use but ships with **no built-in plugins yet**. The engine supports both Go-native plugins and external plugins (any language, communicating over JSON-RPC 2.0 stdio).

## Writing Your Own

See the [Development Guide](./development.md) for how to create plugins. External plugins live in `~/.config/lazymd/plugins/<name>/` with a `manifest.json`.

## Planned Built-In Plugins

The following plugins are planned for future releases:

| Plugin | Type | Description |
|--------|------|-------------|
| word-count | Backend | Word, line, and character counts |
| outline | Frontend | Document outline / table of contents panel |
| zen-mode | Frontend | Distraction-free writing mode |
| pomodoro | Frontend | Focus timer in the status bar |
| git-sync | Backend | Git-based note synchronization |
| export-html | Backend | Export to HTML and PDF |

Community contributions welcome! If you build a plugin, open a PR to get it listed here.
