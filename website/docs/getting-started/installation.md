---
title: Installation
sidebar_position: 1
description: How to install LazyMD, the terminal-based markdown editor with vim keybindings. Build from source with Zig or download pre-built binaries for Linux and macOS.
keywords: [LazyMD install, terminal markdown editor install, zig markdown editor, LazyMD setup, build LazyMD]
---

# Installation

## Prerequisites

lm requires [Zig](https://ziglang.org/download/) version 0.15.1 or later.

## Build from source

```bash
git clone https://github.com/EME130/lazymd.git
cd lazymd
zig build
```

The compiled binary is at `zig-out/bin/lm`. Move it to your PATH:

```bash
cp zig-out/bin/lm /usr/local/bin/
```

## Pre-built binaries

Download pre-built binaries from [GitHub Releases](https://github.com/EME130/lazymd/releases). Available for:

- Linux x86_64
- macOS x86_64
- macOS ARM64 (Apple Silicon)
