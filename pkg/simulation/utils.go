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

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
}

// AngleTo8DirectionsSector takes in an angle in radians between [0, 2pi] and returns an int between [0, 7]
func AngleTo8DirectionsSector(angle float64) int {
	// Calculate the sector based on the angle
	sector := int(math.Floor((angle + math.Pi/8) / (math.Pi / 4)))
	// Map sector to a valid range [0, 7]
	return (sector) % 8
}

// VectToAngle takes in a vector and returns the angle in radians between [0, 2pi]
func VectToAngle(x, y float64) float64 {
	angle := math.Atan2(y, x)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Adjust the angle to be in the range [0, 2*pi)
	angle = math.Mod(angle, 2*math.Pi)
	return angle
}
