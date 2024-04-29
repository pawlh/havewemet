package quadtree

// maxObjects the maximum number of objects in a quadtree node before it splits into 4 sub nodes
const maxObjects = 10

type Quadtree struct {
	node
	allPoints []*Point
}

func NewQuadTree() Quadtree {
	return Quadtree{}
}

func (q *Quadtree) Insert(p *Point) {
	q.allPoints = append(q.allPoints, p)

	if !q.contains(*p) {
		q.grow(*p)
	}

	q.node.insert(p)
}

func (q *Quadtree) resort() {
	q.clearSubNodes()
	for _, p := range q.allPoints {
		q.insert(p)
	}
}

func (q *Quadtree) grow(point Point) {
	if point.x < q.x1 {
		leftXDiff := q.x1 - point.x
		q.x1 -= leftXDiff * 2
	}

	if point.x > q.x2 {
		rightXDiff := point.x - q.x2
		q.x2 += rightXDiff * 2
	}

	if point.y < q.y1 {
		bottomYDiff := q.y1 - point.y
		q.y1 -= bottomYDiff * 2
	}

	if point.y > q.y2 {
		topYDiff := point.y - q.y2
		q.y2 += topYDiff * 2
	}
}

func (q *Quadtree) contains(point Point) bool {
	return point.x >= q.x1 && point.x <= q.x2 && point.y >= q.y1 && point.y <= q.y2
}

func (q *Quadtree) clearSubNodes() {
	q.subNodes = [4]*node{}
}

type node struct {
	x1, y1, x2, y2 float64
	points         []*Point
	subNodes       [4]*node
}

func (n *node) insert(p *Point) {
	if !n.isLeaf() {
		subNodeIndex := n.pickSubNode(*p)
		n.subNodes[subNodeIndex].insert(p)
		return
	}

	n.points = append(n.points, p)

	if len(n.points) > maxObjects {
		n.split()
	}
}

func (n *node) isLeaf() bool {
	return n.subNodes[0] == nil
}

func (n *node) pickSubNode(p Point) int {
	if p.x <= (n.x1+n.x2)/2 {
		if p.y < (n.y1+n.y2)/2 {
			return 0
		} else {
			return 1
		}
	} else {
		if p.y < (n.y1+n.y2)/2 {
			return 2
		} else {
			return 3
		}
	}
}

func (n *node) split() {

	xMid := (n.x1 + n.x2) / 2
	yMid := (n.y1 + n.y2) / 2

	n.subNodes[0] = &node{x1: n.x1, y1: n.y1, x2: xMid, y2: yMid}
	n.subNodes[1] = &node{x1: xMid, y1: n.y1, x2: n.x2, y2: yMid}
	n.subNodes[2] = &node{x1: n.x1, y1: yMid, x2: xMid, y2: n.y2}
	n.subNodes[3] = &node{x1: xMid, y1: yMid, x2: n.x2, y2: n.y2}

	for _, p := range n.points {
		subNodeIndex := n.pickSubNode(*p)
		n.subNodes[subNodeIndex].insert(p)
	}

	n.points = nil
}

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}
