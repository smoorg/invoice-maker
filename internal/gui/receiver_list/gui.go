package receiver_list

import (
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/gui/receiver_add"
	"invoice-maker/internal/gui/receiver_edit"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var selectedRow int = 0

func Render(tui *types.TUI) {
	tui.Pages.AddPage(types.PageReceiverList, ListReceivers(tui), true, false)
}

func ListReceivers(tui *types.TUI) tview.Primitive {
	table := tview.NewTable().SetSelectable(true, false).SetBorders(true)

	table.SetFixed(0, 2)
	table.SetCell(0, 0, tview.NewTableCell("Name"))
	table.SetCell(0, 1, tview.NewTableCell("Address"))
	table.SetCell(0, 2, tview.NewTableCell("Tax ID"))
	if len(tui.Config.Receivers) > 0 {
		for i, r := range tui.Config.Receivers {
			table.SetCell(i+1, 0, tview.NewTableCell(r.Name))
			table.SetCell(i+1, 1, tview.NewTableCell(r.Address))
			table.SetCell(i+1, 2, tview.NewTableCell(r.TaxID))
		}
	}

	if selectedRow == 0 {
		selectedRow = 1
	}
	table.Select(selectedRow, 0)

	table.SetSelectedFunc(func(row int, column int) {
		if row > 0 {
			// we try to ignore table heading row
			tui.ActivePage = types.PageReceiverEdit
			receiver_edit.Render(tui, row-1)
		}
	})

	table.SetSelectionChangedFunc(func(row int, column int) {
		selectedRow = row
	})

	grid := tview.NewFrame(table).
		AddText("ESC/h: go back, l/enter: edit selected receiver, a: add Receiver", false, tview.AlignLeft, tcell.Color100)

	return grid
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if vimkeys.Back(eventKey) {
		tui.SwitchToPage(types.PageDefault)
	}

	if vimkeys.Down(eventKey) {
		return tcell.NewEventKey(tcell.KeyDown, tcell.RuneDArrow, tcell.ModNone)
	}
	if vimkeys.Up(eventKey) {
		return tcell.NewEventKey(tcell.KeyUp, tcell.RuneUArrow, tcell.ModNone)
	}

	if eventKey.Rune() == 'd' && selectedRow > 0 {
		// we try to ignore table heading row there
		tui.Config.Receivers = config.RemoveCompany(tui.Config.Receivers, selectedRow-1)
		if err := tui.Config.WriteConfig(); err != nil {
			modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui) })
		}
		tui.SwitchToPage(types.PageReceiverList)
	}

	if eventKey.Rune() == 'a' {
		receiver_add.Render(tui)
		return nil
	}

	return eventKey
}
