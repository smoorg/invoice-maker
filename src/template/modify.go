package template

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"invoice-maker/config"
	"invoice-maker/pdf"

	"github.com/shopspring/decimal"
)

func replaceField(t string, label string, value string) string {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]`)

	var result = t

	allSubmatches := re.FindAllStringSubmatch(t, -1)
	if len(allSubmatches) == 0 {
		return result
	}

	for _, submatches := range allSubmatches {
		if len(submatches) == 0 {
			continue
		}
		submatch := submatches[0]

		offset := utf8.RuneCountInString(submatch) - utf8.RuneCountInString(value)
		padding := strings.Repeat(" ", offset)
		finalLabel := value + padding

		result = strings.ReplaceAll(result, submatch, finalLabel)
	}

	return result
}

func insertRows(t string, label string, value string) string {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]\n`)

	return re.ReplaceAllString(t, value)
}

func ApplyInvoice(template string, rowTemplate string, cfg config.Invoice) (result string) {
	result = template

	result = replaceField(result, "IssuerName", cfg.Issuer.Name)
	result = replaceField(result, "IssuerAddress", cfg.Issuer.Address)
	result = replaceField(result, "IssuerTaxID", cfg.Issuer.TaxID)
	result = replaceField(result, "AccountNo", cfg.Issuer.Account)
	result = replaceField(result, "IssuerBankName", cfg.Issuer.BankName)
	result = replaceField(result, "IssuerBic", cfg.Issuer.BIC)

	result = replaceField(result, "ReceiverName", cfg.Receiver.Name)
	result = replaceField(result, "ReceiverAddress", cfg.Receiver.Address)
	result = replaceField(result, "ReceiverTaxID", cfg.Receiver.TaxID)

	result = replaceField(result, "InvoiceNo", cfg.InvoiceNo)
	result = replaceField(result, "InvoiceDate", cfg.InvoiceDate)
	result = replaceField(result, "PaymentType", cfg.PaymentType)
	result = replaceField(result, "DueDate", cfg.DueDate)

	amountSum := decimal.Zero
	for _, v := range cfg.Items {
		amountSum.Add(v.Amount)
	}
	result = replaceField(result, "ASum", amountSum.String())

	vatSum := decimal.Zero
	for _, v := range cfg.Items {
		vatSum.Add(v.VatAmount)
	}
	result = replaceField(result, "TaxSum", vatSum.String())

	totSum := decimal.Zero
	for _, v := range cfg.Items {
		totSum.Add(v.VatAmount)
	}
	result = replaceField(result, "TotSum", totSum.String())

	if rowTemplate != "" && len(cfg.Items) > 0 {
		items := ""
		for _, item := range cfg.Items {
			items += rowTemplate
			item.CalculateTotal()

			items = replaceField(items, "Title", fmt.Sprint(item.Title))
			items = replaceField(items, "Qty", fmt.Sprint(item.Quantity))
			items = replaceField(items, "Unit", fmt.Sprint(item.Unit))
			items = replaceField(items, "Price", item.Price.String())
			items = replaceField(items, "Amount", item.Amount.String())
			items = replaceField(items, "VR", fmt.Sprintf("%d%%", item.VatRate))
			items = replaceField(items, "VA", fmt.Sprint(item.VatAmount.String()))
			items = replaceField(items, "Total", fmt.Sprint(item.Total))
		}
		result = insertRows(result, "Items", items)
	}

	return result
}

func saveFile(dirname string, filename string, content string) error {
	if err := os.MkdirAll(dirname, 0744); err != nil {
		return err
	}

	mddir := filepath.Join(dirname, filename)

	file, err := os.Create(mddir)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(content); err != nil {
		log.Fatal("write string err", err)
		return err
	}
	return nil
}

func SaveInvoice(invoice string, dirname string) (string, error) {
	name := time.Now().Format("2006-01-02 15:04:05")
	mdName := name + ".md"

	if err := saveFile(dirname, mdName, invoice); err != nil {
		return "", errors.New("issue while writting markdown file: " + err.Error())
	}

	bytes, err := os.ReadFile(filepath.Join(dirname, mdName))
	if err != nil {
		return "", errors.New("issue while reading markdown file: " + err.Error())
	}

	htmlBytes := pdf.ToHTML(bytes)
	htmlName := name + ".html"
	if err := os.WriteFile(filepath.Join(dirname, htmlName), htmlBytes, 0744); err != nil {
		return "", errors.New("issue while writting html file: " + err.Error())
	}

	pdfName := name + ".pdf"
	re := regexp.MustCompile(`<?.pre>`)
	pdfContent := re.ReplaceAllString(invoice, "")
	err = pdf.PrintInvoice(pdfContent, filepath.Join(dirname, pdfName))
	if err != nil {
		return "", errors.New("issue while writting pdf file: " + err.Error())
	}

	return dirname, nil
}
