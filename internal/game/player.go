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
			Damage:       5,
			Speed:        1.0,
			ActionPoints: 0,
			SightRange:   10,
		},
	}
}

// Move moves the player to a new position
func (p *Player) Move(level *Level, to Pos) {
	p.Pos = to
}
