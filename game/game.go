package game

import (
	"bufio"
	"math"
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
	Search // temp
	None
)

// Input represents the key board input
type Input struct {
	Type InputType
}

// Tile represents the representation of an element in a map
type Tile rune

// Enum of differetn space types
const (
	SpaceStone      Tile = '#'
	SpaceDirt            = '.'
	SpaceClosedDoor      = '|'
	SpaceOpenedDoor      = '/'
	SpaceBlank           = 0
	SpacePlayer          = 'P'
	SpacePending         = -1
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

// Run runs the game user interface
func Run(ui UI) {
	level := loadLevelFromFile("game/maps/level1.map")

	for {
		ui.Draw(level)
		input := ui.GetInput()

		if input != nil && input.Type == Quit {
			return
		}

		handleInput(ui, level, input)
	}
}

func handleInput(ui UI, level *Level, input *Input) {
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
		// bfs(ui, level, level.Player.Pos)
		astar(ui, level, level.Player.Pos, Pos{4, 2})
	}
}

func canWalk(level *Level, pos Pos) bool {
	tile := level.Tiles[pos.Y][pos.X]
	switch tile {
	case SpaceBlank, SpaceClosedDoor, SpaceStone:
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Pos) {
	tile := level.Tiles[pos.Y][pos.X]
	if tile == SpaceClosedDoor {
		level.Tiles[pos.Y][pos.X] = SpaceOpenedDoor
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
				level.Tiles[y][x] = SpaceClosedDoor
			case '/':
				level.Tiles[y][x] = SpaceOpenedDoor
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

func bfs(ui UI, level *Level, start Pos) {
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
				ui.Draw(level)
			}
		}
	}

}

func astar(ui UI, level *Level, start Pos, goal Pos) []Pos {
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
				ui.Draw(level)
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
