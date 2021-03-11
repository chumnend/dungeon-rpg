package main

import (
	"github.com/chumnend/simple-rpg/internal/game"
	"github.com/chumnend/simple-rpg/internal/ui"
)

func main() {
	game := game.NewGame("internal/game/maps/level1.map")
	app := ui.NewApp(game.LevelCh, game.InputCh)

	go app.Run()
	game.Run()
}
