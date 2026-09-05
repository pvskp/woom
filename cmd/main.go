package main

import (
	"context"
	"log"
	"strings"

	"github.com/pvskp/woom/internal/overlay"
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

	overlay := overlay.New()

	err = overlay.Load(strings.TrimPrefix(uri, "file://"))

	if err != nil {
		log.Fatal(err)
	}

}
