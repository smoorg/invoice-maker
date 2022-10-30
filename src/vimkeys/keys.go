package vimkeys

import (
	"github.com/gdamore/tcell/v2"
)

func Back(eventKey *tcell.EventKey) bool {
	return eventKey.Key() == tcell.KeyESC || eventKey.Rune() == 'h'
}

func Up(eventKey *tcell.EventKey) bool {
	return eventKey.Key() == tcell.KeyUp || eventKey.Rune() == 'k'
}

func Down(eventKey *tcell.EventKey) bool {
	return eventKey.Key() == tcell.KeyDown || eventKey.Rune() == 'j'
}

func Enter(eventKey *tcell.EventKey) bool {
	return eventKey.Rune() == 'l' || eventKey.Key() == tcell.KeyEnter
}
