package template

import (
	"invoice-maker/config"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplyAddress(t *testing.T) {
	template := "| [ InvoiceAddress                     ] |"
	length := len(template)
	i := config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	sut := ApplyInvoice(template, "", i)

	if length != len(sut) {
		t.Error("length mismatch")
	}
}

func TestApplyInvoiceRows(t *testing.T) {
	template := "| [ InvoiceAddress                     ] |\n [ Items ]"
	i := config.Invoice{
		Items: []config.InvoiceItem{
			{
				Title:    "Cheese",
				Quantity: 2,
				Unit:     "kg",
				Price:    decimal.NewFromInt(20),
				VatRate:  15,
			},
			{
				Title:    "Cheese",
				Quantity: 2,
				Unit:     "kg",
				Price:    decimal.NewFromInt(20),
				VatRate:  15,
			},
		},
	}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	rowTemplate := "│ [ Title                ] │ [Qty] │ [Unit]  │ [Price] │ [Amount] │ [VR]  │ [VA   ] │ [Total]  │"

	sut := ApplyInvoice(template, rowTemplate, i)

	t.Error(sut)

	lines := strings.Split(sut, "\n")
	if len(lines) > 2 {
		t.Error("too many lines")
	}

}

func TestTotalCalculations(t *testing.T) {
	template := "| [ InvoiceAddress                     ] |\n[ Items                                     ]"
	rowTemplate := "| [ Total                     ] |"

	i := config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	item := config.InvoiceItem{}
	item.VatRate = 22
	item.Price = decimal.NewFromInt32(25000)
	item.Quantity = 1

	i.Items = append(i.Items, item)

	sut := ApplyInvoice(template, rowTemplate, i)
	total := "30500"

	if !strings.Contains(sut, total) {
		t.Error("total invalid", sut)
	}
}
