package ui

import (
	"bufio"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/chumnend/dungeon-rpg/internal/game"
	"github.com/veandco/go-sdl2/mix"
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

	err = mix.Init(mix.INIT_OGG)
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

	r         *rand.Rand
	game      *game.Game
	lastLevel *game.Level

	window              *sdl.Window
	renderer            *sdl.Renderer
	textureAtlas        *sdl.Texture
	textureIndex        map[rune][]sdl.Rect
	eventBackground     *sdl.Texture
	inventoryBackground *sdl.Texture
	str2TexSmall        map[string]*sdl.Texture
	str2TexMedium       map[string]*sdl.Texture
	str2TexLarge        map[string]*sdl.Texture
	smallFont           *ttf.Font
	mediumFont          *ttf.Font
	largeFont           *ttf.Font
	footstepSounds      []*mix.Chunk
	doorOpenSounds      []*mix.Chunk
}

// NewApp returns an App struct
func NewApp(game *game.Game, width, height int32) *App {
	window, err := sdl.CreateWindow("RPG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	r := rand.New(rand.NewSource(1))

	smallFont, err := ttf.OpenFont("internal/ui/assets/fonts/Kingthings.ttf", int(float64(width)*0.015))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	mediumFont, err := ttf.OpenFont("internal/ui/assets/fonts/Kingthings.ttf", int(float64(width)*0.025))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	largeFont, err := ttf.OpenFont("internal/ui/assets/fonts/Kingthings.ttf", int(float64(width)*0.05))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		panic(err)
	}

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open audio: %s\n", err)
		panic(err)
	}

	music, err := mix.LoadMUS("internal/ui/assets/sound/ambient.ogg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to play music: %s\n", err)
		panic(err)
	}
	music.Play(-1)

	footstepSounds := make([]*mix.Chunk, 0)
	footstepBase := "internal/ui/assets/sound/footstep0"
	for i := 0; i < 6; i++ {
		path := footstepBase + strconv.Itoa(i) + ".ogg"
		wav, err := mix.LoadWAV(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open sound file: %s\n", err)
			panic(err)
		}

		footstepSounds = append(footstepSounds, wav)
	}

	doorOpenSounds := make([]*mix.Chunk, 0)
	doorOpenBase := "internal/ui/assets/sound/doorOpen_"
	for i := 1; i < 3; i++ {
		path := doorOpenBase + strconv.Itoa(i) + ".ogg"
		wav, err := mix.LoadWAV(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open sound file: %s\n", err)
			panic(err)
		}

		doorOpenSounds = append(doorOpenSounds, wav)
	}

	app := &App{
		width:          width,
		height:         height,
		centerX:        -1,
		centerY:        -1,
		r:              r,
		game:           game,
		lastLevel:      nil,
		window:         window,
		renderer:       renderer,
		str2TexSmall:   make(map[string]*sdl.Texture),
		str2TexMedium:  make(map[string]*sdl.Texture),
		str2TexLarge:   make(map[string]*sdl.Texture),
		smallFont:      smallFont,
		mediumFont:     mediumFont,
		largeFont:      largeFont,
		footstepSounds: footstepSounds,
		doorOpenSounds: doorOpenSounds,
	}

	app.textureAtlas = app.imgFileToTexture("internal/ui/assets/tiles/tiles.png")
	app.textureIndex = app.loadTextureIndex("internal/ui/assets/atlas-index.txt")
	app.eventBackground = app.getSinglePixelTexture(sdl.Color{R: 0, G: 0, B: 0, A: 128})
	app.inventoryBackground = app.getSinglePixelTexture(sdl.Color{R: 255, G: 0, B: 0, A: 128})

	return app
}

// Start starts the application window
func (a *App) Start() {
	go a.game.Run()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				a.game.InputCh <- &game.Input{Type: game.QuitGame}
				return

			case *sdl.MouseButtonEvent:
				var input game.Input

				if e.Type == sdl.MOUSEBUTTONUP {
					mx, my, _ := sdl.GetMouseState()

					item := a.checkForItem(a.lastLevel, mx, my)
					if item != nil {
						input.Type = game.TakeItem
						input.Item = item

						a.game.InputCh <- &input
					}
				}

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
					case sdl.SCANCODE_T:
						input.Type = game.TakeAll
					default:
						input.Type = game.None
					}

					a.game.InputCh <- &input
				}
			}
		}

		select {
		case loadedLevel, ok := <-a.game.LevelCh:
			if ok {
				a.lastLevel = loadedLevel

				switch loadedLevel.LastEvent {
				case game.Move:
					playRandomSound(a.footstepSounds, 64)
				case game.DoorOpen:
					playRandomSound(a.doorOpenSounds, 64)
				default:
					// do nothing
				}

				a.draw(loadedLevel)
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
			rects = append(rects, sdl.Rect{X: int32(x * 32), Y: int32(y * 32), W: 32, H: 32})
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

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}

func (a *App) draw(level *game.Level) {
	a.renderer.Clear()
	a.r.Seed(1)

	if a.centerX == -1 && a.centerY == -1 {
		a.centerX = level.Player.X
		a.centerY = level.Player.Y
	}

	// move the camera with the player
	limit := 5
	if level.Player.X > a.centerX+limit {
		diff := level.Player.X - (a.centerX + limit)
		a.centerX += diff
	} else if level.Player.X < a.centerX-limit {
		diff := (a.centerX - limit) - level.Player.X
		a.centerX -= diff
	}

	if level.Player.Y > a.centerY+limit {
		diff := level.Player.Y - (a.centerY + limit)
		a.centerY += diff
	} else if level.Player.Y < a.centerY-limit {
		diff := (a.centerY - limit) - level.Player.Y
		a.centerY -= diff
	}

	offsetX := (a.width / 2) - int32(a.centerX*32)
	offsetY := (a.height / 2) - int32(a.centerY*32)

	// draw floor tiles
	for y, row := range level.Tiles {
		for x, tile := range row {
			if tile.Symbol == game.EmptyTile {
				continue
			}

			srcRects := a.textureIndex[tile.Symbol]
			srcRect := srcRects[a.r.Intn(len(srcRects))]

			if tile.Visible || tile.Seen {
				destRect := sdl.Rect{
					X: int32(x*32) + offsetX,
					Y: int32(y*32) + offsetY,
					W: 32,
					H: 32,
				}

				pos := game.Pos{X: x, Y: y}
				if level.Debug[pos] {
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
		if level.Tiles[pos.Y][pos.X].Visible {
			monsterSrcRect := a.textureIndex[monster.Symbol][0]
			monsterDestRect := sdl.Rect{X: int32(pos.X)*32 + offsetX, Y: int32(pos.Y)*32 + offsetY, W: 32, H: 32}
			a.renderer.Copy(a.textureAtlas, &monsterSrcRect, &monsterDestRect)
		}
	}

	// draw items
	for pos, items := range level.Items {
		if level.Tiles[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcRect := a.textureIndex[item.Symbol][0]
				itemDestRect := sdl.Rect{X: int32(pos.X)*32 + offsetX, Y: int32(pos.Y)*32 + offsetY, W: 32, H: 32}
				a.renderer.Copy(a.textureAtlas, &itemSrcRect, &itemDestRect)
			}
		}
	}

	// draw event log
	textStart := int32(float64(a.height) * 0.75)
	a.renderer.Copy(a.eventBackground, nil, &sdl.Rect{
		X: 0,
		Y: textStart,
		W: int32(float64(a.width) * 0.25),
		H: int32(float64(a.height) * 0.75),
	})

	_, fontSizeY, _ := a.smallFont.SizeUTF8("A")

	i := level.EventPos
	count := 0
	for {
		event := level.Events[i]
		if event != "" {
			tex := a.stringToTexture(event, smallFont, sdl.Color{R: 255, G: 0, B: 0})
			_, _, w, h, err := tex.Query()
			if err != nil {
				fmt.Println("Problem loading event: " + event)
			}
			a.renderer.Copy(tex, nil, &sdl.Rect{X: 0, Y: int32(count*fontSizeY) + textStart, W: w, H: h})
		}

		i = (i + 1) % (len(level.Events))
		count++

		if i == level.EventPos {
			break
		}
	}

	// draw inventory
	inventoryStart := int32(float64(a.width) * 0.9)
	inventoryWIdth := a.width - inventoryStart

	a.renderer.Copy(a.inventoryBackground, nil, &sdl.Rect{
		X: inventoryStart,
		Y: a.height - 32,
		W: inventoryWIdth,
		H: 32,
	})

	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemSrcRect := &a.textureIndex[item.Symbol][0]
		itemDestRect := a.getItemRect(i)
		a.renderer.Copy(a.textureAtlas, itemSrcRect, itemDestRect)
	}

	a.renderer.Present()
}

func (a *App) getItemRect(i int) *sdl.Rect {
	return &sdl.Rect{
		X: a.width - 32 - int32(i*32),
		Y: a.height - 32,
		W: 32,
		H: 32,
	}
}

func (a *App) checkForItem(level *game.Level, mx int32, my int32) *game.Item {
	mouseRect := &sdl.Rect{
		X: mx,
		Y: my,
		W: 1,
		H: 1,
	}

	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemRect := a.getItemRect(i)
		if itemRect.HasIntersection(mouseRect) {
			return item
		}
	}

	return nil
}
