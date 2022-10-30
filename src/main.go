package main

import (
	"invoice-maker/gui"
	"invoice-maker/types"
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
