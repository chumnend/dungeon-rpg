package main

import (
	"github.com/chumnend/simple-rpg/game"
	"github.com/chumnend/simple-rpg/ui"
)

func main() {
	ui := &ui.UI{}
	game.Run(ui)
}
