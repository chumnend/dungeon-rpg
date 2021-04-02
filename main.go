package main

import (
	"github.com/chumnend/dungeon-rpg/internal/game"
	"github.com/chumnend/dungeon-rpg/internal/ui"
)

func main() {
	// setup app
	game := game.NewGame("internal/game/maps/level1.map")
	app := ui.NewApp(game, 1280, 730)

	// start the app
	app.Start()
}
