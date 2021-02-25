package game

// Monster represents a monster in the game
type Monster struct {
	Pos

	Symbol    Tile
	Name      string
	Hitpoints int
	Damage    int
	Speed     float64
}

// NewRat creates a Rat monster
func NewRat(p Pos) *Monster {
	return &Monster{
		Pos:       p,
		Symbol:    'R',
		Name:      "Rat",
		Hitpoints: 5,
		Damage:    1,
		Speed:     1.0,
	}
}

// NewSpider creates a Spider monster
func NewSpider(p Pos) *Monster {
	return &Monster{
		Pos:       p,
		Symbol:    'S',
		Name:      "Spider",
		Hitpoints: 10,
		Damage:    2,
		Speed:     2.0,
	}
}

// Update updates the monsters position relative to the player
func (m *Monster) Update(level *Level) {
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)

	if len(positions) > 1 {
		m.Move(level, positions[1])
	}
}

// Move moves the monster to a given position
func (m *Monster) Move(level *Level, to Pos) {
	delete(level.Monsters, m.Pos)
	level.Monsters[to] = m
	m.Pos = to
}
