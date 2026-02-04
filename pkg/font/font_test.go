package font_test

import (
	f "invoice-maker/pkg/font"
	"testing"
)

func TestFindRobotoFont(t *testing.T) {
	fonts, err := f.FindFonts("FreeSerif", "italic")

	if err != nil || len(fonts) == 0 {
		t.Error("should find roboto font")
	}

	if len(fonts) > 1 {
		t.Error("should find single font")
	}
}

func TestFindStyles(t *testing.T) {
	fonts, err := f.GetFontStyles("FreeSans")

	if err != nil {
		t.Error(err)
	}

	if len(fonts) <= 1 {
		t.Error("should have more than one style")
	}
}

func TestFindFamilies(t *testing.T) {
	fonts, err := f.GetFontFamilies()

	if err != nil {
		t.Error(err)
	}

	if len(fonts) == 0 {
		t.Error("should have more than one style")
	}
}
