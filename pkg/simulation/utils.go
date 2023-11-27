package simulation

import (
	"fmt"
	"math"
	"math/rand"
)

func signedAcos(x float64) float64 {
	unsignedAcos := math.Acos(x)
	if x >= 0 {
		return unsignedAcos
	} else {
		return -unsignedAcos
	}
}

// generateValidCoordinates generates random map coordinates (int) that are not inside a wall
func GenerateValidCoordinates(walls [][2]int, width, height int) (int, int) {
	x := rand.Intn(height)
	y := rand.Intn(width)
	coordsOk := false
	// while agent is inside a wall, generate new coordinates
	for !coordsOk {
		coordsOk = true
		for _, wall := range walls {
			if wall[1] == x && wall[0] == y {
				x = rand.Intn(width)
				y = rand.Intn(height)
				coordsOk = false
			}
			if wall[1] == x && wall[0] == y {
				fmt.Println("Maybe that's the problem")
			}
		}
	}
	return x, y
}
