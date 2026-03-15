# Go Programming

Notes on Go as used in LazyMD.

## Why Go?

- Fast compilation, single binary output
- Excellent concurrency primitives
- Strong standard library

## Key Packages

- `bubbletea` — Terminal UI framework
- `lipgloss` — Style definitions
- `glamour` — Markdown rendering

## Code Example

```go
func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Related

- [[markdown]] — The file format we edit
- [[wiki-links]] — How notes connect
