package invoice_add

import (
	"fmt"
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/template"
	"invoice-maker/internal/types"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shopspring/decimal"
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
	tui.SetDefaultStyle(form.Box)

	receivers := []string{}
	paymentTypes := []string{"Cash", "Transfer"}
	for _, r := range tui.Config.Receivers {
		receivers = append(receivers, fmt.Sprint(r.Name))
	}

	data.Issuer = tui.Config.Issuer

	pickedReceiver := updateReceiver(tui, data)

	form.
		AddInputField("Invoice No.", data.InvoiceNo, 50, nil, func(text string) {
			data.InvoiceNo = text
			formChanged()
		}).
		AddInputField("Invoice Date", data.InvoiceDate, 50, nil,
			func(text string) {
				data.InvoiceDate = text
				formChanged()
			}).
		AddInputField("Delivery Date", data.DeliveryDate, 50, nil,
			func(text string) {
				data.DeliveryDate = text
				formChanged()
			}).
		AddInputField("Due Date", data.DueDate, 50, nil,
			func(text string) {
				data.DueDate = text
				formChanged()
			}).
		AddDropDown("Receiver", receivers, pickedReceiver,
			func(option string, optionIndex int) {
				if optionIndex >= 0 {
					data.Receiver = tui.Config.Receivers[optionIndex]
				}
				formChanged()
			}).
		AddDropDown("PaymentType", paymentTypes, 0, func(option string, optionIndex int) {
			data.PaymentType = option
			formChanged()
		})

	//TODO: make it possible to have multiple issuer companies to pick

	//That part probably supposed to be moved to separated invoice items form in future, currently app supports just a single item
	if len(data.Items) == 0 {
		data.Items = append(data.Items, config.InvoiceItem{})
	}

	// invoice items start
	form.
		AddInputField("Unit", data.Items[0].Unit, 50, nil, func(text string) {
			data.Items[0].Unit = text
			formChanged()
		}).
		AddInputField("Price/unit", data.Items[0].Price.String(), 50, nil, func(text string) {
			decimal.DivisionPrecision = 2
			if p, err := decimal.NewFromString(text); err == nil {
				data.Items[0].Price = p
				formChanged()
			}
		}).
		AddInputField("Quantity", fmt.Sprint(data.Items[0].Quantity), 50, nil, func(text string) {
			if q, err := strconv.Atoi(text); err == nil {
				data.Items[0].Quantity = int32(q)
				formChanged()
			}
		}).
		AddInputField("Vat Rate", fmt.Sprint(data.Items[0].VatRate), 50, nil, func(text string) {
			if vr, err := strconv.ParseUint(text, 10, 32); err == nil {
				data.Items[0].VatRate = int32(vr)
				formChanged()
			}
		}).
		AddInputField("Title", fmt.Sprint(data.Items[0].Title), 50, nil, func(text string) {
			data.Items[0].Title = text
			formChanged()
		})

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
		preview.SetText("something went wrong")
	} else {
		preview.SetText(txt)
	}

	return preview
}

func AddOrEditInvoice(tui *types.TUI, data *config.Invoice, save func(data *config.Invoice), cancel func()) tview.Primitive {
	preview := renderPreview(tui, data)
	changed := func() {
		defer func() {
			txt, err := template.GetContent(data)
			if err == nil {
				preview.Clear()
				fmt.Fprint(preview, txt)
			}
		}()
	}
	form := createForm(tui, data, changed)

	form.AddButton("Save", func() { save(data) }).
		AddButton("Cancel", cancel).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1)

	g := tview.NewFlex().
		AddItem(form, 70, 1, true).
		AddItem(preview, 0, 1, false)

	return g
}

func insertInvoice(tui *types.TUI, data *config.Invoice) {
	tui.Config.AddInvoice(*data)
	if err := tui.Config.WriteConfig(); err != nil {
		modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui) })
	}
	goBack(tui)
}

func goBack(tui *types.TUI) {
	tui.SwitchToPage(types.PageInvoiceList)
}

func addInvoice(tui *types.TUI) tview.Primitive {
	data := &config.Invoice{}
	return AddOrEditInvoice(tui, data,
		func(data *config.Invoice) { insertInvoice(tui, data) },
		func() { goBack(tui) })
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		goBack(tui)
	}

	return eventKey
}
