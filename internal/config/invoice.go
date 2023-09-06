package config

import (
	"github.com/shopspring/decimal"
)

type Invoice struct {
	// Arbitrary number for identification purposes.
	InvoiceNo string `yaml:"invoiceNo"`

	// Date of raise of the invoice.
	InvoiceDate string `yaml:"issueDate"`

	// Date of delivery of the goods or service.
	DeliveryDate string `yaml:"deliveryDate"`

	// Due day of payment for the receiver.
	DueDate string `yaml:"dueDate"`

	// Receiver of the invoice.
	Receiver Company `yaml:"receiver"`

	// Issuer of the invoice, in short your company.
	Issuer Issuer `yaml:"issuer"`

	// The way receiver will pay.
	PaymentType string `yaml:"paymentType"`

	// Items on the invoice.
	Items []InvoiceItem `yaml:"items"`
}

type InvoiceItem struct {
	// Name of product or service on the invoice
	Title string `yaml:"title"`

	// Price per good.
	Price decimal.Decimal `yaml:"price"`

	// Number of units sold. In case of service that's usually set to 1.
	Quantity int32 `yaml:"quantity"`

	// Unit of measure of quantity, say a crate, gallon, kilos and so on.
	// Not all that important for services where usually just "unit" applies.
	Unit string `yaml:"unit"`

	// Vat rate for the goods, usually 1 or 2 digits treated as percentage later on.
	// Distinct per item due to complicated vat law in Poland where multiple items
	// can be treated on different vat rates.
	VatRate int32 `yaml:"vatRate"`

	// Net amount from which taxes are deducted. Calculated based on Price * Qty.
	Amount decimal.Decimal `yaml:"amount"`

	// Vat amount to be added to `Total`. Calculated based on `VatRate * Price / 100`.
	VatAmount decimal.Decimal

	// Total calculated by Amount + VatAmount
	Total decimal.Decimal
}

func (item *Invoice) CalculateInvoice() {
	for i := range item.Items {
		item.Items[i].CalculateItemTotal()
	}
}

func (i *InvoiceItem) CalculateItemTotal() {
	i.CalculateAmount()
	i.CalculateVatAmount()

	i.Total = i.Amount.Add(i.VatAmount)
}

func (i *InvoiceItem) CalculateAmount() {
	i.Amount = i.Price.Mul(decimal.NewFromInt32(i.Quantity))
}

func (i *InvoiceItem) CalculateVatAmount() {
	i.CalculateAmount()
	if i.VatRate > 0 {
		i.VatAmount = i.Amount.Mul(decimal.NewFromInt32(i.VatRate).Div(decimal.NewFromInt32(100)))
	}
}
