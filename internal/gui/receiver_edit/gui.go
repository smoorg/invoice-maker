package receiver_edit

import (
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/gui/receiver_add"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI, row int) {
	tui.AddAndSwitchToPage(
		types.PageReceiverEdit,
		modal.Modal(tui, types.PageReceiverEdit, types.PageReceiverList, editReceiver(tui, row), 50, 13, ""),
	)
}

func goBack(tui *types.TUI) {
	tui.Pages.RemovePage(types.PageReceiverEdit)
	tui.SwitchToPage(types.PageReceiverList)
}

func editReceiver(tui *types.TUI, row int) tview.Primitive {
	r := tui.Config.Receivers[row]

	save := func() {
		tui.Config.WriteReceiver(r, row)
		goBack(tui)
	}

	cancel := func() {
		goBack(tui)
	}

	return receiver_add.AddOrEditReceiver(&r, "Edit Receiver", save, cancel)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		goBack(tui)
	}
	return eventKey
}
