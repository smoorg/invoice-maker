package template

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"invoice-maker/internal/config"

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

		var padding string

		// this is when amount of characters for a field value is less than field in the template
		if offset < 0 {
			return errors.New(fmt.Sprintf("offset for field '%s' is negative", label))
		}

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

func SumUp(items *[]config.InvoiceItem) (string, string, string) {
	if len(*items) == 0 {
		return "0", "0", "0"
	}

	amountSum := decimal.NewFromInt32(0)
	vatSum := decimal.NewFromInt32(0)
	totSum := decimal.NewFromInt32(0)

	for _, item := range *items {
		item.CalculateItemTotal()

		amountSum = amountSum.Add(item.Amount)
		vatSum = vatSum.Add(item.VatAmount)
		totSum = totSum.Add(item.Total)
	}

	return amountSum.String(), vatSum.String(), totSum.String()
}

func ApplyInvoice(templateStr string, rowTemplate string, cfg *config.Invoice) (*string, error) {
	result := &templateStr

	amount, tax, total := SumUp(&cfg.Items)
	if rowTemplate != "" && len(cfg.Items) > 0 {
		empty := ""
		var itemsStr *string = &empty

		for i := range cfg.Items {
			*itemsStr += rowTemplate
			cfg.Items[i].CalculateItemTotal()

			itemFields := &map[string]string{
				"Title":  fmt.Sprint(cfg.Items[i].Title),
				"Qty":    fmt.Sprint(cfg.Items[i].Quantity),
				"Unit":   fmt.Sprint(cfg.Items[i].Unit),
				"Price":  fmt.Sprint(cfg.Items[i].Price.String()),
				"Amount": fmt.Sprint(cfg.Items[i].Amount.String()),
				"VR":     fmt.Sprintf("%d%%", cfg.Items[i].VatRate),
				"VA":     fmt.Sprint(cfg.Items[i].VatAmount.String()),
				"Total":  fmt.Sprint(cfg.Items[i].Total),
			}

			for k, v := range *itemFields {
				if err := replaceField(itemsStr, k, v); err != nil {
					return result, err
				}
			}

		}
		*result = InsertRows(*result, "Items", *itemsStr)

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
		"ASum":            amount,
		"TaxSum":          tax,
		"TotSum":          total,
	}
	for k, v := range *fields {
		if err := replaceField(result, k, v); err != nil {
			return result, err
		}
	}
	return result, nil
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

func ToHTML(invoice string) ([]byte, error) {
	htmlBytes := markdown.ToHTML([]byte(invoice), nil, nil)
	return htmlBytes, nil
}
