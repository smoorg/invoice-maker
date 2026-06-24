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
	localResult := *result
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]`)

	allSubmatches := re.FindAllStringSubmatch(*result, -1)
	if len(allSubmatches) == 0 {
		return nil
	}

	submatch := allSubmatches[0][0]
	offset := utf8.RuneCountInString(submatch) - utf8.RuneCountInString(value)

	if label == "Title" {
		// this is when amount of characters for a field value is less than field in the template
		if offset < 0 {
			// add first row
			val := value[0:(len(submatch) - 1)]
			offset = utf8.RuneCountInString(submatch) - utf8.RuneCountInString(val)
			val = val + strings.Repeat(" ", offset)
			localResult = strings.Replace(localResult, submatch, val, 1)

			// add second row
			val = value[(len(submatch) - 1):(len(value) - 1)]
			offset = utf8.RuneCountInString(submatch) - utf8.RuneCountInString(val)
			val = val + strings.Repeat(" ", offset)
			localResult = strings.Replace(localResult, submatch, val, 2)
		} else {
			// add first row
			offset = utf8.RuneCountInString(submatch) - utf8.RuneCountInString(value)
			val := value + strings.Repeat(" ", offset)
			localResult = strings.Replace(localResult, submatch, val, 1)
		}

		*result = localResult
		return nil
	}

	padding := ""
	if offset > 0 {
		padding = strings.Repeat(" ", offset)
	}
	finalLabel := value + padding
	final := strings.ReplaceAll(*result, submatch, finalLabel)
	*result = final

	return nil
}

func InsertRows(t string, label string, value string) string {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]\n`)

	return re.ReplaceAllString(t, value)
}

func SumUp(items *[]config.InvoiceItem) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	decimal.DivisionPrecision = 2
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

type ItemField struct {
	Label   string
	Value   string
	MaxSize int
}

func GetItemTemplateSize(label, rowTemplate string) int {
	re := regexp.MustCompile(`\[\s*` + label + `\s*\]`)

	allSubmatches := re.FindAllStringSubmatch(rowTemplate, -1)
	if len(allSubmatches) == 0 {
		return 0
	}
	if len(allSubmatches[0]) == 0 {
		return 0
	}
	submatch := allSubmatches[0][0]
	return utf8.RuneCountInString(submatch)
}

func ParseItems(i *config.Invoice, rowTemplate string) []*ItemField {
	itemFields := []*ItemField{}
	for _, item := range i.Items {
		itemFields = append(
			itemFields,
			&ItemField{Label: "Title", Value: item.Title},
			&ItemField{Label: "Qty", Value: fmt.Sprint(item.Quantity)},
			&ItemField{Label: "Unit", Value: fmt.Sprint(item.Unit)},
			&ItemField{Label: "Price", Value: item.Price},
			&ItemField{Label: "Amount", Value: item.Amount},
			&ItemField{Label: "VR", Value: fmt.Sprint(item.VatRate)},
			&ItemField{Label: "VA", Value: item.CalculateVatAmount().StringFixed(2)},
			&ItemField{Label: "Total", Value: item.CalculateItemTotal().StringFixed(2)},
		)

		for _, field := range itemFields {
			field.MaxSize = GetItemTemplateSize(field.Label, rowTemplate)
		}
	}

	return itemFields
}

func ApplyInvoice(templateStr *string, rowTemplate string, cfg *config.Invoice) error {
	amount, tax, total := SumUp(&cfg.Items)
	if rowTemplate == "" || len(cfg.Items) <= 0 {
		return nil
	}
	itemsStr := rowTemplate

	items := ParseItems(cfg, rowTemplate)

	// if any field is bigger than a row template then double the lines of a row template
	for _, v := range items {
		if utf8.RuneCountInString(v.Value) > v.MaxSize {
			itemsStr += rowTemplate
			break
		}
	}

	// string representation of item to aply
	for _, v := range items {
		if err := replaceField(&itemsStr, v.Label, v.Value); err != nil {
			return err
		}
	}

	*templateStr = InsertRows(*templateStr, "Items", itemsStr)

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
