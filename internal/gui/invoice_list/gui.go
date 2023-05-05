package invoice_list

import (
	"invoice-maker/internal/gui/invoice_add"
	"invoice-maker/internal/gui/invoice_edit"
	"invoice-maker/internal/template"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.Pages.AddPage(types.PageInvoiceList, invoiceList(tui), true, false)
}

var selectedInvoice = 1

func invoiceList(tui *types.TUI) tview.Primitive {
	table := tview.NewTable().SetSelectable(true, false).SetBorders(true)

	table.SetCell(0, 0, tview.NewTableCell("Invoice No.")).
		SetCell(0, 1, tview.NewTableCell("InvoiceDate")).
		SetCell(0, 2, tview.NewTableCell("DueDate")).
		SetCell(0, 3, tview.NewTableCell("Receiver")).
		SetCell(0, 4, tview.NewTableCell("Issuer")).
		SetCell(0, 5, tview.NewTableCell("PaymentType"))

	for i, invoice := range tui.Config.Invoices {
		table.
			SetCell(i+1, 0, tview.NewTableCell(invoice.InvoiceNo)).
			SetCell(i+1, 1, tview.NewTableCell(invoice.InvoiceDate)).
			SetCell(i+1, 2, tview.NewTableCell(invoice.DueDate)).
			SetCell(i+1, 3, tview.NewTableCell(invoice.Receiver.Name)).
			SetCell(i+1, 4, tview.NewTableCell(invoice.Issuer.Name)).
			SetCell(i+1, 5, tview.NewTableCell(invoice.PaymentType))
	}

	table.Select(selectedInvoice, 0)

	table.SetSelectedFunc(func(row int, column int) {
		if row > 0 {
			// we try to ignore table heading row
			tui.ActivePage = types.PageInvoiceEdit
			invoice_edit.Render(tui, row-1)
		}
	})

	table.SetSelectionChangedFunc(func(row int, column int) {
		selectedInvoice = row
	})

	view := tview.NewFrame(table).
		AddText("ESC, h: go back    l, e, enter: edit    a: add invoice    p: print",
			false,
			tview.AlignLeft,
			tcell.Color100,
		)

	return view
}

func printInvoice(tui *types.TUI) {
	i := &tui.Config.Invoices[selectedInvoice-1]
	dir, err := tui.Config.GetInvoiceDirectory()

	tmpl, err := template.GetTemplate()
	if err != nil {
		tui.Fatal("missing template")
	}

	rowTemplate, err := template.GetRowTemplate()
	if err != nil {
		tui.Fatal("unable to locate row template")
	}

	inv := template.ApplyInvoice(string(tmpl), string(rowTemplate), i)

	_, err = template.SaveInvoice(inv, dir)
	if err != nil {
		tui.Fatal("missing row template", err)
	}
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
		return tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	}

	if eventKey.Rune() == 'a' {
		tui.ActivePage = types.PageInvoiceAdd
		invoice_add.Render(tui)
		return nil
	}

	if eventKey.Rune() == 'p' {
		printInvoice(tui)
	}

	return nil
}
