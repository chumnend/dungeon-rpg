package game

// Item struct declaration
type Item struct {
	Entity
}

// NewSword creates a sword entity
func NewSword(p Pos) *Item {
	return &Item{
		Entity: Entity{
			Pos:    p,
			Name:   "Sword",
			Symbol: 's',
		},
	}
}

// NewHelmet creates a helmet entity
func NewHelmet(p Pos) *Item {
	return &Item{
		Entity: Entity{
			Pos:    p,
			Name:   "Helmet",
			Symbol: 'h',
		},
	}
}
