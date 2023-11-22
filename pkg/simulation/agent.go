package simulation

import (
	"errors"
	"fmt"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	"math"
	"time"
)

type Agent struct {
	X, Y, vx, vy, gx, gy, speed, reactivity float64 // je pense qu'on peut retirer les vx, vy, gx, gy, tx, ty des attributs
	tx, ty                                  float64
	controllable                            bool
	Path                                    []*astar.Node
	currentWayPoint                         int
	Goal                                    *astar.Node
	start                                   *astar.Node
	channelAgent                            chan []*Agent
	perceptChannel                          chan PerceptRequest
	picMap                                  *astar.Map
}

type PerceptRequest struct {
	Agt             *Agent
	ResponseChannel chan []*Agent
}

func NewAgent(xStart, yStart float64, xGoal, yGoal int, picMap *astar.Map, perceptChannel chan PerceptRequest) *Agent {
	return &Agent{
		X:              xStart,
		Y:              yStart,
		speed:          2,
		reactivity:     1,
		controllable:   true,
		channelAgent:   make(chan []*Agent, 1),
		perceptChannel: perceptChannel,
		start:          &astar.Node{Pos: astar.Position{X: int(xStart) / 7, Y: int(yStart) / 7}},
		Goal:           &astar.Node{Pos: astar.Position{X: xGoal, Y: yGoal}},
		picMap:         picMap,
	}
}

func (a *Agent) Run() {
	for {
		if a.Path == nil {
			goalX, goalY := generateValidCoordinates(a.picMap.GetListWalls())
			a.Goal = &astar.Node{Pos: astar.Position{X: goalX, Y: goalY}}
			a.calculatePath()
		}
		a.calculatePosition()
		time.Sleep(16 * time.Millisecond)
	}
}

func (a *Agent) calculatePath() error {
	// find route to goal
	walls := a.picMap.GetListWalls()
	closeWalls := make([][2]int, 0)
	for _, wall := range walls {
		normeEucli := math.Sqrt((float64(wall[0]*7)+3.5-a.X)*(float64(wall[0]*7)+3.5-a.X) + (float64(wall[1]*7)+3.5-a.Y)*(float64(wall[1]*7)+3.5-a.Y))
		if normeEucli < 50 {
			closeWalls = append(closeWalls, wall)
		}
	}
	searcher := astar.NewJumpPointSearch(a.picMap, a.start, a.Goal)
	path, found := searcher.Search()
	if !found {
		return errors.New("no path found")
	}
	a.Path = path
	a.currentWayPoint = 1
	fmt.Printf("goal is now %d,%d, path length is %d from position %d,%d\n", a.Goal.Pos.X, a.Goal.Pos.Y, len(a.Path), a.start.Pos.X, a.start.Pos.Y)
	return nil
}

func (a *Agent) calculatePosition() error {
	var wayPoint *astar.Node

	// move agent towards current waypoint at a speed of 2px per frame
	if a.currentWayPoint < len(a.Path) {
		wayPoint = a.Path[a.currentWayPoint]
		a.vx = float64(wayPoint.Pos.X*7) + 3.5 - a.X
		a.vy = float64(wayPoint.Pos.Y*7) + 3.5 - a.Y
		vNorm := math.Sqrt(a.vx*a.vx + a.vy*a.vy)

		a.vx = a.vx / vNorm
		a.vy = a.vy / vNorm

		//prise en compte des murs
		var vectx, vecty, normeEucli, reactionMurX, reactionMurY float64
		listeMur := a.picMap.GetListWalls()
		for _, mur := range listeMur {
			vectx = float64(mur[0])*7 + 3.5 - a.X
			vecty = float64(mur[1])*7 + 3.5 - a.Y
			normeEucli = math.Sqrt((float64(mur[0]*7)+3.5-a.X)*(float64(mur[0]*7)+3.5-a.X) + (float64(mur[1]*7)+3.5-a.Y)*(float64(mur[1]*7)+3.5-a.Y))
			vectx = vectx / normeEucli
			vecty = vecty / normeEucli

			reactionMurX = vectx * (3 * math.Exp(-normeEucli/2)) * 10
			reactionMurY = vecty * (3 * math.Exp(-normeEucli/2)) * 10

			a.vx -= reactionMurX
			a.vy -= reactionMurY
		}

		// change velocity to avoid other agents following moussaïd 2009
		var closeAgents []*Agent
		a.perceptChannel <- PerceptRequest{Agt: a, ResponseChannel: a.channelAgent}
		closeAgents = <-a.channelAgent
		for _, otherAgent := range closeAgents {
			lambda := 2.0
			A := 4.5
			gamma := 0.35
			n := 2.0
			np := 3.0
			factor := 0.15

			//pour ne pas recalculer la distance on pourrait la passé dans le channel via un dictionnaire ? A discuter
			dist := math.Sqrt((a.X*factor-otherAgent.X*factor)*(a.X*factor-otherAgent.X*factor) + (a.Y*factor-otherAgent.Y*factor)*(a.Y*factor-otherAgent.Y*factor))

			ex := (otherAgent.X*factor - a.X*factor) / dist
			ey := (otherAgent.Y*factor - a.Y*factor) / dist
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
			theta := signedAcos(math.Min(math.Max(ex*tx+ey*ty, -1), 1))
			B := gamma * DNorm
			addedToVX := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*tx + math.Exp(-math.Pow(n*B*theta, 2))*nx) / factor
			addedToVY := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*ty + math.Exp(-math.Pow(n*B*theta, 2))*ny) / factor
			// safeguard against too big values
			if math.Abs(addedToVX) > 10 {
				addedToVX = 0
			}
			if math.Abs(addedToVY) > 10 {
				addedToVY = 0
			}
			a.vx += addedToVX
			a.vy += addedToVY

			//fmt.Println(addedToVX, addedToVY)

		}

		//fmt.Println(a.vx, a.vy)

		a.X += a.vx
		a.Y += a.vy

		//fmt.Println(a.X, a.Y)

		//passage à l'étape d'après :
		if math.Sqrt((float64(wayPoint.Pos.X*7)+3.5-a.X)*(float64(wayPoint.Pos.X*7)+3.5-a.X)+(float64(wayPoint.Pos.Y*7)+3.5-a.Y)*(float64(wayPoint.Pos.Y*7)+3.5-a.Y)) < 2 {
			a.currentWayPoint++
		}
	}
	return nil
}
