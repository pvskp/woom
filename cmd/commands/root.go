package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "woom",
	Short: "Simple tool for screenshot, OCR and overlay zooming",
}

func Execute() error {
	return rootCmd.Execute()
}
