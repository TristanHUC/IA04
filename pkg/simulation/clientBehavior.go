package simulation

import (
	"github.com/ankurjha7/jps"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/exp/slices"
	"math"
	"math/rand"
	"time"
)

type ClientBehavior struct{}

func (ClientBehavior) CoordinatesGenerator(m _map.Map) (float64, float64) {
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

func (ClientBehavior) Reflect(a *Agent) {
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

func (ClientBehavior) Act(a *Agent) {
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
	if a.action == WaitForBeer && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
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

// WaitForBeer listen to the Beer channel, if a Beer is received, drink it
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
