package overlay

import (
	"fmt"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
)

const (
	windowName   = "woom"
	windowWidth  = 1280
	windowHeight = 720

	zoomStep = 0.1
)

type SDLOverlay struct{}

func New() SDLOverlay {
	return SDLOverlay{}
}

func (s SDLOverlay) Load(imagePath string) error {
	loader := binsdl.Load()
	defer loader.Unload()

	imgLoader := binimg.Load()
	defer imgLoader.Unload()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(
		windowName,
		windowWidth,
		windowHeight,
		sdl.WINDOW_BORDERLESS,
	)
	if err != nil {
		return err
	}

	defer window.Destroy()
	defer renderer.Destroy()

	if err := window.SetFullscreen(true); err != nil {
		return err
	}

	texture, err := img.LoadTexture(
		renderer,
		imagePath,
	)

	if err != nil {
		return err
	}

	defer texture.Destroy()

	running := true

	zoom := float32(1)
	ox := float32(0)
	oy := float32(0)

	for running {
		w, h, err := renderer.CurrentOutputSize()
		if err != nil {
			return err
		}

		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type {

			case sdl.EVENT_QUIT:
				running = false

			case sdl.EVENT_KEY_DOWN:
				if event.KeyboardEvent().Key == sdl.K_ESCAPE {
					running = false
				}

			case sdl.EVENT_MOUSE_WHEEL:
				ev := event.MouseWheelEvent()

				previous := zoom
				zoom += ev.Y * zoomStep
				zoom = max(1, zoom)
				fmt.Println(zoom)

				k := zoom / previous
				ox = ev.MouseX - (ev.MouseX-ox)*k
				oy = ev.MouseY - (ev.MouseY-oy)*k
			}
		}

		renderer.Clear()

		dst := &sdl.FRect{
			X: ox,
			Y: oy,
			W: float32(w) * zoom,
			H: float32(h) * zoom,
		}

		renderer.RenderTexture(
			texture,
			nil,
			dst,
		)

		renderer.Present()
	}

	return nil
}
