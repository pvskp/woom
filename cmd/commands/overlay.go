package commands

import (
	"github.com/spf13/cobra"
	"context"
	"log"
	"strings"

	"github.com/pvskp/woom/internal/overlay"
	"github.com/pvskp/woom/internal/portal"
)

var overlayCmd = &cobra.Command{
	Use: "overlay",
	Short: "Place an overlay on top of the screen, which allows zooming and temporary edits",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runOverlay()
	},
}

func runOverlay() error {
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

	overlay := overlay.New(overlay.WithZoomStep(0.8))

	err = overlay.Load(strings.TrimPrefix(uri, "file://"))

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(overlayCmd)
}
