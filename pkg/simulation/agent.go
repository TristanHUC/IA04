package simulation

import (
	"errors"
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
	"math"
	"math/rand"
	"time"
)

type Agent struct {
	X, Y, vx, vy, gx, gy, speed, reactivity float64 // je pense qu'on peut retirer les vx, vy, gx, gy, tx, ty des attributs
	tx, ty                                  float64
	controllable                            bool
	Path                                    []*astar.Node
	CurrentWayPoint                         int
	Goal                                    *astar.Node
	start                                   *astar.Node
	channelAgent                            chan []*Agent
	perceptChannel                          chan PerceptRequest
	picMap                                  *astar.Map
	rollingMeanMovement                     float64
	lastExecutionTime                       time.Time
}

type PerceptRequest struct {
	Agt             *Agent
	ResponseChannel chan []*Agent
}

func NewAgent(xStart, yStart float64, xGoal, yGoal int, picMap *astar.Map, perceptChannel chan PerceptRequest) *Agent {
	return &Agent{
		X:                 xStart,
		Y:                 yStart,
		speed:             float64(rand.Intn(5)+5) / 13,
		reactivity:        1,
		controllable:      true,
		channelAgent:      make(chan []*Agent, 1),
		perceptChannel:    perceptChannel,
		start:             &astar.Node{Pos: astar.Position{X: int(xStart) / 7, Y: int(yStart) / 7}},
		Goal:              &astar.Node{Pos: astar.Position{X: xGoal, Y: yGoal}},
		picMap:            picMap,
		lastExecutionTime: time.Now(),
	}
}

func (a *Agent) Run() {
	for {
		if a.lastExecutionTime.Add(17 * time.Millisecond).Before(time.Now()) {
			a.lastExecutionTime = time.Now()
			if a.Path == nil {
				goalX, goalY := generateValidCoordinates(a.picMap.GetListWalls())
				a.Goal = &astar.Node{Pos: astar.Position{X: goalX, Y: goalY}}
				a.calculatePath()
				//for err != nil {
				//	goalX, goalY := generateValidCoordinates(a.picMap.GetListWalls())
				//	a.Goal = &astar.Node{Pos: astar.Position{X: goalX, Y: goalY}}
				//	err = a.calculatePath()
				//}
			}
			// if rolling mean movement is too low, recalculate path
			if a.rollingMeanMovement < 0.1 {
				a.calculatePath()
				//fmt.Println("stuck, recalculating path")
			}
			a.calculatePosition()
		} else {
			time.Sleep(1 * time.Millisecond)
		}
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
	start := &astar.Node{Pos: astar.Position{X: int(a.X) / 7, Y: int(a.Y) / 7}}
	searcher := astar.NewJumpPointSearch(a.picMap, start, a.Goal)
	path, found := searcher.Search()
	if !found {
		return errors.New("no path found")
	}
	a.Path = path
	a.CurrentWayPoint = 1
	return nil
}

func (a *Agent) calculatePosition() error {
	var wayPoint *astar.Node

	// move agent towards current waypoint at a speed of 2px per frame
	if a.CurrentWayPoint < len(a.Path) {
		wayPoint = a.Path[a.CurrentWayPoint]
		// compute goal velocity (norm = agent speed, direction= towards goal)
		gvx := float64(wayPoint.Pos.X*7) + 3.5 - a.X
		gvy := float64(wayPoint.Pos.Y*7) + 3.5 - a.Y
		gvNorm := math.Sqrt(gvx*gvx + gvy*gvy)
		gvx /= gvNorm
		gvy /= gvNorm
		gvx *= a.speed
		gvy *= a.speed

		// change velocity towards goal velocity at a rate of reactivity
		a.vx += (gvx - a.vx) * a.reactivity
		a.vy += (gvy - a.vy) * a.reactivity

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

			a.vx += addedToVX
			a.vy += addedToVY
		}

		// safeguard against too big values
		if norm := math.Sqrt(a.vx*a.vx + a.vy*a.vy); norm > a.speed {
			a.vx = a.vx / norm * a.speed
			a.vy = a.vy / norm * a.speed
		}

		a.X += a.vx
		a.Y += a.vy

		a.rollingMeanMovement = (a.rollingMeanMovement + math.Sqrt(a.vx*a.vx+a.vy*a.vy)) / 2

		//passage à l'étape d'après :
		if math.Sqrt((float64(wayPoint.Pos.X*7)+3.5-a.X)*(float64(wayPoint.Pos.X*7)+3.5-a.X)+(float64(wayPoint.Pos.Y*7)+3.5-a.Y)*(float64(wayPoint.Pos.Y*7)+3.5-a.Y)) < 2 {
			a.CurrentWayPoint++
		}
		if a.CurrentWayPoint >= len(a.Path) {
			a.Path = nil
			a.CurrentWayPoint = 0
		}
	}
	return nil
}
