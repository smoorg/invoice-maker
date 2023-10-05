package help

import (
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {

	t := helpPage(tui)
	tui.AddAndSwitchToPage(types.PageHelp, t)
}

func helpPage(tui *types.TUI) tview.Primitive {
	t := tview.NewTable()
	tui.SetDefaultStyle(t.Box)

	t.SetCellSimple(1, 0, "?")
	t.SetCellSimple(1, 1, " - opens this help dialog")

	t.SetCellSimple(2, 0, "←↑→↓")
	t.SetCellSimple(3, 0, "hjkl")

	t.SetCellSimple(3, 1, " - navigate on lists and tables")

	t.SetCellSimple(4, 0, "Tab")
	t.SetCellSimple(4, 1, " - navigate forward in forms")

	t.SetCellSimple(4, 0, "Shift + Tab")
	t.SetCellSimple(4, 1, " - navigate backwards in forms")

	t.SetCellSimple(5, 0, "Esc")
	t.SetCellSimple(5, 1, " - cancel or go back")

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
		tui.SwitchToPage(tui.PreviousPage)
		return nil
	}

	return eventKey
}
