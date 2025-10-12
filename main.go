package main

import (
	"log"

	"github.com/rivo/tview"

	"invoice-maker/internal/gui"
	"invoice-maker/internal/types"
)

var tui *types.TUI

func main() {
	tui = &types.TUI{}
	tui.App = tview.NewApplication()
	tui.Pages = tview.NewPages()

	err := gui.Run(tui)
	if err != nil {
		log.Fatal(err)
	}
}
