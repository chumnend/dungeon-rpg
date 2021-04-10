package ui

const spriteHeight = 32
const itemSizeRatio = 0.033

func (a *App) getItemSize() int32 {
	return int32(itemSizeRatio * float32(a.width))
}

func (a *App) getSlotSize() int32 {
	itemSize := a.getItemSize()
	return int32(float32(itemSize) * 1.5)
}
