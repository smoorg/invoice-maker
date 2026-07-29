package preview

import (
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	pkg_help "invoice-maker/pkg/help"
	"invoice-maker/pkg/template"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type InvoicePreviewModel struct {
	invoice     *config.Invoice
	keys        Keymap
	helpContent string
}

type Keymap struct {
	Back  key.Binding
	Next  key.Binding
	Print key.Binding
}

func NewPreviewModel() InvoicePreviewModel {
	m := InvoicePreviewModel{}

	m.keys = Keymap{
		Back: key.NewBinding(
			key.WithKeys("h", "q", tea.KeyEsc.String()),
			key.WithHelp("h/q/esc", "go back"),
		),
		Next: key.NewBinding(
			key.WithKeys("L", "l", tea.KeyRight.String()),
			key.WithHelp("l/L", "edit"),
		),
		Print: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "print"),
		),
	}

	m.helpContent = help.New().ShortHelpView(pkg_help.MapToBindingsList(m.keys))

	return m
}

func (m *InvoicePreviewModel) SetInvoice(v config.Invoice) {
	m.invoice = &v
}

func (m InvoicePreviewModel) Init() {}

func (m InvoicePreviewModel) Update(msg tea.Msg) (InvoicePreviewModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Next):
			cmds = append(cmds, pkg.GoInvoiceEdit())
		case key.Matches(msg, m.keys.Back):
			cmds = append(cmds, pkg.GoInvoiceList())
		case key.Matches(msg, m.keys.Print):
			log.Println("go invoice print")
			cmds = append(cmds, pkg.GoInvoicePrint())
		}
	}
	return m, tea.Batch(cmds...)
}

func (m InvoicePreviewModel) View() string {
	b := strings.Builder{}

	content, err := template.GetContent(m.invoice)
	if err != nil {
		return err.Error()
	}
	b.WriteString(content)
	b.WriteString("\n")
	b.WriteString(m.helpContent)

	return b.String()
}
