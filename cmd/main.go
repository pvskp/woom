package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/pvskp/woom/internal/portal"
)


func main() {
	client, err := portal.New()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	ctx := context.Background()

	opts := portal.ScreenshotOptions{
		Mode:      portal.ScreenshotModeFull,
		Selection: portal.Coordinate{},
	}

	uri, err := client.Screenshot(ctx, opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(uri)

	loader := binsdl.Load()
	defer loader.Unload()

	imgLoader := binimg.Load()
	defer imgLoader.Unload()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(
		"woom",
		1280,
		720,
		sdl.WINDOW_BORDERLESS,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer window.Destroy()
	defer renderer.Destroy()

	if err := window.SetFullscreen(true); err != nil {
		log.Fatal(err)
	}

	texture, err := img.LoadTexture(
		renderer,
		strings.TrimPrefix(uri, "file://"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer texture.Destroy()

	running := true

	for running {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type {

			case sdl.EVENT_QUIT:
				running = false

			case sdl.EVENT_KEY_DOWN:
				if event.KeyboardEvent().Key == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		renderer.Clear()

		renderer.RenderTexture(
			texture,
			nil,
			nil,
		)

		renderer.Present()
	}

}
