package game

// Character represents a type of character in the game (ie. player, monster)
type Character struct {
	Entity
	Hitpoints    int
	Damage       int
	Speed        float64
	ActionPoints float64
}

// Attack ...
func Attack(c1, c2 *Character) {
	c2.Hitpoints -= c1.Damage

	if c2.Hitpoints > 0 {
		c2.ActionPoints--
		c1.Hitpoints -= c2.Damage
	}
}
