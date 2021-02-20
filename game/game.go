package game

import (
	"bufio"
	"os"
)

// UI interface
type UI interface {
	Draw(*Level)
	GetInput() *Input
}

// InputType used for input enumeration
type InputType int

// Input Enum for getting input for the game
const (
	Up InputType = iota
	Down
	Left
	Right
	Quit
)

// Input represents the key board input
type Input struct {
	Type InputType
}

// Tile represents the representation of an element in a map
type Tile rune

// Enum of differetn space types
const (
	SpaceStone   Tile = '#'
	SpaceDirt         = '.'
	SpaceDoor         = '|'
	SpaceBlank        = 0
	SpacePlayer       = 'P'
	SpacePending      = -1
)

// Entity represnts an object
type Entity struct {
	X, Y int
}

// Player represents a player object
type Player struct {
	Entity
}

// Level represents the mapping of a level
type Level struct {
	Tiles  [][]Tile
	Player Player
}

// Run runs the game user interface
func Run(ui UI) {
	level := loadLevelFromFile("game/maps/level1.map")

	for {
		ui.Draw(level)
		input := ui.GetInput()
		if input != nil {
			switch input.Type {
			case Up:
				level.Player.Y--
			case Down:
				level.Player.Y++
			case Left:
				level.Player.X--
			case Right:
				level.Player.X++
			case Quit:
				return
			}
		}
	}
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

	for y := range level.Tiles {
		line := lines[y]
		for x, c := range line {
			switch c {
			case ' ', '\t', '\n', '\r':
				level.Tiles[y][x] = SpaceBlank
			case '#':
				level.Tiles[y][x] = SpaceStone
			case '|':
				level.Tiles[y][x] = SpaceDoor
			case '.':
				level.Tiles[y][x] = SpaceDirt
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				level.Tiles[y][x] = SpacePending
			default:
				panic("Invalid Character: " + string(c))
			}
		}
	}

	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile == SpacePending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Tiles[searchY][searchX]
						switch searchTile {
						case SpaceDirt:
							level.Tiles[y][x] = SpaceDirt
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}
