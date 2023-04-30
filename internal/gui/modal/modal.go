package modal

import (
	"invoice-maker/internal/types"

	"github.com/rivo/tview"
)

// Returns a new primitive which puts the provided primitive in the center and
// sets its size to the given width and height.
func Modal(tui *types.TUI, base string, content tview.Primitive, width, height int, title string) tview.Primitive {
	tui.Pages.SwitchToPage(base)

	_, page := tui.Pages.GetFrontPage()

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	return tview.NewPages().
		AddPage(base, page, true, true).
		AddPage(types.PageModal, modal, true, true)
}
