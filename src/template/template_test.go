package template

import (
	"invoice-maker/config"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplyInvoice(t *testing.T) {
	template := "| [ InvoiceAddress                     ] |"
	length := len(template)
	i := config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	sut := ApplyInvoice(template, "", i)

	if length != len(sut) {
		t.Error("length mismatch")
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
