package main

import (
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/image/colornames"
	"log"
	"math"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700


func signedAcos(x float64) float64 {
	unsignedAcos := math.Acos(x)
	if x >= 0 {
		return unsignedAcos
	} else {
		return -unsignedAcos
	}
}


type Agent struct {
	x, y, vx, vy, gx, gy, speed, reactivity float64 // je pense qu'on peut retirer les vx, vy, gx, gy, tx, ty des attributs 
	tx, ty                                  float64
	controllable                            bool
	path            						[]*astar.Node
	currentWayPoint 						int
	goal 									*astar.Node
	start									*astar.Node	
	channelAgent							chan []*Agent
	channelMur								chan [][2]int
}


func (a *Agent)calculePosition ()error{

	var wayPoint *astar.Node

	// move agent towards current waypoint at a speed of 2px per frame
	if a.currentWayPoint < len(a.path) {
		wayPoint = a.path[a.currentWayPoint]
		a.vx = float64(wayPoint.Pos.X*7) + 3.5 - a.x
		a.vy = float64(wayPoint.Pos.Y*7) + 3.5 - a.y
		vNorm := math.Sqrt(a.vx*a.vx + a.vy*a.vy)

		a.vx = a.vx / vNorm
		a.vy = a.vy / vNorm

		//prise en compte des murs
		var vectx, vecty, normeEucli, reactionMurX, reactionMurY  float64
		var listeMur [][2]int
		listeMur = <- a.channelMur
		for _, mur := range listeMur {
			vectx = float64(mur[0])*7 + 3.5 - a.x
			vecty = float64(mur[1])*7 + 3.5 - a.y
			normeEucli = math.Sqrt((float64(mur[0]*7)+3.5-a.x)*(float64(mur[0]*7)+3.5-a.x) + (float64(mur[1]*7)+3.5-a.y)*(float64(mur[1]*7)+3.5-a.y))
			vectx = vectx / normeEucli
			vecty = vecty / normeEucli

			reactionMurX = vectx * (3 * math.Exp(-normeEucli/2)) * 10
			reactionMurY = vecty * (3 * math.Exp(-normeEucli/2)) * 10

			a.vx -= reactionMurX
			a.vy -= reactionMurY
		}


		// change velocity to avoid other agents following moussaïd 2009
		var listeAgent []*Agent
		listeAgent = <- a.channelAgent
		for _,otherAgent := range listeAgent{
			lambda := 2.0
			A := 4.5
			gamma := 0.35
			n := 2.0
			np := 3.0
			factor := 0.15

			//pour ne pas recalculer la distance on pourrait la passé dans le channel via un dictionnaire ? A discuter
			dist := math.Sqrt((a.x*factor-otherAgent.x*factor)*(a.x*factor-otherAgent.x*factor) + (a.y*factor-otherAgent.y*factor)*(a.y*factor-otherAgent.y*factor))
			
			ex := (otherAgent.x*factor - a.x*factor) / dist
			ey := (otherAgent.y*factor - a.y*factor) / dist
			Dx := lambda*(a.vx*factor-otherAgent.vx*factor) + ex
			Dy := lambda*(a.vy*factor-otherAgent.vy*factor) + ey
			DNorm := math.Sqrt(Dx*Dx + Dy*Dy)
			tx := Dx / DNorm
			ty := Dy / DNorm
			// nx, ny is the normal vector to tx,ty pointing to the left
			nx := ty
			ny := -tx
			a.tx = nx
			a.ty = ny
			//theta := math.Acos(math.Min(math.Max(ex*tx+ey*ty, -1), 1)) / (2 * math.Pi) * 360
			theta := signedAcos(math.Min(math.Max(ex*tx+ey*ty, -1), 1))
			B := gamma * DNorm
			addedToVX := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*tx + math.Exp(-math.Pow(n*B*theta, 2))*nx) / factor
			addedToVY := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*ty + math.Exp(-math.Pow(n*B*theta, 2))*ny) / factor
			// safeguard against too big values
			if addedToVX > 10 {
				addedToVX = 0
			}
			if addedToVY > 10 {
				addedToVY = 0
			}
			a.vx += addedToVX
			a.vy += addedToVY

		}
		a.x += a.vx
		a.y += a.vy

		//passage à l'étape d'après :
		if math.Sqrt((float64(wayPoint.Pos.X*7)+3.5-a.x)*(float64(wayPoint.Pos.X*7)+3.5-a.x)+(float64(wayPoint.Pos.Y*7)+3.5-a.y)*(float64(wayPoint.Pos.Y*7)+3.5-a.y)) < 2 {
			a.currentWayPoint++
		}
	} else {
		<- a.channelMur
		<- a.channelAgent

	}
	return nil
}

type Simulation struct {
	agents  		[]*Agent
	nAgents 		int
	walls           [][2]int
	m 				*astar.Map
}

func (s *Simulation) Update() error {

	//initialization
	if s.m == nil {
		// find route to goal
		s.m = astar.NewMap(100, 100)
		for _, wall := range s.walls {
			s.m.SetCell(astar.Position{X: wall[0], Y: wall[1]}, astar.WallCell)
		}

		var searcher *astar.JumpPointSearch
		var path []*astar.Node
		var found bool

		for i:=0;i<s.nAgents; i++{
		
			s.agents[i].channelAgent = make(chan []*Agent,1)
			s.agents[i].channelMur = make(chan [][2]int,1)

			//s.m.SetCell(astar.Position{X: int(s.agentX / 7), Y: int(s.agentY / 7)}, astar.EmptyCell)
			s.agents[i].start = &astar.Node{Pos: astar.Position{X: int(s.agents[i].x / 7), Y: int(s.agents[i].y / 7)}}
			s.agents[i].goal = &astar.Node{Pos: astar.Position{X: 99, Y: 99}}
	
			searcher = astar.NewJumpPointSearch(s.m, s.agents[i].start, s.agents[i].goal)
			path, found = searcher.Search()
			if !found {
				return errors.New("no path found")
			}
			s.agents[i].path = path
			fmt.Printf("path calculated, len=%d\n", len(s.agents[i].path))
			s.agents[i].currentWayPoint = 1
	
		}

	}

	//update
	for i:=0;i<s.nAgents; i++{

		listeAgent := make([]*Agent, 0, s.nAgents)
		listeMur := make([][2]int, 0)

		//loading agent'channel with walls around them 
		for _,mur := range s.walls{
			normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-s.agents[i].x)*(float64(mur[0]*7)+3.5-s.agents[i].x) + (float64(mur[1]*7)+3.5-s.agents[i].y)*(float64(mur[1]*7)+3.5-s.agents[i].y))
			if normeEucli < 50 {
				listeMur = append(listeMur, mur)
			}
		}
		s.agents[i].channelMur <- listeMur		

		//loading agent'channel with other agents around them 
		for _,otherAgent := range s.agents{
			if otherAgent == s.agents[i] {
				continue
			}
			factor := 0.15
			dist := math.Sqrt((s.agents[i].x*factor-otherAgent.x*factor)*(s.agents[i].x*factor-otherAgent.x*factor) + (s.agents[i].y*factor-otherAgent.y*factor)*(s.agents[i].y*factor-otherAgent.y*factor))
			if dist < 10 {
				listeAgent = append(listeAgent, otherAgent)
			}
		}
		s.agents[i].channelAgent <- listeAgent

		go s.agents[i].calculePosition()

		
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

	//draw agents, their position and their goals
	for i:=0;i<s.nAgents; i++{

		// draw red circle for goal (99,99)
		ebitenvector.DrawFilledCircle(screen, float32(s.agents[i].goal.Pos.X*7), float32(s.agents[i].goal.Pos.Y*7), 4, colornames.Red, false)

		// draw agent
		ebitenvector.DrawFilledCircle(screen, float32(s.agents[i].x), float32(s.agents[i].y), 4, colornames.Blue, false)
		
		// draw lines between waypoints
		for j := 0; j < len(s.agents[i].path)-1; j++ {
			ebitenvector.StrokeLine(screen, float32(s.agents[i].path[j].Pos.X*7)+3.5, float32(s.agents[i].path[j].Pos.Y*7)+3.5, float32(s.agents[i].path[j+1].Pos.X*7)+3.5, float32(s.agents[i].path[j+1].Pos.Y*7)+3.5, 1, colornames.Green, false)
		}
		for _, mur := range s.walls {
			normeEucli := math.Sqrt((float64(mur[0]*7)+3.5-s.agents[i].x)*(float64(mur[0]*7)+3.5-s.agents[i].x) + (float64(mur[1]*7)+3.5-s.agents[i].y)*(float64(mur[1]*7)+3.5-s.agents[i].y))
			if normeEucli < 50 {
				color := colornames.Blue
				color.A = 50
				//color.R -= uint8(normeEucli / 5)
				//color.G -= uint8(normeEucli / 5)
				//color.B -= uint8(normeEucli / 5)
				ebitenvector.StrokeLine(screen, float32(s.agents[i].x), float32(s.agents[i].y), float32(mur[0])*7+3.5, float32(mur[1])*7+3.5, 1, color, false)
			}
		}
	}

	
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

	// load map from file
	testmap := _map.Map{}
	err := testmap.LoadFromFile("testmap")
	if err != nil {
		return
	}

	nAgents := 10
	agents := make([]*Agent,0,5)
	for i:=0;i<nAgents; i++{
		agents = append(agents,&Agent{x: float64(70+35*i), y: float64(70)})
	}

	sim := Simulation{
		agents :		agents,
		nAgents :		nAgents,
		walls:  testmap.Walls,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
