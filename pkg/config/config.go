package config

type Config struct {
	Issuer     Issuer    `yaml:"issuer"`
	Receivers  []Company `yaml:"receivers"`
	Invoices   []Invoice `yaml:"invoices"`
	Config     FontCfg   `yaml:"font"`
	DateFormat string    `yaml:"date_format"`
}
