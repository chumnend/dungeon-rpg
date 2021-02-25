package main

import (
	"github.com/chumnend/simple-rpg/game"
	"github.com/chumnend/simple-rpg/ui"
)

func main() {
	game := game.NewGame("game/maps/level1.map")
	go func() {
		app := ui.NewApp(game.LevelCh, game.InputCh)
		app.Run()
	}()

	game.Run()
}
