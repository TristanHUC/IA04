package astar

import (
	"fmt"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
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

func MapToAstarGrid(m _map.Map) [][]uint8 {
	// Create an empty grid
	grid := make([][]uint8, m.Height)
	fmt.Println("Height : ", m.Height)
	fmt.Println("Width : ", m.Width)
	for i := range grid {
		grid[i] = make([]uint8, m.Width)
	}
	// Fill the grid with walls
	for _, wall := range m.Walls {
		grid[wall[1]][wall[0]] = 1
	}
	return grid
}
