package ui

import (
	"bufio"
	"fmt"
	"image/png"
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
var textureIndex map[game.Tile]sdl.Rect

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
	file, err := os.Open("ui/assets/atlas-index.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		tile := game.Tile(line[0])
		pos := strings.Split(line[1:], ",")

		posX := strings.TrimSpace(pos[0])
		x, err := strconv.ParseInt(posX, 10, 64)
		if err != nil {
			panic(err)
		}

		posY := strings.TrimSpace(pos[1])
		y, err := strconv.ParseInt(posY, 10, 64)
		if err != nil {
			panic(err)
		}

		fmt.Println(tile, x, y)
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
	fmt.Println("Hello World")

	renderer.Copy(textureAtlas, nil, nil)
	renderer.Present()

	sdl.Delay(5000)
}
