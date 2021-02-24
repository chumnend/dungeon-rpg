package main

import (
	"runtime"

	"github.com/chumnend/simple-rpg/game"
	"github.com/chumnend/simple-rpg/ui"
)

func main() {
	numWindows := 1

	game := game.NewGame(numWindows, "game/maps/level1.map")
	for i := 0; i < numWindows; i++ {
		go func(i int) {
			runtime.LockOSThread()
			app := ui.NewApp(game.LevelChs[i], game.InputCh)
			app.Run()
		}(i)
	}

	game.Run()
}
