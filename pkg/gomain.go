package pkg

import (
	"invoice-maker/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

type JumpMainView struct{}
type JumpInvoicePreview struct{}

type SetInvoiceRows struct {
	Rows []config.Invoice
}

func GoMain() tea.Cmd {
	return func() tea.Msg {
		return JumpMainView{}
	}

}
func GoInvoicePreview() tea.Cmd {
	return func() tea.Msg {
		return JumpInvoicePreview{}
	}
}
