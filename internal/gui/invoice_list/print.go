package invoice_list

import (
	"invoice-maker/internal/gui/invoice_print_modal"
	"invoice-maker/internal/template"
	"invoice-maker/internal/types"
)

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

	if path, err := template.SaveInvoice(inv, dir); err != nil {
		tui.Fatal("missing row template", err)
	} else {
		invoiceprintmodal.Render(tui, path)
	}
}
