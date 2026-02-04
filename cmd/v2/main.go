package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"invoice-maker/cmd/v2/root"
)

func main() {
	m := root.NewRootModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
