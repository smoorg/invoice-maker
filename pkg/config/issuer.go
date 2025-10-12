package config

type Issuer struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	TaxID    string `yaml:"taxID"`
	Account  string `yaml:"accountNo"`
	BankName string `yaml:"bankName"`
	BIC      string `yaml:"bic"`
}
