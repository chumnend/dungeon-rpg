package game

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// InputType used for input enumeration
type InputType int

// Input Enum for getting input for the game
const (
	QuitGame InputType = iota
	Up
	Down
	Left
	Right
	Take
	None
)

// Input represents the key board input
type Input struct {
	Type InputType
}

// Pos reprsents the x an y coordinate
type Pos struct {
	X, Y int
}

// Entity represents the identity of game entity (ie. player, monsters, items)
type Entity struct {
	Pos
	Name   string
	Symbol rune
}

// Character represents a type of character in the game (ie. player, monster)
type Character struct {
	Entity
	Hitpoints    int
	Damage       int
	Speed        float64
	ActionPoints float64
	SightRange   int
	Items        []*Item
}

// Game represents the RPG game state
type Game struct {
	LevelCh      chan *Level
	InputCh      chan *Input
	Levels       map[string]*Level
	CurrentLevel *Level
}

// NewGame creates a new Game struct
func NewGame(path string) *Game {

	// load levels from maps directory
	levels := loadLevels()

	// set initial level
	currentLevel := levels["level1"]
	currentLevel.lineOfSight()

	game := &Game{
		LevelCh:      make(chan *Level),
		InputCh:      make(chan *Input),
		Levels:       levels,
		CurrentLevel: currentLevel,
	}

	game.loadWorld()

	return game
}

// Run runs the game user interface
func (game *Game) Run() {
	game.LevelCh <- game.CurrentLevel

	for {
		input := <-game.InputCh
		if input.Type == QuitGame {
			return
		}

		game.handleInput(input)
		game.LevelCh <- game.CurrentLevel
	}
}

func (game *Game) handleInput(input *Input) {
	level := game.CurrentLevel
	var pos Pos
	newPos := false

	switch input.Type {
	case Up:
		pos = Pos{level.Player.X, level.Player.Y - 1}
		newPos = true
	case Down:
		pos = Pos{level.Player.X, level.Player.Y + 1}
		newPos = true
	case Left:
		pos = Pos{level.Player.X - 1, level.Player.Y}
		newPos = true
	case Right:
		pos = Pos{level.Player.X + 1, level.Player.Y}
		newPos = true
	case Take:
		items := level.Items[level.Player.Pos]
		if len(items) > 0 {
			for _, item := range items {
				level.moveItem(item, &level.Player.Character)
				level.AddEvent("Player picked up " + item.Name)
			}
		} else {
			level.AddEvent("Nothing to take!")
		}
	default:
		// do nothing
	}

	// if player did move, resolve that movement
	if newPos {
		// check if at portal
		nextLevel := level.Portals[pos]
		if nextLevel != nil {
			level.LastEvent = Portal
			level.Player.Pos = nextLevel.Pos
			game.CurrentLevel = nextLevel.Level
			game.CurrentLevel.lineOfSight()
		}

		level.resolveMove(pos)
	}
}

// LevelPos represents the starting location of a level
type LevelPos struct {
	Level *Level
	Pos   Pos
}

func (game *Game) loadWorld() {
	file, err := os.Open("internal/game/maps/world.txt")
	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	for rowIdx, row := range rows {
		if rowIdx == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			continue
		}

		levelWithPortal := game.Levels[row[0]]
		if levelWithPortal == nil {
			fmt.Println("Couldn't find level name in the world file")
			panic(nil)
		}

		x, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			panic(err)
		}
		pos := Pos{X: int(x), Y: int(y)}

		levelToGo := game.Levels[row[3]]
		if levelToGo == nil {
			fmt.Println("Couldn't find level name in the world file")
			panic(nil)
		}

		x, err = strconv.ParseInt(row[4], 10, 64)
		if err != nil {
			panic(err)
		}
		y, err = strconv.ParseInt(row[5], 10, 64)
		if err != nil {
			panic(err)
		}
		posToGo := Pos{X: int(x), Y: int(y)}

		levelWithPortal.Portals[pos] = &LevelPos{Level: levelToGo, Pos: posToGo}
	}
}
