package template_test

import (
	"fmt"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/template"
	"log"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestApplyAddress(t *testing.T) {
	sut := "| [ InvoiceAddress                     ] |"
	initialLength := len(sut)
	i := &config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	err := template.ApplyInvoice(&sut, "", i)

	if err != nil {
		t.Error("error thrown", err)
	}

	if initialLength != len(sut) {
		t.Error("length mismatch")
	}
}

func TestApplyInvoiceRows(t *testing.T) {
	sut := "| [ Items ] |"
	i := &config.Invoice{
		Items: []config.InvoiceItem{
			{
				Title:    "Cheese",
				Quantity: 2,
				Unit:     "kg",
				Price:    "20",
				VatRate:  15,
			},
			{
				Title:    "Cottage Cheese",
				Quantity: 4,
				Unit:     "kg",
				Price:    "20",
				VatRate:  15,
			},
		},
	}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	rowTemplate := "│ [ Title                ] │ [Qty] │ [Unit]  │ [Price] │ [Amount] │ [VR]  │ [VA   ] │ [Total]  │"

	err := template.ApplyInvoice(&sut, rowTemplate, i)

	if err != nil {
		t.Error("error thrown", err)
	}

	lines := strings.Split(sut, "\n")
	if len(lines) != 2 {
		t.Error("invalid number of lines", len(lines))
	}

	// expected1 := fmt.Sprintf("│ %s │ %s │ %s  │ %s │ %s │ %s  │ %s │ %s  │",i.Items[0].ToInvoice())
	// expected1 := fmt.Sprintf("│ %s │ %s │ %s  │ %s │ %s │ %s  │ %s │ %s  │",)

	// if sut != expected
}

func TestTotalCalculations(t *testing.T) {
	sut := "| [ IssuerAddress                     ] |\n[ Items                                     ]\n"
	rowTemplate := "| [ Total                     ] |"

	i := &config.Invoice{}
	i.Issuer.Address = "ul. Narutowicza 14B/2, 80-501 Gdańsk"

	item := config.InvoiceItem{}
	item.VatRate = 23
	item.Price = "25000"
	item.Quantity = 1
	price, err := decimal.NewFromString(item.Price)
	if err != nil {
		t.Error(err)
	}
	total := price.Mul(decimal.NewFromInt32(item.Quantity)).Mul(decimal.NewFromFloat32(1.23))

	i.Items = append(i.Items, item)

	err = template.ApplyInvoice(&sut, rowTemplate, i)
	if err != nil {
		t.Error("error thrown", err)
	}

	if !strings.Contains(sut, total.String()) {
		t.Error("total invalid", sut, total)
	}
}

func TestApplySingleItem(t *testing.T) {
	templ := fmt.Sprintf("| [ %s ] | [ %s ] |",
		template.TemplFPrice,
		template.TemplFQuantity,
	)
	log.Println("templ", templ)

	row := &template.Item{
		Fields: []*template.ItemField{
			{
				Label:  template.TemplFPrice,
				Value:  "206",
				MaxLen: template.GetItemTemplateSize(template.TemplFPrice, templ),
			},
			{
				Label:  template.TemplFQuantity,
				Value:  "4",
				MaxLen: template.GetItemTemplateSize(template.TemplFQuantity, templ),
			},
		},
	}

	sut, err := template.AppendItem(row, string(templ))
	if err != nil {
		t.Error(err)
	}

	log.Printf("max len, qty: %d; price: %d", row.Fields[0].MaxLen, row.Fields[1].MaxLen)

	expectedTempl := fmt.Sprintf("| %s       | %s       |", row.Fields[0].Value, row.Fields[1].Value)

	result := strings.Compare(sut, expectedTempl)
	if sut != expectedTempl {
		t.Error("invalid value", result, []byte(sut), []byte(expectedTempl), "\n", sut, "\n", expectedTempl)
	}
}

func TestApplyTwoItems(t *testing.T) {
	templ := fmt.Sprintf("| [ %s ] | [ %s ] |",
		template.TemplFPrice,
		template.TemplFQuantity,
	)
	log.Println("templ", templ)

	rows := []*template.Item{
		{
			Fields: []*template.ItemField{
				{
					Label:  template.TemplFPrice,
					Value:  "206",
					MaxLen: template.GetItemTemplateSize(template.TemplFPrice, templ),
				},
				{
					Label:  template.TemplFQuantity,
					Value:  "7",
					MaxLen: template.GetItemTemplateSize(template.TemplFQuantity, templ),
				},
			},
		},
		{
			Fields: []*template.ItemField{
				{
					Label:  template.TemplFPrice,
					Value:  "206",
					MaxLen: template.GetItemTemplateSize(template.TemplFPrice, templ),
				},
				{
					Label:  template.TemplFQuantity,
					Value:  "4",
					MaxLen: template.GetItemTemplateSize(template.TemplFQuantity, templ),
				},
			},
		},
	}

	invoiceRows := []string{}
	expectedTempl := ""
	for idx, row := range rows {
		str, err := template.AppendItem(row, string(templ))
		if err != nil {
			t.Error(err)
		}

		log.Printf("max len, qty: %d; price: %d", row.Fields[0].MaxLen, row.Fields[1].MaxLen)

		invoiceRows = append(invoiceRows, str)

		expectedTempl += fmt.Sprintf("| %s       | %s       |", row.Fields[0].Value, row.Fields[1].Value)
		if idx != (len(rows) - 1) {
			log.Printf("triggered idx %v templ %s rows %s", idx, invoiceRows, expectedTempl)
			expectedTempl += "\n"
		}
	}
	log.Printf("rows %v", invoiceRows)

	sut := strings.Join(invoiceRows, "\n")
	result := strings.Compare(sut, expectedTempl)
	if sut != expectedTempl {
		t.Error("invalid value", result, "\n", sut, "\n", expectedTempl, "\n")
	}
}
