package ui

import (
	"bufio"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/chumnend/simple-rpg/game"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
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

	window       *sdl.Window
	renderer     *sdl.Renderer
	textureAtlas *sdl.Texture
	textureIndex map[game.Tile][]sdl.Rect

	levelCh chan *game.Level
	inputCh chan *game.Input

	r *rand.Rand
}

// NewApp returns an App struct
func NewApp(levelCh chan *game.Level, inputCh chan *game.Input) *App {
	app := new(App)
	app.width = 1280
	app.height = 730
	app.centerX = -1
	app.centerY = -1

	var err error

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

	app.levelCh = levelCh
	app.inputCh = inputCh

	app.r = rand.New(rand.NewSource(1))

	return app
}

// Run starts the application window
func (a *App) Run() {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				a.inputCh <- &game.Input{Type: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					a.inputCh <- &game.Input{Type: game.CloseWindow, LevelCh: a.levelCh}
				}
			case *sdl.KeyboardEvent:
				code := e.Keysym.Scancode
				if e.Type == sdl.KEYUP && code == sdl.SCANCODE_UP {
					a.inputCh <- &game.Input{Type: game.Up}
				}
				if e.Type == sdl.KEYUP && code == sdl.SCANCODE_DOWN {
					a.inputCh <- &game.Input{Type: game.Down}
				}
				if e.Type == sdl.KEYUP && code == sdl.SCANCODE_LEFT {
					a.inputCh <- &game.Input{Type: game.Left}
				}
				if e.Type == sdl.KEYUP && code == sdl.SCANCODE_RIGHT {
					a.inputCh <- &game.Input{Type: game.Right}
				}
				if e.Type == sdl.KEYUP && code == sdl.SCANCODE_S {
					a.inputCh <- &game.Input{Type: game.Search}
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

// Draw draws a level for the game
func (a *App) draw(level *game.Level) {
	if a.centerX == -1 {
		a.centerX = level.Player.X
		a.centerY = level.Player.Y
	}

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

	a.r.Seed(1)
	a.renderer.Clear()

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

	srcRect := &sdl.Rect{X: 21 * 32, Y: 59 * 32, W: 32, H: 32}
	destRect := &sdl.Rect{X: int32(level.Player.X*32) + offsetX, Y: int32(level.Player.Y*32) + offsetY, W: 32, H: 32}
	a.renderer.Copy(a.textureAtlas, srcRect, destRect)
	a.renderer.Present()
}
