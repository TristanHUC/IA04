package _map

import "os"

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
