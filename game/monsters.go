package game

// Monster represents a monster in the game
type Monster struct {
	Symbol    Tile
	Name      string
	Hitpoints int
	Damage    int
	Speed     float64
}

// NewRat creates a Rat monster
func NewRat() *Monster {
	return &Monster{
		Symbol:    'R',
		Name:      "Rat",
		Hitpoints: 5,
		Damage:    1,
		Speed:     2.0,
	}
}

// NewSpider creates a Spider monster
func NewSpider() *Monster {
	return &Monster{
		Symbol:    'S',
		Name:      "Spider",
		Hitpoints: 10,
		Damage:    2,
		Speed:     2.0,
	}
}
