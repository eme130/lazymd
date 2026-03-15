# Wiki-Links

Wiki-links are the backbone of LazyMD's knowledge graph.

## Syntax

Use double brackets to link notes:

    [[note-name]]

LazyMD resolves links by filename (case-insensitive, ignores path prefixes).

## How It Works

1. The scanner finds all `[[...]]` patterns in your vault
2. Each unique note becomes a node in the graph
3. Each wiki-link becomes a directed edge
4. The Brain panel (Alt+3) visualizes the result

## Features

- **Backlinks**: See what links TO a note
- **Neighbors**: Find connected notes within N hops
- **Path finding**: Shortest path between any two notes
- **Orphan detection**: Notes with zero connections
- **Hub notes**: Most connected notes ranked by link count

## Related

- [[zettelkasten]] — The method behind linked notes
- [[markdown]] — Base syntax for everything
- [[welcome]] — Back to the start
