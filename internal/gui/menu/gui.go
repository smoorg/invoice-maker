package menu

import (
	"invoice-maker/internal/gui/config_page"
	"invoice-maker/internal/gui/help"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.Pages.AddPage(types.PageDefault, MenuPage(tui), true, true)
}

func MenuPage(tui *types.TUI) *tview.List {
	page := tview.NewList()

	page.AddItem("Invoices", "List, add and modify invoices", 'i',
		func() { tui.SwitchToPage(types.PageInvoiceList) })

	page.AddItem("Issuer", "Modify your company data, bank account etc.", 'm',
		func() { tui.SwitchToPage(types.PageIssuerEdit) })

	page.AddItem("Receivers", "List, add and modify your invoice obtainers", 'r',
		func() { tui.SwitchToPage(types.PageReceiverList) })

	page.AddItem("Config", "General configuration of the program", 'c',
		func() {
			config_page.Render(tui)
		})

	page.AddItem("Help", "Help page", 'c',
		func() {
			help.Render(tui)
		})

	tui.SetDefaultStyle(page.Box)

	return page
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Rune() == 'j' {
		return HandleEvents(tcell.NewEventKey(tcell.KeyDown, tcell.RuneDArrow, tcell.ModNone), tui)
	}
	if eventKey.Rune() == 'k' {
		return HandleEvents(tcell.NewEventKey(tcell.KeyUp, tcell.RuneUArrow, tcell.ModNone), tui)
	}
	if eventKey.Rune() == 'l' {
		return tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	}

	return eventKey
}
