package simulation

import (
	"errors"
	"fmt"
	"github.com/ankurjha7/jps"
	"github.com/go-faker/faker/v4"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"math"
	"math/rand"
	"reflect"
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
	GoToExit
)

//enum for the type of agent

type TypeAgent int

const (
	ClientTypeAgent TypeAgent = iota
	BarmanTypeAgent
)

type Behavior interface {
	CoordinatesGenerator(m _map.Map, isLaterGenerated bool) (float64, float64)
	Reflect(a *Agent)
	Act(a *Agent)
}
type Agent struct {
	ID                                      int
	X, Y, Vx, Vy, gx, gy, Speed, reactivity float64
	tx, ty                                  float64
	controllable                            bool
	Path                                    []jps.Node
	CurrentWayPoint                         int
	Goal                                    *jps.Node
	start                                   *jps.Node
	channelAgent                            chan []*Agent
	perceptChannel                          chan PerceptRequest
	PerceptExitChannel                      chan Action
	PerceptPeeChannel                       chan bool
	BeerChannel                             chan bool
	BeerCounterChan                         chan bool
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
	Action                                  Action
	Behavior                                Behavior
	closeAgents                             []*Agent
	client                                  *Agent
	hasABarman                              bool
	endOfLife                               bool
	Paused                                  bool
	Name                                    string
	justPie                                 bool
	woman                                   bool
}

type PerceptRequest struct {
	Agt             *Agent
	ResponseChannel chan []*Agent
}

func NewAgent(ID int, behavior Behavior, picMapDense [][]uint8, picMapSparse *_map.Map, perceptChannel chan PerceptRequest, isLaterGenerated bool, BeerChanCounter chan bool) *Agent {
	agent := &Agent{
		ID:                 ID,
		Speed:              float64(rand.Intn(1)+1) / 30,
		reactivity:         0.1,
		controllable:       true,
		channelAgent:       make(chan []*Agent, 1),
		BeerChannel:        make(chan bool, 1),
		perceptChannel:     perceptChannel,
		PerceptExitChannel: make(chan Action, 1),
		PerceptPeeChannel:  make(chan bool, 1),
		picMapDense:        picMapDense,
		picMapSparse:       picMapSparse,
		lastExecutionTime:  time.Now(),
		DrinkContents:      0, // in milliliters
		timeBetweenDrinks:  time.Duration(rand.Intn(2)) * time.Second,
		drinkSpeed:         0.1,
		drinkEmptyTime:     time.Now(),
		BladderContents:    0, // in milliliters
		BloodAlcoholLevel:  0,
		Action:             None,
		closeAgents:        make([]*Agent, 0),
		client:             nil,
		hasABarman:         false,
		endOfLife:          false,
		Behavior:           behavior,
		Name:               faker.FirstName() + " " + faker.LastName(),
		justPie:            false,
	}
	if rand.Intn(2) == 1 {
		agent.woman = true
	}
	agent.X, agent.Y = agent.Behavior.CoordinatesGenerator(*picMapSparse, isLaterGenerated)
	return agent
}

func (a *Agent) Percept() {
	a.perceptChannel <- PerceptRequest{Agt: a, ResponseChannel: a.channelAgent}
	a.closeAgents = <-a.channelAgent
}

func (a *Agent) PerceptOrderExit() {
	a.Action = <-a.PerceptExitChannel
}

func (a *Agent) PerceptEndOfPieEnhanced() {
	a.justPie = <-a.PerceptPeeChannel
	if a.justPie == true {
		time.Sleep(6 * time.Second)
		a.justPie = false
	}
	go a.PerceptEndOfPieEnhanced()
}

func (a *Agent) Run() {
	pathNotCalculatedYet := true
	go a.PerceptOrderExit()
	go a.PerceptEndOfPieEnhanced()
	for !a.endOfLife {
		if a.Paused {
			time.Sleep(1 * time.Millisecond)
			continue
		}
		if a.lastExecutionTime.Add(17 * time.Millisecond).Before(time.Now()) {
			a.Percept()
			a.lastExecutionTime = time.Now()

			if a.BloodAlcoholLevel > 0.0001 {
				a.BloodAlcoholLevel -= 0.0001
			}
			a.Behavior.Reflect(a)
			a.Behavior.Act(a)

			if (a.Action == GoToExit) && pathNotCalculatedYet {
				pathNotCalculatedYet = false
				err := a.calculatePath()
				if err != nil {
					fmt.Errorf("error calculating path: %v", err)
				}
			}

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
			wallInteractionDistanceMultiplier := 1.
			if reflect.TypeOf(a.Behavior) == reflect.TypeOf(BarmanBehavior{}) {
				wallInteractionDistanceMultiplier = 1.5
			}
			agentStrengthMultiplier := 1.
			if a.Action == WaitForBeer {
				agentStrengthMultiplier = 20
			}
			if a.justPie == true {
				agentStrengthMultiplier = 2
			}
			err := a.calculatePosition(wallInteractionDistanceMultiplier, 1, agentStrengthMultiplier)
			if err != nil {
				fmt.Errorf("error calculating position: %v", err)
			}
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

// GetClosestBarmenArea returns the closest barmen area to a given client
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

func (a *Agent) calculatePosition(
	wallInteractionDistanceMultiplier float64,
	wallInteractionStrengthMultiplier float64,
	agentStrengthMultiplier float64,
) error {
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
		factor := 1.0 * agentStrengthMultiplier

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

		reactionMurX = vectx * (wallInteractionStrengthMultiplier * 3e5 * math.Exp(-(normeEucli+0.5*wallInteractionDistanceMultiplier)/0.1))
		reactionMurY = vecty * (wallInteractionStrengthMultiplier * 3e5 * math.Exp(-(normeEucli+0.5*wallInteractionDistanceMultiplier)/0.1))

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
