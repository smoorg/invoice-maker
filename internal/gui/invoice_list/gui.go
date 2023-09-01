package invoice_list

import (
	"invoice-maker/internal/gui/invoice_add"
	"invoice-maker/internal/gui/invoice_edit"
	"invoice-maker/internal/gui/invoice_print_modal"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"
	"log"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.Pages.AddPage(types.PageInvoiceList, invoiceList(tui), true, false)
}

var selectedInvoice = 1

func invoiceList(tui *types.TUI) tview.Primitive {
	table := tview.NewTable().SetSelectable(true, false).SetBorders(true)

	table.
		SetCellSimple(0, 0, "Invoice No.").
		SetCellSimple(0, 1, "InvoiceDate").
		SetCellSimple(0, 2, "DueDate").
		SetCellSimple(0, 3, "Receiver").
		SetCellSimple(0, 4, "Issuer").
		SetCellSimple(0, 5, "PaymentType")

	for i, invoice := range tui.Config.Invoices {
		table.
			SetCellSimple(i+1, 0, invoice.InvoiceNo).
			SetCellSimple(i+1, 1, invoice.InvoiceDate).
			SetCellSimple(i+1, 2, invoice.DueDate).
			SetCellSimple(i+1, 3, invoice.Receiver.Name).
			SetCellSimple(i+1, 4, invoice.Issuer.Name).
			SetCellSimple(i+1, 5, invoice.PaymentType)
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
	if vimkeys.Forward(eventKey) || eventKey.Rune() == 'e' {
		return tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	}

	if eventKey.Rune() == 'a' {
		invoice_add.Render(tui)
		return nil
	}

	if eventKey.Rune() == 'p' {
		path, err := printInvoice(tui)
		if err != nil {
			log.Fatal("print invoide failed", err)
			return nil
		}
		if path == "" {
			log.Fatal("missing pdf path")
			return nil
		}

		invoiceprintmodal.Render(tui, path)
	}

	return nil
}
