package ui

import (
	"bufio"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func (a *App) imgFileToTexture(filename string) *sdl.Texture {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	pixels := make([]byte, w*h*4)
	idx := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[idx] = byte(r / 256)
			idx++
			pixels[idx] = byte(g / 256)
			idx++
			pixels[idx] = byte(b / 256)
			idx++
			pixels[idx] = byte(a / 256)
			idx++
		}
	}

	tex, err := a.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic(err)
	}

	tex.Update(nil, pixels, w*4)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}

func (a *App) loadTextureIndex(filename string) map[rune][]sdl.Rect {
	textureIndex := make(map[rune][]sdl.Rect)

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		tile := rune(line[0])
		tileInfo := strings.Split(line[1:], ",") // x, y, variation

		x, err := strconv.ParseInt(strings.TrimSpace(tileInfo[0]), 10, 64)
		if err != nil {
			panic(err)
		}

		y, err := strconv.ParseInt(strings.TrimSpace(tileInfo[1]), 10, 64)
		if err != nil {
			panic(err)
		}

		variation, err := strconv.ParseInt(strings.TrimSpace(tileInfo[2]), 10, 64)
		if err != nil {
			panic(err)
		}

		rects := make([]sdl.Rect, 0)
		for i := 0; i < int(variation); i++ {
			rects = append(rects, sdl.Rect{
				X: int32(x * spriteHeight),
				Y: int32(y * spriteHeight),
				W: spriteHeight,
				H: spriteHeight,
			})
			x = (x + 1)
			if x > 62 {
				x = 0
				y++
			}
		}

		textureIndex[tile] = rects
	}

	return textureIndex
}

type fontSize int

const (
	smallFont fontSize = iota
	mediumFont
	largeFont
)

func (a *App) stringToTexture(s string, size fontSize, color sdl.Color) *sdl.Texture {
	var font *ttf.Font
	switch size {
	case smallFont:
		font = a.smallFont
		if tex, exists := a.str2TexSmall[s]; exists {
			return tex
		}
	case mediumFont:
		font = a.mediumFont
		if tex, exists := a.str2TexMedium[s]; exists {
			return tex
		}
	case largeFont:
		font = a.largeFont
		if tex, exists := a.str2TexLarge[s]; exists {
			return tex
		}
	}

	surface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	tex, err := a.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}

	switch size {
	case smallFont:
		a.str2TexSmall[s] = tex
	case mediumFont:
		a.str2TexMedium[s] = tex
	case largeFont:
		a.str2TexLarge[s] = tex
	}

	return tex
}

func (a *App) getSinglePixelTexture(color sdl.Color) *sdl.Texture {
	tex, err := a.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}

	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A

	tex.Update(nil, pixels, 4)
	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}
