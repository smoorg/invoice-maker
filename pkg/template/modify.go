package template

import (
	"fmt"
	"log"
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
	if len(allSubmatches) == 0 || len(allSubmatches[0]) == 0 {
		return nil
	}

	submatch := allSubmatches[0][0]
	offset := utf8.RuneCountInString(submatch) - utf8.RuneCountInString(value)

	if label == "Title" {
		// this is when amount of characters for a field value is less than field in the template
		if offset < 0 {
			row1, row2 := SplitIntoTwoByWord(value, submatch)
			localResult = strings.Replace(localResult, submatch, row1, 1)
			localResult = strings.Replace(localResult, submatch, row2, 2)
		} else {
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
	final := strings.Replace(*result, submatch, finalLabel, 1)
	*result = final

	return nil
}

func SplitIntoTwoByWord(value string, match string) (string, string) {
	var row1, row2 string

	valueWords := strings.SplitSeq(value, " ")

	for word := range valueWords {
		newVal := row1 + " " + word
		offset := utf8.RuneCountInString(match) - utf8.RuneCountInString(newVal)
		if offset < 0 {
			row2 = value[len(row1)-1:]
			break
		}
		row1 = newVal
	}

	offsetRow1 := utf8.RuneCountInString(match) - utf8.RuneCountInString(row1)
	row1 = row1 + strings.Repeat(" ", offsetRow1)

	offsetRow2 := utf8.RuneCountInString(match) - utf8.RuneCountInString(row2)
	row2 = row2 + strings.Repeat(" ", offsetRow2)
	return row1, row2
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

type Item struct {
	Fields []*ItemField
}

type ItemField struct {
	Label  string
	Value  string
	MaxLen int
}

func GetItemTemplateSize(label, rowTemplate string) int {
	re := regexp.MustCompile(`\[\s+` + label + `\s+\]`)

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

const (
	TemplFTitle           = "Title"
	TemplFQuantity        = "Qty"
	TemplFUnit            = "Unit"
	TemplFPrice           = "Price"
	TemplFAmount          = "Amount"
	TemplFVatRate         = "VR"
	TemplFVatAmount       = "VA"
	TemplFTotal           = "Total"
	TemplFIssuerName      = "IssuerName"
	TemplFIssuerAddress   = "IssuerAddress"
	TemplFIssuerTaxID     = "IssuerTaxID"
	TemplFAccountNo       = "AccountNo"
	TemplFIssuerBankName  = "IssuerBankName"
	TemplFIssuerBic       = "IssuerBic"
	TemplFReceiverName    = "ReceiverName"
	TemplFReceiverAddress = "ReceiverAddress"
	TemplFReceiverTaxID   = "ReceiverTaxID"
	TemplFPaymentType     = "PaymentType"
	TemplFInvoiceNo       = "InvoiceNo"
	TemplFInvoiceDate     = "InvoiceDate"
	TemplFDueDate         = "DueDate"
	TemplFASum            = "ASum"
	TemplFTaxSum          = "TaxSum"
	TemplFTotSum          = "TotSum"
)

func ParseItems(i *config.Invoice, rowTemplate string) []*Item {
	items := []*Item{}
	for _, v := range i.Items {
		itemFields := []*ItemField{}
		itemFields = append(
			itemFields,
			&ItemField{
				Label: TemplFTitle,
				Value: v.Title,
			},
			&ItemField{
				Label: TemplFQuantity,
				Value: fmt.Sprint(v.Quantity),
			},
			&ItemField{
				Label: TemplFUnit,
				Value: fmt.Sprint(v.Unit),
			},
			&ItemField{
				Label: TemplFPrice,
				Value: v.Price,
			},
			&ItemField{
				Label: TemplFAmount,
				Value: v.Amount,
			},
			&ItemField{
				Label: TemplFVatRate,
				Value: fmt.Sprint(v.VatRate, "%"),
			},
			&ItemField{
				Label: TemplFVatAmount,
				Value: v.CalculateVatAmount().StringFixed(2),
			},
			&ItemField{
				Label: TemplFTotal,
				Value: v.CalculateItemTotal().StringFixed(2),
			},
		)

		for _, field := range itemFields {
			field.MaxLen = GetItemTemplateSize(field.Label, rowTemplate)
		}

		items = append(items, &Item{Fields: itemFields})
	}

	return items
}

func ApplyInvoice(templateStr *string, rowTemplate string, cfg *config.Invoice) error {
	amount, tax, total := SumUp(&cfg.Items)
	if rowTemplate == "" || len(cfg.Items) <= 0 {
		return nil
	}

	items := ParseItems(cfg, rowTemplate)

	filledRows := []string{}
	log.Printf("items len %d", len(items))
	for idx, v := range items {
		str, err := AppendItem(v, rowTemplate)
		if err != nil {
			panic(err)
		}

		// trim newline on last item
		if idx == (len(items) - 1) {
			str = str[0 : len(str)-2]
		}
		filledRows = append(filledRows, str)
	}

	log.Println("inserting rows to templateStr", filledRows)

	if err := replaceField(templateStr, "Items", strings.Join(filledRows, "")); err != nil {
		panic(err)
	}

	fields := &map[string]string{
		TemplFIssuerName:      cfg.Issuer.Name,
		TemplFIssuerAddress:   cfg.Issuer.Address,
		TemplFIssuerTaxID:     cfg.Issuer.TaxID,
		TemplFAccountNo:       cfg.Issuer.Account,
		TemplFIssuerBankName:  cfg.Issuer.BankName,
		TemplFIssuerBic:       cfg.Issuer.BIC,
		TemplFReceiverName:    cfg.Receiver.Name,
		TemplFReceiverAddress: cfg.Receiver.Address,
		TemplFReceiverTaxID:   cfg.Receiver.TaxID,
		TemplFPaymentType:     cfg.PaymentType,
		TemplFInvoiceNo:       cfg.InvoiceNo,
		TemplFInvoiceDate:     cfg.InvoiceDate,
		TemplFDueDate:         cfg.DueDate,
		TemplFASum:            amount.StringFixed(2),
		TemplFTaxSum:          tax.StringFixed(2),
		TemplFTotSum:          total.StringFixed(2),
	}
	for k, v := range *fields {
		if err := replaceField(templateStr, k, v); err != nil {
			return err
		}
	}
	return nil
}

func AppendItem(v *Item, rowTemplate string) (string, error) {
	var itemsStr = ""
	itemsStr += rowTemplate

	// if any field is bigger than a row template then double the lines of a row template
	// exclude MaxLen == 0 from that rule; row has no such field so we will ignore it anyway
	for _, f := range v.Fields {
		if f.MaxLen > 0 && utf8.RuneCountInString(f.Value) > f.MaxLen {
			log.Printf("value %s is bigger than len %d", f.Value, f.MaxLen)
			itemsStr += rowTemplate
			break
		}
	}

	// string representation of item to apply
	for _, f := range v.Fields {
		if err := replaceField(&itemsStr, f.Label, f.Value); err != nil {
			return "", err
		}
	}

	// go again to clear not replaced crap like double lines for expand title purpose
	for _, f := range v.Fields {
		if err := replaceField(&itemsStr, f.Label, ""); err != nil {
			return "", err
		}
	}

	log.Printf("append item: %s", itemsStr)

	return itemsStr, nil
}

func ToHTML(invoice string) ([]byte, error) {
	htmlBytes := markdown.ToHTML([]byte(invoice), nil, nil)
	return htmlBytes, nil
}
