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
