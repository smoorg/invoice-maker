package template

import (
	"errors"
	"os"
	"path/filepath"

	"invoice-maker/internal/config"
)

type String string

func GetTemplate() ([]byte, error) {
	return getTemplate("template")
}

func GetRowTemplate() ([]byte, error) {
	return getTemplate("row-template")
}

func getTemplate(name string) ([]byte, error) {
	dir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	if template, err := os.ReadFile(filepath.Join(dir, name+".md")); err == nil {
		return template, nil
	}
	if template, err := os.ReadFile(filepath.Join(dir, name+".html")); err == nil {
		return template, nil
	}

	return nil, errors.New("no template file found")
}


func GetContent(i *config.Invoice) (string, error) {
	tmpl, err := GetTemplate()
	if err != nil {
		return "", errors.New("missing template")
	}

	rowTemplate, err := GetRowTemplate()
	if err != nil {
		return "", errors.New("unable to locate row template")
	}

	inv, err := ApplyInvoice(string(tmpl), string(rowTemplate), i)
	if err != nil {
		return "", err
	}
	return *inv, err
}
