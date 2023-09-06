package config_test

import (
	"fmt"
	"testing"

	"invoice-maker/internal/config"

	"github.com/shopspring/decimal"
)

func TestPriceCalculation(t *testing.T) {
	i := config.InvoiceItem{
		Price:    decimal.NewFromInt32(1000),
		Quantity: 23,
	}

	i.CalculateAmount()
	sut := i.Amount.String()

	if sut != "23000" {
		t.Error("invalid sum", sut)
	}
}

func TestVatCalculation(t *testing.T) {

	{

		i := config.InvoiceItem{
			Price:    decimal.NewFromInt32(1000),
			Quantity: 1,
			VatRate:  23,
		}
		i.CalculateVatAmount()

		sut := i.VatAmount.String()
		if sut != "230" {
			t.Error("invalid vat", sut)
		}
	}

	{
		i := config.InvoiceItem{
			Price:    decimal.NewFromInt32(1),
			Quantity: 1,
			VatRate:  23,
		}

		i.CalculateVatAmount()
		sut := i.VatAmount.String()

		if sut != "0.23" {
			t.Error("invalid vat", sut)
		}
	}
}

func TestCalculateTotal(t *testing.T) {
	inv := &config.Invoice{}
	inv.Items = append(inv.Items, config.InvoiceItem{
		Price:    decimal.NewFromInt32(10),
		Quantity: 2,
	})
	inv.CalculateInvoice()

	if len(inv.Items) != 1 {
		t.Error("should be single Item")
	}

	if !inv.Items[0].Amount.Equal(decimal.NewFromInt32(20)) {
		fmt.Print(inv.Items)
		t.Error(inv.Items[0].Amount, "should be 20")
	}
}
