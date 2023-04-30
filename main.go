package main

import (
	"invoice-maker/internal/gui"
	"invoice-maker/internal/types"
	"log"
)

var tui *types.TUI

func main() {
	tui = &types.TUI{}

	err := gui.Run(tui)
	if err != nil {
		log.Fatal(err)
	}
}
