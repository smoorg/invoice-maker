package invoices

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/font"
	pkg_help "invoice-maker/pkg/help"
	"invoice-maker/pkg/pdf"
	"invoice-maker/pkg/template"

	"github.com/76creates/stickers/flexbox"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type InvoiceView uint64

const (
	ViewList InvoiceView = iota
	ViewPreview
	ViewPrint
	ViewEdit
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Next  key.Binding
	Print key.Binding
	Back  key.Binding
	Edit  key.Binding
}

type InvoicesModel struct {
	invoices  []config.Invoice `yaml:"invoices"`
	invoice   InvoicePreviewModel
	edit      InvoiceEditModel
	directory string
	font      config.FontCfg

	table table.Model
	cols  *InvoiceListColumns
	flex  flexbox.HorizontalFlexBox

	helpContent  string
	printContent string
	printPath    string
	view         InvoiceView
	keys         KeyMap
}

func (m *InvoicesModel) SetConfig(cfg config.Config) {
	m.invoices = cfg.Invoices
	m.directory = cfg.Config.InvoiceDirectory
	m.font = cfg.Config
}

func (m *InvoicesModel) SetRows(rows []table.Row) {
	m.table.SetRows(rows)
}

func (m *InvoicesModel) SetSize(width int, height int) {
	m.flex.SetWidth(width)
	m.flex.SetHeight(height)
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

const (
	ColumnDate = iota
	ColumnPaymentDate
	ColumnInvoiceNo
	ColumnReceiver
	ColumnNet
	ColumnGross
)

type InvoiceListColumns struct {
	columns []table.Column
}

func NewCols() *InvoiceListColumns {
	c := make([]table.Column, 6)

	c[ColumnDate] = table.Column{Title: "Date", Width: 10}
	c[ColumnPaymentDate] = table.Column{Title: "Payment Date", Width: 10}
	c[ColumnInvoiceNo] = table.Column{Title: "Invoice No.", Width: 10}
	c[ColumnReceiver] = table.Column{Title: "Receiver", Width: 10}
	c[ColumnNet] = table.Column{Title: "Net", Width: 8}
	c[ColumnGross] = table.Column{Title: "Gross", Width: 8}

	return &InvoiceListColumns{
		columns: c,
	}
}

func (c InvoiceListColumns) Get() []table.Column {
	return c.columns
}

func (c InvoiceListColumns) GetInvoiceNo() table.Column {
	return c.columns[ColumnInvoiceNo]
}

func New(config config.Config) InvoicesModel {
	m := InvoicesModel{}
	m.cols = NewCols()
	m.table = table.New(
		table.WithColumns(m.cols.Get()),
	)
	m.table.SetHeight(5)
	m.table.SetWidth(100)

	m.flex = *flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.flex.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.flex.AddColumns(columns)
	m.directory = config.Config.InvoiceDirectory
	m.font = config.Config
	m.edit = NewEditModel()

	m.keys = KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", tea.KeyUp.String()),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", tea.KeyDown.String()),
			key.WithHelp("↓/j", "down"),
		),
		Next: key.NewBinding(
			key.WithKeys("l", tea.KeyRight.String()),
			key.WithHelp("→/l", "preview"),
		),
		Back: key.NewBinding(
			key.WithKeys("h", tea.KeyLeft.String()),
			key.WithHelp("←/h", "go back"),
		),
		Print: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "print"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
	}

	helpBubble := help.New()
	m.helpContent = helpBubble.ShortHelpView(pkg_help.MapToBindingsList(m.keys))

	return m
}

func (m InvoicesModel) Init() tea.Cmd {
	// invoices
	return nil
}
func (m InvoicesModel) Update(msg tea.Msg) (InvoicesModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
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
		m.table.SetRows(rows)
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height)
		m.flex.SetWidth(msg.Width)
		m.flex.SetHeight(msg.Height)
	case pkg.JumpInvoiceEdit:
		m.view = ViewEdit

		picked := m.table.SelectedRow()
		inv, _, err := getInvoice(m.invoices, picked[ColumnInvoiceNo], picked[ColumnNet])
		if err != nil {
			panic(err)
		}
		m.edit.invoice = inv

	case pkg.JumpInvoicePreview:
		m.view = ViewPreview
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			switch m.view {
			case ViewList:
				cmds = append(cmds, pkg.GoMain())
			case ViewPreview:
				m.view = ViewList
			case ViewPrint:
				m.view = ViewList
			}
		case key.Matches(msg, m.keys.Down):
			switch m.view {
			case ViewList:
				m.table.MoveDown(1)
			}
		case key.Matches(msg, m.keys.Up):
			switch m.view {
			case ViewList:
				m.table.MoveUp(1)
			}
		case key.Matches(msg, m.keys.Next):
			switch m.view {
			case ViewList:
				cmd = pkg.GoInvoicePreview()
				cmds = append(cmds, cmd)
			case ViewPreview:
				cmd = pkg.GoInvoiceEdit()
				cmds = append(cmds, cmd)
			}
		case key.Matches(msg, m.keys.Print):
			switch m.view {
			case ViewList:
				m.view = ViewPrint

				picked := m.table.SelectedRow()
				invContent, err := getInvoiceContent(m.invoices, picked[ColumnInvoiceNo], picked[ColumnNet])
				if err != nil {
					panic(err)
				}
				m.printPath = m.printInvoice(invContent)

				go func(file string) {
					cmd := exec.Command("xdg-open", file)
					_, err := cmd.Output()
					if err != nil {
						log.Fatal(err)
					}
				}(m.printPath)

				return m, pkg.GoInvoicePreview()
			}
		case key.Matches(msg, m.keys.Edit):
			m.view = ViewEdit
		}
	}

	switch m.view {
	case ViewEdit:
		m.edit, cmd = m.edit.Update(msg)
		cmds = append(cmds, cmd)
	case ViewPreview:
		m.invoice, cmd = m.invoice.Update(msg)
		cmds = append(cmds, cmd)
	case ViewList:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		m.flex.GetColumn(0).GetCell(0).SetContent(m.table.View())

		picked := m.table.SelectedRow()
		if picked == nil {
			panic("no invoice selected")
		}
		for _, v := range m.invoices {
			if v.InvoiceNo == picked[2] && v.NetSum() == picked[4] {
				m.invoice.SetInvoice(v)
				m.edit.SetInvoiceModel(&v)
			}
		}
	case ViewPrint:
	}

	return m, tea.Batch(cmds...)
}
func (m InvoicesModel) View() string {
	content := ""
	switch m.view {
	case ViewPreview:
		return m.invoice.View()
	case ViewEdit:
		return m.edit.View()
	case ViewPrint:
		if m.printPath == "" {
			return "Missing invoice print path..."

		}
		content = fmt.Sprintf("Your invoice got print at:\n%s", m.printPath)

	case ViewList:
		availHeight := m.flex.GetHeight()
		availHeight -= lipgloss.Height(m.helpContent)
		m.flex.SetHeight(availHeight)
		content = m.flex.Render()

		var sections []string
		sections = append(sections, content)
		sections = append(sections, m.helpContent)
		content = lipgloss.JoinVertical(lipgloss.Left, sections...)
	}
	return content
}

func (m *InvoicesModel) printInvoice(invContent string) string {
	dir, err := config.GetInvoicePath(m.directory)
	if err != nil {
		panic(err)
	}
	fonts, err := font.FindFonts(m.font.Family, m.font.Style)
	if err != nil {
		panic(err)
	}

	if len(fonts) == 0 {
		errMsg := fmt.Sprint(
			"font from the config could not be found in the system, font-family: ",
			m.font.Family, "font-style: ", m.font.Style)
		panic(errMsg)
	}

	htmlBytes, err := template.ToHTML(invContent)

	name := time.Now().Format("2006-01-02 15:04:05")
	mdName := name + ".md"
	htmlName := name + ".html"
	pdfName := name + ".pdf"

	if err := saveFile(dir, mdName, []byte(invContent)); err != nil {
		panic("issue while writting markdown file: " + err.Error())
	}
	if err := saveFile(dir, htmlName, htmlBytes); err != nil {
		panic("issue while writting html file: " + err.Error())
	}

	re := regexp.MustCompile(`<?.pre>`)
	pdfContent := re.ReplaceAllString(invContent, "")

	pdf.InitializePdf("")

	err = pdf.SetFont(
		m.font.Family,
		m.font.Style,
		m.font.Filepath,
		8,
	)
	if err != nil {
		panic(err)
	}

	pdf.SetText(pdfContent, 0, 4)

	path := filepath.Join(dir, pdfName)
	if err := pdf.Output(path); err != nil {
		panic("pdf output: " + err.Error())
	}
	return path
}
