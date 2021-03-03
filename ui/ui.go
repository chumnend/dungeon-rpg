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
	"github.com/veandco/go-sdl2/ttf"
)

func init() {
	var err error

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}
}

// App represents the application window that runs the RPG game
type App struct {
	width   int32
	height  int32
	centerX int
	centerY int

	r *rand.Rand

	window          *sdl.Window
	renderer        *sdl.Renderer
	textureAtlas    *sdl.Texture
	eventBackground *sdl.Texture
	str2TexSmall    map[string]*sdl.Texture
	str2TexMedium   map[string]*sdl.Texture
	str2TexLarge    map[string]*sdl.Texture

	textureIndex map[game.Tile][]sdl.Rect
	fontSmall    *ttf.Font
	fontMedium   *ttf.Font
	fontLarge    *ttf.Font

	levelCh chan *game.Level
	inputCh chan *game.Input
}

// NewApp returns an App struct
func NewApp(levelCh chan *game.Level, inputCh chan *game.Input) *App {
	var err error

	app := new(App)
	app.width = 1280
	app.height = 730
	app.centerX = -1
	app.centerY = -1

	app.r = rand.New(rand.NewSource(1))

	app.window, err = sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1280, 720, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	app.renderer, err = sdl.CreateRenderer(app.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	app.textureAtlas = app.imgFileToTexture("ui/assets/tiles.png")
	app.loadTextureIndex()

	app.eventBackground = app.getSinglePixelTexture(sdl.Color{R: 0, G: 0, B: 0, A: 128})

	app.str2TexSmall = make(map[string]*sdl.Texture)
	app.str2TexMedium = make(map[string]*sdl.Texture)
	app.str2TexLarge = make(map[string]*sdl.Texture)

	app.fontSmall, err = ttf.OpenFont("ui/assets/Kingthings.ttf", 24)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	app.fontMedium, err = ttf.OpenFont("ui/assets/Kingthings.ttf", 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	app.fontLarge, err = ttf.OpenFont("ui/assets/Kingthings.ttf", 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	app.levelCh = levelCh
	app.inputCh = inputCh

	return app
}

// Run starts the application window
func (a *App) Run() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				a.inputCh <- &game.Input{Type: game.QuitGame}
			case *sdl.KeyboardEvent:
				var input game.Input

				if e.Type == sdl.KEYUP {
					switch e.Keysym.Scancode {
					case sdl.SCANCODE_UP:
						input.Type = game.Up
					case sdl.SCANCODE_DOWN:
						input.Type = game.Down
					case sdl.SCANCODE_LEFT:
						input.Type = game.Left
					case sdl.SCANCODE_RIGHT:
						input.Type = game.Right
					case sdl.SCANCODE_S:
						input.Type = game.Search
					}

					a.inputCh <- &input
				}
			}
		}

		select {
		case newLevel, ok := <-a.levelCh:
			if ok {
				a.draw(newLevel)
			}
		default:
			// do nothing
		}

		sdl.Delay(10)
	}
}

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

func (a *App) loadTextureIndex() {
	a.textureIndex = make(map[game.Tile][]sdl.Rect)

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
		a.textureIndex[tile] = rects
	}
}

type fontSize int

const (
	fontSmall fontSize = iota
	fontMedium
	fontLarge
)

func (a *App) stringToTexture(s string, size fontSize, color sdl.Color) *sdl.Texture {
	var font *ttf.Font
	switch size {
	case fontSmall:
		font = a.fontSmall
		if tex, exists := a.str2TexSmall[s]; exists {
			return tex
		}
	case fontMedium:
		font = a.fontMedium
		if tex, exists := a.str2TexMedium[s]; exists {
			return tex
		}
	case fontLarge:
		font = a.fontLarge
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
	case fontSmall:
		a.str2TexSmall[s] = tex
	case fontMedium:
		font = a.fontMedium
		a.str2TexMedium[s] = tex
	case fontLarge:
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

func (a *App) draw(level *game.Level) {
	a.r.Seed(1)
	a.renderer.Clear()
	if a.centerX == -1 {
		a.centerX = level.Player.X
		a.centerY = level.Player.Y
	}

	// move the camera with the player
	limit := 5
	if level.Player.X > a.centerX+limit {
		a.centerX++
	} else if level.Player.X < a.centerX-limit {
		a.centerX--
	}

	if level.Player.Y > a.centerY+limit {
		a.centerY++
	} else if level.Player.Y < a.centerY-limit {
		a.centerY--
	}

	offsetX := (a.width / 2) - int32(a.centerX*32)
	offsetY := (a.height / 2) - int32(a.centerY*32)

	// draw floor tiles
	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile == game.EmptyTile {
				continue
			}

			srcRects := a.textureIndex[tile]
			srcRect := srcRects[a.r.Intn(len(srcRects))]
			destRect := sdl.Rect{
				X: int32(x*32) + offsetX,
				Y: int32(y*32) + offsetY,
				W: 32,
				H: 32,
			}

			pos := game.Pos{X: x, Y: y}
			if level.Debug[pos] {
				a.textureAtlas.SetColorMod(128, 0, 0)
			} else {
				a.textureAtlas.SetColorMod(255, 255, 255)
			}

			a.renderer.Copy(a.textureAtlas, &srcRect, &destRect)
		}
	}

	// draw player
	playerSrcRect := a.textureIndex[level.Player.Symbol][0]
	playerDestRect := sdl.Rect{
		X: int32(level.Player.X*32) + offsetX,
		Y: int32(level.Player.Y*32) + offsetY,
		W: 32,
		H: 32,
	}
	a.renderer.Copy(a.textureAtlas, &playerSrcRect, &playerDestRect)

	// draw monsters
	for pos, monster := range level.Monsters {
		monsterSrcRect := a.textureIndex[monster.Symbol][0]
		monsterDestRect := sdl.Rect{X: int32(pos.X)*32 + offsetX, Y: int32(pos.Y)*32 + offsetY, W: 32, H: 32}
		a.renderer.Copy(a.textureAtlas, &monsterSrcRect, &monsterDestRect)
	}

	// draw event log
	textStart := int32(float64(a.height) * 0.75)
	a.renderer.Copy(a.eventBackground, nil, &sdl.Rect{
		X: 0,
		Y: textStart,
		W: int32(float64(a.width) * 0.25),
		H: int32(float64(a.height) * 0.75),
	})

	for idx, event := range level.Events {
		if event != "" {
			tex := a.stringToTexture(event, fontMedium, sdl.Color{R: 255, G: 0, B: 0})
			_, _, w, h, err := tex.Query()
			if err != nil {
				fmt.Println("Problem loading event: " + event)
			}
			a.renderer.Copy(tex, nil, &sdl.Rect{X: 0, Y: int32(idx*32) + textStart, W: w, H: h})
		}
	}

	a.renderer.Present()
}
