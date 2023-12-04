package simulation

import (
	_map "gitlab.utc.fr/royhucheradorni/ia04.git/pkg/map"
)

type Environment struct {
	MapSparse      _map.Map
	MapDense       [][]uint8
	PerceptChannel chan PerceptRequest
	Agents         []*Agent
}

func NewEnvironment(sparseMap _map.Map, denseMap [][]uint8, nAgents int) *Environment {
	env := Environment{
		MapSparse:      sparseMap,
		MapDense:       denseMap,
		PerceptChannel: make(chan PerceptRequest, 1),
		Agents:         make([]*Agent, nAgents),
	}
	for i := 0; i < nAgents; i++ {
		x, y := GenerateValidCoordinates(sparseMap.Walls, sparseMap.Width, sparseMap.Height)
		env.Agents[i] = NewAgent(float64(x), float64(y), denseMap, &sparseMap, env.PerceptChannel)
	}
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
