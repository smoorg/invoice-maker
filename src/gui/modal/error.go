package modal

import (
	"invoice-maker/types"

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
			tui.Pages.RemovePage(types.PageErrModal)
			back()
		})

	tui.AddAndSwitchToPage(types.PageErrModal, Modal(tui, bgPageName, m, width, height, title))
}
