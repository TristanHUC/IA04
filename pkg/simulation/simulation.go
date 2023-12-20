package simulation

import (
	"fmt"
	"math"
)

type Simulation struct {
	NAgents     int
	NBarmans    int
	Environment *Environment
	Paused      bool
}

func (s *Simulation) Start() {
	go s.Environment.Counter.Run()
	for i := 0; i < s.NAgents; i++ {
		go s.Environment.Agents[i].Run()
	}
}

func (s *Simulation) TogglePause() {
	s.Paused = !s.Paused
	for i := 0; i < s.NAgents; i++ {
		s.Environment.Agents[i].Paused = !s.Environment.Agents[i].Paused
	}
}

func (s *Simulation) SetSpeed(speed float32) {
	for i := 0; i < s.NAgents; i++ {
		s.Environment.Agents[i].SimulationSpeed = float32(math.Max(0, float64(speed)))
	}
	fmt.Println("Set speed to", s.Environment.Agents[0].SimulationSpeed)
}
