package game

import (
	"bufio"
	"math"
	"os"
)

// Tile represents the representation of an element in a map
type Tile rune

// Enum of differetn space types
const (
	StoneTile      Tile = '#'
	DirtTile            = '.'
	ClosedDoorTile      = '|'
	OpenedDoorTile      = '/'
	EmptyTile           = 0
	PlayerTile          = '@'
	PendingTile         = -1
)

// Level represents the mapping of a level
type Level struct {
	Tiles    [][]Tile
	Player   *Player
	Monsters map[Pos]*Monster
	Debug    map[Pos]bool
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
	level.Monsters = make(map[Pos]*Monster)
	for i := range level.Tiles {
		level.Tiles[i] = make([]Tile, longestRow)
	}

	for y := range level.Tiles {
		line := lines[y]
		for x, c := range line {
			t := level.Tiles[y][x]
			pos := Pos{x, y}

			switch c {
			case ' ', '\t', '\n', '\r':
				t = EmptyTile
			case '#':
				t = StoneTile
			case '|':
				t = ClosedDoorTile
			case '/':
				t = OpenedDoorTile
			case '.':
				t = DirtTile
			case '@':
				level.Player = NewPlayer(pos)
				t = PendingTile
			case 'R':
				level.Monsters[pos] = NewRat(pos)
				t = PendingTile
			case 'S':
				level.Monsters[pos] = NewSpider(pos)
				t = PendingTile
			default:
				panic("Invalid Character: " + string(c))
			}

			level.Tiles[y][x] = t
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

func (level *Level) canWalk(pos Pos) bool {
	tile := level.Tiles[pos.Y][pos.X]
	switch tile {
	case EmptyTile, ClosedDoorTile, StoneTile:
		return false
	default:
		return true
	}
}

func (level *Level) checkDoor(pos Pos) {
	tile := level.Tiles[pos.Y][pos.X]
	if tile == ClosedDoorTile {
		level.Tiles[pos.Y][pos.X] = OpenedDoorTile
	}
}

func (level *Level) bfs(start Pos) {
	queue := make([]Pos, 0, 8)
	queue = append(queue, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	level.Debug = visited

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range level.getNeighbors(current) {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}

}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	queue := make(posPriorityQueue, 0, 8)
	queue = queue.push(start, 1)

	from := make(map[Pos]Pos)
	from[start] = start

	cost := make(map[Pos]int)
	cost[start] = 0

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

			return path
		}

		for _, neighbor := range level.getNeighbors(current) {
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

func (level *Level) getNeighbors(pos Pos) []Pos {
	neighbors := make([]Pos, 0, 8)

	right := Pos{X: pos.X + 1, Y: pos.Y}
	if level.canWalk(right) {
		neighbors = append(neighbors, right)
	}

	left := Pos{X: pos.X - 1, Y: pos.Y}
	if level.canWalk(left) {
		neighbors = append(neighbors, left)
	}

	up := Pos{X: pos.X, Y: pos.Y - 1}
	if level.canWalk(up) {
		neighbors = append(neighbors, up)
	}

	down := Pos{X: pos.X, Y: pos.Y + 1}
	if level.canWalk(down) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}
