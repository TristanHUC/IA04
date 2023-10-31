package main

import (
	"fmt"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
)

func main() {
	// Create a new map
	width, height := 10, 10
	m := astar.NewMap(width, height)

	// Set the map with some obstacles (walls)
	m.SetCell(astar.Position{X: 3, Y: 3}, astar.WallCell)
	m.SetCell(astar.Position{X: 3, Y: 4}, astar.WallCell)
	m.SetCell(astar.Position{X: 3, Y: 5}, astar.WallCell)
	m.SetCell(astar.Position{X: 3, Y: 6}, astar.WallCell)

	// Set the map with some obstacles (dense cells)
	m.SetCell(astar.Position{X: 5, Y: 3}, astar.DenseCell)
	m.SetCell(astar.Position{X: 5, Y: 4}, astar.DenseCell)
	m.SetCell(astar.Position{X: 5, Y: 5}, astar.DenseCell)
	m.SetCell(astar.Position{X: 5, Y: 6}, astar.DenseCell)
	m.SetCell(astar.Position{X: 7, Y: 4}, astar.DenseCell)
	m.SetCell(astar.Position{X: 7, Y: 5}, astar.DenseCell)
	m.SetCell(astar.Position{X: 7, Y: 7}, astar.DenseCell)

	// Define the start and goal astar.Positions
	start := &astar.Node{Pos: astar.Position{X: 1, Y: 1}}
	goal := &astar.Node{Pos: astar.Position{X: 8, Y: 8}}

	// Find the path using A* algorithm
	path, found := astar.AStar(m, start, goal)

	if found {
		// Create a grid to visualize the map
		grid := make([][]rune, height)
		for i := range grid {
			grid[i] = make([]rune, width)
			for j := range grid[i] {
				switch m.GetCell(astar.Position{X: j, Y: i}) {
				case astar.EmptyCell:
					grid[i][j] = '.'
				case astar.WallCell:
					grid[i][j] = '#'
				case astar.DenseCell:
					grid[i][j] = 'X'
				}
			}
		}
		// Mark the path on the grid
		for _, node := range path {
			grid[node.Pos.Y][node.Pos.X] = '*'
		}

		// Print the map with the path
		for i := range grid {
			fmt.Println(string(grid[i]))
		}
		fmt.Println("Path found:")
		for i := len(path) - 1; i >= 0; i-- {
			fmt.Printf("X: %d, Y: %d\n", path[i].Pos.X, path[i].Pos.Y)
		}
	} else {
		fmt.Println("Path not found")
	}

}
