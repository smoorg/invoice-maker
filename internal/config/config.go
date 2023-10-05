package config


type Config struct {
	Issuer           Issuer    `yaml:"issuer"`
	Receivers        []Company `yaml:"receivers"`
	Invoices         []Invoice `yaml:"invoices"`
	InvoiceDirectory string    `yaml:"invoiceDirectory"`
	Font             FontCfg   `yaml:"font"`
}


func (c *Config) AddInvoice(i Invoice) {
	c.Invoices = append(c.Invoices, i)
}

func (c *Config) UpdateInvoice(index int, invoice Invoice) {
	c.Invoices[index] = invoice
}

func (c *Config) GetInvoice(index int) *Invoice {
	return &c.Invoices[index]
}
