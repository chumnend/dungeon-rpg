package game

// InputType used for input enumeration
type InputType int

// Input Enum for getting input for the game
const (
	QuitGame InputType = iota
	Up
	Down
	Left
	Right
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
}

// Game represents the RPG game state
type Game struct {
	LevelCh chan *Level
	InputCh chan *Input
	Level   *Level
}

// NewGame creates a new Game struct
func NewGame(path string) *Game {
	return &Game{
		LevelCh: make(chan *Level),
		InputCh: make(chan *Input),
		Level:   loadLevelFromFile(path),
	}
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
	default:
		// do nothing
	}

	if newPos {
		level.resolveMove(pos)
	}
}
