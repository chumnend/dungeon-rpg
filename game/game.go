package game

import (
	"bufio"
	"math"
	"os"
)

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

// Tile represents the representation of an element in a map
type Tile rune

// Enum of differetn space types
const (
	StoneTile      Tile = '#'
	DirtTile            = '.'
	ClosedDoorTile      = '|'
	OpenedDoorTile      = '/'
	EmptyTile           = 0
	PlayerTile          = 'P'
	PendingTile         = -1
)

// Pos reprsents the x an y coordinate
type Pos struct {
	X, Y int
}

// Entity represnts an object
type Entity struct {
	Pos
}

// Player represents a player object
type Player struct {
	Entity
}

// Level represents the mapping of a level
type Level struct {
	Tiles  [][]Tile
	Player Player
	Debug  map[Pos]bool
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

	for input := range game.InputCh {
		if input.Type == QuitGame {
			return
		}

		game.handleInput(input)

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

	switch input.Type {
	case Up:
		if canWalk(level, Pos{level.Player.X, level.Player.Y - 1}) {
			level.Player.Y--
		} else {
			checkDoor(level, Pos{level.Player.X, level.Player.Y - 1})
		}
	case Down:
		if canWalk(level, Pos{level.Player.X, level.Player.Y + 1}) {
			level.Player.Y++
		} else {
			checkDoor(level, Pos{level.Player.X, level.Player.Y + 1})
		}
	case Left:
		if canWalk(level, Pos{level.Player.X - 1, level.Player.Y}) {
			level.Player.X--
		} else {
			checkDoor(level, Pos{level.Player.X - 1, level.Player.Y})
		}
	case Right:
		if canWalk(level, Pos{level.Player.X + 1, level.Player.Y}) {
			level.Player.X++
		} else {
			checkDoor(level, Pos{level.Player.X + 1, level.Player.Y})
		}
	case Search:
		//game.bfs(ui, level, level.Player.Pos)
		game.astar(level.Player.Pos, Pos{4, 2})
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

}

func canWalk(level *Level, pos Pos) bool {
	tile := level.Tiles[pos.Y][pos.X]
	switch tile {
	case EmptyTile, ClosedDoorTile, StoneTile:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Pos) {
	tile := level.Tiles[pos.Y][pos.X]
	if tile == ClosedDoorTile {
		level.Tiles[pos.Y][pos.X] = OpenedDoorTile
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
				level.Tiles[y][x] = EmptyTile
			case '#':
				level.Tiles[y][x] = StoneTile
			case '|':
				level.Tiles[y][x] = ClosedDoorTile
			case '/':
				level.Tiles[y][x] = OpenedDoorTile
			case '.':
				level.Tiles[y][x] = DirtTile
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				level.Tiles[y][x] = PendingTile
			default:
				panic("Invalid Character: " + string(c))
			}
		}
	}

	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile == PendingTile {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Tiles[searchY][searchX]
						switch searchTile {
						case DirtTile:
							level.Tiles[y][x] = DirtTile
							break SearchLoop
						}
					}
				}
			}
		}
	}

	return level
}

func getNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 8)

	right := Pos{X: pos.X + 1, Y: pos.Y}
	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}

	left := Pos{X: pos.X - 1, Y: pos.Y}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}

	up := Pos{X: pos.X, Y: pos.Y - 1}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}

	down := Pos{X: pos.X, Y: pos.Y + 1}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}

func (game *Game) bfs(start Pos) {
	level := game.Level

	queue := make([]Pos, 0, 8)
	queue = append(queue, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range getNeighbors(level, current) {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}

}

func (game *Game) astar(start Pos, goal Pos) []Pos {
	level := game.Level

	queue := make(posPriorityQueue, 0, 8)
	queue = queue.push(start, 1)

	from := make(map[Pos]Pos)
	from[start] = start

	cost := make(map[Pos]int)
	cost[start] = 0

	level.Debug = make(map[Pos]bool)

	var current Pos
	for len(queue) > 0 {
		queue, current = queue.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = from[p]
			}

			path = append(path, p)
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			for _, pos := range path {
				level.Debug[pos] = true
			}

			return path
		}

		for _, neighbor := range getNeighbors(level, current) {
			newCost := cost[current] + 1 // always 1 for now
			if _, exists := cost[neighbor]; !exists || newCost < cost[neighbor] {
				cost[neighbor] = newCost
				xDist := int(math.Abs(float64(goal.X - neighbor.X)))
				yDist := int(math.Abs(float64(goal.Y - neighbor.Y)))
				priority := newCost + xDist + yDist
				queue = queue.push(neighbor, priority)
				from[neighbor] = current
			}
		}
	}

	return nil
}
