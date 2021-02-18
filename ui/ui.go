package ui

import (
	"bufio"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/chumnend/simple-rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 1280
	windowHeight = 720
)

var window *sdl.Window
var renderer *sdl.Renderer
var textureAtlas *sdl.Texture
var textureIndex map[game.Tile][]sdl.Rect

func imgFileToTexture(filename string) *sdl.Texture {
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

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
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

func loadTextureIndex() {
	textureIndex = make(map[game.Tile][]sdl.Rect)

	file, err := os.Open("ui/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		tile := game.Tile(line[0])
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
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
			x = (x + 1)
			if x > 62 {
				x = 0
				y++
			}
		}

		// rect := sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32}
		textureIndex[tile] = rects
	}
}

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}

	window, err = sdl.CreateWindow("Simple RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	textureAtlas = imgFileToTexture("ui/assets/tiles.png")
	loadTextureIndex()
}

// UI struct represents the user interface for the game
type UI struct {
	game.UI
}

// Draw draws a level for the game
func (ui *UI) Draw(level *game.Level) {
	rand.Seed(1)

	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile == game.SpaceBlank {
				continue
			}

			srcRects := textureIndex[tile]
			srcRect := srcRects[rand.Intn(len(srcRects))]

			destRect := sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32}
			renderer.Copy(textureAtlas, &srcRect, &destRect)
		}
	}

	// 21, 59
	renderer.Copy(textureAtlas, &sdl.Rect{21 * 32, 59 * 32, 32, 32}, &sdl.Rect{int32(level.Player.X * 32), int32(level.Player.Y * 32), 32, 32})

	renderer.Present()
	sdl.Delay(3000)
}
