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

func (m *Map) GetListWalls() [][2]int {
	walls := make([][2]int, 0)
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.GetCell(Position{X: x, Y: y}) == WallCell {
				walls = append(walls, [2]int{x, y})
			}
		}
	}
	return walls
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

type JumpPointSearch struct {
	Map         *Map
	OpenSet     NodeHeap
	Start, Goal *Node
	GScore      map[Position]float64
}

func NewJumpPointSearch(m *Map, start, goal *Node) *JumpPointSearch {
	jps := &JumpPointSearch{
		Map:     m,
		OpenSet: make(NodeHeap, 0),
		Start:   start,
		Goal:    goal,
		GScore:  make(map[Position]float64),
	}
	heap.Init(&jps.OpenSet)
	jps.GScore[start.Pos] = 0
	jps.OpenSet.Push(start)
	return jps
}

func (jps *JumpPointSearch) Search() ([]*Node, bool) {
	for jps.OpenSet.Len() > 0 {
		current := heap.Pop(&jps.OpenSet).(*Node)
		if current.Pos == jps.Goal.Pos {
			return reconstructPath(current)
		}

		successors := jps.successors(current)
		for _, successor := range successors {
			tentativeG := jps.GScore[current.Pos] + jps.Map.GetCostToMoveFactory(current.Pos)(successor.Pos)
			if tentativeG < jps.GScore[successor.Pos] || jps.GScore[successor.Pos] == 0 {
				jps.GScore[successor.Pos] = tentativeG
				successor.Heuristic = calculateHeuristic(successor, jps.Goal)
				heap.Push(&jps.OpenSet, successor)
			}
		}
	}
	return nil, false
}

func (jps *JumpPointSearch) successors(node *Node) []*Node {
	successors := make([]*Node, 0)
	for _, dir := range Directions {
		neighborPos := node.Pos
		for {
			neighborPos = Position{X: neighborPos.X + dir[0], Y: neighborPos.Y + dir[1]}
			if !jps.Map.IsOkToMoveTo(neighborPos) || neighborPos == jps.Goal.Pos {
				break
			}
			jumpPoint := jps.jump(neighborPos, dir)
			if jumpPoint != nil {
				successor := &Node{
					Pos:    *jumpPoint,
					Parent: node,
					Cost:   node.Cost + jps.Map.GetCostToMoveFactory(node.Pos)(*jumpPoint),
				}
				successors = append(successors, successor)
			}
		}
	}
	return successors
}

func (jps *JumpPointSearch) jump(neighborPos Position, dir [2]int) *Position {
	if !jps.Map.IsOkToMoveTo(neighborPos) {
		return nil
	}
	if neighborPos == jps.Goal.Pos {
		return &neighborPos
	}

	nextPos := Position{X: neighborPos.X + dir[0], Y: neighborPos.Y + dir[1]}
	if !jps.Map.IsOkToMoveTo(nextPos) || !jps.Map.IsOkToMoveTo(Position{X: nextPos.X - dir[1], Y: nextPos.Y - dir[0]}) {
		return &neighborPos
	}

	return jps.jump(nextPos, dir)
}

func reconstructPath(node *Node) ([]*Node, bool) {
	path := make([]*Node, 0)
	for node != nil {
		path = append(path, node)
		node = node.Parent
	}
	return reversePath(path), true
}

func reversePath(path []*Node) []*Node {
	reversed := make([]*Node, len(path))
	for i, j := 0, len(path)-1; j >= 0; i, j = i+1, j-1 {
		reversed[i] = path[j]
	}
	return reversed
}
