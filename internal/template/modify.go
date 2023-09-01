package template

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"invoice-maker/internal/config"

	"github.com/gomarkdown/markdown"
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

		var padding string

		// this is when amount of characters for a field value is less than field in the template
		if offset < 0 {
			log.Fatalf("offset for field '%s' is negative", label)
		}

		if offset > 0 {
			padding = strings.Repeat(" ", offset)
		}

		finalLabel := value + padding

		result = strings.ReplaceAll(result, submatch, finalLabel)
	}

	return result
}

func InsertRows(t string, label string, value string) string {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]\n`)

	return re.ReplaceAllString(t, value)
}

func ApplyInvoice(templateStr string, rowTemplate string, cfg *config.Invoice) (result string) {
	result = templateStr

	result = replaceField(result, "IssuerName", cfg.Issuer.Name)
	result = replaceField(result, "IssuerAddress", cfg.Issuer.Address)
	result = replaceField(result, "IssuerTaxID", cfg.Issuer.TaxID)
	result = replaceField(result, "AccountNo", cfg.Issuer.Account)
	result = replaceField(result, "IssuerBankName", cfg.Issuer.BankName)
	result = replaceField(result, "IssuerBic", cfg.Issuer.BIC)

	result = replaceField(result, "ReceiverName", cfg.Receiver.Name)
	result = replaceField(result, "ReceiverAddress", cfg.Receiver.Address)
	result = replaceField(result, "ReceiverTaxID", cfg.Receiver.TaxID)
	result = replaceField(result, "PaymentType", cfg.PaymentType)

	result = replaceField(result, "InvoiceNo", cfg.InvoiceNo)
	result = replaceField(result, "InvoiceDate", cfg.InvoiceDate)
	result = replaceField(result, "PaymentType", cfg.PaymentType)
	result = replaceField(result, "DueDate", cfg.DueDate)

	amountSum := decimal.NewFromInt32(0)
	vatSum := decimal.NewFromInt32(0)
	totSum := decimal.NewFromInt32(0)
	if rowTemplate != "" && len(cfg.Items) > 0 {
		itemsStr := ""

		for i := range cfg.Items {
			itemsStr += rowTemplate
			cfg.Items[i].CalculateItemTotal()

			itemsStr = replaceField(itemsStr, "Title", fmt.Sprint(cfg.Items[i].Title))
			itemsStr = replaceField(itemsStr, "Qty", fmt.Sprint(cfg.Items[i].Quantity))
			itemsStr = replaceField(itemsStr, "Unit", fmt.Sprint(cfg.Items[i].Unit))
			itemsStr = replaceField(itemsStr, "Price", cfg.Items[i].Price.String())
			itemsStr = replaceField(itemsStr, "Amount", cfg.Items[i].Amount.String())
			itemsStr = replaceField(itemsStr, "VR", fmt.Sprintf("%d%%", cfg.Items[i].VatRate))
			itemsStr = replaceField(itemsStr, "VA", fmt.Sprint(cfg.Items[i].VatAmount.String()))
			itemsStr = replaceField(itemsStr, "Total", fmt.Sprint(cfg.Items[i].Total))

			amountSum = amountSum.Add(cfg.Items[i].Amount)
			vatSum = vatSum.Add(cfg.Items[i].VatAmount)
			totSum = totSum.Add(cfg.Items[i].Total)

		}
		result = InsertRows(result, "Items", itemsStr)
	}

	result = replaceField(result, "ASum", amountSum.String())
	result = replaceField(result, "TaxSum", vatSum.String())
	result = replaceField(result, "TotSum", totSum.String())

	return result
}

func ToHTML(invoice string) ([]byte, error) {
	htmlBytes := markdown.ToHTML([]byte(invoice), nil, nil)
	return htmlBytes, nil
}
