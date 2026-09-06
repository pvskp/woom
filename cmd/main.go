package main

import (
	"log"

	"github.com/pvskp/woom/cmd/commands"
)

func main() {
	err := commands.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
