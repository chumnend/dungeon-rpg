package game

// InputType used for input enumeration
type InputType int

// Input Enum for getting input for the game
const (
	QuitGame InputType = iota
	CloseWindow
	Up
	Down
	Left
	Right
	Search // temp
)

// Input represents the key board input
type Input struct {
	Type    InputType
	LevelCh chan *Level
}

// Pos reprsents the x an y coordinate
type Pos struct {
	X, Y int
}

// Game represents the RPG game state
type Game struct {
	LevelChs []chan *Level
	InputCh  chan *Input
	Level    *Level
}

// NewGame creates a new Game struct
func NewGame(numWindows int, path string) *Game {
	game := new(Game)

	game.LevelChs = make([]chan *Level, numWindows)
	for i := range game.LevelChs {
		game.LevelChs[i] = make(chan *Level)
	}
	game.InputCh = make(chan *Input)
	game.Level = loadLevelFromFile(path)

	return game
}

// Run runs the game user interface
func (game *Game) Run() {
	for _, lch := range game.LevelChs {
		lch <- game.Level
	}

	for {
		input := <-game.InputCh
		if input.Type == QuitGame {
			return
		}

		game.handleInput(input)

		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
		}

		if len(game.LevelChs) == 0 {
			return
		}

		for _, lch := range game.LevelChs {
			lch <- game.Level
		}
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
	case CloseWindow:
		close(input.LevelCh)
		chanIdx := 0
		for i, c := range game.LevelChs {
			if c == input.LevelCh {
				chanIdx = i
				break
			}
		}

		game.LevelChs = append(game.LevelChs[:chanIdx], game.LevelChs[chanIdx+1:]...)
	}

	if newPos {
		if level.canWalk(pos) {
			level.Player.Move(level, pos)
		} else {
			level.checkDoor(pos)
		}
	}
}
