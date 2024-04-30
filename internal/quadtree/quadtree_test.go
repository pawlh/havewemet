package quadtree

import "testing"

func TestQuadTree_Insert(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(10)
	point := &Point[string]{x: 1, y: 1, Value: "A"}

	quadTree.Insert(point)

	verifyNodePointsLength(t, quadTree.node, 1)
	verifyAllPointsLength(t, quadTree, 1)
}

func TestQuadTree_Insert_NoSplit(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(3)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "A"},
		{x: 2, y: 2, Value: "B"},
		{x: 3, y: 3, Value: "C"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	verifyNodePointsLength(t, quadTree.node, 3)
	verifyAllPointsLength(t, quadTree, 3)

	for _, subNode := range quadTree.node.subNodes {
		if subNode != nil {
			t.Errorf("Expected nil subNode, got %v", subNode)
		}
	}
}

func TestQuadTree_Insert_Split(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(2)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "A"},
		{x: 2, y: 2, Value: "B"},
		{x: 3, y: 3, Value: "C"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	verifyNodePointsLength(t, quadTree.node, 0)
	verifyAllPointsLength(t, quadTree, 3)

	for _, subNode := range quadTree.node.subNodes {
		if subNode == nil {
			t.Errorf("Expected non-nil subNode, got %v", subNode)
		}
	}
}

func TestQuadTree_Insert_CorrectSubNodes(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(2)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "Quad 1"},
		{x: -1, y: 1, Value: "Quad 2"},
		{x: -1, y: -1, Value: "Quad 3"},
		{x: 1, y: -1, Value: "Quad 4"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	verifyNodePointsLength(t, quadTree.node, 0)
	verifyNonNilSubNodes(t, quadTree.node.subNodes)

	if quadTree.node.subNodes[0].points[0] != points[0] {
		t.Errorf("Expected %s in subNode 0, got %s", points[0].Value, quadTree.node.subNodes[0].points[0].Value)
	}

	if quadTree.node.subNodes[1].points[0] != points[1] {
		t.Errorf("Expected %s in subNode 1, got %s", points[1].Value, quadTree.node.subNodes[1].points[0].Value)
	}

	if quadTree.node.subNodes[2].points[0] != points[2] {
		t.Errorf("Expected %s in subNode 2, got %s", points[2].Value, quadTree.node.subNodes[2].points[0].Value)
	}

	if quadTree.node.subNodes[3].points[0] != points[3] {
		t.Errorf("Expected %s in subNode 3, got %s", points[3].Value, quadTree.node.subNodes[3].points[0].Value)
	}
}

func TestQuadTree_QueryRadius(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(10)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "Close"},
		{x: -1, y: 1, Value: "Close"},
		{x: -1, y: -1, Value: "Close"},
		{x: 1, y: -1, Value: "Close"},
		{x: 10, y: 10, Value: "Far"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	result := quadTree.QueryRadius(0, 0, 2)

	if len(result) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result))
	}

	for _, point := range result {
		if point.Value == "Far" {
			t.Errorf("Expected only close points, got %s", point.Value)
		}

	}
}

func TestQuadTree_QueryRadius_NoResults(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(10)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "Close"},
		{x: -1, y: 1, Value: "Close"},
		{x: -1, y: -1, Value: "Close"},
		{x: 1, y: -1, Value: "Close"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	result := quadTree.QueryRadius(10, 10, 2)

	if len(result) != 0 {
		t.Errorf("Expected 0 results, got %d", len(result))
	}
}

func TestQuadTree_QueryRadius_SubNodes(t *testing.T) {
	quadTree := newQuadTreeWithMaxObjects(3)
	points := []*Point[string]{
		{x: 1, y: 1, Value: "Close"},
		{x: 2, y: 1, Value: "Close"},
		{x: 1, y: 2, Value: "Close"},
		{x: 2, y: 2, Value: "Close"},
		{x: -10, y: -10, Value: "Far"},
	}

	for _, point := range points {
		quadTree.Insert(point)
	}

	result := quadTree.QueryRadius(0, 0, 3)
	if len(result) != 4 {
		t.Errorf("Expected 4 results, got %d", len(result))
	}

	for _, point := range result {
		if point.Value == "Far" {
			t.Errorf("Expected only close points, got %s", point.Value)
		}
	}
}

func verifyNonNilSubNodes(t *testing.T, subNodes [4]*node[string]) {
	for _, subNode := range subNodes {
		if subNode == nil {
			t.Errorf("Expected non-nil subNode, got %v", subNode)
		}
	}

}

func verifyNodePointsLength(t *testing.T, node node[string], expectedLength int) {
	if len(node.points) != expectedLength {
		t.Errorf("Expected %d points, got %d", expectedLength, len(node.points))
	}
}

func verifyAllPointsLength(t *testing.T, quadTree QuadTree[string], expectedLength int) {
	if len(quadTree.allPoints) != expectedLength {
		t.Errorf("Expected %d in allPoints, got %d", expectedLength, len(quadTree.allPoints))
	}
}

func newQuadTreeWithMaxObjects(maxObjects int) QuadTree[string] {
	return QuadTree[string]{
		node: node[string]{
			maxObjects: maxObjects,
		},
	}
}
