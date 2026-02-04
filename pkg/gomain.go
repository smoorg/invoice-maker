package pkg

import (
	"invoice-maker/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

type JumpMainView struct{}
type JumpReceivers struct {}
type JumpReceiverEdit struct {
	Receiver config.Company
}
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

func GoReceivers() tea.Cmd {
	return func() tea.Msg {
		return JumpReceivers{}
	}
}

func GoReceiverEdit(v config.Company) tea.Cmd {
	return func() tea.Msg {
		return JumpReceiverEdit{
			Receiver: v,
		}
	}
}
