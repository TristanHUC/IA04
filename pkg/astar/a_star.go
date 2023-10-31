package astar

import (
	"container/heap"
	"math"
)

type CellType int

const (
	EmptyCell CellType = iota
	WallCell
	DenseCell
)

type Position struct {
	X int
	Y int
}

type Map struct {
	Width  int
	Height int
	Grid   [][]CellType
}

func (m *Map) GetCell(position Position) CellType {
	return m.Grid[position.Y][position.X]
}

func (m *Map) SetCell(position Position, cellType CellType) {
	m.Grid[position.Y][position.X] = cellType
}

func (m *Map) IsOkToMoveTo(position Position) bool {
	return position.X >= 0 && position.X < m.Width && position.Y >= 0 && position.Y < m.Height && m.GetCell(position) != WallCell
}

type Node struct {
	Pos       Position
	Parent    *Node
	Cost      float64
	Heuristic float64
}

type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].Cost+h[i].Heuristic < h[j].Cost+h[j].Heuristic }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*Node))
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

var Directions = [][2]int{
	{-1, 0},  // North
	{-1, 1},  // NorthEast
	{0, 1},   // East
	{1, 1},   // SouthEast
	{1, 0},   // South
	{1, -1},  // SouthWest
	{0, -1},  // West
	{-1, -1}, // NorthWest
}

func calculateHeuristic(n1, n2 *Node) float64 {
	return math.Sqrt(math.Pow(float64(n1.Pos.X-n2.Pos.X), 2) + math.Pow(float64(n1.Pos.Y-n2.Pos.Y), 2))
}

func (m *Map) GetCostToMoveFactory(from Position) func(to Position) float64 {
	return func(to Position) float64 {
		var costFactor float64
		switch m.GetCell(to) {
		case EmptyCell:
			costFactor = 1
		case DenseCell:
			costFactor = 2
		default:
			costFactor = 1
		}
		return costFactor * math.Sqrt(math.Pow(float64(from.X-to.X), 2)+math.Pow(float64(from.Y-to.Y), 2))
	}
}

func CreateNewNodes(m *Map, n *Node, goal *Node) ([]*Node, error) {
	costFunction := m.GetCostToMoveFactory(n.Pos)
	newNodes := make([]*Node, 0)

	for _, direction := range Directions {
		newPos := Position{X: n.Pos.X + direction[0], Y: n.Pos.Y + direction[1]}

		if m.IsOkToMoveTo(newPos) {
			newNode := &Node{
				Pos:       newPos,
				Cost:      n.Cost + costFunction(newPos),
				Heuristic: calculateHeuristic(&Node{Pos: newPos}, goal),
				Parent:    n,
			}
			newNodes = append(newNodes, newNode)
		}
	}
	return newNodes, nil
}

func NewMap(width, height int) *Map {
	grid := make([][]CellType, height)
	for i := range grid {
		grid[i] = make([]CellType, width)
	}
	return &Map{
		Width:  width,
		Height: height,
		Grid:   grid,
	}
}

func AStar(m *Map, start, goal *Node) ([]*Node, bool) {
	openSet := &NodeHeap{}
	heap.Init(openSet)
	heap.Push(openSet, start)

	start.Cost = 0
	start.Heuristic = calculateHeuristic(start, goal)

	closedSet := make(map[Position]bool)

	for openSet.Len() > 0 {
		currentNode := heap.Pop(openSet).(*Node)
		if currentNode.Pos == goal.Pos {
			path := make([]*Node, 0)
			for currentNode != nil {
				path = append(path, currentNode)
				currentNode = currentNode.Parent
			}
			return path, true
		}

		closedSet[currentNode.Pos] = true

		newNodes, err := CreateNewNodes(m, currentNode, goal)
		if err != nil {
			continue
		}

		for _, newNode := range newNodes {
			if _, ok := closedSet[newNode.Pos]; ok {
				continue
			}
			heap.Push(openSet, newNode)
		}
	}

	return nil, false
}