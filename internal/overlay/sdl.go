package overlay

import (
	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
)

const (
	defaultWindowName   = "woom"
	defaultZoomStep = 0.5

	windowWidth  = 1280
	windowHeight = 720

)

type SDLOverlay struct {
	windowName string
	zoomStep   float32
}

func New(options ...func(*SDLOverlay)) *SDLOverlay {
	ov := &SDLOverlay{
		windowName: defaultWindowName,
		zoomStep:   defaultZoomStep,
	}

	for _, opt := range options {
		opt(ov)
	}

	return ov
}

func WithWindowName(name string) func(*SDLOverlay) {
	return func(s *SDLOverlay)  {
		s.windowName = name
	}
}

func WithZoomStep(step float32) func(*SDLOverlay) {
	return func(s *SDLOverlay)  {
		s.zoomStep = step
	}
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

	dragging := false

	window, renderer, err := sdl.CreateWindowAndRenderer(
		defaultWindowName,
		windowWidth,
		windowHeight,
		sdl.WINDOW_FULLSCREEN,
	)

	if err != nil {
		return err
	}

	defer window.Destroy()
	defer renderer.Destroy()

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

	handCursor, err := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_MOVE)
	if err != nil {
		return err
	}

	defer handCursor.Destroy()

	defaultCursor, err := sdl.GetDefaultCursor()
	if err != nil {
		return err
	}

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
				key := event.KeyboardEvent().Key
				if key == sdl.K_ESCAPE || key == sdl.K_Q {
					running = false
				}

			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				ev := event.MouseButtonEvent()
				if ev.Button == uint8(sdl.BUTTON_LEFT) {
					dragging = true
					sdl.SetCursor(handCursor)
				}

			case sdl.EVENT_MOUSE_BUTTON_UP:
				ev := event.MouseButtonEvent()
				if ev.Button == uint8(sdl.BUTTON_LEFT) {
					dragging = false
					sdl.SetCursor(defaultCursor)
				}

			case sdl.EVENT_MOUSE_MOTION:
				ev := event.MouseMotionEvent()
				if dragging {
					ox += ev.Xrel
					oy += ev.Yrel

					ox = min(max(ox, float32(w)*(1-zoom)), 0)
					oy = min(max(oy, float32(h)*(1-zoom)), 0)
				}

			case sdl.EVENT_MOUSE_WHEEL:
				ev := event.MouseWheelEvent()

				previous := zoom
				zoom += ev.Y * s.zoomStep
				zoom = max(1, zoom)

				k := zoom / previous
				ox = ev.MouseX - (ev.MouseX-ox)*k
				oy = ev.MouseY - (ev.MouseY-oy)*k

				ox = min(max(ox, float32(w)*(1-zoom)), 0)
				oy = min(max(oy, float32(h)*(1-zoom)), 0)
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
