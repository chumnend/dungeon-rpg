package game

// Player represents a player object
type Player struct {
	Pos

	Symbol Tile
}

// NewPlayer creates player struct
func NewPlayer(p Pos) *Player {
	return &Player{
		Pos:    p,
		Symbol: '@',
	}
}

// Move moves the player to a new position
func (p *Player) Move(level *Level, to Pos) {
	// check if valid tile
	if _, exists := level.Monsters[to]; !exists {
		p.Pos = to
	}
}
