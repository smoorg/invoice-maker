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

func TestSumUp(t *testing.T) {
	items := []config.InvoiceItem{}
	items = append(items, config.InvoiceItem{
		Price:     decimal.NewFromInt32(100),
		VatAmount: decimal.NewFromInt(2300),
		Total:     decimal.NewFromInt(12300),
		Quantity:  100,
		VatRate:   23,
	})

	sum, tax, total := template.SumUp(&items)

	if sum != "10000" {
		t.Error("amount of 100 * 100 should be 10000 but is", sum)
	}

	if tax != "2300" {
		t.Error("amount of 10000 * 0.23 should be 2300 but is ", tax)
	}

	if total != "12300" {
		t.Error("total should sum up amount and tax but is", total)
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

	sut, err := template.ApplyInvoice(templateStr, rowStr, invoice)

	if err != nil {
		t.Error("should not return error", err)
	}

	if !strings.Contains(*sut, "| 2000   | 460     | 2460") {
		t.Error("Invalid amount sum", sut)
	}
}
