package simulation

import (
	"gitlab.utc.fr/royhucheradorni/ia04.git/pkg/astar"
)

type Environment struct {
	Walls          [][2]int
	m              *astar.Map
	PerceptChannel chan PerceptRequest
	Agents         []*Agent
}

func NewEnvironment(walls [][2]int, m *astar.Map, nAgents int) *Environment {
	env := Environment{
		Walls:          walls,
		m:              m,
		PerceptChannel: make(chan PerceptRequest, 1),
		Agents:         make([]*Agent, nAgents),
	}
	for i := 0; i < nAgents; i++ {
		env.Agents[i] = NewAgent(float64(70+35*i), float64(70), 99, 99, m, env.PerceptChannel)
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
