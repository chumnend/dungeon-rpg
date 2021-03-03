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
	Search // temp
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
	Symbol Tile
}

// Character represents a type of character in the game (ie. player, monster)
type Character struct {
	Entity
	Hitpoints    int
	Damage       int
	Speed        float64
	ActionPoints float64
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

		for pos, monster := range game.Level.Monsters {
			monster.Update(game.Level)
			if monster.Hitpoints <= 0 {
				delete(game.Level.Monsters, pos)
			}
		}

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
	case Search:
		level.astar(level.Player.Pos, Pos{4, 2})
	}

	if newPos {
		if level.canWalk(pos) {
			level.Player.Move(level, pos)
		} else {
			level.checkDoor(pos)
		}
	}
}
