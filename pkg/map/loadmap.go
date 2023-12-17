//go:build !js
// +build !js

package _map

import "fmt"

func LoadMap(mapName string) (Map, error) {
	testmap := Map{}
	err := testmap.LoadFromFile(mapName)
	if err != nil {
		fmt.Println(err)
		return testmap, err
	}
	return testmap, nil
}
