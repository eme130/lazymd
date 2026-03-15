# MCP Server

LazyMD is an MCP (Model Context Protocol) server. AI agents connect via stdio using JSON-RPC 2.0.

## Starting the Server

    lm --mcp-server              # Start MCP server
    lm --mcp-server myfile.md    # Start with file preloaded

## Tools (22 total)

### Document Tools
`open_file`, `read_document`, `write_document`, `list_headings`, `edit_section`, `insert_text`, `delete_lines`, `search_content`, `get_structure`

### Navigation Tools
`read_section`, `list_tasks`, `update_task`, `get_breadcrumb`, `move_section`, `read_section_range`

### Brain Tools
`list_links`, `get_backlinks`, `get_graph`, `get_neighbors`, `find_path`, `get_orphans`, `get_hub_notes`

## Claude Code Setup

    claude mcp add LazyMD -- /path/to/lm --mcp-server

## Gemini CLI Setup

Add to `~/.gemini/settings.json`:

    {
      "mcpServers": {
        "LazyMD": {
          "command": "/path/to/lm",
          "args": ["--mcp-server"]
        }
      }
    }
