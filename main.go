package main

import (
	"github.com/chumnend/simple-rpg/game"
	"github.com/chumnend/simple-rpg/ui"
)

func main() {
	numWindows := 1

	game := game.NewGame(numWindows, "game/maps/level1.map")
	for i := 0; i < numWindows; i++ {
		go func(i int) {
			app := ui.NewApp(game.LevelChs[i], game.InputCh)
			app.Run()
		}(i)
	}

	game.Run()
}
