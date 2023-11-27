package main

import (
	"fmt"
	"github.com/ankurjha7/jps"
	"github.com/hajimehoshi/ebiten/v2"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/simulation"
	"golang.org/x/image/colornames"
	"log"
)

type View struct {
	Environment [][]uint8
	start       jps.Node
	goal        jps.Node
	waypoints   []jps.Node
	Walls       [][2]int
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

func PrintGrid(grid [][]uint8) {
	fmt.Println("Grid : ")
	for _, row := range grid {
		fmt.Println(row)
	}
}

func MaxIndexWalls(walls [][2]int) (int, int) {
	maxX := 0
	maxY := 0
	for _, wall := range walls {
		if wall[0] > maxX {
			maxX = wall[0]
		}
		if wall[1] > maxY {
			maxY = wall[1]
		}
	}
	return maxX, maxY
}

func (v *View) Update() error {
	return nil
}

func (v *View) Draw(screen *ebiten.Image) {
	// fill white
	screen.Fill(colornames.White)

	// draw walls (7px thick)
	for _, wall := range v.Walls {
		ebitenvector.DrawFilledRect(screen, float32(wall[0]*7), float32(wall[1]*7), 7, 7, colornames.Black, false)
	}

	// draw start and goal
	ebitenvector.DrawFilledRect(screen, float32(v.start.GetCol()*7), float32(v.start.GetRow()*7), 7, 7, colornames.Red, false)
	ebitenvector.DrawFilledRect(screen, float32(v.goal.GetCol()*7), float32(v.goal.GetRow()*7), 7, 7, colornames.Blue, false)
	// draw lines between future waypoints
	for j := 0; j < len(v.waypoints)-1; j++ {
		ebitenvector.StrokeLine(screen, float32(v.waypoints[j].GetCol()*7)+3.5, float32(v.waypoints[j].GetRow()*7)+3.5, float32(v.waypoints[j+1].GetCol()*7)+3.5, float32(v.waypoints[j+1].GetRow()*7)+3.5, 1, colornames.Black, false)
	}

}

func (v *View) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 700, 700
}

func main() {

	//grid := [][]uint8{
	//	{0, 1, 0, 0, 0},
	//	{0, 0, 1, 0, 0},
	//	{1, 1, 1, 1, 0},
	//	{0, 0, 0, 0, 0},
	//}
	//start := jps.GetNode(0, 0)
	//end := jps.GetNode(3, 2)
	//path, err := jps.AStarWithJump(grid, start, end, 1)
	//if err == nil {
	//	fmt.Printf("Path is : ")
	//	for _, node := range path.Nodes {
	//		fmt.Printf("%d %d -> ", node.GetRow(), node.GetCol())
	//	}
	//	fmt.Printf("\nTotal distance %f", path.Weight)
	//}

	maptest := _map.Map{}
	errLoad := maptest.LoadFromFile("testmap")
	if errLoad != nil {
		fmt.Println(errLoad)
	}
	fmt.Println(MaxIndexWalls(maptest.Walls))
	astarMap := MapToAstarGrid(maptest)
	//fmt.Println("Astar map : ")
	//PrintGrid(astarMap)
	//start := jps.GetNode(0, 0)
	//end := jps.GetNode(3, 2)
	var Nodes []jps.Node
	fmt.Println(Nodes)
	randX, randY := simulation.GenerateValidCoordinates(maptest.Walls, maptest.Width, maptest.Height)
	start := jps.GetNode(randX, randY)
	randX, randY = simulation.GenerateValidCoordinates(maptest.Walls, maptest.Width, maptest.Height)
	end := jps.GetNode(randX, randY)
	fmt.Println("Start : ", start.GetRow(), start.GetCol())
	fmt.Println("End : ", end.GetRow(), end.GetCol())
	path, err := jps.AStarWithJump(astarMap, start, end, 1)
	if err == nil {
		Nodes = path.Nodes
		for _, node := range path.Nodes {
			fmt.Printf("%d %d -> ", node.GetRow(), node.GetCol())
		}
	} else {
		Nodes = nil
		log.Fatalf("error: %v", err)
	}
	//run ebiten
	ebiten.SetWindowSize(700, 700)
	ebiten.SetWindowTitle("JPS")
	view := &View{
		Environment: astarMap,
		start:       start,
		goal:        end,
		waypoints:   Nodes,
		Walls:       maptest.Walls,
	}
	if err := ebiten.RunGame(view); err != nil {
		panic(err)
	}
}
