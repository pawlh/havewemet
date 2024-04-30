package quadtree

import "math"

// maxObjects the maximum number of objects in a quadtree node before it splits into 4 sub nodes
const maxObjects = 10

type QuadTree[V any] struct {
	node[V]
	allPoints []*Point[V]
}

func NewQuadTree[V any]() QuadTree[V] {
	return QuadTree[V]{}
}

func (q *QuadTree[V]) Insert(p *Point[V]) {
	if !q.contains(*p) {
		q.grow(*p)
		q.resort()
	}

	q.allPoints = append(q.allPoints, p)

	q.node.insert(p)
}

func (q *QuadTree[V]) QueryRadius(x, y, radius float64) []*Point[V] {
	var points []*Point[V]
	q.queryRadius(x, y, radius, &points)
	return points
}

func (q *QuadTree[V]) resort() {
	q.clearSubNodes()
	for _, p := range q.allPoints {
		q.insert(p)
	}
}

func (q *QuadTree[V]) grow(point Point[V]) {
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

func (q *QuadTree[V]) contains(point Point[V]) bool {
	return point.x >= q.x1 && point.x <= q.x2 && point.y >= q.y1 && point.y <= q.y2
}

func (q *QuadTree[V]) clearSubNodes() {
	q.points = nil
	q.subNodes = [4]*node[V]{}
}

type node[V any] struct {
	x1, y1, x2, y2 float64
	points         []*Point[V]
	subNodes       [4]*node[V]
}

func (n *node[V]) insert(p *Point[V]) {
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

func (n *node[V]) isLeaf() bool {
	return n.subNodes[0] == nil
}

func (n *node[V]) pickSubNode(p Point[V]) int {
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

func (n *node[V]) split() {

	xMid := (n.x1 + n.x2) / 2
	yMid := (n.y1 + n.y2) / 2

	n.subNodes[0] = &node[V]{x1: n.x1, y1: n.y1, x2: xMid, y2: yMid}
	n.subNodes[1] = &node[V]{x1: xMid, y1: n.y1, x2: n.x2, y2: yMid}
	n.subNodes[2] = &node[V]{x1: n.x1, y1: yMid, x2: xMid, y2: n.y2}
	n.subNodes[3] = &node[V]{x1: xMid, y1: yMid, x2: n.x2, y2: n.y2}

	for _, p := range n.points {
		subNodeIndex := n.pickSubNode(*p)
		n.subNodes[subNodeIndex].insert(p)
	}

	n.points = nil
}

func (n *node[V]) queryRadius(x, y, radius float64, points *[]*Point[V]) {
	if !n.intersectsCircle(x, y, radius) {
		return
	}

	for _, p := range n.points {
		if p.distance(x, y) <= radius {
			*points = append(*points, p)
		}
	}

	if !n.isLeaf() {
		for _, subNode := range n.subNodes {
			subNode.queryRadius(x, y, radius, points)
		}
	}
}

func (n *node[V]) intersectsCircle(x, y, radius float64) bool {
	if x+radius < n.x1 || x-radius > n.x2 || y+radius < n.y1 || y-radius > n.y2 {
		return false
	}

	return true

}

type Point[V any] struct {
	x, y  float64
	Value V
}

func NewPoint[V any](x, y float64, value V) *Point[V] {
	return &Point[V]{x: x, y: y, Value: value}
}

func (p *Point[V]) distance(x, y float64) float64 {
	return math.Sqrt(math.Pow(p.x-x, 2) + math.Pow(p.y-y, 2))
}
