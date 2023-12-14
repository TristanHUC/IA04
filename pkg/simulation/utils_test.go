package simulation

import (
	"math"
	"testing"
)

func TestAngleTo8DirectionsSector(t *testing.T) {
	// Test the 8 directions
	for i := 0; i < 8; i++ {
		sector := AngleTo8DirectionsSector(float64(i) * math.Pi / 4)
		if sector != i {
			t.Errorf("AngleTo8DirectionsSector(%f) = %d, want %d", float64(i)*math.Pi/4, sector, i)
		}
	}
	// Test the 8 directions + 1/16
	for i := 0; i < 8; i++ {
		sector := AngleTo8DirectionsSector(float64(i)*math.Pi/4 + math.Pi/16)
		if sector != i {
			t.Errorf("AngleTo8DirectionsSector(%f) = %d, want %d", float64(i)*math.Pi/4+math.Pi/16, sector, i)
		}
	}
	// Test the 8 directions - 1/16
	for i := 0; i < 8; i++ {
		sector := AngleTo8DirectionsSector(float64(i)*math.Pi/4 - math.Pi/16)
		if sector != i {
			t.Errorf("AngleTo8DirectionsSector(%f) = %d, want %d", float64(i)*math.Pi/4-math.Pi/16, sector, i)
		}
	}
}

func TestVectToAngle(t *testing.T) {
	// Test the 8 directions
	rightX, rightY := 1.0, 0.0
	upRightX, upRightY := 1.0, 1.0
	upX, upY := 0.0, 1.0
	upLeftX, upLeftY := -1.0, 1.0
	leftX, leftY := -1.0, 0.0
	downLeftX, downLeftY := -1.0, -1.0
	downX, downY := 0.0, -1.0
	downRightX, downRightY := 1.0, -1.0

	// Test the 8 directions
	angleRight := VectToAngle(rightX, rightY)
	if angleRight != 0 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", rightX, rightY, angleRight, 0.0)
	}
	angleUpRight := VectToAngle(upRightX, upRightY)
	if angleUpRight != math.Pi/4 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", upRightX, upRightY, angleUpRight, math.Pi/4)
	}
	angleUp := VectToAngle(upX, upY)
	if angleUp != math.Pi/2 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", upX, upY, angleUp, math.Pi/2)
	}
	angleUpLeft := VectToAngle(upLeftX, upLeftY)
	if angleUpLeft != 3*math.Pi/4 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", upLeftX, upLeftY, angleUpLeft, 3*math.Pi/4)
	}
	angleLeft := VectToAngle(leftX, leftY)
	if angleLeft != math.Pi {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", leftX, leftY, angleLeft, math.Pi)
	}
	angleDownLeft := VectToAngle(downLeftX, downLeftY)
	if angleDownLeft != 5*math.Pi/4 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", downLeftX, downLeftY, angleDownLeft, 5*math.Pi/4)
	}
	angleDown := VectToAngle(downX, downY)
	if angleDown != 3*math.Pi/2 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", downX, downY, angleDown, 3*math.Pi/2)
	}
	angleDownRight := VectToAngle(downRightX, downRightY)
	if angleDownRight != 7*math.Pi/4 {
		t.Errorf("VectToAngle(%f, %f) = %f, want %f", downRightX, downRightY, angleDownRight, 7*math.Pi/4)
	}
}
