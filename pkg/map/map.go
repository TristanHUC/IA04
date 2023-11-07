package _map

import (
	"math"
	"os"
)

// LoadFromFile loads a map from a file.
func (m *Map) LoadFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	m.Width = 0
	m.Height = 0
	for _, c := range content {
		if c == '\n' {
			m.Height++
			m.Width = 0
		} else {
			m.Width++
		}
		if c == 'w' {
			m.Walls = append(m.Walls, [2]int{m.Width, m.Height})
		}
	}
	return nil
}

// SaveToFile saves a map to a file.
func (m *Map) SaveToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	minXWall := math.Inf(1)
	minYWall := math.Inf(1)
	maxXWall := math.Inf(-1)
	maxYWall := math.Inf(-1)
	for _, wall := range m.Walls {
		if float64(wall[0]) < minXWall {
			minXWall = float64(wall[0])
		}
		if float64(wall[1]) < minYWall {
			minYWall = float64(wall[1])
		}
		if float64(wall[0]) > maxXWall {
			maxXWall = float64(wall[0])
		}
		if float64(wall[1]) > maxYWall {
			maxYWall = float64(wall[1])
		}
	}

	for y := int(minYWall); y <= int(maxYWall); y++ {
		for x := int(minXWall); x <= int(maxXWall); x++ {
			isWall := false
			for _, wall := range m.Walls {
				if wall[0] == x && wall[1] == y {
					isWall = true
					break
				}
			}
			if isWall {
				f.Write([]byte{'w'})
			} else {
				f.Write([]byte{' '})
			}
		}
		f.Write([]byte{'\n'})
	}

	return nil
}
