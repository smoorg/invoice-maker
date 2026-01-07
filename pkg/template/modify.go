package template

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"invoice-maker/pkg/config"

	"github.com/gomarkdown/markdown"
	"github.com/shopspring/decimal"
)

func replaceField(result *string, label string, value string) error {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]`)

	allSubmatches := re.FindAllStringSubmatch(*result, -1)
	if len(allSubmatches) == 0 {
		return nil
	}

	for _, submatches := range allSubmatches {
		if len(submatches) == 0 {
			continue
		}
		submatch := submatches[0]
		offset := utf8.RuneCountInString(submatch) - utf8.RuneCountInString(value)

		// this is when amount of characters for a field value is less than field in the template
		if offset < 0 {
			return fmt.Errorf("offset for field '%s' is negative", label)
		}

		padding := ""
		if offset > 0 {
			padding = strings.Repeat(" ", offset)
		}
		finalLabel := value + padding
		final := strings.ReplaceAll(*result, submatch, finalLabel)
		*result = final
	}

	return nil
}

func InsertRows(t string, label string, value string) string {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]\n`)

	return re.ReplaceAllString(t, value)
}

func SumUp(items *[]config.InvoiceItem) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	if len(*items) == 0 {
		return decimal.NewFromInt32(0),
			decimal.NewFromInt32(0),
			decimal.NewFromInt32(0)
	}

	amountSum := decimal.NewFromInt32(0)
	vatSum := decimal.NewFromInt32(0)
	totSum := decimal.NewFromInt32(0)

	for _, item := range *items {
		item.CalculateItemTotal()

		amountSum = amountSum.Add(item.CalculateAmount())
		vatSum = vatSum.Add(item.CalculateVatAmount())
		totSum = totSum.Add(item.CalculateItemTotal())
	}

	return amountSum, vatSum, totSum
}

func ApplyInvoice(templateStr *string, rowTemplate string, cfg *config.Invoice) error {
	amount, tax, total := SumUp(&cfg.Items)
	if rowTemplate != "" && len(cfg.Items) > 0 {
		itemsStr := ""

		for i := range cfg.Items {
			itemsStr += rowTemplate
			cfg.Items[i].CalculateItemTotal()

			itemFields := &map[string]string{
				"Title":  fmt.Sprint(cfg.Items[i].Title),
				"Qty":    fmt.Sprint(cfg.Items[i].Quantity),
				"Unit":   fmt.Sprint(cfg.Items[i].Unit),
				"Price":  fmt.Sprint(cfg.Items[i].Price),
				"Amount": fmt.Sprint(cfg.Items[i].Amount),
				"VR":     fmt.Sprintf("%d%%", cfg.Items[i].VatRate),
				"VA":     fmt.Sprint(cfg.Items[i].CalculateVatAmount().StringFixed(2)),
				"Total":  fmt.Sprint(cfg.Items[i].CalculateItemTotal().StringFixed(2)),
			}

			for k, v := range *itemFields {
				if err := replaceField(&itemsStr, k, v); err != nil {
					return err
				}
			}

		}
		*templateStr = InsertRows(*templateStr, "Items", itemsStr)

	}
	fields := &map[string]string{
		"IssuerName":      cfg.Issuer.Name,
		"IssuerAddress":   cfg.Issuer.Address,
		"IssuerTaxID":     cfg.Issuer.TaxID,
		"AccountNo":       cfg.Issuer.Account,
		"IssuerBankName":  cfg.Issuer.BankName,
		"IssuerBic":       cfg.Issuer.BIC,
		"ReceiverName":    cfg.Receiver.Name,
		"ReceiverAddress": cfg.Receiver.Address,
		"ReceiverTaxID":   cfg.Receiver.TaxID,
		"PaymentType":     cfg.PaymentType,
		"InvoiceNo":       cfg.InvoiceNo,
		"InvoiceDate":     cfg.InvoiceDate,
		"DueDate":         cfg.DueDate,
		"ASum":            amount.StringFixed(2),
		"TaxSum":          tax.StringFixed(2),
		"TotSum":          total.StringFixed(2),
	}
	for k, v := range *fields {
		if err := replaceField(templateStr, k, v); err != nil {
			return err
		}
	}
	return nil
}

func ToHTML(invoice string) ([]byte, error) {
	htmlBytes := markdown.ToHTML([]byte(invoice), nil, nil)
	return htmlBytes, nil
}
