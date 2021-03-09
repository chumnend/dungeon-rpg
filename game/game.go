package game

import "math"

// InputType used for input enumeration
type InputType int

// Input Enum for getting input for the game
const (
	QuitGame InputType = iota
	Up
	Down
	Left
	Right
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
}

// Game represents the RPG game state
type Game struct {
	LevelCh chan *Level
	InputCh chan *Input
	Level   *Level
}

// NewGame creates a new Game struct
func NewGame(path string) *Game {
	game := new(Game)

	game.LevelCh = make(chan *Level)
	game.InputCh = make(chan *Input)
	game.Level = loadLevelFromFile(path)

	return game
}

// Run runs the game user interface
func (game *Game) Run() {
	game.LevelCh <- game.Level

	for {
		input := <-game.InputCh
		if input.Type == QuitGame {
			return
		}

		game.handleInput(input)

		game.LevelCh <- game.Level
	}
}

func (game *Game) handleInput(input *Input) {
	level := game.Level
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
	}

	if newPos {
		level.resolveMove(pos)
	}
}

func bresenham(start Pos, end Pos) []Pos {
	result := make([]Pos, 0)
	isSteep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))
	if isSteep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}

	if start.X > end.X {
		start.X, end.X = end.X, start.X
		start.Y, end.Y = end.Y, start.Y
	}

	deltaX := end.X - start.X
	deltaY := int(math.Abs(float64(end.Y - start.Y)))

	err := 0
	y := start.Y
	yStep := 1
	if start.Y >= end.Y {
		yStep = -1
	}

	for x := start.X; x < end.X; x++ {
		if isSteep {
			result = append(result, Pos{y, x})
		} else {
			result = append(result, Pos{x, y})
		}

		err += deltaY
		if 2*err >= deltaX {
			y += yStep
			err -= deltaX
		}
	}

	return result
}
