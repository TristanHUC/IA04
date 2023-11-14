package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	ebitenvector "github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"log"
	"math"
	"math/rand"
)

const SCREEN_WIDTH = 700
const SCREEN_HEIGHT = 700

var wallsMargin = 100

type Agent struct {
	x, y, vx, vy, gx, gy, speed, reactivity float64
	tx, ty                                  float64
	controllable                            bool
}

type Simulation struct {
	agents  []*Agent
	nAgents int
}

func signedAcos(x float64) float64 {
	unsignedAcos := math.Acos(x)
	if x >= 0 {
		return unsignedAcos
	} else {
		return -unsignedAcos
	}
}

func (s *Simulation) Update() error {

	// if we press up arrow, increase margin
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		wallsMargin++
		if wallsMargin > SCREEN_WIDTH/2 {
			wallsMargin = SCREEN_WIDTH/2 - 100
		}
	}
	// if we press down arrow, reduce margin
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		wallsMargin--
		if wallsMargin < 0 {
			wallsMargin = 0
		}
	}

	// add agents until we have nAgents
	for len(s.agents) < s.nAgents {
		if rand.Int()%2 == 0 {
			// top to bottom
			s.agents = append(s.agents, &Agent{
				x:          float64(wallsMargin + rand.Intn(SCREEN_WIDTH-wallsMargin*2)),
				y:          0,
				vx:         0,
				vy:         0,
				gx:         float64(wallsMargin + rand.Intn(SCREEN_WIDTH-wallsMargin*2)),
				gy:         SCREEN_HEIGHT - 50,
				speed:      float64(rand.Intn(5)+5) / 10,
				reactivity: 0.2,
			})
		} else {
			// bottom to top
			s.agents = append(s.agents, &Agent{
				x:          float64(wallsMargin + rand.Intn(SCREEN_WIDTH-wallsMargin*2)),
				y:          SCREEN_HEIGHT,
				vx:         0,
				vy:         0,
				gx:         float64(wallsMargin + rand.Intn(SCREEN_WIDTH-wallsMargin*2)),
				gy:         50,
				speed:      float64(rand.Intn(5)+5) / 10,
				reactivity: 0.2,
			})
		}
	}

	// update agents
	for _, agent := range s.agents {
		// if left click, move goal of controllable agent
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && agent.controllable {
			x, y := ebiten.CursorPosition()
			agent.gx = float64(x)
			agent.gy = float64(y)
		}

		// compute goal velocity (norm = agent speed, direction= towards goal)
		gvx := agent.gx - agent.x
		gvy := agent.gy - agent.y
		gvNorm := math.Sqrt(gvx*gvx + gvy*gvy)
		gvx /= gvNorm
		gvy /= gvNorm
		gvx *= agent.speed
		gvy *= agent.speed

		// change velocity towards goal velocity at a rate of reactivity
		agent.vx += (gvx - agent.vx) * agent.reactivity
		agent.vy += (gvy - agent.vy) * agent.reactivity

		// change velocity to avoid walls following johannson 2007 and the values provided by moussaïd 2009
		if agent.y > 100 && agent.y < SCREEN_HEIGHT-100 {
			dLeft := agent.x - float64(wallsMargin)
			agent.vx += 3 * math.Exp(-dLeft/2)
			dRight := float64(SCREEN_WIDTH-wallsMargin) - agent.x
			agent.vx -= 3 * math.Exp(-dRight/2)
		}

		// change velocity to avoid other agents following moussaïd 2009
		for _, otherAgent := range s.agents {
			if otherAgent == agent {
				continue
			}
			lambda := 2.0
			A := 4.5
			gamma := 0.35
			n := 2.0
			np := 3.0
			factor := 0.15
			//if agent.controllable {
			//	factor = 10
			//}
			dist := math.Sqrt((agent.x*factor-otherAgent.x*factor)*(agent.x*factor-otherAgent.x*factor) + (agent.y*factor-otherAgent.y*factor)*(agent.y*factor-otherAgent.y*factor))
			if dist > 10 {
				continue
			}
			ex := (otherAgent.x*factor - agent.x*factor) / dist
			ey := (otherAgent.y*factor - agent.y*factor) / dist
			Dx := lambda*(agent.vx*factor-otherAgent.vx*factor) + ex
			Dy := lambda*(agent.vy*factor-otherAgent.vy*factor) + ey
			DNorm := math.Sqrt(Dx*Dx + Dy*Dy)
			tx := Dx / DNorm
			ty := Dy / DNorm
			// nx, ny is the normal vector to tx,ty pointing to the left
			nx := ty
			ny := -tx
			agent.tx = nx
			agent.ty = ny
			//theta := math.Acos(math.Min(math.Max(ex*tx+ey*ty, -1), 1)) / (2 * math.Pi) * 360
			theta := signedAcos(math.Min(math.Max(ex*tx+ey*ty, -1), 1))
			B := gamma * DNorm
			addedToVX := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*tx + math.Exp(-math.Pow(n*B*theta, 2))*nx) / factor
			addedToVY := -A * math.Exp(-dist/B) * (math.Exp(-math.Pow(np*B*theta, 2))*ty + math.Exp(-math.Pow(n*B*theta, 2))*ny) / factor
			// safeguard against too big values
			if addedToVX > 10 {
				addedToVX = 0
			}
			if addedToVY > 10 {
				addedToVY = 0
			}
			agent.vx += addedToVX
			agent.vy += addedToVY
		}
	}

	// update positions
	i := 0
	for _ = range s.agents {
		agent := s.agents[i]
		// integrate velocity
		agent.x += agent.vx
		agent.y += agent.vy

		// forbid agents to go outside of the screen
		if agent.x < 0 {
			agent.x = 5
		}
		if agent.x > SCREEN_WIDTH {
			agent.x = SCREEN_WIDTH - 5
		}
		if agent.y < 0 {
			agent.y = 5
		}
		if agent.y > SCREEN_HEIGHT {
			agent.y = SCREEN_HEIGHT - 5
		}

		// if agent is way outside the walls, put it back in the middle
		if agent.x < float64(wallsMargin-10) {
			agent.x = float64(wallsMargin + 30)
			agent.vx = 0
			agent.vy = 0
		}
		if agent.x > float64(SCREEN_WIDTH-wallsMargin+10) {
			agent.x = float64(SCREEN_WIDTH - wallsMargin - 30)
			agent.vx = 0
			agent.vy = 0
		}

		// if agent is close to goal, remove it
		if agent.x < agent.gx+10 && agent.x > agent.gx-10 && agent.y < agent.gy+10 && agent.y > agent.gy-10 && !agent.controllable {
			s.agents = append(s.agents[:i], s.agents[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)

	// draw walls
	ebitenvector.DrawFilledRect(screen, float32(wallsMargin-10), 100, 10, float32(SCREEN_HEIGHT-200), colornames.Black, false)
	ebitenvector.DrawFilledRect(screen, float32(SCREEN_WIDTH-wallsMargin), 100, 10, float32(SCREEN_HEIGHT-200), colornames.Black, false)

	for _, agent := range s.agents {
		color := colornames.Black
		if agent.controllable {
			color = colornames.Red
		}
		ebitenvector.DrawFilledCircle(screen, float32(agent.x), float32(agent.y), 5, color, true)
		// draw line to goal if controllable
		if agent.controllable {
			ebitenvector.StrokeLine(screen, float32(agent.x), float32(agent.y), float32(agent.gx), float32(agent.gy), 1, colornames.Blue, true)
		}
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
		nAgents: 500,
	}
	sim.agents = append(sim.agents, &Agent{
		x:            400,
		y:            400,
		speed:        float64(rand.Intn(5)+5) / 5,
		reactivity:   0.2,
		controllable: true,
	})

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
