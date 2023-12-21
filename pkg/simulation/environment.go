package simulation

import (
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
)

type Environment struct {
	MapSparse       _map.Map
	MapDense        [][]uint8
	PerceptChannel  chan PerceptRequest
	Agents          []*Agent
	Counter         *Counter
	SimulationSpeed *float32
}

func NewEnvironment(sparseMap _map.Map, denseMap [][]uint8, nAgents int, nBarmans int, SimulationSpeed *float32) *Environment {
	env := Environment{
		MapSparse:       sparseMap,
		MapDense:        denseMap,
		PerceptChannel:  make(chan PerceptRequest, 1),
		Agents:          make([]*Agent, nAgents),
		Counter:         NewCounter(),
		SimulationSpeed: SimulationSpeed,
	}
	//for i := 0; i < nAgents; i++ {
	//	x, y := GenerateValidCoordinates(sparseMap.Walls, sparseMap.Width, sparseMap.Height)
	//	env.Agents[i] = NewAgent(float64(x), float64(y), denseMap, &sparseMap, env.PerceptChannel)
	//}
	for iClient := 0; iClient < nAgents-nBarmans; iClient++ {
		env.Agents[iClient] = NewAgent(iClient, ClientBehavior{}, denseMap, &sparseMap, env.PerceptChannel, false, env.Counter.GetChannelCounter(), SimulationSpeed, GoToBar)
	}
	for iBarman := nAgents - nBarmans; iBarman < nAgents; iBarman++ {
		env.Agents[iBarman] = NewAgent(iBarman, BarmanBehavior{}, denseMap, &sparseMap, env.PerceptChannel, false, env.Counter.GetChannelCounter(), SimulationSpeed, None)
	}

	//
	return &env
}

func (e *Environment) GetNearbyAgents(agt *Agent) []*Agent {
	nearbyAgents := make([]*Agent, 0)
	for _, agent := range e.Agents {
		if agent != agt {
			nearbyAgents = append(nearbyAgents, agent)
		}
	}
	return nearbyAgents
}

func (e *Environment) PerceptRequestsHandler() {
	for {
		select {
		case perceptRequest := <-e.PerceptChannel:
			perceptRequest.ResponseChannel <- e.GetNearbyAgents(perceptRequest.Agt)

		}
	}
}

func (e *Environment) Update() {
	for i, agent := range e.Agents {
		if agent.endOfLife == true {
			e.Agents = append(e.Agents[:i], e.Agents[i+1:]...)
		}
	}
}
