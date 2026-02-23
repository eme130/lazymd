---
title: Building from Source
sidebar_position: 2
description: Build LazyMD from source with Zig — debug builds, release builds, and running directly. Compile the terminal markdown editor yourself.
keywords: [build LazyMD, zig build, compile from source, debug build, release build, zig project]
---

# Building from Source

```bash
# Debug build
zig build

# Release build
zig build -Doptimize=ReleaseSafe

# Run directly
zig build run -- myfile.md
```
