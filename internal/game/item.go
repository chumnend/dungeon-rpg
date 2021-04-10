package game

// Item struct declaration
type Item struct {
	Entity
	Type  ItemType
	Power float64
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
		Type:  Weapon,
		Power: 2.0,
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
		Type:  Armor,
		Power: 0.8,
	}
}
