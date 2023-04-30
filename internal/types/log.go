package types

import (
	"fmt"
	"log"
)

func (tui *TUI) Fatal(v ...any) {
	err := fmt.Sprintf("%s\n", v...)
	tui.App.Stop()
	log.Fatal(err)
}
