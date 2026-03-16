package ui

// quadtree provides spatial indexing for efficient viewport culling and
// Barnes-Hut force approximation on large node sets.

const qtMaxItems = 8 // items per leaf before splitting

// rect2D is an axis-aligned bounding rectangle.
type rect2D struct {
	minX, minY, maxX, maxY float64
}

func (r rect2D) contains(x, y float64) bool {
	return x >= r.minX && x <= r.maxX && y >= r.minY && y <= r.maxY
}

func (r rect2D) intersects(o rect2D) bool {
	return r.minX <= o.maxX && r.maxX >= o.minX && r.minY <= o.maxY && r.maxY >= o.minY
}

func (r rect2D) width() float64  { return r.maxX - r.minX }
func (r rect2D) height() float64 { return r.maxY - r.minY }
func (r rect2D) midX() float64   { return (r.minX + r.maxX) / 2 }
func (r rect2D) midY() float64   { return (r.minY + r.maxY) / 2 }

// qtItem is a point with an associated node index.
type qtItem struct {
	x, y float64
	id   int
}

// quadNode is a node in the quadtree.
type quadNode struct {
	bounds   rect2D
	items    []qtItem // leaf items (nil if subdivided)
	children [4]*quadNode
	// Barnes-Hut aggregate: center of mass and total count
	comX, comY float64
	count      int
}

// quadTree is a spatial index.
type quadTree struct {
	root *quadNode
}

// newQuadTree builds a quadtree from positions.
func newQuadTree(positions []pos) *quadTree {
	if len(positions) == 0 {
		return &quadTree{}
	}

	// Compute bounds
	b := rect2D{positions[0].x, positions[0].y, positions[0].x, positions[0].y}
	for _, p := range positions[1:] {
		if p.x < b.minX {
			b.minX = p.x
		}
		if p.x > b.maxX {
			b.maxX = p.x
		}
		if p.y < b.minY {
			b.minY = p.y
		}
		if p.y > b.maxY {
			b.maxY = p.y
		}
	}
	// Add small padding to avoid zero-size bounds
	b.minX -= 1
	b.minY -= 1
	b.maxX += 1
	b.maxY += 1

	root := &quadNode{bounds: b}
	for i, p := range positions {
		qtInsert(root, qtItem{x: p.x, y: p.y, id: i})
	}
	qtComputeMass(root)
	return &quadTree{root: root}
}

func qtInsert(node *quadNode, item qtItem) {
	if !node.bounds.contains(item.x, item.y) {
		return
	}

	// Leaf with room
	if node.children[0] == nil && len(node.items) < qtMaxItems {
		node.items = append(node.items, item)
		return
	}

	// Need to subdivide
	if node.children[0] == nil {
		qtSubdivide(node)
	}

	for _, child := range node.children {
		if child.bounds.contains(item.x, item.y) {
			qtInsert(child, item)
			return
		}
	}
}

func qtSubdivide(node *quadNode) {
	b := node.bounds
	mx, my := b.midX(), b.midY()

	node.children = [4]*quadNode{
		{bounds: rect2D{b.minX, b.minY, mx, my}},
		{bounds: rect2D{mx, b.minY, b.maxX, my}},
		{bounds: rect2D{b.minX, my, mx, b.maxY}},
		{bounds: rect2D{mx, my, b.maxX, b.maxY}},
	}

	// Re-insert existing items
	for _, item := range node.items {
		for _, child := range node.children {
			if child.bounds.contains(item.x, item.y) {
				qtInsert(child, item)
				break
			}
		}
	}
	node.items = nil
}

// qtComputeMass computes center-of-mass for Barnes-Hut.
func qtComputeMass(node *quadNode) {
	if node == nil {
		return
	}

	if node.children[0] == nil {
		// Leaf
		node.count = len(node.items)
		if node.count > 0 {
			var sx, sy float64
			for _, item := range node.items {
				sx += item.x
				sy += item.y
			}
			node.comX = sx / float64(node.count)
			node.comY = sy / float64(node.count)
		}
		return
	}

	node.count = 0
	var sx, sy float64
	for _, child := range node.children {
		qtComputeMass(child)
		if child.count > 0 {
			node.count += child.count
			sx += child.comX * float64(child.count)
			sy += child.comY * float64(child.count)
		}
	}
	if node.count > 0 {
		node.comX = sx / float64(node.count)
		node.comY = sy / float64(node.count)
	}
}

// queryRange returns all node IDs within the given rectangle.
func (qt *quadTree) queryRange(area rect2D) []int {
	if qt.root == nil {
		return nil
	}
	var results []int
	qtQueryRange(qt.root, area, &results)
	return results
}

func qtQueryRange(node *quadNode, area rect2D, results *[]int) {
	if node == nil || !node.bounds.intersects(area) {
		return
	}

	if node.children[0] == nil {
		for _, item := range node.items {
			if area.contains(item.x, item.y) {
				*results = append(*results, item.id)
			}
		}
		return
	}

	for _, child := range node.children {
		qtQueryRange(child, area, results)
	}
}
