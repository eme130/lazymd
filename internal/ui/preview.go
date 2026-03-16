package ui

import (
	"regexp"
	"strings"

	"github.com/EME130/lazymd/internal/buffer"
	"github.com/EME130/lazymd/internal/themes"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// PreviewModel renders a markdown preview using Glamour.
type PreviewModel struct {
	rendered  string
	dirty     bool
	scrollOff int
	lastWidth int
	renderer  *glamour.TermRenderer
}

// NewPreview creates a new preview model.
func NewPreview() PreviewModel {
	return PreviewModel{dirty: true}
}

// Invalidate marks the preview as needing re-render.
func (p *PreviewModel) Invalidate() {
	p.dirty = true
}

// ScrollUp scrolls the preview up.
func (p *PreviewModel) ScrollUp(n int) {
	p.scrollOff -= n
	if p.scrollOff < 0 {
		p.scrollOff = 0
	}
}

// ScrollDown scrolls the preview down.
func (p *PreviewModel) ScrollDown(n int) {
	p.scrollOff += n
}

// headingBorder returns the lipgloss border style for a heading level.
func headingBorder(level int) lipgloss.Border {
	switch level {
	case 1:
		return lipgloss.ThickBorder()
	case 2:
		return lipgloss.DoubleBorder()
	case 3:
		return lipgloss.NormalBorder()
	default:
		return lipgloss.RoundedBorder()
	}
}

// headingColor returns the theme color for a heading level.
func headingColor(c *themes.ThemeColors, level int) themes.Color {
	switch level {
	case 1:
		return c.H1
	case 2:
		return c.H2
	case 3:
		return c.H3
	case 4:
		return c.H4
	case 5:
		return c.H5
	default:
		return c.H6
	}
}

// renderHeading renders a heading as a bordered box.
func renderHeading(text string, level, width int) string {
	c := themes.CurrentColors()
	color := headingColor(c, level)
	innerW := width - 4 // account for border + padding
	if innerW < 4 {
		innerW = 4
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color)).
		Border(headingBorder(level)).
		BorderForeground(lipgloss.Color(color)).
		Width(innerW).
		Padding(0, 1)

	return style.Render(text)
}

// segment represents a chunk of markdown content — either a heading or a block of regular text.
type segment struct {
	isHeading bool
	level     int    // heading level (1-6), only set if isHeading
	text      string // heading text (stripped of #) or raw markdown block
}

// splitSegments splits markdown into heading and non-heading segments.
func splitSegments(content string) []segment {
	lines := strings.Split(content, "\n")
	var segments []segment
	var block []string
	inFencedBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track fenced code blocks to avoid treating # inside them as headings
		if strings.HasPrefix(trimmed, "```") {
			inFencedBlock = !inFencedBlock
		}

		if !inFencedBlock {
			if m := headingRe.FindStringSubmatch(line); m != nil {
				// Flush accumulated non-heading block
				if len(block) > 0 {
					segments = append(segments, segment{text: strings.Join(block, "\n")})
					block = nil
				}
				segments = append(segments, segment{
					isHeading: true,
					level:     len(m[1]),
					text:      strings.TrimSpace(m[2]),
				})
				continue
			}
		}

		block = append(block, line)
	}

	// Flush remaining block
	if len(block) > 0 {
		segments = append(segments, segment{text: strings.Join(block, "\n")})
	}

	return segments
}

// RenderNow performs the render immediately (called from Update on debounce tick).
func (p *PreviewModel) RenderNow(buf *buffer.Buffer, rect Rect) {
	if buf.Length() == 0 {
		p.rendered = ""
		p.dirty = false
		return
	}

	wrapWidth := rect.W - 4
	if wrapWidth < 20 {
		wrapWidth = 20
	}

	// Recreate renderer if width changed or first use
	if p.renderer == nil || p.lastWidth != wrapWidth {
		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(wrapWidth),
		)
		if err == nil {
			p.renderer = r
			p.lastWidth = wrapWidth
		}
	}

	content := buf.Content()
	segs := splitSegments(content)

	var parts []string
	for _, seg := range segs {
		if seg.isHeading {
			parts = append(parts, renderHeading(seg.text, seg.level, rect.W-2))
		} else {
			// Render non-heading markdown with Glamour
			if p.renderer != nil {
				rendered, err := p.renderer.Render(seg.text)
				if err == nil {
					parts = append(parts, rendered)
				} else {
					parts = append(parts, seg.text)
				}
			} else {
				parts = append(parts, seg.text)
			}
		}
	}

	p.rendered = strings.Join(parts, "\n")
	p.dirty = false
}

// View renders the markdown preview panel.
func (p *PreviewModel) View(buf *buffer.Buffer, rect Rect) string {
	c := themes.CurrentColors()

	if buf.Length() == 0 {
		style := lipgloss.NewStyle().
			Width(rect.W).Height(rect.H).
			Foreground(lipgloss.Color(c.TextMuted))
		return style.Render("  No content to preview")
	}

	if p.dirty && p.rendered == "" {
		style := lipgloss.NewStyle().
			Width(rect.W).Height(rect.H).
			Foreground(lipgloss.Color(c.TextMuted))
		return style.Render("  Rendering...")
	}

	// Split into lines and apply scrolling
	lines := strings.Split(p.rendered, "\n")

	// Clamp scroll
	if p.scrollOff > len(lines)-1 {
		p.scrollOff = len(lines) - 1
	}
	if p.scrollOff < 0 {
		p.scrollOff = 0
	}

	contentH := rect.H - 1 // account for panel header
	if contentH < 1 {
		contentH = 1
	}

	var visible []string
	for i := 0; i < contentH; i++ {
		idx := p.scrollOff + i
		if idx >= len(lines) {
			visible = append(visible, "")
			continue
		}
		line := lines[idx]
		if lipgloss.Width(line) > rect.W-2 {
			runes := []rune(line)
			if len(runes) > rect.W-2 {
				line = string(runes[:rect.W-2])
			}
		}
		visible = append(visible, line)
	}

	return strings.Join(visible, "\n")
}
