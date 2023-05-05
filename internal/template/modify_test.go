package template_test

import (
	"strings"
	"testing"

	"invoice-maker/internal/config"
	"invoice-maker/internal/template"

	"github.com/shopspring/decimal"
)

func TestInsertItems(t *testing.T) {
	row := "[ Items                       ]\n"

	sut := template.InsertRows(row, "Items", "test")

	if !strings.Contains(sut, "test") {
		t.Error("InsertRows did not add value properly", sut)
	}
}

func TestTotal(t *testing.T) {
	templateStr := "[ Items                                                   ]\n| [ASum] | [TaxSum]| [TotSum]"
	rowStr := "[Qty][Price][Amount][VR][VA][Total]"

	invoice := &config.Invoice{}
	invoice.Items = append(invoice.Items, config.InvoiceItem{
		Title:    "Test",
		Price:    decimal.NewFromInt32(1000),
		VatRate:  23,
		Quantity: 2,
	})

	if len(invoice.Items) == 0 {
		t.Error("invoice items should be greater than 0")
	}

	sut := template.ApplyInvoice(templateStr, rowStr, invoice)

	if !strings.Contains(sut, "| 2000   | 460     | 2460") {
		t.Error("Invalid amount sum", sut)
	}
}
