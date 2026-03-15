# Panels

LazyMD has four panels: File Tree, Editor, Preview, and Brain.

## File Tree (Alt+1)

Browse files in the current directory.

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Open selected file |
| `r` | Refresh file listing |
| `Esc` | Return focus to editor |

Click a file entry to select it; click again to open.

## Editor (always visible)

Vim-style modal editor with syntax-aware rendering.

- **Normal mode**: Navigation and commands
- **Insert mode**: Text input (press `i` to enter)
- **Command mode**: Ex commands (press `:` to enter)

Line numbers in the gutter. Active line is highlighted.

## Preview (Alt+2)

Rendered markdown preview using Glamour.

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll down/up |
| `d` / `u` | Page down/up |
| `Esc` | Return focus to editor |

Preview updates automatically with a 250ms debounce after edits.

## Brain (Alt+3)

Force-directed ASCII graph of wiki-link connections.

| Key | Action |
|-----|--------|
| `j` / `k` | Select next/previous node |
| `f` | Toggle local mode |
| `+` / `-` | Zoom in/out |
| `h` / `l` | Pan left/right |
| `Esc` | Return focus to editor |

Local mode shows only direct neighbors of the selected node.
