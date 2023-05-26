package modal

import (
	"invoice-maker/internal/types"

	"github.com/rivo/tview"
)

type ErrorModal struct {
	*tview.Modal
}

func Error(tui *types.TUI, errMsg string, bgPageName string, width int, height int, title string, back func()) {
	m := tview.NewModal().
		SetText(errMsg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			tui.Pages.RemovePage(types.PageModal)
			back()
		})

	tui.AddAndSwitchToPage(types.PageModal, Modal(tui, types.PageModal, bgPageName, m, width, height, title))
}
