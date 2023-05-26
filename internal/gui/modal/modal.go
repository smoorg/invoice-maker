package modal

import (
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Returns a new primitive which puts the provided primitive in the center and
// sets its size to the given width and height.
func Modal(tui *types.TUI, pageName string, previousPage string, content tview.Primitive, width, height int, title string) tview.Primitive {
	_, page := tui.Pages.GetFrontPage()

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	return tview.NewPages().
		AddPage(previousPage, page, true, true).
		AddPage(pageName, modal, true, true)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if vimkeys.Back(eventKey) ||
		eventKey.Key() == tcell.KeyESC ||
		eventKey.Key() == tcell.KeyEnter {
		tui.Pages.RemovePage(tui.ActivePage)
		tui.SwitchToPrevious()
	}

	return nil
}
