package help

import (
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.AddAndSwitchToPage(
		types.PageHelp,
		modal.Modal(tui, tui.ActivePage, HelpPage(tui), 50, 11, "Add Receiver"),
	)
}

func HelpPage(tui *types.TUI) tview.Primitive {
	t := tview.NewTable()

	t.SetCell(1, 0, tview.NewTableCell("?"))
	t.SetCell(1, 1, tview.NewTableCell(" - opens this help dialog"))

	t.SetCell(2, 0, tview.NewTableCell("←↑→↓"))
	t.SetCell(3, 0, tview.NewTableCell("hjkl"))

	t.SetCell(3, 1, tview.NewTableCell(" - navigate on lists and tables"))

	t.SetCell(4, 0, tview.NewTableCell("Tab"))
	t.SetCell(4, 1, tview.NewTableCell(" - navigate form inputs"))

	t.SetCell(5, 0, tview.NewTableCell("Esc"))
	t.SetCell(5, 1, tview.NewTableCell(" - cancel or go back"))

	f := tview.NewGrid()
	f.SetTitle(" MANUAL ")
	f.SetBorder(true)
	v := tview.NewTextArea()
	v.SetText("General keyboard settings. Most of the view specific keyboard shortcuts lands on the bottom of the screen. For more info use `man im`.", false)
	f.AddItem(v, 0, 0, 1, 1, 1, 1, false)
	f.AddItem(t, 1, 0, 1, 1, 1, 1, true)
	return f
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.SwitchToPage(types.PageReceiverList)
		return nil
	}

	return eventKey
}
