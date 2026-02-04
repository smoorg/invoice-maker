package config_test

import (
	"fmt"
	"testing"

	"invoice-maker/pkg/config"

	"github.com/shopspring/decimal"
)

func TestPriceCalculation(t *testing.T) {
	i := config.InvoiceItem{
		Price:    "1000",
		Quantity: 23,
	}

	i.CalculateAmount()
	sut := i.Amount

	if sut != "23000.00" {
		t.Error("invalid sum", sut)
	}
}

func TestVatCalculation(t *testing.T) {

	{

		i := config.InvoiceItem{
			Price:    "1000",
			Quantity: 1,
			VatRate:  23,
		}
		i.CalculateVatAmount()

		sut := i.CalculateVatAmount().StringFixed(2)
		if sut != "230.00" {
			t.Error("invalid vat", sut)
		}
	}

	{
		i := config.InvoiceItem{
			Price:    "1",
			Quantity: 1,
			VatRate:  23,
		}

		i.CalculateVatAmount()
		sut := i.CalculateVatAmount().StringFixed(2)

		if sut != "0.23" {
			t.Error("invalid vat", sut)
		}
	}
}

func TestCalculateTotal(t *testing.T) {
	inv := &config.Invoice{}
	inv.Items = append(inv.Items, config.InvoiceItem{
		Price:    "10",
		Quantity: 2,
	})
	inv.CalculateInvoice()

	if len(inv.Items) != 1 {
		t.Error("should be single Item")
	}

	amount, err := decimal.NewFromString(inv.Items[0].Amount)
	if err != nil {
		t.Error(err)
	}
	if !amount.Equal(decimal.NewFromInt32(20)) {
		fmt.Print(inv.Items)
		t.Error(inv.Items[0].Amount, "should be 20")
	}
}

func TestAddInvoiceItem(t *testing.T) {
	inv := &config.Invoice{}
	inv.AddNewItem()

	if len(inv.Items) != 1 {
		t.Fatal("item was not added properly")
	}
}

func TestDeleteInvoiceItem(t *testing.T) {
	inv := &config.Invoice{
		Items: []config.InvoiceItem{
			{Title: "abc"},
		},
	}

	if len(inv.Items) != 1 {
		t.Fatal("item was not added properly")
	}

	inv.DeleteInvoiceItem(0)
	if len(inv.Items) != 0 {
		t.Fatal("item was not deleted properly")
	}
}
