package config

type Company struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	TaxID   string `yaml:"taxID"`
}

func RemoveCompany(s []Company, index int) []Company {
	return append(s[:index], s[index+1:]...)
}
