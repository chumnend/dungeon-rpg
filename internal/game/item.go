package game

// Item struct declaration
type Item struct {
	Entity
	Type ItemType
}

// ItemType declaration
type ItemType int

// ItemType enum declaration
const (
	Weapon ItemType = iota
	Armor
	Other
)

// NewSword creates a sword entity
func NewSword(p Pos) *Item {
	return &Item{
		Entity: Entity{
			Pos:    p,
			Name:   "Sword",
			Symbol: 's',
		},
		Type: Weapon,
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
		Type: Armor,
	}
}
