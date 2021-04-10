package ui

import (
	"github.com/chumnend/dungeon-rpg/internal/game"
)

func (a *App) toggleInventory() {
	if a.state == mainState {
		a.state = inventoryState
	} else if a.state == inventoryState {
		a.dragged = nil
		a.state = mainState
	}
}

func (a *App) checkForFloorItem(mx int32, my int32) *game.Item {
	mouseRect := a.getMouseRect(mx, my)

	level := a.loadedLevel
	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := a.getPickupItemRect(i)
		if itemRect.HasIntersection(mouseRect) {
			return item
		}
	}

	return nil
}

func (a *App) checkForInventoryItem(mx int32, my int32) *game.Item {
	mouseRect := a.getMouseRect(mx, my)

	level := a.loadedLevel
	items := level.Player.Items
	for i, item := range items {
		itemRect := a.getInventoryItemRect(i)
		if itemRect.HasIntersection(mouseRect) {
			return item
		}
	}

	return nil
}

func (a *App) checkForDropItem(mx int32, my int32) bool {
	mouseRect := a.getMouseRect(mx, my)

	inventoryRect := a.getInventoryBackdropRect()
	if !inventoryRect.HasIntersection(mouseRect) {
		return true
	}

	return false
}

func (a *App) checkForEquipItem(mx int32, my int32) *game.Item {
	mouseRect := a.getMouseRect(mx, my)

	if a.dragged.Type == game.Weapon {
		weaponRect := a.getWeaponSlotRect()
		if weaponRect.HasIntersection(mouseRect) {
			return a.dragged
		}
	}

	if a.dragged.Type == game.Armor {
		armorRect := a.getArmorSlotRect()
		if armorRect.HasIntersection(mouseRect) {
			return a.dragged
		}
	}

	return nil
}
