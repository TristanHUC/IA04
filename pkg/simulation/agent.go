package simulation

import (
	"errors"
	"fmt"
	"github.com/ankurjha7/jps"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/exp/slices"
	"math"
	"math/rand"
	"time"
)

type Action int

const (
	None Action = iota
	GoToRandomSpot
	GoToToilet
	GoToBar
	GoToBeerTap
	WaitForBeer
	WaitForClient
	GoToClient
)

//enum for the type of agent

type TypeAgent int

const (
	ClientTypeAgent TypeAgent = iota
	BarmanTypeAgent
)

type Behavior struct {
	CoordinatesGenerator func(m _map.Map) (float64, float64)
	Reflect              func()
	Act                  func()
}

func (a *Agent) ClientCoordinatesGenerator(m _map.Map) (float64, float64) {
	x := rand.Intn(m.Width)
	y := rand.Intn(m.Height)
	coordsOk := false
	// while agent is inside a wall, generate new coordinates
	for !coordsOk {
		coordsOk = true
		for _, wall := range m.Walls {
			if wall[0] == x && wall[1] == y {
				x = rand.Intn(m.Width)
				y = rand.Intn(m.Height)
				coordsOk = false
			}
		}
		for _, counter := range m.BarmenArea {
			if counter[0] == x && counter[1] == y {
				x = rand.Intn(m.Width)
				y = rand.Intn(m.Height)
				coordsOk = false
			}
		}
	}
	xFloat := float64(x) + rand.Float64()
	yFloat := float64(y) + rand.Float64()
	return xFloat, yFloat
}

func (a *Agent) BarmanCoordinatesGenerator(m _map.Map) (float64, float64) {
	// Take a random point in the bar area
	counterPoints := m.BarmenArea[rand.Intn(len(m.BarmenArea))]
	return float64(counterPoints[0]) + rand.Float64(), float64(counterPoints[1]) + rand.Float64()
}

type Agent struct {
	X, Y, Vx, Vy, gx, gy, Speed, reactivity float64 // je pense qu'on peut retirer les Vx, vy, gx, gy, tx, ty des attributs
	tx, ty                                  float64
	controllable                            bool
	Path                                    []jps.Node
	CurrentWayPoint                         int
	Goal                                    *jps.Node
	start                                   *jps.Node
	channelAgent                            chan []*Agent
	perceptChannel                          chan PerceptRequest
	BeerChannel                             chan bool
	picMapDense                             [][]uint8
	picMapSparse                            *_map.Map
	rollingMeanMovement                     float64
	lastExecutionTime                       time.Time
	DrinkContents                           float64
	timeBetweenDrinks                       time.Duration
	drinkEmptyTime                          time.Time
	BladderContents                         float64
	drinkSpeed                              float64
	BloodAlcoholLevel                       float64
	action                                  Action
	typeAgent                               TypeAgent
	behavior                                Behavior
	closeAgents                             []*Agent
	client                                  *Agent
	hasABarman                              bool
}

type PerceptRequest struct {
	Agt             *Agent
	ResponseChannel chan []*Agent
}

func NewAgent(typeAgent TypeAgent, picMapDense [][]uint8, picMapSparse *_map.Map, perceptChannel chan PerceptRequest) *Agent {
	agent := &Agent{
		typeAgent:         typeAgent,
		Speed:             float64(rand.Intn(1)+1) / 30,
		reactivity:        0.1,
		controllable:      true,
		channelAgent:      make(chan []*Agent, 1),
		BeerChannel:       make(chan bool, 1),
		perceptChannel:    perceptChannel,
		picMapDense:       picMapDense,
		picMapSparse:      picMapSparse,
		lastExecutionTime: time.Now(),
		DrinkContents:     0, // in milliliters
		timeBetweenDrinks: time.Duration(rand.Intn(2)) * time.Second,
		drinkSpeed:        0.1,
		drinkEmptyTime:    time.Now(),
		BladderContents:   0, // in milliliters
		BloodAlcoholLevel: 0,
		action:            None,
		closeAgents:       make([]*Agent, 0),
		client:            nil,
		hasABarman:        false,
	}
	if typeAgent == ClientTypeAgent {
		agent.behavior = Behavior{CoordinatesGenerator: agent.ClientCoordinatesGenerator, Reflect: agent.ClientReflect, Act: agent.ClientAct}
	} else {
		agent.behavior = Behavior{CoordinatesGenerator: agent.BarmanCoordinatesGenerator, Reflect: agent.BarmanReflect, Act: agent.BarmanAct}
	}
	agent.X, agent.Y = agent.behavior.CoordinatesGenerator(*picMapSparse)
	return agent
}

func (a *Agent) Percept() {
	a.perceptChannel <- PerceptRequest{Agt: a, ResponseChannel: a.channelAgent}
	a.closeAgents = <-a.channelAgent
}

func (a *Agent) Run() {
	for {
		if a.lastExecutionTime.Add(17 * time.Millisecond).Before(time.Now()) {
			a.Percept()
			a.lastExecutionTime = time.Now()
			a.BloodAlcoholLevel -= 0.0001
			a.behavior.Reflect()
			a.behavior.Act()
			// if agent has no path but a goal is set, calculate path
			if a.Goal != nil && a.Path == nil {
				err := a.calculatePath()
				if err != nil {
					fmt.Errorf("error calculating path: %v", err)
				}
			}
			// if rolling mean movement is too low, recalculate path (anti-stuck)
			if a.Path != nil && a.rollingMeanMovement < 0.02 {
				err := a.calculatePath()
				if err != nil {
					fmt.Errorf("error calculating path: %v", err)
				}
			}
			// agent reflexes
			a.calculatePosition()
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (a *Agent) Drink() {
	if a.DrinkContents >= a.drinkSpeed {
		a.DrinkContents -= a.drinkSpeed
		a.BladderContents += a.drinkSpeed
		// 1000 for ml -> l, 0.07 for alcohol percentage, 0.78 alcohol density, 5 for liters in the body
		a.BloodAlcoholLevel += (a.drinkSpeed * 1000) * 0.07 * 0.78 / 5
	} else if a.drinkEmptyTime.IsZero() {
		// if drink just finished, set time
		a.drinkEmptyTime = time.Now()
	}
}

func (a *Agent) BarmanReflect() {
	if a.action != None {
		return
	}
	//lauching the search for a client
	a.action = GoToBeerTap

}

// get the closest barmen area in euclidean distance from the attribute client
func (a *Agent) GetClosestBarmenArea(client Agent) jps.Node {
	var closestBarmenArea [2]int
	var minDistance float64 = 100000
	var distance float64
	for _, barmenArea := range a.picMapSparse.BarmenArea {
		distance = distanceInt(barmenArea[0], barmenArea[1], int(client.X), int(client.Y))
		if distance < minDistance {
			minDistance = distance
			closestBarmenArea = barmenArea
		}
	}
	return jps.GetNode(closestBarmenArea[1], closestBarmenArea[0])
}

// may not find a client if there is none
func (a *Agent) SearchForClient() {
	for _, agent := range a.closeAgents {
		if agent.typeAgent == ClientTypeAgent && agent.action == WaitForBeer && !agent.hasABarman {
			a.client = agent
			g := a.GetClosestBarmenArea(*agent)
			a.Goal = &g
			// notify the client that he has a barman
			a.client.BeerChannel <- false
			a.action = GoToClient
			break
		}
	}
}

func (a *Agent) GiveABeer() {
	a.client.BeerChannel <- true
	a.DrinkContents = 0
}

func (a *Agent) BarmanAct() {
	// if agent want to go to beer tap, and current goal does not reflect that, change goal
	if a.action == GoToBeerTap && (a.Goal == nil || !slices.Contains(a.picMapSparse.BeerTaps, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		beerTap := a.picMapSparse.BeerTaps[rand.Intn(len(a.picMapSparse.BeerTaps))]
		g := jps.GetNode(beerTap[1], beerTap[0])
		a.Goal = &g
	}

	// if agent is waiting for a client, he should find one
	if a.action == WaitForClient && a.client == nil {
		a.SearchForClient()
	}

	// if goal is reached
	if a.action != None && a.Goal != nil && distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		a.Path = nil
		a.CurrentWayPoint = 0
		a.Goal = nil
		if a.action == WaitForClient {
			a.action = GoToClient
		} else if a.action == GoToBeerTap {
			a.DrinkContents = 300
			a.action = WaitForClient
		} else if a.action == GoToClient {
			a.GiveABeer()
			a.action = GoToBeerTap
			a.client = nil
		}
	}
}

func (a *Agent) ClientReflect() {
	if a.action != None { // doucement cabron, une action à la fois
		return
	}
	if a.BladderContents > 450 {
		// go to toilet
		a.action = GoToToilet
	}
	if !a.drinkEmptyTime.IsZero() && a.drinkEmptyTime.Add(a.timeBetweenDrinks).Before(a.lastExecutionTime) {
		// go to bar
		a.action = GoToBar
	}
}

// listen to the Beer channel, if a Beer is received, drink it
func (a *Agent) WaitForBeer() {
	var response bool
	response = <-a.BeerChannel
	// a barman has chosen this client
	if !response {
		a.hasABarman = true
		a.WaitForBeer()
	} else {
		a.DrinkContents = 300
		a.hasABarman = false
		a.action = GoToRandomSpot
		goalX, goalY := GenerateValidCoordinates(a.picMapSparse.Walls, a.picMapSparse.Width, a.picMapSparse.Height)
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}
}

func (a *Agent) ClientAct() {
	a.Drink()
	// if agent want to go to toilet, and current goal does not reflect that, change goal
	if a.action == GoToToilet && (a.Goal == nil || !slices.Contains(a.picMapSparse.ManToiletPoints, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		toilet := a.picMapSparse.ManToiletPoints[rand.Intn(len(a.picMapSparse.ManToiletPoints))]
		g := jps.GetNode(toilet[1], toilet[0])
		a.Goal = &g
	}

	// if agent want to go to bar, and current goal does not reflect that, change goal
	if a.action == GoToBar && (a.Goal == nil || !slices.Contains(a.picMapSparse.BarPoints, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		bar := a.picMapSparse.BarPoints[rand.Intn(len(a.picMapSparse.BarPoints))]
		g := jps.GetNode(bar[1], bar[0])
		a.Goal = &g
	}

	// if agent wants to go to random spot, and current goal does not reflect that, change goal
	if a.action == GoToRandomSpot && a.Goal == nil {
		goalX, goalY := GenerateValidCoordinates(a.picMapSparse.Walls, a.picMapSparse.Width, a.picMapSparse.Height)
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}

	// if agent has nothing to do, try to stay still
	if a.action == None && a.Goal == nil {
		goalX, goalY := a.X, a.Y
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}

	if a.action == None && a.Goal != nil {
		vecToGoalX := float64(a.Goal.GetCol()) - a.X
		vecToGoalY := float64(a.Goal.GetRow()) - a.Y
		distToGoal := math.Sqrt(vecToGoalX*vecToGoalX + vecToGoalY*vecToGoalY)
		if distToGoal > 2 {
			goalX, goalY := a.X, a.Y
			g := jps.GetNode(int(goalY), int(goalX))
			a.Goal = &g
			a.Path = nil
		}
	}

	// if agent is waiting for a Beer, doesnt move even if he has reached his goal
	if a.action == WaitForBeer && distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		return
	}

	// if goal is reached
	if a.action != None && a.Goal != nil && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		a.Path = nil
		a.CurrentWayPoint = 0
		a.Goal = nil
		if a.action == GoToToilet {
			a.BladderContents = 0
			a.action = GoToRandomSpot
		} else if a.action == GoToBar {
			a.action = WaitForBeer
			go a.WaitForBeer()
			// try to stay still
			goalX, goalY := a.X, a.Y
			g := jps.GetNode(int(goalY), int(goalX))
			a.Goal = &g
		} else if a.action == GoToRandomSpot {
			a.action = None
		}
	}
}

func (a *Agent) calculatePath() error {
	defer func() { // jps sometimes panics randomly, this is a safeguard
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	// find route to goal
	start := jps.GetNode(int(math.Floor(a.Y)), int(math.Floor(a.X)))
	path, err := jps.AStarWithJump(a.picMapDense, start, *a.Goal, 1)
	if err != nil {
		return errors.New("no path found")
	}
	a.Path = path.Nodes
	// append goal
	a.Path = append(a.Path, *a.Goal)
	a.CurrentWayPoint = 1
	return nil
}

func (a *Agent) calculatePosition() error {
	var wayPoint *jps.Node

	// move agent towards current waypoint at a speed of 2px per frame
	if a.CurrentWayPoint < len(a.Path) {
		wayPoint = &a.Path[a.CurrentWayPoint]
		// compute goal velocity (norm = agent speed, direction= towards goal)
		gvx := float64(wayPoint.GetCol()) - a.X
		gvy := float64(wayPoint.GetRow()) - a.Y
		gvNorm := math.Sqrt(gvx*gvx + gvy*gvy)
		gvx /= gvNorm
		gvy /= gvNorm
		gvx *= a.Speed
		gvy *= a.Speed

		// change velocity towards goal velocity at a rate of reactivity
		a.Vx += (gvx - a.Vx) * a.reactivity
		a.Vy += (gvy - a.Vy) * a.reactivity
	}

	// change velocity to avoid other agents following moussaïd 2009
	for _, otherAgent := range a.closeAgents {
		lambda := 2.0
		A := 4.5
		gamma := 0.2
		n := 2.0
		np := 3.0
		factor := 1.0

		// pour ne pas recalculer la distance on pourrait la passer dans le channel via un dictionnaire ? A discuter
		dist := math.Sqrt((a.X*factor-otherAgent.X*factor)*(a.X*factor-otherAgent.X*factor) + (a.Y*factor-otherAgent.Y*factor)*(a.Y*factor-otherAgent.Y*factor))

		ex := (otherAgent.X*factor - a.X*factor) / dist
		ey := (otherAgent.Y*factor - a.Y*factor) / dist
		Dx := lambda*(a.Vx*factor-otherAgent.Vx*factor) + ex
		Dy := lambda*(a.Vy*factor-otherAgent.Vy*factor) + ey
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

		a.Vx += addedToVX
		a.Vy += addedToVY
	}

	//prise en compte des murs
	var vectx, vecty, normeEucli, reactionMurX, reactionMurY float64
	listeMur := a.picMapSparse.Walls
	for _, mur := range listeMur {
		vectx = float64(mur[0]) + 0.5 - a.X
		vecty = float64(mur[1]) + 0.5 - a.Y
		normeEucli = math.Sqrt(vectx*vectx + vecty*vecty)

		vectx = vectx / normeEucli
		vecty = vecty / normeEucli

		reactionMurX = vectx * (3e5 * math.Exp(-(normeEucli+0.5)/0.1))
		reactionMurY = vecty * (3e5 * math.Exp(-(normeEucli+0.5)/0.1))

		a.Vx -= reactionMurX
		a.Vy -= reactionMurY
	}

	//safeguard against too big values
	if norm := math.Sqrt(a.Vx*a.Vx + a.Vy*a.Vy); norm > a.Speed*1.5 {
		a.Vx = a.Vx / norm * a.Speed
		a.Vy = a.Vy / norm * a.Speed
	}

	a.X += a.Vx
	a.Y += a.Vy

	a.rollingMeanMovement = (a.rollingMeanMovement + math.Sqrt(a.Vx*a.Vx+a.Vy*a.Vy)) / 2

	// passage à l'étape d'après :
	if wayPoint != nil && a.CurrentWayPoint < len(a.Path)-1 && math.Sqrt((float64(wayPoint.GetCol())-a.X)*(float64(wayPoint.GetCol())-a.X)+(float64(wayPoint.GetRow())-a.Y)*(float64(wayPoint.GetRow())-a.Y)) < 1 {
		a.CurrentWayPoint++
	}

	return nil
}
