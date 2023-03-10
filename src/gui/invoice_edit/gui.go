package invoice_edit

import (
	"invoice-maker/config"
	"invoice-maker/gui/invoice_add"
	"invoice-maker/gui/modal"
	"invoice-maker/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI, row int) {
	tui.AddAndSwitchToPage(
		types.PageInvoiceEdit,
		modal.Modal(tui, types.PageInvoiceList, editInvoice(tui, row), 50, 20, ""),
	)
}

func saveInvoice(tui *types.TUI, row int, data *config.Invoice) {
	tui.Config.Invoices[row] = *data
	if err := tui.Config.WriteConfig(); err != nil {
		modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui, row) })
	}
	tui.Rerender()
	goBack(tui)
}

func goBack(tui *types.TUI) {
	tui.SwitchToPage(types.PageInvoiceList)
	tui.Pages.RemovePage(types.PageInvoiceAdd)
}

func editInvoice(tui *types.TUI, row int) tview.Primitive {
	invoice := tui.Config.Invoices[row]
	return invoice_add.AddOrEditInvoice(
		tui,
		&invoice,
		func(data *config.Invoice) { saveInvoice(tui, row, data) },
		func() { goBack(tui) },
	)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyESC {
		goBack(tui)
	}
	return eventKey
}
