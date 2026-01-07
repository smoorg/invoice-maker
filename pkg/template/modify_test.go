package template_test

import (
	"strings"
	"testing"

	"invoice-maker/pkg/config"
	"invoice-maker/pkg/template"

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
		Price:     "100",
		Quantity:  100,
		VatRate:   23,
	})

	sum, tax, total := template.SumUp(&items)

	expectedSum := decimal.NewFromInt(10000)
	if !expectedSum.Equal(sum) {
		t.Error("amount of 100 * 100 should be 10000 but is", sum, expectedSum.Cmp(sum))
	}

	expectedTax := decimal.NewFromInt(2300)
	if !expectedTax.Equal(tax) {
		t.Error("amount of 10000 * 0.23 should be 2300 but is ", tax, expectedTax.Cmp(tax))
	}

	expectedTotal := decimal.NewFromInt(12300)
	if !expectedTotal.Equal(total) {
		t.Error("total should sum up amount and tax but is", total)
	}
}

func TestTotal(t *testing.T) {
	sut := "[ Items                                                   ]\n|[ ASum ]| [TaxSum]| [TotSum]"
	rowStr := "[Qty]|[Price]|[Amount]|[VR]|[  VA  ]|[Total]"

	invoice := &config.Invoice{}
	invoice.Items = append(invoice.Items, config.InvoiceItem{
		Title:    "Test",
		Price:    "1000",
		VatRate:  23,
		Quantity: 2,
	})

	if len(invoice.Items) == 0 {
		t.Error("invoice items should be greater than 0")
	}

	err := template.ApplyInvoice(&sut, rowStr, invoice)

	if err != nil {
		t.Error("should not return error", err)
	}

	expected := "2    |1000.00|2000.00 |23% |460.00  |2460.00"
	if !strings.Contains(sut, expected) {
	    t.Error("Invalid amount sum:", sut, "compared to expected:", expected)
	}
}
