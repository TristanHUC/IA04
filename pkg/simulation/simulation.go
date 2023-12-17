package simulation

type Simulation struct {
	NAgents     int
	NBarmans    int
	Environment *Environment
}

func (s *Simulation) Start() {
	for i := 0; i < s.NAgents; i++ {
		go s.Environment.Agents[i].Run()
	}
}
