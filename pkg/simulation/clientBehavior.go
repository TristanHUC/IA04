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

// CoordinatesGenerator for the client behavior generates coordinates inside the walls of the map
func (ClientBehavior) CoordinatesGenerator(m _map.Map, isLaterGenerated bool) (float64, float64) {
	var (
		xFloat, yFloat float64
		x, y           int
	)

	minWallX := math.Inf(1)
	minWallY := math.Inf(1)
	maxWallX := math.Inf(-1)
	maxWallY := math.Inf(-1)

	for _, wall := range m.Walls {
		if float64(wall[0]) > maxWallX {
			maxWallX = float64(wall[0])
		} else if float64(wall[0]) < minWallX {
			minWallX = float64(wall[0])
		}
		if float64(wall[1]) > maxWallY {
			maxWallY = float64(wall[1])
		} else if float64(wall[1]) < minWallY {
			minWallY = float64(wall[1])
		}
	}

	if isLaterGenerated {
		xFloat = float64(m.Enter[rand.Intn(len(m.Enter))][0])
		yFloat = float64(m.Enter[rand.Intn(len(m.Enter))][1])
	} else {
		x = int(minWallX) + rand.Intn(int(maxWallX-minWallX))
		y = int(minWallY) + rand.Intn(int(maxWallY-minWallY))
		coordsOk := false
		// while agent is inside a wall or other things, generate new coordinates
		for !coordsOk {
			x = int(minWallX) + rand.Intn(int(maxWallX-minWallX))
			y = int(minWallY) + rand.Intn(int(maxWallY-minWallY))
			coordsOk = true
			for _, wall := range m.Walls {
				if wall[0] == x && wall[1] == y {
					coordsOk = false
				}
			}
			for _, counter := range m.BarmenArea {
				if counter[0] == x && counter[1] == y {
					coordsOk = false
				}
			}
			for _, beerTap := range m.BeerTaps {
				if beerTap[0] == x && beerTap[1] == y {
					coordsOk = false
				}
			}
		}
		xFloat = float64(x) + rand.Float64()
		yFloat = float64(y) + rand.Float64()
	}
	return xFloat, yFloat
}

func (ClientBehavior) Reflect(a *Agent) {
	if a.Action != None { // doucement cabron, une action à la fois
		return
	}
	if a.BladderContents > 450 {
		// go to toilet
		a.Action = GoToToilet
	} else {
		if !a.drinkEmptyTime.IsZero() && a.drinkEmptyTime.Add(a.timeBetweenDrinks).Before(a.lastExecutionTime) {
			// go to bar
			a.Action = GoToBar
		}
	}
}

func (ClientBehavior) Act(a *Agent) {
	a.Drink()

	// if agent have to leave
	if a.Action == GoToExit && (a.Goal == nil || !slices.Contains(a.picMapSparse.Exit, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		exit := a.picMapSparse.Exit[rand.Intn(len(a.picMapSparse.Exit))]
		g := jps.GetNode(exit[1], exit[0])
		a.Goal = &g
	}

	// if agent want to go to toilet, and current goal does not reflect that, change goal
	if a.Action == GoToToilet && (a.Goal == nil || !(slices.Contains(a.picMapSparse.ManToiletPoints, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())}) || slices.Contains(a.picMapSparse.WomanToiletPoints, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())}))) {
		if a.woman == true {
			toilet := a.picMapSparse.WomanToiletPoints[rand.Intn(len(a.picMapSparse.WomanToiletPoints))]
			g := jps.GetNode(toilet[1], toilet[0])
			a.Goal = &g
		} else {
			toilet := a.picMapSparse.ManToiletPoints[rand.Intn(len(a.picMapSparse.ManToiletPoints))]
			g := jps.GetNode(toilet[1], toilet[0])
			a.Goal = &g
		}
	}

	// if agent want to go to bar, and current goal does not reflect that, change goal
	if a.Action == GoToBar && (a.Goal == nil || !slices.Contains(a.picMapSparse.BarPoints, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		bar := a.picMapSparse.BarPoints[rand.Intn(len(a.picMapSparse.BarPoints))]
		g := jps.GetNode(bar[1], bar[0])
		a.Goal = &g
	}

	// if agent wants to go to random spot, and current goal does not reflect that, change goal
	if a.Action == GoToRandomSpot && a.Goal == nil {
		goalX, goalY := a.Behavior.CoordinatesGenerator(*a.picMapSparse, false)
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}

	// if agent has nothing to do, try to stay still
	if a.Action == None && a.Goal == nil {
		goalX, goalY := a.X, a.Y
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}

	if a.Action == None && a.Goal != nil {
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
	if a.Action == WaitForBeer && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		return
	}

	// if goal is reached
	if a.Action != None && a.Goal != nil && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		a.Path = nil
		a.CurrentWayPoint = 0
		a.Goal = nil
		if a.Action == GoToExit {
			a.endOfLife = true
		} else if a.Action == GoToToilet {
			a.BladderContents = 0
			a.Action = GoToRandomSpot
			a.PerceptPeeChannel <- true
		} else if a.Action == GoToBar {
			a.Action = WaitForBeer
			go a.WaitForBeer()
			// try to stay still
			goalX, goalY := a.X, a.Y
			g := jps.GetNode(int(goalY), int(goalX))
			a.Goal = &g
		} else if a.Action == GoToRandomSpot {
			a.Action = None
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
	// check if the client didn't decide to go home in the meantime
	if a.Action != WaitForBeer {
		return
	}
	if !response {
		a.hasABarman = true
		a.WaitForBeer()
	} else {
		a.DrinkContents = 300
		a.hasABarman = false
		a.Action = GoToRandomSpot
		goalX, goalY := a.Behavior.CoordinatesGenerator(*a.picMapSparse, false)
		g := jps.GetNode(int(goalY), int(goalX))
		a.Goal = &g
	}
}
