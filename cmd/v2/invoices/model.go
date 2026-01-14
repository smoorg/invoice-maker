package invoices

import (
	"errors"
	"fmt"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/font"
	"invoice-maker/pkg/pdf"
	"invoice-maker/pkg/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InvoiceView uint64

const (
	ViewMain InvoiceView = iota
	ViewPreview
	ViewPrint
)

type InvoicesModel struct {
	invoices         []config.Invoice `yaml:"invoices"`
	table            table.Model
	InvoiceViewFlex  flexbox.HorizontalFlexBox
	invoiceDirectory string
	font             config.FontCfg

	printContent     string
	invoicePrintPath string
	view             InvoiceView
}

func (m *InvoicesModel) SetConfig(cfg config.Config) {
	m.invoices = cfg.Invoices
	m.invoiceDirectory = cfg.InvoiceDirectory
	m.font = cfg.Font
}

func (m *InvoicesModel) SetRows(rows []table.Row) {
	m.table.SetRows(rows)
}

func (m *InvoicesModel) SetSize(width int, height int) {
	m.InvoiceViewFlex.SetWidth(width)
	m.InvoiceViewFlex.SetHeight(height)
	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func New(config config.Config) InvoicesModel {
	m := InvoicesModel{}
	m.table = table.New(
		table.WithColumns([]table.Column{
			{Title: "Date", Width: 10},
			{Title: "Payment Date", Width: 10},
			{Title: "Invoice No.", Width: 11},
			{Title: "Receiver", Width: 20},
			{Title: "Net", Width: 8},
			{Title: "Gross", Width: 8},
		}),
	)
	m.table.SetHeight(5)
	m.table.SetWidth(100)

	m.InvoiceViewFlex = *flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.InvoiceViewFlex.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.InvoiceViewFlex.AddColumns(columns)
	m.invoiceDirectory = config.InvoiceDirectory
	m.font = config.Font

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
		m.table.SetRows(rows)
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height)
		m.InvoiceViewFlex.SetWidth(msg.Width)
		m.InvoiceViewFlex.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			switch m.view {
			case ViewMain:
				cmd = pkg.GoMain()
			case ViewPreview:
				m.view = ViewMain
			case ViewPrint:
				m.view = ViewMain
			}
		case "j":
			m.table.MoveDown(1)
		case "k":
			m.table.MoveUp(1)
		case "l":
			switch m.view {
			case ViewMain:
				m.view = ViewPreview
			}
		case "p":
			m.view = ViewPrint

			picked := m.table.SelectedRow()
			invContent, err := getInvoiceContent(m.invoices, picked[2], picked[4])
			if err != nil {
				panic(err)
			}
			m.invoicePrintPath = m.printInvoice(invContent)

			go func(file string) {
				cmd := exec.Command("xdg-open", file)
				_, err := cmd.Output()
				if err != nil {
					log.Fatal(err)
				}
			}(m.invoicePrintPath)
		}
	}

	m.table, cmd = m.table.Update(msg)
	m.InvoiceViewFlex.GetColumn(0).GetCell(0).SetContent(m.table.View())

	return m, cmd
}
func (m InvoicesModel) View() string {
	content := ""
	switch m.view {
	case ViewPrint:
		if m.invoicePrintPath == "" {
			return "Missing invoice print path..."

		}
		content = fmt.Sprintf("Your invoice got print at %s", m.invoicePrintPath)
	case ViewPreview:
		picked := m.table.SelectedRow()
		if picked == nil {
			return ""
		}
		for _, v := range m.invoices {
			if v.InvoiceNo == picked[2] && v.NetSum() == picked[4] {
				c, err := template.GetContent(&v)
				if err != nil {
					panic(err)
				}
				content = c
			}
		}
	case ViewMain:
		availHeight := m.InvoiceViewFlex.GetHeight()

		helpBubble := help.New()
		kb := []key.Binding{
			key.NewBinding(
				key.WithKeys("l"),
				key.WithHelp("p", "preview"),
			),
			key.NewBinding(
				key.WithKeys("p"),
				key.WithHelp("p", "print"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit"),
			),
		}
		helpContent := helpBubble.ShortHelpView(kb)
		availHeight -= lipgloss.Height(helpContent)

		m.InvoiceViewFlex.SetHeight(availHeight)
		content = m.InvoiceViewFlex.Render()

		var sections []string
		sections = append(sections, content)
		sections = append(sections, helpContent)
		content = lipgloss.JoinVertical(lipgloss.Left, sections...)
	}
	return content
}

func (m *InvoicesModel) printInvoice(invContent string) string {
	dir, err := config.GetInvoicePath(m.invoiceDirectory)
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

func getInvoice(invoices []config.Invoice, invoiceNo string, netSum string) (*config.Invoice, int, error) {
	for i, v := range invoices {
		if v.InvoiceNo == invoiceNo && v.NetSum() == netSum {
			return &v, i, nil
		}
	}

	return nil, 0, errors.New("no such invoice")
}
func getInvoiceContent(invoices []config.Invoice, invoiceNo string, netSum string) (string, error) {
	invoice, _, err := getInvoice(invoices, invoiceNo, netSum)
	if err != nil {
		return "", err
	}

	content, err := template.GetContent(invoice)
	if err != nil {
		return "", err
	}

	return content, nil
}

func saveFile(dirname string, filename string, content []byte) error {
	if err := os.MkdirAll(dirname, 0744); err != nil {
		return err
	}

	mddir := filepath.Join(dirname, filename)

	file, err := os.Create(mddir)
	if err != nil {
		return err
	}
	if _, err := file.Write(content); err != nil {
		log.Fatal("write string err", err)
		return err
	}
	return nil
}
