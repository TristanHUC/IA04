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

// GenerateValidCoordinates generates random map coordinates (int) that are not inside a wall
func GenerateValidCoordinates(walls [][2]int, width, height int) (float32, float32) {
	x := rand.Intn(width)
	y := rand.Intn(height)
	coordsOk := false
	// while agent is inside a wall, generate new coordinates
	for !coordsOk {
		coordsOk = true
		for _, wall := range walls {
			if wall[0] == x && wall[1] == y {
				x = rand.Intn(width)
				y = rand.Intn(height)
				coordsOk = false
			}
		}
	}
	xFloat := float32(x) + rand.Float32()
	yFloat := float32(y) + rand.Float32()
	return xFloat, yFloat
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}
