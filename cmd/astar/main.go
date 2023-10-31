package main

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	"golang.org/x/image/colornames"
	"log"
	"math"
	"math/rand"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

type Simulation struct {
	agentX          float64
	agentY          float64
	path            []*astar.Node
	currentWayPoint int
	walls           [][2]int
}

func (s *Simulation) Update() error {
	if s.path == nil {
		// find route to goal
		m := astar.NewMap(100, 100)
		for _, wall := range s.walls {
			m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
		}
		//m.SetCell(astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}, astar.EmptyCell)
		start := &astar.Node{Pos: astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}}
		goal := &astar.Node{Pos: astar.Position{X: 99, Y: 99}}

		searcher := astar.NewJumpPointSearch(m, start, goal)
		path, found := searcher.Search()
		if !found {
			return errors.New("no path found")
		}
		s.path = path
		fmt.Printf("path calculated, len=%d\n", len(s.path))
		s.currentWayPoint = 1
	}

	// move agent towards current waypoint at a speed of 2px per frame
	if s.currentWayPoint < len(s.path) {
		wayPoint := s.path[s.currentWayPoint]
		vx := float64(wayPoint.Pos.X*7) - s.agentX
		vy := float64(wayPoint.Pos.Y*7) - s.agentY
		vNorm := math.Sqrt(vx*vx + vy*vy)
		s.agentX += vx / vNorm
		s.agentY += vy / vNorm
		if math.Abs(s.agentX-float64(wayPoint.Pos.X*7)) < 1 && math.Abs(s.agentY-float64(wayPoint.Pos.Y*7)) < 1 {
			s.currentWayPoint++
		}
	}
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// draw walls (7px thick)
	for _, wall := range s.walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0]*7), float32(wall[1]*7), 7, 7, colornames.Black, false)
	}

	// draw red circle for goal (99,99)
	ebitenvector.DrawFilledCircle(screen, 99*7, 99*7, 4, colornames.Red, false)

	// draw agent
	ebitenvector.DrawFilledCircle(screen, float32(s.agentX), float32(s.agentY), 4, colornames.Blue, false)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (s *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Pic")

	// choose 3000 random walls positions
	walls := make([][2]int, 3000)
	for i := range walls {
		walls[i] = [2]int{10 + rand.Intn(90), 10 + rand.Intn(90)}
	}
	sim := Simulation{
		agentX: 0,
		agentY: 0,
		walls:  walls,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}

//
//func main() {
//	// Create a new map
//	width, height := 10, 10
//	m := astar.NewMap(width, height)
//
//	// Set the map with some obstacles (walls)
//	m.SetCell(astar.Position{X: 3, Y: 3}, astar.WallCell)
//	m.SetCell(astar.Position{X: 3, Y: 4}, astar.WallCell)
//	m.SetCell(astar.Position{X: 3, Y: 5}, astar.WallCell)
//	m.SetCell(astar.Position{X: 3, Y: 6}, astar.WallCell)
//
//	// Set the map with some obstacles (dense cells)
//	m.SetCell(astar.Position{X: 5, Y: 3}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 5, Y: 4}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 5, Y: 5}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 5, Y: 6}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 7, Y: 4}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 7, Y: 5}, astar.DenseCell)
//	m.SetCell(astar.Position{X: 7, Y: 7}, astar.DenseCell)
//
//	// Define the start and goal astar.Positions
//	start := &astar.Node{Pos: astar.Position{X: 1, Y: 1}}
//	goal := &astar.Node{Pos: astar.Position{X: 8, Y: 8}}
//
//	// Find the path using A* algorithm
//	path, found := astar.AStar(m, start, goal)
//
//	if found {
//		// Create a grid to visualize the map
//		grid := make([][]rune, height)
//		for i := range grid {
//			grid[i] = make([]rune, width)
//			for j := range grid[i] {
//				switch m.GetCell(astar.Position{X: j, Y: i}) {
//				case astar.EmptyCell:
//					grid[i][j] = '.'
//				case astar.WallCell:
//					grid[i][j] = '#'
//				case astar.DenseCell:
//					grid[i][j] = 'X'
//				}
//			}
//		}
//		// Mark the path on the grid
//		for _, node := range path {
//			grid[node.Pos.Y][node.Pos.X] = '*'
//		}
//
//		// Print the map with the path
//		for i := range grid {
//			fmt.Println(string(grid[i]))
//		}
//		fmt.Println("Path found:")
//		for i := len(path) - 1; i >= 0; i-- {
//			fmt.Printf("X: %d, Y: %d\n", path[i].Pos.X, path[i].Pos.Y)
//		}
//	} else {
//		fmt.Println("Path not found")
//	}
//
//}
