package _map

import (
	"log"
	"math"
	"os"
)

func (m *Map) IsWall(x, y int) bool {
	for _, wall := range m.Walls {
		if wall[0] == x && wall[1] == y {
			return true
		}
	}
	return false
}

// LoadFromString loads a map from a string.
func (m *Map) LoadFromString(content string) error {
	m.Height = 0
	// mesure width on first line
	m.Width = 0
	column := 0
	for _, c := range content {
		if c == 'w' {
			m.Walls = append(m.Walls, [2]int{column, m.Height})
		}
		if c == 'm' {
			m.ManToiletPoints = append(m.ManToiletPoints, [2]int{column, m.Height})
		}
		if c == 'f' {
			m.WomanToiletPoints = append(m.WomanToiletPoints, [2]int{column, m.Height})
		}
		if c == 'b' {
			m.BarPoints = append(m.BarPoints, [2]int{column, m.Height})
		}
		if c == 'c' {
			m.BarmenArea = append(m.BarmenArea, [2]int{column, m.Height})
		}
		if c == 't' {
			m.BeerTaps = append(m.BeerTaps, [2]int{column, m.Height})
		}
		if c == 'a' {
			m.CounterArea = append(m.CounterArea, [2]int{column, m.Height})
			m.Walls = append(m.Walls, [2]int{column, m.Height})
		}
		if c == 'e' {
			m.Exit = append(m.Exit, [2]int{column, m.Height})
		}
		if c == 'd' {
			m.Enter = append(m.Enter, [2]int{column, m.Height})
		}
		if c == '\n' {
			m.Height++
			if column > m.Width {
				m.Width = column
			}
			column = 0
		} else {
			column++
		}
	}
	if m.Width == 0 || m.Height == 0 {
		log.Fatalf("invalid map: %dx%d", m.Width, m.Height)
	}
	return nil
}

// LoadFromFile loads a map from a file.
func (m *Map) LoadFromFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return m.LoadFromString(string(content))
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
	for _, entry := range m.Enter {
		if float64(entry[0]) < minXWall {
			minXWall = float64(entry[0])
		}
		if float64(entry[1]) < minYWall {
			minYWall = float64(entry[1])
		}
		if float64(entry[0]) > maxXWall {
			maxXWall = float64(entry[0])
		}
		if float64(entry[1]) > maxYWall {
			maxYWall = float64(entry[1])
		}
	}
	for _, exit := range m.Exit {
		if float64(exit[0]) < minXWall {
			minXWall = float64(exit[0])
		}
		if float64(exit[1]) < minYWall {
			minYWall = float64(exit[1])
		}
		if float64(exit[0]) > maxXWall {
			maxXWall = float64(exit[0])
		}
		if float64(exit[1]) > maxYWall {
			maxYWall = float64(exit[1])
		}
	}

	for y := int(minYWall); y <= int(maxYWall); y++ {
		for x := int(minXWall); x <= int(maxXWall); x++ {
			isWall := false
			isBar := false
			isWomanToilet := false
			isManToilet := false
			isBarmenArea := false
			isBeerTap := false
			isCounterArea := false
			isExit := false
			isEnter := false
			isEmpty := true

			//adding wall to the map
			for _, wall := range m.Walls {
				if wall[0] == x && wall[1] == y {
					isWall = true
					break
				}
			}
			//adding CounterArea to the map
			for _, CounterArea := range m.CounterArea {
				if CounterArea[0] == x && CounterArea[1] == y {
					isCounterArea = true
					break
				}
			}
			if isCounterArea {
				f.Write([]byte{'a'})
				isEmpty = false
			} else if isWall {
				f.Write([]byte{'w'})
				isEmpty = false
			}

			//adding Bar to the map
			for _, Bar := range m.BarPoints {
				if Bar[0] == x && Bar[1] == y {
					isBar = true
					break
				}
			}
			if isBar {
				f.Write([]byte{'b'})
				isEmpty = false
			}

			//adding WomanToilet to the map
			for _, WomanToilet := range m.WomanToiletPoints {
				if WomanToilet[0] == x && WomanToilet[1] == y {
					isWomanToilet = true
					break
				}
			}
			if isWomanToilet {
				f.Write([]byte{'f'})
				isEmpty = false
			}

			//adding ManToilet to the map
			for _, ManToilet := range m.ManToiletPoints {
				if ManToilet[0] == x && ManToilet[1] == y {
					isManToilet = true
					break
				}
			}
			if isManToilet {
				f.Write([]byte{'m'})
				isEmpty = false
			}

			//adding BarmenArea to the map
			for _, BarmenArea := range m.BarmenArea {
				if BarmenArea[0] == x && BarmenArea[1] == y {
					isBarmenArea = true
					break
				}
			}
			if isBarmenArea {
				f.Write([]byte{'c'})
				isEmpty = false
			}

			//adding BeerTaps to the map
			for _, BeerTap := range m.BeerTaps {
				if BeerTap[0] == x && BeerTap[1] == y {
					isBeerTap = true
					break
				}
			}
			if isBeerTap {
				f.Write([]byte{'t'})
				isEmpty = false
			}

			//adding Exit to the map
			for _, Exit := range m.Exit {
				if Exit[0] == x && Exit[1] == y {
					isExit = true
					break
				}
			}
			if isExit {
				f.Write([]byte{'e'})
				isEmpty = false
			}

			//adding Enter to the map
			for _, Enter := range m.Enter {
				if Enter[0] == x && Enter[1] == y {
					isEnter = true
					break
				}
			}
			if isEnter {
				f.Write([]byte{'d'})
				isEmpty = false
			}

			if isEmpty == true {
				f.Write([]byte{' '})
			}

		}
		f.Write([]byte{'\n'})
	}

	return nil
}
