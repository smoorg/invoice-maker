package invoice_item_list

import (
	"fmt"
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/invoice_item_edit"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	//"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RenderItemTable(tui *types.TUI, data *config.Invoice, invoiceRow *int) {
	tui.AddAndSwitchToPage(
		types.PageInvoiceItemList,
		invoiceItemList(tui, data, invoiceRow),
	)
}

var selectedInvoiceItem = 0
var invoiceIndex = 0

func invoiceItemList(tui *types.TUI, invoice *config.Invoice, invoiceRow *int) tview.Primitive {
	if invoiceRow == nil {
		lastItem := len(tui.Config.Invoices) - 1
		invoiceRow = &lastItem
	}
	invoiceIndex = *invoiceRow
	table := tview.NewTable().SetSelectable(true, false).SetBorders(true)
	tui.SetDefaultStyle(table.Box)

	table.
		SetCellSimple(0, 0, config.FieldTitle).
		SetCellSimple(0, 1, config.FieldUnit).
		SetCellSimple(0, 2, config.FieldPrice).
		SetCellSimple(0, 3, config.FieldQuantity).
		SetCellSimple(0, 4, config.FieldVatRate)

	if invoice != nil && invoice.Items != nil {
		for i, invoice := range invoice.Items {
			index := i + 1
			table.
				SetCellSimple(index, 0, invoice.Title).
				SetCellSimple(index, 1, invoice.Unit).
				SetCellSimple(index, 2, invoice.Price.String()).
				SetCellSimple(index, 3, fmt.Sprint(invoice.Quantity)).
				SetCellSimple(index, 4, fmt.Sprint(invoice.VatRate))
		}
	}

	table.Select(selectedInvoiceItem, 0)

	table.SetSelectedFunc(func(itemRow int, column int) {
		if itemRow > 0 {
			selectItem(tui, invoice, invoiceRow, itemRow-1)
		}
	})

	table.SetSelectionChangedFunc(func(row int, column int) {
		selectedInvoiceItem = row
	})

	view := tview.NewFrame(table).
		AddText("ESC, h: go back    l, e, enter: edit invoice item    a: add invoice item",
			false,
			tview.AlignLeft,
			tcell.Color100,
		)

	return view
}

func selectItem(tui *types.TUI, invoice *config.Invoice, invoiceRow *int, itemRow int) {
	if invoiceRow == nil {
		panic("invoiceRow nil")
	}
	tui.RefreshConfig()
	tui.Pages.RemovePage(types.PageInvoiceItemList)
	invoice_item_edit.Render(tui, invoice, invoiceRow, &itemRow, func() {
		tui.Config.WriteConfig()
		tui.RefreshConfig()
		tui.Rerender()
		RenderItemTable(tui, tui.Config.GetInvoice(*invoiceRow), invoiceRow)
	})
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if vimkeys.Back(eventKey) {
		tui.Pages.RemovePage(types.PageInvoiceItemList)
		tui.SwitchToPage(types.PageInvoiceList)
	}

	if vimkeys.Forward(eventKey) || eventKey.Rune() == 'e' {
		return tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	}

	if eventKey.Rune() == 'a' {
		invoice := tui.Config.GetInvoice(invoiceIndex)
		i := invoice.AddNewItem()
		tui.Config.WriteConfig()
		selectItem(tui, tui.Config.GetInvoice(invoiceIndex), &selectedInvoiceItem, i)
		tui.Config.WriteConfig()
		//selectItem(tui, tui.Config.GetInvoice(invoiceIndex), selectedInvoiceItem, newItemIndex)
	}

	if eventKey.Rune() == 'd' {
		invoice := tui.Config.GetInvoice(invoiceIndex)
		invoice.DeleteInvoiceItem(selectedInvoiceItem - 1)
		tui.Config.UpdateInvoice(invoiceIndex, *invoice)
		tui.Config.WriteConfig()
		tui.RefreshConfig()

		tui.Pages.RemovePage(types.PageInvoiceItemList)

		invoice = tui.Config.GetInvoice(invoiceIndex)
		RenderItemTable(tui, invoice, &invoiceIndex)

	}

	return eventKey
}
