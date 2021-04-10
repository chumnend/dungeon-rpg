package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (a *App) getPickupItemRect(i int) *sdl.Rect {
	itemSize := int32(itemSizeRatio * float32(a.width))

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
	itemSize := int32(itemSizeRatio * float32(a.width))

	return &sdl.Rect{
		X: inventoryRect.X + int32(i)*itemSize,
		Y: inventoryRect.Y + inventoryRect.H - itemSize,
		W: itemSize,
		H: itemSize,
	}
}
