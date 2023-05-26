package invoice_list

import (
	"errors"
	"invoice-maker/internal/template"
	"invoice-maker/internal/types"
)

func printInvoice(tui *types.TUI) (string, error) {
	i := &tui.Config.Invoices[selectedInvoice-1]
	dir, err := tui.Config.GetInvoiceDirectory()

	tmpl, err := template.GetTemplate()
	if err != nil {
		return "", errors.New("missing template")
	}

	rowTemplate, err := template.GetRowTemplate()
	if err != nil {
		return "", errors.New("unable to locate row template")
	}

	inv := template.ApplyInvoice(string(tmpl), string(rowTemplate), i)
	path, err := template.SaveInvoice(inv, dir)

	if err != nil {
		return "", err
	}
	return path, nil
}
