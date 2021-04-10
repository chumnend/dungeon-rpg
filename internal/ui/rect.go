package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (a *App) getMouseRect(mx int32, my int32) *sdl.Rect {
	return &sdl.Rect{
		X: mx,
		Y: my,
		W: 1,
		H: 1,
	}
}

func (a *App) getPickupItemRect(i int) *sdl.Rect {
	itemSize := a.getItemSize()

	return &sdl.Rect{
		X: a.width - itemSize - int32(i)*itemSize,
		Y: a.height - itemSize,
		W: itemSize,
		H: itemSize,
	}
}

func (a *App) getInventoryBackdropRect() *sdl.Rect {
	inventoryWidth := int32(a.width / 2)
	inventoryHeight := int32(a.height * 3 / 4)
	offsetX := (a.width - inventoryWidth) / 2
	offsetY := (a.height - inventoryHeight) / 2

	return &sdl.Rect{
		X: offsetX,
		Y: offsetY,
		W: inventoryWidth,
		H: inventoryHeight,
	}
}

func (a *App) getInventoryItemRect(i int) *sdl.Rect {
	inventoryRect := a.getInventoryBackdropRect()
	itemSize := a.getItemSize()

	return &sdl.Rect{
		X: inventoryRect.X + int32(i)*itemSize,
		Y: inventoryRect.Y + inventoryRect.H - itemSize,
		W: itemSize,
		H: itemSize,
	}
}

func (a *App) getWeaponSlotRect() *sdl.Rect {
	inventoryRect := a.getInventoryBackdropRect()
	slotSize := a.getSlotSize()

	return &sdl.Rect{
		X: inventoryRect.X + inventoryRect.X/2 - slotSize/2,
		Y: inventoryRect.Y + inventoryRect.H/3,
		W: slotSize,
		H: slotSize,
	}
}

func (a *App) getArmorSlotRect() *sdl.Rect {
	inventoryRect := a.getInventoryBackdropRect()
	slotSize := a.getSlotSize()

	return &sdl.Rect{
		X: inventoryRect.X + 3*inventoryRect.X/2 + slotSize/2,
		Y: inventoryRect.Y + inventoryRect.H/3,
		W: slotSize,
		H: slotSize,
	}
}
