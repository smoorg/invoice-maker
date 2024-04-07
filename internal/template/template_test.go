package template_test

import (
	"invoice-maker/internal/config"
	"invoice-maker/internal/template"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplyAddress(t *testing.T) {
	sut := "| [ InvoiceAddress                     ] |"
	initialLength := len(sut)
	i := &config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	err := template.ApplyInvoice(&sut, "", i)

	if err != nil {
		t.Error("error thrown", err)
	}

	if initialLength != len(sut) {
		t.Error("length mismatch")
	}
}

func TestApplyInvoiceRows(t *testing.T) {
	sut := "| [ InvoiceAddress                     ] |\n [ Items ]\n"
	i := &config.Invoice{
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

	err := template.ApplyInvoice(&sut, rowTemplate, i)

	if err != nil {
		t.Error("error thrown", err)
	}

	lines := strings.Split(sut, "\n")
	if len(lines) > 2 {
		t.Error("too many lines")
	}

}

func TestTotalCalculations(t *testing.T) {
	sut := "| [ IssuerAddress                     ] |\n[ Items                                     ]\n"
	rowTemplate := "| [ Total                     ] |"

	i := &config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	item := config.InvoiceItem{}
	item.VatRate = 23
	item.Price = decimal.NewFromInt32(25000)
	item.Quantity = 1
	total := item.Price.Mul(decimal.NewFromInt32(item.Quantity)).Mul(decimal.NewFromFloat32(1.23))

	i.Items = append(i.Items, item)

	err := template.ApplyInvoice(&sut, rowTemplate, i)
	if err != nil {
		t.Error("error thrown", err)
	}

	if !strings.Contains(sut, total.String()) {
		t.Error("total invalid", sut, total)
	}
}
