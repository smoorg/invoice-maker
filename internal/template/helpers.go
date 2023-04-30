package template

import (
	"log"
	"os"
	"path/filepath"

	"invoice-maker/internal/config"
)

func GetInvoiceTemplate() ([]byte, error) {
	cfgDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}
	invTemplate := filepath.Join(cfgDir, "template.md")

	templateFile, err := os.ReadFile(invTemplate)
	if err != nil {
		file, fileCreateErr := os.Create(invTemplate)
		defer file.Close()
		if fileCreateErr != nil {
			log.Fatalln("unable to create invoice template markdown file at", invTemplate)
		}

		_, err = file.Read(templateFile)
	}

	//TODO: fallback to /etc/invoice-maker/template.md as should always exist

	return templateFile, nil
}
