package game

import (
	"bufio"
	"os"
)

// UI interface
type UI interface {
	Draw(*Level)
}

// Run runs the game user interface
func Run(ui UI) {
	level1 := loadLevelFromFile("game/maps/level1.map")
	ui.Draw(level1)
}

// Tile represents the representation of an element in a map
type Tile rune

const (
	stone Tile = '#'
	dirt       = '.'
	door       = '|'
)

// Level represents the mapping of a level
type Level struct {
	Tiles [][]Tile
}

func loadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	longestRow := 0
	index := 0
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines[index]) > longestRow {
			longestRow = len(lines[index])
		}
		index++
	}

	level := &Level{}
	level.Tiles = make([][]Tile, len(lines))
	for i := range level.Tiles {
		level.Tiles[i] = make([]Tile, longestRow)
	}

	return level
}
