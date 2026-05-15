package invoices

import (
	"errors"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/template"
	"log"
	"os"
	"path/filepath"
)

func getInvoice(invoices []config.Invoice, invoiceNo string, netSum string) (*config.Invoice, int, error) {
	for i, v := range invoices {
		if v.InvoiceNo == invoiceNo && v.NetSum() == netSum {
			return &v, i, nil
		}
	}

	return nil, 0, errors.New("no such invoice")
}
func getInvoiceContent(invoices []config.Invoice, invoiceNo string, netSum string) (string, error) {
	invoice, _, err := getInvoice(invoices, invoiceNo, netSum)
	if err != nil {
		return "", err
	}

	content, err := template.GetContent(invoice)
	if err != nil {
		return "", err
	}

	return content, nil
}

func saveFile(dirname string, filename string, content []byte) error {
	if err := os.MkdirAll(dirname, 0744); err != nil {
		return err
	}

	mddir := filepath.Join(dirname, filename)

	file, err := os.Create(mddir)
	if err != nil {
		return err
	}
	if _, err := file.Write(content); err != nil {
		log.Fatal("write string err", err)
		return err
	}
	return nil
}
