package ui

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

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

	state       appState
	r           *rand.Rand
	game        *game.Game
	loadedLevel *game.Level
	dragged     *game.Item

	window       *sdl.Window
	renderer     *sdl.Renderer
	textureAtlas *sdl.Texture
	textureIndex map[rune][]sdl.Rect

	eventBackground     *sdl.Texture
	inventoryBackground *sdl.Texture
	slotBackground      *sdl.Texture

	str2TexSmall  map[string]*sdl.Texture
	str2TexMedium map[string]*sdl.Texture
	str2TexLarge  map[string]*sdl.Texture
	smallFont     *ttf.Font
	mediumFont    *ttf.Font
	largeFont     *ttf.Font

	footstepSounds []*mix.Chunk
	doorOpenSounds []*mix.Chunk
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

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

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

	a := &App{
		width:          width,
		height:         height,
		centerX:        -1,
		centerY:        -1,
		state:          mainState,
		r:              r,
		game:           game,
		loadedLevel:    nil,
		dragged:        nil,
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

	a.textureAtlas = a.imgFileToTexture("internal/ui/assets/tiles/tiles.png")
	a.textureIndex = a.loadTextureIndex("internal/ui/assets/atlas-index.txt")

	a.eventBackground = a.getSinglePixelTexture(sdl.Color{R: 0, G: 0, B: 0, A: 128})
	a.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	a.inventoryBackground = a.getSinglePixelTexture(sdl.Color{R: 149, G: 84, B: 19, A: 200})
	a.inventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	a.slotBackground = a.getSinglePixelTexture(sdl.Color{R: 0, G: 0, B: 0, A: 255})

	return a
}

// Start starts the application window
func (a *App) Start() {

	// run the game engine
	go a.game.Run()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				input := game.Input{
					Type: game.QuitGame,
				}

				a.game.InputCh <- &input
				return

			// check mouse events
			case *sdl.MouseButtonEvent:
				switch a.state {
				case mainState:
					input := game.Input{
						Type: game.None,
					}

					if e.Type == sdl.MOUSEBUTTONUP {
						item := a.checkForFloorItem(a.loadedLevel, e.X, e.Y)
						if item != nil {
							input.Type = game.TakeItem
							input.Item = item
							a.game.InputCh <- &input
						}
					}

				case inventoryState:
					input := game.Input{
						Type: game.None,
					}

					if e.Type == sdl.MOUSEBUTTONDOWN {
						// look for drag event if in inventory
						item := a.checkForInventoryItem(a.loadedLevel, e.X, e.Y)
						if item != nil {
							a.dragged = item
						}
					}

					if e.Type == sdl.MOUSEBUTTONUP {
						if a.dragged != nil {
							shouldDrop := a.checkForDropItem(a.loadedLevel, e.X, e.Y)
							if shouldDrop {
								input.Type = game.DropItem
								input.Item = a.dragged

								a.game.InputCh <- &input
							}
							a.dragged = nil
						}
					}
				}

			// check keyboard events
			case *sdl.KeyboardEvent:
				switch a.state {
				case mainState:
					input := game.Input{
						Type: game.None,
					}

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
						case sdl.SCANCODE_I:
							a.toggleInventory()
						case sdl.SCANCODE_T:
							input.Type = game.TakeAll
						default:
							// do nothing
						}

						a.game.InputCh <- &input
					}

				case inventoryState:
					input := game.Input{
						Type: game.None,
					}

					if e.Type == sdl.KEYUP {
						switch e.Keysym.Scancode {
						case sdl.SCANCODE_I:
							a.toggleInventory()
						default:
							// do nothing
						}

						a.game.InputCh <- &input
					}
				}
			}
		}

		// check for level update
		select {
		case loadedLevel, ok := <-a.game.LevelCh:
			if ok {
				a.loadedLevel = loadedLevel // keep track of the loaded level

				switch loadedLevel.LastEvent {
				case game.Move:
					playRandomSound(a.footstepSounds, 64)
				case game.DoorOpen:
					playRandomSound(a.doorOpenSounds, 64)
				default:
					// do nothing
				}
			}
		default:
			// do nothing
		}

		// draw to the window
		a.draw()

		sdl.Delay(10)
	}
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	chunkIndex := rand.Intn(len(chunks))
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}
