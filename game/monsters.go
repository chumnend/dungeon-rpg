package game

// Monster represents a monster in the game
type Monster struct {
	Character
}

// NewRat creates a Rat monster
func NewRat(p Pos) *Monster {
	return &Monster{
		Character: Character{
			Entity: Entity{
				Pos:    p,
				Name:   "Rat",
				Symbol: 'R',
			},
			Hitpoints:    5,
			Damage:       1,
			Speed:        2.0,
			ActionPoints: 0,
		},
	}
}

// NewSpider creates a Spider monster
func NewSpider(p Pos) *Monster {
	return &Monster{
		Character: Character{
			Entity: Entity{
				Pos:    p,
				Name:   "Spider",
				Symbol: 'S',
			},
			Hitpoints:    10,
			Damage:       2,
			Speed:        1.0,
			ActionPoints: 0,
		},
	}
}

// Update updates the monsters position relative to the player
func (m *Monster) Update(level *Level) {
	playerPos := level.Player.Pos
	positions := level.astar(m.Pos, playerPos)
	moveIndex := 1

	if len(positions) > 0 {
		m.ActionPoints += m.Speed
	}
	ap := int(m.ActionPoints)

	for i := 0; i < ap; i++ {
		if moveIndex < len(positions) {
			m.Move(level, positions[moveIndex])
			moveIndex++
			m.ActionPoints--
		}
	}
}

// Move moves the monster to a given position
func (m *Monster) Move(level *Level, to Pos) {
	// check if valid tile
	if _, exists := level.Monsters[to]; !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	} else {
		level.AddEvent(m.Name + " attacked Player")
		m.Attack(level.Player)
	}
}

// Attack ...
func (m *Monster) Attack(p *Player) {
	p.Hitpoints -= m.Damage
	m.ActionPoints--
}
