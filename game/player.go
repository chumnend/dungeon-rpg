package game

// Player represents a player object
type Player struct {
	Character
}

// NewPlayer creates player struct
func NewPlayer(p Pos) *Player {
	return &Player{
		Character: Character{
			Entity: Entity{
				Pos:    p,
				Name:   "Player",
				Symbol: '@',
			},
			Hitpoints:    20,
			Damage:       10,
			Speed:        1.0,
			ActionPoints: 0,
		},
	}
}

// Move moves the player to a new position
func (p *Player) Move(level *Level, to Pos) {
	// check if valid tile
	if monster, exists := level.Monsters[to]; !exists {
		p.Pos = to
	} else {
		level.AddEvent("Player attacked " + monster.Name)
		p.Attack(monster)
		if p.Hitpoints <= 0 {
			panic("YOU DIED!")
		}
	}
}

// Attack ...
func (p *Player) Attack(m *Monster) {
	m.Hitpoints -= p.Damage
	p.ActionPoints--
}
