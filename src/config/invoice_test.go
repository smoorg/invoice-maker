package config_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"invoice-maker/config"
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
