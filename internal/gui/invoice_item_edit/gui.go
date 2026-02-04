package invoice_item_edit

import (
	"fmt"
	"strconv"

	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"invoice-maker/pkg/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
)

func Render(tui *types.TUI, invoice *config.Invoice, invoiceIndex *int, itemIndex *int, postsave func()) {
	tui.AddAndSwitchToPage(types.PageInvoiceItemEdit,
		modal.Modal(tui, types.PageInvoiceItemEdit, types.PageInvoiceItemList,
			renderItem(tui, invoice, invoiceIndex, itemIndex, postsave),
			80, 16, "Edit Invoice Item"),
	)
}

func renderItem(tui *types.TUI, invoice *config.Invoice, invoiceIndex *int, invoiceItemIndex *int, postsave func()) tview.Primitive {
	var item config.InvoiceItem
	if invoiceItemIndex == nil {
		index := invoice.AddNewItem()
		item = invoice.Items[index]
		invoiceItemIndex = &index
	} else {
		item = invoice.Items[*invoiceItemIndex]
	}

	i := tview.NewForm().
		AddInputField(config.FieldTitle, item.Title, 20, nil, func(text string) {
			item.Title = text
		}).
		AddInputField(config.FieldUnit, item.Unit, 10, nil, func(text string) {
			item.Unit = text
		}).
		AddInputField(config.FieldPrice, item.Price, 5, nil, func(text string) {
			if val, err := decimal.NewFromString(text); err == nil {
				item.Price = val.StringFixed(2)
			}
		}).
		AddInputField(config.FieldQuantity, fmt.Sprint(item.Quantity), 5, nil, func(text string) {
			if val, err := strconv.ParseInt(text, 0, 32); err == nil {
				item.Quantity = int32(val)
			}
		}).
		AddInputField(config.FieldVatRate, fmt.Sprint(item.VatRate), 3, nil, func(text string) {
			if val, err := strconv.ParseInt(text, 0, 32); err == nil {
				item.VatRate = int32(val)
			}
		}).
		AddButton("Save", func() {
			//TODO: check error
			if invoiceIndex == nil {
				panic("invoiceIndex nil")
			}
			tui.Config.WriteInvoiceItem(item, *invoiceIndex, *invoiceItemIndex)
			tui.Pages.RemovePage(types.PageInvoiceItemEdit)
			tui.SetActivePage(types.PageInvoiceItemList)
			postsave()
		}).
		AddButton("Cancel", func() {
			tui.Pages.RemovePage(types.PageInvoiceItemEdit)
			tui.SetActivePage(types.PageInvoiceItemList)
		})

	return i
}
func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.Pages.RemovePage(types.PageInvoiceItemEdit)
	}

	return eventKey
}
