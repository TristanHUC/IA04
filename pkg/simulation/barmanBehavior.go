package simulation

import (
	"fmt"
	"github.com/ankurjha7/jps"
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
	"golang.org/x/exp/slices"
	"math/rand"
)

type BarmanBehavior struct{}

func (BarmanBehavior) Reflect(a *Agent) {
	if a.Action != None {
		return
	}
	//lauching the search for a client
	a.Action = GoToBeerTap

}

func (BarmanBehavior) Act(a *Agent) {
	// if agent want to go to beer tap, and current goal does not reflect that, change goal
	if a.Action == GoToBeerTap && (a.Goal == nil || !slices.Contains(a.picMapSparse.BeerTaps, [2]int{int(a.Goal.GetCol()), int(a.Goal.GetRow())})) {
		beerTap := a.picMapSparse.BeerTaps[rand.Intn(len(a.picMapSparse.BeerTaps))]
		g := jps.GetNode(beerTap[1], beerTap[0])
		a.Goal = &g
	}

	// if agent is waiting for a client, he should find one
	if a.Action == WaitForClient && a.client == nil {
		a.SearchForClient()
	}

	// if agent want to go to client, and the client is not waiting for a barman or the client doesn't exist anymore or he must go, change the Action to WaitForClient
	if a.Action == GoToClient && (a.client == nil || (a.client != nil && (a.client.Action != WaitForBeer || a.client.endOfLife))) {
		a.Action = WaitForClient
		fmt.Println("client is not waiting for a barman anymore")
		a.client = nil
	}

	// if agent want to go to client, and current goal is not the closest barmen area to the client, change goal
	if a.Action == GoToClient && a.client != nil && (a.Goal == nil || distanceNode(*a.Goal, a.GetClosestBarmenArea(*a.client)) > 1.5) {
		g := a.GetClosestBarmenArea(*a.client)
		a.Goal = &g
	}

	// if goal is reached
	if a.Action != None && a.Goal != nil && Distance(a.X, a.Y, float64(a.Goal.GetCol()), float64(a.Goal.GetRow())) < 1 {
		a.Path = nil
		a.CurrentWayPoint = 0
		a.Goal = nil
		if a.Action == WaitForClient {
			a.Action = GoToClient
		} else if a.Action == GoToBeerTap {
			a.DrinkContents = 300
			a.Action = WaitForClient
		} else if a.Action == GoToClient {
			a.GiveABeer()
			a.Action = GoToBeerTap
			a.client = nil
		}
	}
}

func (BarmanBehavior) CoordinatesGenerator(m _map.Map, isLaterGenerated bool) (float64, float64) {
	// Take a random point in the bar area
	counterPoints := m.BarmenArea[rand.Intn(len(m.BarmenArea))]
	return float64(counterPoints[0]) + rand.Float64(), float64(counterPoints[1]) + rand.Float64()
}

// SearchForClient may not find a client if there is none
func (a *Agent) SearchForClient() {
	for _, agent := range a.closeAgents {
		if agent.Action == WaitForBeer && !agent.hasABarman {
			a.client = agent
			g := a.GetClosestBarmenArea(*agent)
			a.Goal = &g
			// notify the client that he has a barman
			a.client.BeerChannel <- false
			a.Action = GoToClient
			break
		}
	}
}

func (a *Agent) GiveABeer() {
	a.client.BeerChannel <- true
	a.DrinkContents = 0
	a.BeerCounterChan <- a.Age
}
