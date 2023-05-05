package config_test

import (
	"fmt"
	"testing"

	"invoice-maker/internal/config"

	"github.com/shopspring/decimal"
)

func TestPriceCalculation(t *testing.T) {

	price := decimal.NewFromInt32(1000)
	result := config.CalculateAmount(price, 22)
	sut := result.String()

	if sut != "22000" {
		t.Error("invalid sum", sut)
	}
}

func TestVatCalculation(t *testing.T) {

	{
		price := decimal.NewFromInt32(1000)
		result := config.CalculateVatAmount(price, 22)
		sut := result.String()

		if sut != "220" {
			t.Error("invalid vat", sut)
		}
	}

	{
		price := decimal.NewFromInt32(1)
		result := config.CalculateVatAmount(price, 22)
		sut := result.String()

		if sut != "0.22" {
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
