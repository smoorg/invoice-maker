package pdf

import (
	"github.com/phpdave11/gofpdf"
	"os"
)

func PrintInvoice(txt string, pdfFileName string) error {
	f := gofpdf.New("P", "mm", "A4", "/usr/share/fonts/TTF/")
	f.SetFont("Arial", "B", 12)

	if fontBytes, err := getCustomFont(); err == nil {
		f.AddUTF8FontFromBytes("JetBrainsMono", "", fontBytes)
		f.SetFont("JetBrainsMono", "", 8)
	}

	f.AddPage()
	f.MultiCell(0, 4, txt, "", "", false)

	err := f.OutputFileAndClose(pdfFileName)
	if err != nil {
		return err
	}

	return nil
}

func getCustomFont() ([]byte, error) {
	fontPath := "/usr/share/fonts/TTF/JetBrainsMono-Regular.ttf"
	fontBytes, err := os.ReadFile(fontPath)

	return fontBytes, err
}
