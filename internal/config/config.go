package config

type Config struct {
	Issuer           Issuer    `yaml:"issuer"`
	Receivers        []Company `yaml:"receivers"`
	Invoices         []Invoice `yaml:"invoices"`
	InvoiceDirectory string    `yaml:"invoiceDirectory"`
	Font             FontCfg   `yaml:"font"`
}
