package simulation

type Simulation struct {
	NAgents     int
	NBarmans    int
	Environment *Environment
	Paused      bool
}

func (s *Simulation) Start() {
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
