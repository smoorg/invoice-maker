package receiver_edit

import (
	"invoice-maker/gui/modal"
	"invoice-maker/gui/receiver_add"
	"invoice-maker/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI, row int) {
	tui.AddAndSwitchToPage(
		types.PageReceiverEdit,
		modal.Modal(tui, types.PageReceiverList, editReceiver(tui, row), 50, 13, ""),
	)
}

func editReceiver(tui *types.TUI, row int) tview.Primitive {
	r := tui.Config.Receivers[row]

	save := func() {
		tui.Config.WriteReceiver(r, row)
		tui.Rerender()
		tui.SwitchToPage(types.PageReceiverList)
		tui.Pages.RemovePage(types.PageReceiverEdit)
	}

	cancel := func() {
		tui.SwitchToPage(types.PageReceiverList)
		tui.Pages.RemovePage(types.PageReceiverEdit)
	}

	return receiver_add.AddOrEditReceiver(&r, "Edit Receiver", save, cancel)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.SwitchToPage(types.PageReceiverList)
		tui.Pages.RemovePage(types.PageReceiverEdit)
	}
	return eventKey
}
