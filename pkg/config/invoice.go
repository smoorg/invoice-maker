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

const (
	FieldInvoiceNo    = "Invoice No."
	FieldInvoiceDate  = "Invoice Date"
	FieldDeliveryDate = "Delivery Date"
	FieldDueDate      = "Payout Due Date"
	FieldPaymentType  = "Payment Type"
	FieldUnit         = "Unit of Measure"
	FieldTitle        = "Title"
	FieldPrice        = "Price / unit"
	FieldReceiver     = "Receiver"
	FieldQuantity     = "Quantity"
	FieldVatRate      = "Vat Rate"
	FieldAmount       = "Amount"
	FieldTotal        = "Total"
)

type InvoiceItem struct {
	// Name of product or service on the invoice
	Title string `yaml:"title"`

	// Price per good.
	Price string `yaml:"price"`

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
	Amount string `yaml:"amount"`
}

func (item *Invoice) CalculateInvoice() {
	decimal.DivisionPrecision = 2
	for i := range item.Items {
		item.Items[i].CalculateItemTotal()
	}
}

func (i *Invoice) AddNewItem() int {
	i.Items = append(i.Items, InvoiceItem{})

	return len(i.Items) - 1
}

func (i *InvoiceItem) CalculateItemTotal() decimal.Decimal {
	decimal.DivisionPrecision = 2
	amount := i.CalculateAmount()
	vat := i.CalculateVatAmount()

	return amount.Add(vat)
}

func (i *InvoiceItem) CalculateAmount() decimal.Decimal {
	decimal.DivisionPrecision = 2
	price, err := decimal.NewFromString(i.Price)
	if err != nil {
		panic(err)
	}
	result := price.Mul(decimal.NewFromInt32(i.Quantity))
	i.Amount = result.StringFixed(2)

	return result
}

func (i *InvoiceItem) CalculateVatAmount() decimal.Decimal {
	decimal.DivisionPrecision = 2
	if i.VatRate > 0 {
		amount := i.CalculateAmount()
		return amount.Mul(decimal.NewFromInt32(i.VatRate).Div(decimal.NewFromInt32(100)))
	}
	return decimal.NewFromInt32(0)
}

func (c *Config) AddInvoice(i Invoice) {
	c.Invoices = append(c.Invoices, i)
}
func (c *Config) AddInvoiceItem(index int, i InvoiceItem) int {
	c.Invoices[index].Items = append(c.Invoices[index].Items, i)

	return len(c.Invoices[index].Items) - 1
}

func (c *Config) UpdateInvoice(index int, invoice Invoice) {
	c.Invoices[index] = invoice
}

func (c *Config) GetInvoice(index int) *Invoice {
	return &c.Invoices[index]
}

func (i *Invoice) DeleteInvoiceItem(itemIndex int) {
	i.Items = append(i.Items[:itemIndex], i.Items[(itemIndex+1):]...)
}

func (i *Invoice) NetSum() string {
	amount := decimal.NewFromInt(0)
	for _, v := range i.Items {
		amount = amount.Add(v.CalculateAmount())
	}

	return amount.StringFixed(2)
}

func (i *Invoice) GrossSum() string {
	amount := decimal.NewFromInt(0)
	for _, v := range i.Items {
		amount = amount.Add(v.CalculateItemTotal())
	}

	return amount.StringFixed(2)
}

func (c *Config) RemoveInvoice(index int) {
	c.Invoices = append(c.Invoices[:index], c.Invoices[index+1:]...)
}
