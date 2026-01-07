package invoices

import (
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/template"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InvoicesModel struct {
	Invoices        []config.Invoice `yaml:"invoices"`
	invoicesTable   table.Model
	invoiceViewFlex flexbox.HorizontalFlexBox
	preview         bool
}

func New() InvoicesModel {
	m := InvoicesModel{}
	m.invoiceViewFlex = *flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.invoiceViewFlex.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.invoiceViewFlex.AddColumns(columns)

	m.invoicesTable = table.New(
		table.WithColumns([]table.Column{
			{Title: "Date", Width: 10},
			{Title: "Payment Date", Width: 10},
			{Title: "Invoice No.", Width: 11},
			{Title: "Receiver", Width: 20},
			{Title: "Net", Width: 8},
			{Title: "Gross", Width: 8},
		}),
	)
	m.invoicesTable.SetHeight(5)
	m.invoicesTable.SetWidth(10)

	rows := []table.Row{}
	for _, v := range m.Invoices {
		rows = append(rows, table.Row{
			v.DeliveryDate,
			v.DueDate,
			v.InvoiceNo,
			v.Receiver.Name,
			v.NetSum(),
			v.GrossSum(),
		})
	}
	m.invoicesTable.SetRows(rows)

	return m
}

func (m InvoicesModel) Init() tea.Cmd {
	// invoices
	return nil
}
func (m InvoicesModel) Update(msg tea.Msg) (InvoicesModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pkg.SetInvoiceRows:
		rows := []table.Row{}
		for _, v := range msg.Rows {
			rows = append(rows, table.Row{
				v.DeliveryDate,
				v.DueDate,
				v.InvoiceNo,
				v.Receiver.Name,
				v.NetSum(),
				v.GrossSum(),
			})
		}
		m.invoicesTable.SetRows(rows)
	case tea.WindowSizeMsg:
		m.invoicesTable.SetWidth(msg.Width)
		m.invoicesTable.SetHeight(msg.Height)
		m.invoiceViewFlex.SetWidth(msg.Width)
		m.invoiceViewFlex.SetHeight(msg.Height)
	case tea.Key:
		switch msg.String() {
		case "h":
			if m.preview {
				m.preview = false
			} else {
				cmd = pkg.GoMain()
			}
		case "j":
			m.invoicesTable.MoveDown(1)
		case "k":
			m.invoicesTable.MoveUp(1)
		case "p":
			m.preview = true
		}
	}

	m.invoicesTable, cmd = m.invoicesTable.Update(msg)
	if m.invoiceViewFlex.ColumnsLen() > 0 {
		m.invoiceViewFlex.GetColumn(0).GetCell(0).SetContent(m.invoicesTable.View())
	}

	return m, cmd
}
func (m InvoicesModel) View() string {
	if m.preview {
		picked := m.invoicesTable.SelectedRow()
		if picked == nil {
			return ""
		}
		for _, v := range m.Invoices {
			if v.InvoiceNo == picked[2] && v.NetSum() == picked[4] {
				content, err := template.GetContent(&v)
				if err != nil {
					panic(err)
				}
				return content
			}
		}
		return ""
	}

	availHeight := m.invoiceViewFlex.GetHeight()

	helpBubble := help.New()
	kb := []key.Binding{
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "preview"),
		),
		key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
	}
	helpContent := helpBubble.ShortHelpView(kb)
	availHeight -= lipgloss.Height(helpContent)

	m.invoiceViewFlex.SetHeight(availHeight)
	content := m.invoiceViewFlex.Render()

	var sections []string
	sections = append(sections, content)
	sections = append(sections, helpContent)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
