package invoice_add

import (
	"fmt"
	"invoice-maker/internal/gui/invoice_item_list"
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"invoice-maker/pkg/config"
	"invoice-maker/pkg/template"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.AddAndSwitchToPage(
		types.PageInvoiceAdd,
		addInvoice(tui),
	)
}

// When data stores existing receiver it ensures we can find it from list and save updated data to the invoice
func updateReceiver(tui *types.TUI, data *config.Invoice) int {
	if data.Receiver.Name == "" {
		return -1
	}

	pickedReceiver := -1
	for i, c := range tui.Config.Receivers {
		if c.Name == data.Receiver.Name {
			pickedReceiver = i
			break
		}
	}
	// NOTE: we upgrade receiver data in case this is edit and we have outdated info about it in the invoice.
	// We do that only if there's such receiver on the list as it should be possible to remove receiver and edit old invoice.
	if pickedReceiver > -1 && len(tui.Config.Receivers) > pickedReceiver {
		data.Receiver = tui.Config.Receivers[pickedReceiver]
	}

	return pickedReceiver
}

func createForm(tui *types.TUI, data *config.Invoice, formChanged func()) *tview.Form {
	form := tview.NewForm()

	receivers := []string{}
	paymentTypes := []string{"Cash", "Transfer"}
	for _, r := range tui.Config.Receivers {
		receivers = append(receivers, fmt.Sprint(r.Name))
	}

	data.Issuer = tui.Config.Issuer

	pickedReceiver := updateReceiver(tui, data)

	form.
		AddInputField(config.FieldInvoiceNo, data.InvoiceNo, 50, nil, func(text string) {
			data.InvoiceNo = text
			formChanged()
		}).
		AddInputField(config.FieldInvoiceDate, data.InvoiceDate, 50, nil,
			func(text string) {
				data.InvoiceDate = text
				formChanged()
			}).
		AddInputField(config.FieldDeliveryDate, data.DeliveryDate, 50, nil,
			func(text string) {
				data.DeliveryDate = text
				formChanged()
			}).
		AddInputField(config.FieldDueDate, data.DueDate, 50, nil,
			func(text string) {
				data.DueDate = text
				formChanged()
			}).
		AddDropDown(config.FieldReceiver, receivers, pickedReceiver,
			func(option string, optionIndex int) {
				if optionIndex >= 0 {
					data.Receiver = tui.Config.Receivers[optionIndex]
				}
				formChanged()
			}).
		AddDropDown(config.FieldPaymentType, paymentTypes, 0, func(option string, optionIndex int) {
			data.PaymentType = option
			formChanged()
		})

	//TODO: make it possible to have multiple issuer companies to pick

	//That part probably supposed to be moved to separated invoice items form in future, currently app supports just a single item
	if len(data.Items) == 0 {
		data.Items = append(data.Items, config.InvoiceItem{})
	}

	// invoice items start
	//form.
	//	AddInputField(config.FieldUnit, data.Items[0].Unit, 50, nil, func(text string) {
	//		data.Items[0].Unit = text
	//		formChanged()
	//	}).
	//	AddInputField(config.FieldPrice, data.Items[0].Price.String(), 50, nil, func(text string) {
	//		decimal.DivisionPrecision = 2
	//		if p, err := decimal.NewFromString(text); err == nil {
	//			data.Items[0].Price = p
	//			formChanged()
	//		}
	//	}).
	//	AddInputField(config.FieldQuantity, fmt.Sprint(data.Items[0].Quantity), 50, nil, func(text string) {
	//		if q, err := strconv.Atoi(text); err == nil {
	//			data.Items[0].Quantity = int32(q)
	//			formChanged()
	//		}
	//	}).
	//	AddInputField(config.FieldVatRate, fmt.Sprint(data.Items[0].VatRate), 50, nil, func(text string) {
	//		if vr, err := strconv.ParseUint(text, 10, 32); err == nil {
	//			data.Items[0].VatRate = int32(vr)
	//			formChanged()
	//		}
	//	}).
	//	AddInputField(config.FieldTitle, fmt.Sprint(data.Items[0].Title), 50, nil, func(text string) {
	//		data.Items[0].Title = text
	//		formChanged()
	//	})

	form.SetTitle("Add Invoice")
	if data.Receiver.Name != "" || data.InvoiceNo != "" || data.DueDate != "" || data.InvoiceDate != "" ||
		data.Issuer.Name != "" {
		form.SetTitle("Edit Invoice")
	}

	return form
}

func renderPreview(tui *types.TUI, data *config.Invoice) *tview.TextView {
	preview := tview.NewTextView()
	txt, err := template.GetContent(data)
	if err != nil {
		preview.SetText(fmt.Sprintf("something went wrong: %s", err.Error()))
	} else {
		preview.SetText(txt)
	}

	return preview
}

func save(tui *types.TUI, row *int, data *config.Invoice) {
	if row != nil {
		// update
		tui.Config.UpdateInvoice(*row, *data)
	} else {
		//insert
		tui.Config.AddInvoice(*data)
	}

	if err := tui.Config.WriteConfig(); err != nil {
		modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", nil)
	}
}

func AddOrEditInvoice(tui *types.TUI, data *config.Invoice, row *int, cancel func()) tview.Primitive {
	// preview
	preview := renderPreview(tui, data)
	tui.SetDefaultStyle(preview)

	// replaces preview on form edit
	changed := func() {
		defer func() {
			txt, err := template.GetContent(data)
			if err == nil {
				preview.Clear()
				fmt.Fprint(preview, txt)
			}
		}()
	}

	////editable invoice items table
	//it := tview.NewTable().SetSelectable(true, false)
	//it.SetCellSimple(0, 0, "Id")
	//it.SetCellSimple(0, 1, "Title")
	//it.SetCellSimple(0, 2, "Total")

	//for i, v := range data.Items {
	//	index := i + 1
	//	it.SetCellSimple(index, 0, fmt.Sprint(index))
	//	it.SetCellSimple(index, 1, v.Title)
	//	it.SetCellSimple(index, 2, v.Total.String())
	//}

	// form
	form := createForm(tui, data, changed)
	form.
		AddButton("Save", func() {
			save(tui, row, data)
		}).
		AddButton("Cancel", cancel).
		AddButton("Next", func() {
			// HINT: save can happen on item level only
			save(tui, row, data)
			invoice_item_list.RenderItemTable(tui, data, row)
		}).
		SetBorderPadding(1, 1, 1, 1)

	// footer
	//footer := tview.NewTextView().SetText(
	//    ""
	//)

	// grid
	// ------
	// |form | preview
	// |-----|
	// |items|
	g := tview.NewGrid().SetRows(-30).SetColumns(-2, -4)

	g.
		AddItem(form, 0, 0, 1, 1, 0, 0, true).
		AddItem(preview, 0, 1, 1, 1, 0, 0, false) //.
		//AddItem(footer, 1, 0, 1, 2, 0, 0, false)

	//tui.SetDefaultStyle(footer)
	tui.SetDefaultStyle(form)

	return g
}

func goBack(tui *types.TUI) {
	tui.SwitchToPage(types.PageInvoiceList)
	tui.Pages.RemovePage(types.PageInvoiceAdd)
}

func addInvoice(tui *types.TUI) tview.Primitive {
	data := &config.Invoice{}
	return AddOrEditInvoice(tui, data, nil,
		func() { goBack(tui) })
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		goBack(tui)
	}

	return eventKey
}
