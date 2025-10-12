package receiver_list

import (
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/gui/receiver_add"
	"invoice-maker/internal/gui/receiver_edit"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	"invoice-maker/pkg/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var selectedRow int = 0

func Render(tui *types.TUI) {
	tui.Pages.AddPage(types.PageReceiverList, listReceivers(tui), true, false)
}

func listReceivers(tui *types.TUI) tview.Primitive {
	table := tview.NewTable().SetSelectable(true, false).SetBorders(true)
	tui.SetDefaultStyle(table.Box)

	table.SetFixed(0, 2)
	table.SetCellSimple(0, 0, "Name")
	table.SetCellSimple(0, 1, "Address")
	table.SetCellSimple(0, 2, "Tax ID")
	if len(tui.Config.Receivers) > 0 {
		for i, r := range tui.Config.Receivers {
			table.SetCellSimple(i+1, 0, r.Name)
			table.SetCellSimple(i+1, 1, r.Address)
			table.SetCellSimple(i+1, 2, r.TaxID)
		}
	}

	if selectedRow == 0 {
		selectedRow = 1
	}
	table.Select(selectedRow, 0)

	table.SetSelectedFunc(func(row int, column int) {
		if row > 0 {
			// we try to ignore table heading row
			selectedRow = row - 1
		}
	})

	table.SetSelectionChangedFunc(func(row int, column int) {
		selectedRow = row - 1
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

	if vimkeys.Forward(eventKey) {
		receiver_edit.Render(tui, selectedRow)
		return nil
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
