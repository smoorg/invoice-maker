package invoice_list

import (
	"errors"
	"fmt"
	"invoice-maker/internal/pdf"
	"invoice-maker/internal/template"
	"invoice-maker/internal/types"
	"invoice-maker/pkg/font"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func printInvoice(tui *types.TUI) (string, error) {
	i := &tui.Config.Invoices[selectedInvoice-1]
	dir, err := tui.Config.GetInvoiceDirectory()

	tmpl, err := template.GetTemplate()
	if err != nil {
		return "", errors.New("missing template")
	}

	rowTemplate, err := template.GetRowTemplate()
	if err != nil {
		return "", errors.New("unable to locate row template")
	}

	inv, err := template.ApplyInvoice(string(tmpl), string(rowTemplate), i)
	if err != nil {
		return "", err
	}

	fonts, err := font.FindFonts(tui.Config.Font.Family, tui.Config.Font.Style)
	if err != nil {
		return "", err
	}

	if len(fonts) == 0 {
		errMsg := fmt.Sprint(
			"font from the config could not be found in the system, font-family: ",
			tui.Config.Font.Family, "font-style: ", tui.Config.Font.Style)
		return "", errors.New(errMsg)
	}

	htmlBytes, err := template.ToHTML(*inv)

	name := time.Now().Format("2006-01-02 15:04:05")
	mdName := name + ".md"
	htmlName := name + ".html"
	pdfName := name + ".pdf"

	if err := saveFile(dir, mdName, []byte(*inv)); err != nil {
		return "", errors.New("issue while writting markdown file: " + err.Error())
	}
	if err := saveFile(dir, htmlName, htmlBytes); err != nil {
		return "", errors.New("issue while writting html file: " + err.Error())
	}

	re := regexp.MustCompile(`<?.pre>`)
	pdfContent := re.ReplaceAllString(*inv, "")

	pdf.InitializePdf("")

	if err := pdf.SetFont(tui.Config.Font.Family, tui.Config.Font.Style, tui.Config.Font.Filepath, 8); err != nil {
		return "", err
	}

	pdf.SetText(pdfContent, 0, 4)

	path := filepath.Join(dir, pdfName)
	if err := pdf.Output(path); err != nil {
		panic("pdf output: " + err.Error())
	}

	return path, nil
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
