package main

import (
	"fmt"
	"log"
	"os"

	"invoice-maker/cmd/v2/root"

	tea "github.com/charmbracelet/bubbletea"
)

const logFileName = "logs.txt"

func main() {
	err := os.Remove(logFileName)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("This log message will be written to the file.")

	m := root.NewRootModel()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
