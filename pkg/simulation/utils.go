package simulation

import (
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
func generateValidCoordinates(walls [][2]int) (int, int) {
	x := rand.Intn(100)
	y := rand.Intn(100)
	coordsOk := false
	// while agent is inside a wall, generate new coordinates
	for !coordsOk {
		coordsOk = true
		for _, wall := range walls {
			if wall[0] == x && wall[1] == y {
				x = rand.Intn(100)
				y = rand.Intn(100)
				coordsOk = false
			}
		}
	}
	return x, y
}
