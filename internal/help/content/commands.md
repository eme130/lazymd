# Built-in Commands

Enter command mode by pressing `:` in normal mode.

## File Commands

| Command | Description |
|---------|-------------|
| `:w` | Save current file |
| `:w <path>` | Save to a specific path |
| `:q` | Quit (fails if unsaved changes) |
| `:q!` | Force quit without saving |
| `:wq` or `:x` | Save and quit |
| `:e <path>` | Open a file |

## Theme Commands

| Command | Description |
|---------|-------------|
| `:theme` | Show current theme name |
| `:theme.cycle` | Cycle to next theme |
| `:theme.list` | List all available themes |
| `:theme <name>` | Switch to a specific theme |

## Help Commands

| Command | Description |
|---------|-------------|
| `:help` | Show this help overview |
| `:help keys` | Keybindings reference |
| `:help brain` | Knowledge graph guide |
| `:help plugins` | List all plugins and commands |
| `:help commands` | Built-in commands (this page) |
| `:help mcp` | MCP server setup |
| `:help panels` | Panel usage guide |

## Plugin Commands

Plugins register additional commands. Use `:help plugins` to see all available plugin commands.
