package simulation

import (
	"github.com/ankurjha7/jps"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/exp/slices"
	"math/rand"
)

type BarmanBehavior struct{}

func (BarmanBehavior) Reflect(a *Agent) {
	if a.action != None {
		return
	}
	//lauching the search for a client
	a.action = GoToBeerTap

}

func (BarmanBehavior) Act(a *Agent) {
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
	if a.action != None && a.Goal != nil && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
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

func (BarmanBehavior) CoordinatesGenerator(m _map.Map) (float64, float64) {
	// Take a random point in the bar area
	counterPoints := m.BarmenArea[rand.Intn(len(m.BarmenArea))]
	return float64(counterPoints[0]) + rand.Float64(), float64(counterPoints[1]) + rand.Float64()
}

// SearchForClient may not find a client if there is none
func (a *Agent) SearchForClient() {
	for _, agent := range a.closeAgents {
		if agent.action == WaitForBeer && !agent.hasABarman {
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
