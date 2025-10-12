package pdf

import (
	"os"
	"strings"

	"github.com/phpdave11/gofpdf"
)

var file *gofpdf.Fpdf

func InitializePdf(fontdir string) {
	if fontdir == "" {
		//TODO: make it more generic to work on BSD and Windows
		fontdir = "/usr/share/fonts/TTF/"
	}
	file = gofpdf.New("P", "mm", "A4", fontdir)
	file.AddPage()
}

// gofpdf expects "B" for bold or "I" for italic as a font style, may be combined
func getFpdfStyle(f string) string {
	style := ""
	switch {
	case strings.Contains(f, "Bold"):
		style += "B"
	case strings.Contains(f, "Italic"):
		style += "I"
	}
	return style
}

func SetFont(fontname string, fontstyle string, fontpath string, size float64) error {
	fontBytes, err := os.ReadFile(fontpath)
	if err != nil {
		return err
	}

	fpdf_style := getFpdfStyle(fontstyle)

	file.AddUTF8FontFromBytes(fontname, fpdf_style, fontBytes)
	file.SetFont(fontname, fpdf_style, size)

	return nil
}

func SetText(txt string, width float64, height float64) {
	file.MultiCell(width, height, txt, "", "", false)
}

func Output(path string) error {
	err := file.OutputFileAndClose(path)
	if err != nil {
		return err
	}

	return nil
}
