package ui

import (
	"github.com/chumnend/dungeon-rpg/internal/game"
	"github.com/veandco/go-sdl2/sdl"
)

func (a *App) toggleInventory() {
	if a.state == mainState {
		a.state = inventoryState
	} else if a.state == inventoryState {
		a.dragged = nil
		a.state = mainState
	}
}

func (a *App) checkForFloorItem(level *game.Level, mx int32, my int32) *game.Item {
	mouseRect := &sdl.Rect{
		X: mx,
		Y: my,
		W: 1,
		H: 1,
	}

	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := a.getPickupItemRect(i)
		if itemRect.HasIntersection(mouseRect) {
			return item
		}
	}

	return nil
}

func (a *App) checkForInventoryItem(level *game.Level, mx int32, my int32) *game.Item {
	mouseRect := &sdl.Rect{
		X: mx,
		Y: my,
		W: 1,
		H: 1,
	}

	items := level.Player.Items
	for i, item := range items {
		itemRect := a.getInventoryItemRect(i)
		if itemRect.HasIntersection(mouseRect) {
			return item
		}
	}

	return nil
}

func (a *App) checkForDropItem(level *game.Level, mx int32, my int32) bool {
	mouseRect := &sdl.Rect{
		X: mx,
		Y: my,
		W: 1,
		H: 1,
	}

	inventoryRect := a.getInventoryBackdropRect()
	if !inventoryRect.HasIntersection(mouseRect) {
		return true
	}

	return false
}
