package invoices

import (
	"invoice-maker/pkg/config"
	pkg_help "invoice-maker/pkg/help"
	"invoice-maker/pkg/template"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type InvoicePreviewModel struct {
	invoice *config.Invoice
	keys    keymap
}

type keymap struct {
	Exit key.Binding
}

func NewPreviewModel() InvoicePreviewModel {
	m := InvoicePreviewModel{}
	m.keys = keymap{
		Exit: key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "go back")),
	}

	return m
}

func (m *InvoicePreviewModel) SetInvoice(v config.Invoice) {
	m.invoice = &v
}

func (m InvoicePreviewModel) Init() {}
func (m InvoicePreviewModel) Update(msg tea.Msg) (InvoicePreviewModel, tea.Cmd) {
	return m, nil
}

func (m InvoicePreviewModel) View() string {
	content, err := template.GetContent(m.invoice)
	if err != nil {
		return err.Error()
	}

	content += "\n"
	content += help.New().ShortHelpView(pkg_help.MapToBindingsList(m.keys))

	return content
}
