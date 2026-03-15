# Brain — Knowledge Graph

The Brain panel visualizes connections between your notes using wiki-links.

## Wiki-Links

Use double brackets to link notes:

    [[note-name]]

Links are resolved by filename, case-insensitive. Path prefixes and display aliases are supported:

    [[folder/note]]       resolved by last segment
    [[note|display text]] resolved by part before |

## Brain Panel

Toggle with `Alt+3`. The panel shows a force-directed ASCII graph of your vault.

| Key | Action |
|-----|--------|
| `j` / `k` | Select next/previous node |
| `f` | Toggle local mode (show only neighbors) |
| `+` / `-` | Zoom in / out |
| `h` / `l` | Pan left / right |
| `Esc` | Return focus to editor |

## Graph Features

- **Backlinks**: See which notes link TO a given note
- **Neighbors**: Find connected notes within N hops
- **Path finding**: Shortest path between any two notes
- **Orphan detection**: Notes with zero connections
- **Hub notes**: Most connected notes ranked by link count

## MCP Brain Tools

When running as an MCP server (`lm --mcp-server`), AI agents can query the graph:

- `list_links` — outgoing wiki-links from the current document
- `get_backlinks` — files that link TO a given note
- `get_graph` — full graph as JSON (nodes, edges, stats)
- `get_neighbors` — directly connected notes with link counts
- `find_path` — shortest path between two notes
- `get_orphans` — notes with zero links
- `get_hub_notes` — most connected notes ranked by link count
