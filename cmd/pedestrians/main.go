package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"log"
	"math"
	"math/rand"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

type Agent struct {
	x, y, vx, vy, gx, gy, speed, reactivity float64
}

type Simulation struct {
	agents  []*Agent
	nAgents int
}

func (s *Simulation) Update() error {
	// add agents until we have nAgents
	for len(s.agents) < s.nAgents {
		if rand.Int()%2 == 0 {
			// right to left
			s.agents = append(s.agents, &Agent{
				x:     SCREEN_WIDTH + 30,
				y:     float64(rand.Intn(SCREEN_HEIGHT)),
				vx:    -1,
				vy:    0,
				gx:    -30,
				gy:    float64(rand.Intn(SCREEN_HEIGHT)),
				speed: float64(rand.Intn(5)+5) / 10,
			})
		} else {
			// left to right
			s.agents = append(s.agents, &Agent{
				x:          -30,
				y:          float64(rand.Intn(SCREEN_HEIGHT)),
				vx:         1,
				vy:         0,
				gx:         SCREEN_WIDTH + 30,
				gy:         float64(rand.Intn(SCREEN_HEIGHT)),
				speed:      float64(rand.Intn(5)+5) / 10,
				reactivity: 0.5,
			})
		}
	}

	//fmt.Printf("%+v\n", s.agents[0])

	// update agents
	i := 0
	for _ = range s.agents {
		agent := s.agents[i]
		// find the goal velocity
		// direction
		goalVelX := agent.gx - agent.x
		goalVelY := agent.gy - agent.y
		// normalize
		total := math.Sqrt(goalVelX*goalVelX + goalVelY*goalVelY)
		goalVelX /= total
		goalVelY /= total
		goalVelX *= agent.speed // scale
		goalVelY *= agent.speed

		// compute the acceleration
		accelX := goalVelX - agent.vx
		accelY := goalVelY - agent.vy
		// normalize
		total = math.Sqrt(accelX*accelX + accelY*accelY)
		accelX /= total
		accelY /= total
		accelX *= agent.reactivity // scale
		accelY *= agent.reactivity

		// integrate
		agent.vx += accelX
		agent.vy += accelY

		// collision detection
		for j := 0; j < len(s.agents); j++ {
			if i == j {
				continue
			}
			other := s.agents[j]
			if agent.x < other.x+10 && agent.x > other.x-10 && agent.y < other.y+10 && agent.y > other.y-10 {
				// compute the transfer direction (axis between the two agents)
				axisX := other.x - agent.x
				axisY := other.y - agent.y
				// proportion of transmitted velocity is the dot product of the axis and the agent's velocity
				normalizedAxisX := axisX / math.Sqrt(axisX*axisX+axisY*axisY)
				normalizedAxisY := axisY / math.Sqrt(axisX*axisX+axisY*axisY)
				//normalizedAgentVX := agent.vx / math.Sqrt(agent.vx*agent.vx+agent.vy*agent.vy)
				//normalizedAgentVY := agent.vy / math.Sqrt(agent.vx*agent.vx+agent.vy*agent.vy)
				//prop := 1 - math.Abs(normalizedAxisX*normalizedAgentVX+normalizedAxisY*normalizedAgentVY)
				prop := float64(1)
				//fmt.Println(prop)
				// subtract the scaled repulsion direction from the agent's velocity
				agent.vx -= normalizedAxisX * prop
				agent.vy -= normalizedAxisY * prop
				// add to the other agent's velocity
				other.vx += normalizedAxisX * prop
				other.vy += normalizedAxisY * prop
			}
		}

		agent.x += agent.vx // integrate
		agent.y += agent.vy

		// if agent is close to goal, remove it
		if agent.x < agent.gx+10 && agent.x > agent.gx-10 && agent.y < agent.gy+10 && agent.y > agent.gy-10 {
			s.agents = append(s.agents[:i], s.agents[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)

	for _, agent := range s.agents {
		vector.DrawFilledCircle(screen, float32(agent.x), float32(agent.y), 5, colornames.Black, true)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (s *Simulation) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Pic")

	sim := Simulation{
		nAgents: 100,
	}

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
