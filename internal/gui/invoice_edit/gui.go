package invoice_edit

import (
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/invoice_add"
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI, row int) {
	tui.AddAndSwitchToPage(
		types.PageInvoiceEdit,
		editInvoice(tui, row),
	)
}

func updateInvoice(tui *types.TUI, row int, data *config.Invoice) {
	tui.Config.UpdateInvoice(row, *data)
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
	return invoice_add.AddOrEditInvoice(
		tui,
		tui.Config.GetInvoice(row),
		func(data *config.Invoice) { updateInvoice(tui, row, data) },
		func() { goBack(tui) },
	)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyESC {
		goBack(tui)
	}
	return eventKey
}
