package template_test

import (
	"strings"
	"testing"

	"invoice-maker/internal/template"
)

func TestInsertItems(t *testing.T) {
	row := "[ Items                       ]\n"

	sut := template.InsertRows(row, "Items", "test")

	if !strings.Contains(sut, "test") {
		t.Error("InsertRows did not add value properly", sut)
	}
}
