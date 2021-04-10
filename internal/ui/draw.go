package ui

import (
	"fmt"

	"github.com/chumnend/dungeon-rpg/internal/game"
	"github.com/veandco/go-sdl2/sdl"
)

func (a *App) draw() {
	a.renderer.Clear()
	a.r.Seed(1)

	// move the camera with the player
	a.setCamera()

	// draw floor tiles
	a.drawFloor()

	// draw player
	a.drawPlayer()

	// draw monsters
	a.drawMonsters()

	// draw items on ground
	a.drawFloorItems()

	// draw items on pickup bar
	a.drawPickupBarItems()

	// draw event log
	a.drawEventLog()

	// draw the inventory screen
	if a.state == inventoryState {
		a.drawInventory()
	}

	a.renderer.Present()
}

func (a *App) setCamera() {
	if a.centerX == -1 && a.centerY == -1 {
		a.centerX = a.loadedLevel.Player.X
		a.centerY = a.loadedLevel.Player.Y
	}

	limit := 5
	if a.loadedLevel.Player.X > a.centerX+limit {
		diff := a.loadedLevel.Player.X - (a.centerX + limit)
		a.centerX += diff
	} else if a.loadedLevel.Player.X < a.centerX-limit {
		diff := (a.centerX - limit) - a.loadedLevel.Player.X
		a.centerX -= diff
	}

	if a.loadedLevel.Player.Y > a.centerY+limit {
		diff := a.loadedLevel.Player.Y - (a.centerY + limit)
		a.centerY += diff
	} else if a.loadedLevel.Player.Y < a.centerY-limit {
		diff := (a.centerY - limit) - a.loadedLevel.Player.Y
		a.centerY -= diff
	}
}

func (a *App) drawFloor() {
	offsetX := (a.width / 2) - int32(a.centerX*spriteHeight)
	offsetY := (a.height / 2) - int32(a.centerY*spriteHeight)

	for y, row := range a.loadedLevel.Tiles {
		for x, tile := range row {
			if tile.Symbol == game.EmptyTile {
				continue
			}

			srcRects := a.textureIndex[tile.Symbol]
			srcRect := srcRects[a.r.Intn(len(srcRects))]

			if tile.Visible || tile.Seen {
				destRect := sdl.Rect{
					X: int32(x*spriteHeight) + offsetX,
					Y: int32(y*spriteHeight) + offsetY,
					W: spriteHeight,
					H: spriteHeight,
				}

				pos := game.Pos{X: x, Y: y}
				if a.loadedLevel.Debug[pos] {
					a.textureAtlas.SetColorMod(128, 0, 0)
				} else if tile.Seen && !tile.Visible {
					a.textureAtlas.SetColorMod(128, 128, 128)
				} else {
					a.textureAtlas.SetColorMod(255, 255, 255)
				}

				a.renderer.Copy(a.textureAtlas, &srcRect, &destRect)

				if tile.OverlaySymbol != game.EmptyTile {
					srcRect = a.textureIndex[tile.OverlaySymbol][0]
					a.renderer.Copy(a.textureAtlas, &srcRect, &destRect)
				}

			}
		}
	}

	a.textureAtlas.SetColorMod(255, 255, 255)
}

func (a *App) drawPlayer() {
	offsetX := (a.width / 2) - int32(a.centerX*spriteHeight)
	offsetY := (a.height / 2) - int32(a.centerY*spriteHeight)

	playerSrcRect := a.textureIndex[a.loadedLevel.Player.Symbol][0]
	playerDestRect := sdl.Rect{
		X: int32(a.loadedLevel.Player.X*spriteHeight) + offsetX,
		Y: int32(a.loadedLevel.Player.Y*spriteHeight) + offsetY,
		W: spriteHeight,
		H: spriteHeight,
	}
	a.renderer.Copy(a.textureAtlas, &playerSrcRect, &playerDestRect)
}

func (a *App) drawMonsters() {
	offsetX := (a.width / 2) - int32(a.centerX*spriteHeight)
	offsetY := (a.height / 2) - int32(a.centerY*spriteHeight)

	for pos, monster := range a.loadedLevel.Monsters {
		if a.loadedLevel.Tiles[pos.Y][pos.X].Visible {
			monsterSrcRect := a.textureIndex[monster.Symbol][0]
			monsterDestRect := sdl.Rect{
				X: int32(pos.X)*spriteHeight + offsetX,
				Y: int32(pos.Y)*spriteHeight + offsetY,
				W: spriteHeight,
				H: spriteHeight,
			}
			a.renderer.Copy(a.textureAtlas, &monsterSrcRect, &monsterDestRect)
		}
	}

}

func (a *App) drawFloorItems() {
	offsetX := (a.width / 2) - int32(a.centerX*spriteHeight)
	offsetY := (a.height / 2) - int32(a.centerY*spriteHeight)

	for pos, items := range a.loadedLevel.Items {
		if a.loadedLevel.Tiles[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcRect := a.textureIndex[item.Symbol][0]
				itemDestRect := sdl.Rect{
					X: int32(pos.X)*spriteHeight + offsetX,
					Y: int32(pos.Y)*spriteHeight + offsetY,
					W: spriteHeight,
					H: spriteHeight,
				}
				a.renderer.Copy(a.textureAtlas, &itemSrcRect, &itemDestRect)
			}
		}
	}
}

func (a *App) drawPickupBarItems() {
	inventoryStart := int32(float64(a.width) * 0.9)
	inventoryWIdth := a.width - inventoryStart
	itemSize := a.getItemSize()

	items := a.loadedLevel.Items[a.loadedLevel.Player.Pos]
	if len(items) == 0 {
		return
	}

	a.renderer.Copy(a.inventoryBackground, nil, &sdl.Rect{
		X: inventoryStart,
		Y: a.height - itemSize,
		W: inventoryWIdth,
		H: itemSize,
	})

	for i, item := range items {
		itemSrcRect := &a.textureIndex[item.Symbol][0]
		itemDestRect := a.getPickupItemRect(i)
		a.renderer.Copy(a.textureAtlas, itemSrcRect, itemDestRect)
	}
}

func (a *App) drawEventLog() {
	textStart := int32(float64(a.height) * 0.75)
	a.renderer.Copy(a.eventBackground, nil, &sdl.Rect{
		X: 0,
		Y: textStart,
		W: int32(float64(a.width) * 0.25),
		H: int32(float64(a.height) * 0.75),
	})

	_, fontSizeY, _ := a.smallFont.SizeUTF8("A")

	i := a.loadedLevel.EventPos
	count := 0
	for {
		event := a.loadedLevel.Events[i]
		if event != "" {
			tex := a.stringToTexture(event, smallFont, sdl.Color{R: 255, G: 0, B: 0})
			_, _, w, h, err := tex.Query()
			if err != nil {
				fmt.Println("Problem loading event: " + event)
			}
			a.renderer.Copy(tex, nil, &sdl.Rect{X: 0, Y: int32(count*fontSizeY) + textStart, W: w, H: h})
		}

		i = (i + 1) % (len(a.loadedLevel.Events))
		count++

		if i == a.loadedLevel.EventPos {
			break
		}
	}
}

func (a *App) drawInventory() {
	// draw inventory backdrop
	inventoryRect := a.getInventoryBackdropRect()
	a.renderer.Copy(a.inventoryBackground, nil, inventoryRect)

	// draw equipment bar
	weaponRect := a.getWeaponSlotRect()
	a.renderer.Copy(a.slotBackground, nil, weaponRect)
	if a.loadedLevel.Player.Weapon != nil {
		a.renderer.Copy(a.textureAtlas, &a.textureIndex[a.loadedLevel.Player.Weapon.Symbol][0], weaponRect)
	}

	armorRect := a.getArmorSlotRect()
	a.renderer.Copy(a.slotBackground, nil, armorRect)
	if a.loadedLevel.Player.Armor != nil {
		a.renderer.Copy(a.textureAtlas, &a.textureIndex[a.loadedLevel.Player.Armor.Symbol][0], armorRect)
	}

	// draw player in inventory
	playerSrcRect := a.textureIndex[a.loadedLevel.Player.Symbol][0]
	a.renderer.Copy(a.textureAtlas, &playerSrcRect, &sdl.Rect{
		X: inventoryRect.X + inventoryRect.X/2,
		Y: inventoryRect.Y + inventoryRect.Y/2,
		W: inventoryRect.W / 2,
		H: inventoryRect.H / 2,
	})

	// draw items in inventory
	for i, item := range a.loadedLevel.Player.Items {
		itemSrcRect := &a.textureIndex[item.Symbol][0]

		if item == a.dragged {
			itemSize := a.getItemSize()
			mx, my, _ := sdl.GetMouseState()
			itemDestRect := &sdl.Rect{
				X: mx - itemSize/2,
				Y: my - itemSize/2,
				W: itemSize,
				H: itemSize,
			}
			a.renderer.Copy(a.textureAtlas, itemSrcRect, itemDestRect)
		} else {
			itemDestRect := a.getInventoryItemRect(i)
			a.renderer.Copy(a.textureAtlas, itemSrcRect, itemDestRect)
		}
	}
}
