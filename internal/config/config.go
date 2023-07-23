package config

type FontCfg struct {
	Family   string `yaml:"family"`
	Style    string `yaml:"style"`
	Filepath string `yaml:"filepath"`
}

type Config struct {
	Issuer           Issuer    `yaml:"issuer"`
	Receivers        []Company `yaml:"receivers"`
	Invoices         []Invoice `yaml:"invoices"`
	InvoiceDirectory string    `yaml:"invoiceDirectory"`
	Font             FontCfg   `yaml:"font"`
}
