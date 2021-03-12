package game

import (
	"bufio"
	"math"
	"os"
	"strconv"
)

// Tile represents the representation of an element in a map
type Tile struct {
	Symbol  rune
	Visible bool
}

// Enum of differetn space types
const (
	StoneTile      rune = '#'
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
	Events   []string
	EventPos int
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

	level := &Level{
		Tiles:    make([][]Tile, len(lines)),
		Player:   nil,
		Monsters: make(map[Pos]*Monster),
		Events:   make([]string, 8),
		EventPos: 0,
		Debug:    make(map[Pos]bool),
	}

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
				t.Symbol = EmptyTile
			case '#':
				t.Symbol = StoneTile
			case '|':
				t.Symbol = ClosedDoorTile
			case '/':
				t.Symbol = OpenedDoorTile
			case '.':
				t.Symbol = DirtTile
			case '@':
				level.Player = NewPlayer(pos)
				t.Symbol = PendingTile
			case 'R':
				level.Monsters[pos] = NewRat(pos)
				t.Symbol = PendingTile
			case 'S':
				level.Monsters[pos] = NewSpider(pos)
				t.Symbol = PendingTile
			default:
				panic("Invalid Character: " + string(c))
			}

			level.Tiles[y][x] = t
		}
	}

	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile.Symbol == PendingTile {
				searchPos := Pos{x, y}
				level.Tiles[y][x] = level.searchTile(searchPos)
			}
		}
	}

	level.lineOfSight()

	return level
}

// AddEvent adds a string ti the event slice
func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++
	if level.EventPos >= len(level.Events) {
		level.EventPos = 0
	}
}

func (level *Level) inRange(pos Pos) bool {
	return pos.X < len(level.Tiles[0]) && pos.Y < len(level.Tiles) && pos.X >= 0 && pos.Y >= 0
}

func (level *Level) canWalk(pos Pos) bool {
	if !level.inRange(pos) {
		return false
	}

	tile := level.Tiles[pos.Y][pos.X]
	switch tile.Symbol {
	case EmptyTile, ClosedDoorTile, StoneTile:
		return false
	}

	if _, exists := level.Monsters[pos]; exists {
		return false
	}

	return true
}

func (level *Level) canSee(pos Pos) bool {
	if !level.inRange(pos) {
		return false
	}

	tile := level.Tiles[pos.Y][pos.X]
	switch tile.Symbol {
	case EmptyTile, ClosedDoorTile, StoneTile:
		return false
	default:
		return true
	}
}

func (level *Level) lineOfSight() {
	pos := level.Player.Pos
	dist := level.Player.SightRange

	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y

			d := math.Sqrt(float64(xDelta*xDelta + yDelta*yDelta))
			if d <= float64(dist) {
				level.bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) bresenham(start Pos, end Pos) {
	isSteep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))
	if isSteep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}

	deltaY := int(math.Abs(float64(end.Y - start.Y)))

	err := 0
	y := start.Y
	yStep := 1
	if start.Y >= end.Y {
		yStep = -1
	}

	if start.X > end.X {
		deltaX := start.X - end.X
		for x := start.X; x >= end.X; x-- {
			var pos Pos

			if isSteep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}

			level.Tiles[pos.Y][pos.X].Visible = true
			if !level.canSee(pos) {
				return
			}

			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	} else {
		deltaX := end.X - start.X
		for x := start.X; x < end.X; x++ {
			var pos Pos

			if isSteep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}

			level.Tiles[pos.Y][pos.X].Visible = true
			if !level.canSee(pos) {
				return
			}

			err += deltaY
			if 2*err >= deltaX {
				y += yStep
				err -= deltaX
			}
		}
	}
}

func (level *Level) checkDoor(pos Pos) {
	tile := level.Tiles[pos.Y][pos.X]
	if tile.Symbol == ClosedDoorTile {
		level.Tiles[pos.Y][pos.X].Symbol = OpenedDoorTile
	}
}

func (level *Level) attack(c1 *Character, c2 *Character) {
	c2.Hitpoints -= c1.Damage
	c1.ActionPoints--

	if c2.Hitpoints > 0 {
		level.AddEvent(c1.Name + " attacked " + c2.Name + " for " + strconv.Itoa(c1.Damage))
	} else {
		level.AddEvent(c1.Name + " killed " + c2.Name)
	}
}

func (level *Level) resolveMove(pos Pos) {
	monster, exists := level.Monsters[pos]
	if exists {
		level.attack(&level.Player.Character, &monster.Character)
		if monster.Hitpoints <= 0 {
			delete(level.Monsters, monster.Pos)
		}
		if level.Player.Hitpoints <= 0 {
			panic("You Died!")
		}
	} else if level.canWalk(pos) {
		level.Player.Move(level, pos)
	} else {
		level.checkDoor(pos)
	}

	for _, monster := range level.Monsters {
		monster.Update(level)
	}

	if len(level.Monsters) == 0 {
		panic("Level Complete")
	}
}

func (level *Level) searchTile(start Pos) Tile {
	// utilizes BFS
	queue := make([]Pos, 0, 8)
	queue = append(queue, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		currentTile := level.Tiles[current.Y][current.X]
		switch currentTile.Symbol {
		case DirtTile:
			return Tile{DirtTile, false}
		default:
			// do nothing
		}

		queue = queue[1:]

		for _, neighbor := range level.getNeighbors(current) {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}

	return Tile{DirtTile, false}
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
