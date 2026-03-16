package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/EME130/lazymd/internal/brain"
	"github.com/EME130/lazymd/internal/themes"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// Force-directed layout constants.
const (
	maxIterations = 150
	bhTheta       = 0.8  // Barnes-Hut threshold (higher = faster, less accurate)
	attraction    = 0.02 // edge attraction (lower = more spread)
	damping       = 0.9
	centerPull    = 0.005
	minDist       = 1.0
)

type pos struct{ x, y float64 }

type cell struct {
	ch    rune
	style lipgloss.Style
}

// nodeSpring holds per-node spring animation state.
type nodeSpring struct {
	x, xVel float64
	y, yVel float64
}

// BrainViewModel displays the knowledge graph as a force-directed ASCII layout.
type BrainViewModel struct {
	Graph           *brain.Graph
	positions       []pos // current display positions (animated)
	targetPositions []pos // final layout positions
	SelectedNode    int
	CurrentFile     int // node ID of current file, -1 if none
	LocalMode       bool
	viewportX       float64
	viewportY       float64
	zoom            float64

	// Spring animation
	spring    harmonica.Spring
	springs   []nodeSpring
	animating bool

	// Spatial index for viewport culling (rebuilt after animation settles)
	spatialIndex *quadTree
}

// NewBrainView creates a new brain view model.
func NewBrainView() BrainViewModel {
	return BrainViewModel{
		SelectedNode: 0,
		CurrentFile:  -1,
		zoom:         1.0,
		spring:       harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.7),
	}
}

// SetGraph sets the graph, computes the final layout, and starts spring animation.
func (bv *BrainViewModel) SetGraph(g *brain.Graph) {
	bv.Graph = g
	bv.positions = nil
	bv.targetPositions = nil
	bv.springs = nil
	bv.animating = false
	bv.spatialIndex = nil

	n := len(g.Nodes)
	if n == 0 {
		return
	}

	// Initial circular positions — spread wide
	startPositions := make([]pos, n)
	radius := math.Max(float64(n*6), 30)
	for i := range startPositions {
		angle := float64(i) * (2.0 * math.Pi / float64(n))
		startPositions[i] = pos{
			x: radius * math.Cos(angle),
			y: radius * math.Sin(angle),
		}
	}

	// Compute final layout using Barnes-Hut, then normalize
	bv.targetPositions = bv.computeLayout(g, startPositions)
	normalizePositions(bv.targetPositions)

	// Set display positions to start and init springs
	bv.positions = make([]pos, n)
	bv.springs = make([]nodeSpring, n)
	for i := range n {
		bv.positions[i] = startPositions[i]
		bv.springs[i] = nodeSpring{
			x: startPositions[i].x, xVel: 0,
			y: startPositions[i].y, yVel: 0,
		}
	}
	bv.animating = true

	if bv.SelectedNode >= n {
		bv.SelectedNode = 0
	}
}

// computeLayout runs the Barnes-Hut force-directed algorithm to completion.
func (bv *BrainViewModel) computeLayout(g *brain.Graph, initial []pos) []pos {
	n := len(g.Nodes)
	if n <= 1 {
		result := make([]pos, n)
		copy(result, initial)
		return result
	}

	// Scale repulsion up for small graphs so nodes don't collapse
	repStrength := math.Max(2000.0, 500.0*math.Sqrt(float64(n)))

	positions := make([]pos, n)
	copy(positions, initial)
	forces := make([]pos, n)

	for iter := range maxIterations {
		// Clear forces
		for i := range forces {
			forces[i] = pos{}
		}

		// Build quadtree for Barnes-Hut repulsion
		qt := newQuadTree(positions)

		// Barnes-Hut repulsion: O(n log n) instead of O(n²)
		for i := 0; i < n; i++ {
			bhRepulse(qt.root, i, positions[i].x, positions[i].y, &forces[i], repStrength)
		}

		// Attraction along edges: O(e)
		for _, edge := range g.Edges {
			dx := positions[edge.To].x - positions[edge.From].x
			dy := positions[edge.To].y - positions[edge.From].y
			fx := dx * attraction
			fy := dy * attraction
			forces[edge.From].x += fx
			forces[edge.From].y += fy
			forces[edge.To].x -= fx
			forces[edge.To].y -= fy
		}

		// Center pull + apply with damping
		maxMove := 0.0
		for i := 0; i < n; i++ {
			forces[i].x -= positions[i].x * centerPull
			forces[i].y -= positions[i].y * centerPull
			fx := forces[i].x * damping
			fy := forces[i].y * damping
			positions[i].x += fx
			positions[i].y += fy
			move := math.Abs(fx) + math.Abs(fy)
			if move > maxMove {
				maxMove = move
			}
		}

		if maxMove < 0.1 || iter >= maxIterations-1 {
			break
		}
	}

	return positions
}

// bhRepulse computes Barnes-Hut repulsion force on node i from a quadtree node.
func bhRepulse(node *quadNode, i int, px, py float64, force *pos, repStrength float64) {
	if node == nil || node.count == 0 {
		return
	}

	dx := node.comX - px
	dy := node.comY - py
	dist := math.Sqrt(dx*dx + dy*dy)

	// If it's a leaf with a single item that is this node, skip
	if node.children[0] == nil && node.count == 1 && len(node.items) == 1 && node.items[0].id == i {
		return
	}

	// Barnes-Hut criterion: if the cell is far enough away, treat as one body
	size := node.bounds.maxX - node.bounds.minX
	if node.children[0] == nil || (dist > 0 && size/dist < bhTheta) {
		if dist < minDist {
			dist = minDist
			dx = 0.5
			dy = 0.5
		}
		f := repStrength * float64(node.count) / (dist * dist)
		force.x -= (dx / dist) * f
		force.y -= (dy / dist) * f
		return
	}

	// Recurse into children
	for _, child := range node.children {
		bhRepulse(child, i, px, py, force, repStrength)
	}
}

// Animate steps the spring animation forward one frame. Returns true if still animating.
func (bv *BrainViewModel) Animate() bool {
	if !bv.animating || bv.Graph == nil {
		return false
	}

	n := len(bv.Graph.Nodes)
	if n == 0 || len(bv.springs) != n || len(bv.targetPositions) != n {
		bv.animating = false
		return false
	}

	settled := true
	for i := range n {
		bv.springs[i].x, bv.springs[i].xVel = bv.spring.Update(
			bv.springs[i].x, bv.springs[i].xVel, bv.targetPositions[i].x,
		)
		bv.springs[i].y, bv.springs[i].yVel = bv.spring.Update(
			bv.springs[i].y, bv.springs[i].yVel, bv.targetPositions[i].y,
		)

		bv.positions[i].x = bv.springs[i].x
		bv.positions[i].y = bv.springs[i].y

		dx := math.Abs(bv.springs[i].x - bv.targetPositions[i].x)
		dy := math.Abs(bv.springs[i].y - bv.targetPositions[i].y)
		vel := math.Abs(bv.springs[i].xVel) + math.Abs(bv.springs[i].yVel)
		if dx > 0.5 || dy > 0.5 || vel > 0.1 {
			settled = false
		}
	}

	if settled {
		for i := range n {
			bv.positions[i] = bv.targetPositions[i]
		}
		bv.animating = false
		// Build spatial index now that positions are final
		bv.spatialIndex = newQuadTree(bv.positions)
	}

	return bv.animating
}

// IsAnimating returns whether the brain view is currently animating.
func (bv *BrainViewModel) IsAnimating() bool {
	return bv.animating
}

// SetCurrentFile sets the current file node by name.
func (bv *BrainViewModel) SetCurrentFile(name string) {
	if bv.Graph == nil {
		return
	}
	if id, ok := bv.Graph.Resolve(name); ok {
		bv.CurrentFile = int(id)
	} else {
		bv.CurrentFile = -1
	}
}

// SelectedPath returns the file path of the currently selected node, or "".
func (bv *BrainViewModel) SelectedPath() string {
	if bv.Graph == nil {
		return ""
	}
	n := len(bv.Graph.Nodes)
	if bv.SelectedNode < 0 || bv.SelectedNode >= n {
		return ""
	}
	return bv.Graph.Nodes[bv.SelectedNode].Path
}

// normalizePositions centers and scales positions to fill a reasonable range. O(n).
func normalizePositions(positions []pos) {
	n := len(positions)
	if n <= 1 {
		return
	}

	// Find bounding box — single pass
	minX, maxX := positions[0].x, positions[0].x
	minY, maxY := positions[0].y, positions[0].y
	for _, p := range positions[1:] {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	// Center + scale in one pass
	cx := (minX + maxX) / 2
	cy := (minY + maxY) / 2
	span := math.Max(maxX-minX, maxY-minY)
	if span < 1 {
		span = 1
	}
	target := math.Max(40, float64(n)*8)
	scale := target / span

	for i := range positions {
		positions[i].x = (positions[i].x - cx) * scale
		positions[i].y = (positions[i].y - cy) * scale
	}
}

// MoveSelection moves the selected node by delta.
func (bv *BrainViewModel) MoveSelection(delta int) {
	if bv.Graph == nil {
		return
	}
	n := len(bv.Graph.Nodes)
	if n == 0 {
		return
	}
	bv.SelectedNode = ((bv.SelectedNode+delta)%n + n) % n
}

// PanView pans the viewport by the given delta.
func (bv *BrainViewModel) PanView(dx, dy float64) {
	bv.viewportX += dx
	bv.viewportY += dy
}

// Zoom adjusts zoom level by delta, clamped to [0.1, 5.0].
func (bv *BrainViewModel) Zoom(delta float64) {
	bv.zoom += delta
	if bv.zoom < 0.1 {
		bv.zoom = 0.1
	}
	if bv.zoom > 5.0 {
		bv.zoom = 5.0
	}
}

// ResetView resets viewport and zoom to defaults.
func (bv *BrainViewModel) ResetView() {
	bv.viewportX = 0
	bv.viewportY = 0
	bv.zoom = 1.0
}

// CenterOnSelected pans the viewport to center the selected node.
func (bv *BrainViewModel) CenterOnSelected() {
	if bv.Graph == nil || bv.SelectedNode < 0 || bv.SelectedNode >= len(bv.positions) {
		return
	}
	p := bv.positions[bv.SelectedNode]
	bv.viewportX = p.x * bv.zoom
	bv.viewportY = p.y * bv.zoom * 0.5
}

// nodeHeat returns a 0.0-1.0 connectivity score for a node.
func nodeHeat(node *brain.Node, maxLinks int) float64 {
	if maxLinks <= 0 {
		return 0
	}
	links := len(node.OutLinks) + len(node.InLinks)
	h := float64(links) / float64(maxLinks)
	if h > 1.0 {
		h = 1.0
	}
	return h
}

// visibleWorldRect computes the world-space rectangle visible in the viewport.
func (bv *BrainViewModel) visibleWorldRect(gridW, gridH int) rect2D {
	cx := float64(gridW) / 2.0
	cy := float64(gridH) / 2.0
	z := bv.zoom
	if z < 0.01 {
		z = 0.01
	}
	// Invert worldToScreen: wx = (sx - cx + vpX) / zoom
	return rect2D{
		minX: (0 - cx + bv.viewportX) / z,
		maxX: (float64(gridW) - cx + bv.viewportX) / z,
		minY: (0 - cy + bv.viewportY) / (z * 0.5),
		maxY: (float64(gridH) - cy + bv.viewportY) / (z * 0.5),
	}
}

// View renders the brain view panel.
func (bv *BrainViewModel) View(rect Rect) string {
	c := themes.CurrentColors()

	if bv.Graph == nil {
		style := lipgloss.NewStyle().
			Width(rect.W).Height(rect.H).
			Foreground(lipgloss.Color(c.TextMuted))
		return style.Render("  No graph loaded")
	}

	g := bv.Graph
	n := len(g.Nodes)
	if n == 0 {
		style := lipgloss.NewStyle().
			Width(rect.W).Height(rect.H).
			Foreground(lipgloss.Color(c.TextMuted))
		return style.Render("  No notes found")
	}

	// Find max link count for heat mapping
	maxLinks := 0
	for i := range g.Nodes {
		links := len(g.Nodes[i].OutLinks) + len(g.Nodes[i].InLinks)
		if links > maxLinks {
			maxLinks = links
		}
	}

	// Create a character grid
	gridW := rect.W - 2
	gridH := rect.H - 2
	if gridW < 1 {
		gridW = 1
	}
	if gridH < 1 {
		gridH = 1
	}

	grid := make([][]cell, gridH)
	for i := range grid {
		grid[i] = make([]cell, gridW)
		for j := range grid[i] {
			grid[i][j] = cell{ch: ' ', style: lipgloss.NewStyle()}
		}
	}

	cx := float64(gridW) / 2.0
	cy := float64(gridH) / 2.0

	worldToScreen := func(wx, wy float64) (int, int) {
		sx := int((wx*bv.zoom - bv.viewportX) + cx)
		sy := int((wy*bv.zoom*0.5 - bv.viewportY) + cy)
		return sx, sy
	}

	// Determine visible nodes with viewport culling
	worldRect := bv.visibleWorldRect(gridW, gridH)

	var visible []int
	if bv.LocalMode {
		// Local mode: BFS neighbors (small set, no culling needed)
		center := bv.SelectedNode
		if bv.CurrentFile >= 0 {
			center = bv.CurrentFile
		}
		neighbors := g.GetNeighbors(uint16(center), 2)
		for _, nid := range neighbors {
			visible = append(visible, int(nid))
		}
	} else if bv.spatialIndex != nil && !bv.animating {
		// Use quadtree for O(log n + k) viewport culling
		visible = bv.spatialIndex.queryRange(worldRect)
	} else {
		// Fallback: linear scan with bounds check
		for i := 0; i < n; i++ {
			if worldRect.contains(bv.positions[i].x, bv.positions[i].y) {
				visible = append(visible, i)
			}
		}
	}

	isVis := make([]bool, n)
	for _, v := range visible {
		isVis[v] = true
	}

	// Build set of nodes directly connected to selected node
	connectedToSelected := make([]bool, n)
	if bv.SelectedNode >= 0 && bv.SelectedNode < n {
		sel := g.Nodes[bv.SelectedNode]
		for _, id := range sel.OutLinks {
			connectedToSelected[id] = true
		}
		for _, id := range sel.InLinks {
			connectedToSelected[id] = true
		}
	}

	// Edge styles
	edgeDim := lipgloss.NewStyle().Foreground(lipgloss.Color(c.TextMuted))
	edgeLit := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Border))
	edgeActive := lipgloss.NewStyle().Foreground(lipgloss.Color(c.BorderActive))

	// Draw edges as sparse dots (only if both endpoints visible)
	for _, edge := range g.Edges {
		if !isVis[int(edge.From)] && !isVis[int(edge.To)] {
			continue
		}
		x1, y1 := worldToScreen(bv.positions[edge.From].x, bv.positions[edge.From].y)
		x2, y2 := worldToScreen(bv.positions[edge.To].x, bv.positions[edge.To].y)

		eStyle := edgeDim
		if int(edge.From) == bv.SelectedNode || int(edge.To) == bv.SelectedNode {
			eStyle = edgeActive
		} else if int(edge.From) == bv.CurrentFile || int(edge.To) == bv.CurrentFile {
			eStyle = edgeLit
		}

		drawDottedLine(grid, gridW, gridH, x1, y1, x2, y2, eStyle)
	}

	// Heat color tiers
	heatColors := []string{c.TextMuted, c.Text, string(c.Link), string(c.H2), string(c.H1)}

	// Draw nodes
	for _, nid := range visible {
		sx, sy := worldToScreen(bv.positions[nid].x, bv.positions[nid].y)
		if sy < 0 || sy >= gridH || sx < 0 || sx >= gridW {
			continue
		}

		node := &g.Nodes[nid]
		isSelected := nid == bv.SelectedNode
		isCurrent := nid == bv.CurrentFile

		var nodeColor string
		if isSelected {
			nodeColor = string(c.BorderActive)
		} else if isCurrent {
			nodeColor = string(c.Success)
		} else if connectedToSelected[nid] {
			nodeColor = string(c.H5)
		} else {
			heat := nodeHeat(node, maxLinks)
			idx := int(heat * float64(len(heatColors)-1))
			if idx >= len(heatColors) {
				idx = len(heatColors) - 1
			}
			nodeColor = heatColors[idx]
		}

		nodeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(nodeColor))
		boldStyle := nodeStyle.Bold(true)

		marker := '●'
		if isCurrent {
			marker = '◉'
		} else if isSelected {
			marker = '◈'
		} else if len(node.OutLinks)+len(node.InLinks) == 0 {
			marker = '○'
		}

		grid[sy][sx] = cell{ch: marker, style: boldStyle}

		// Label with link count
		label := node.Name
		linkCount := len(node.OutLinks) + len(node.InLinks)
		badge := ""
		if linkCount > 0 {
			badge = fmt.Sprintf(" (%d)", linkCount)
		}
		fullLabel := " " + label + badge
		maxLen := gridW - sx - 1
		if maxLen > 0 {
			if len(fullLabel) > maxLen {
				fullLabel = fullLabel[:maxLen]
			}
			labelStyle := nodeStyle
			if isSelected || isCurrent {
				labelStyle = boldStyle
			}
			for ci, ch := range fullLabel {
				col := sx + 1 + ci
				if col >= gridW {
					break
				}
				grid[sy][col] = cell{ch: ch, style: labelStyle}
			}
		}
	}

	// Render grid to lines
	lines := make([]string, 0, gridH+1)
	for _, row := range grid {
		var sb strings.Builder
		for _, c := range row {
			sb.WriteString(c.style.Render(string(c.ch)))
		}
		lines = append(lines, sb.String())
	}

	// Status line
	var status string
	if bv.SelectedNode >= 0 && bv.SelectedNode < n {
		node := g.Nodes[bv.SelectedNode]
		mode := "GRAPH"
		if bv.LocalMode {
			mode = "LOCAL"
		}
		hint := "⏎ dive in"
		if bv.LocalMode {
			hint = "⏎ dive in  ⌫ back"
		}
		status = fmt.Sprintf(" %s  ↗%d ↙%d  [%s]  %s",
			node.Name, len(node.OutLinks), len(node.InLinks), mode, hint)
	}
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.TitleFg)).Bold(true)
	lines = append(lines, statusStyle.Render(status))

	content := strings.Join(lines, "\n")
	border := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(c.Border)).
		Width(rect.W).Height(rect.H)
	return border.Render(content)
}

// drawDottedLine draws a sparse dotted line between two points.
func drawDottedLine(grid [][]cell, gridW, gridH, x1, y1, x2, y2 int, style lipgloss.Style) {
	dx := x2 - x1
	dy := y2 - y1
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		return
	}
	if steps > 200 {
		steps = 200
	}

	sx := float64(dx) / float64(steps)
	sy := float64(dy) / float64(steps)
	px := float64(x1)
	py := float64(y1)

	for i := range steps {
		if i%2 == 0 {
			px += sx
			py += sy
			continue
		}
		ix := int(px)
		iy := int(py)
		if ix >= 0 && ix < gridW && iy >= 0 && iy < gridH {
			grid[iy][ix] = cell{ch: '·', style: style}
		}
		px += sx
		py += sy
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
