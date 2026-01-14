package main

import (
	"errors"
	"fmt"
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
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	//"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type View uint64

const (
	ViewMain View = iota
	ViewInvoices
	ViewReceivers
	ViewConfig
	ViewInvoicePreview
	ViewInvoicePrint
)

type Config struct {
	//Font             config.FontCfg   `yaml:"font"`
	//Issuer           config.Issuer    `yaml:"issuer"`
	//Receivers        []config.Company `yaml:"receivers"`
	//Invoices         []config.Invoice `yaml:"invoices"`
	//InvoiceDirectory string           `yaml:"invoiceDirectory"`
	config           config.Config
	viewList         list.Model
	view             View
	invoicesTable    table.Model
	invoicePreview   flexbox.HorizontalFlexBox
	receiversTable   table.Model
	receivers        ReceiversModel
	printContent     string
	invoicePrintPath string
}

type JumpMainView struct{}

func goMain() tea.Cmd {
	return func() tea.Msg {
		return JumpMainView{}
	}

}

type ReceiversModel struct {
	Receivers []config.Company
}

func (m ReceiversModel) Update(msg tea.Msg) (ReceiversModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.Key:
		switch msg.String() {
		case "h":
			return m, goMain()
		}
	}

	return m, nil
}

func (m ReceiversModel) Init() tea.Cmd {
	return nil
}

func (m ReceiversModel) View() string {
	msg := "receivers\n"
	for _, v := range m.Receivers {
		msg += fmt.Sprintf("%s\n", v.Name)
		msg += fmt.Sprintf("    %s\n", v.Address)
		msg += fmt.Sprintf("    %s\n", v.TaxID)
	}
	return msg
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var HelpStyle lipgloss.Style = lipgloss.NewStyle().Padding(1, 0, 0, 2).Foreground(lipgloss.Color("#666"))

type HelpFormat map[string]string

func (hp HelpFormat) RenderHelp() string {
	var content string
	for i, v := range hp {
		content += fmt.Sprintf("%s %s · ", i, v)
	}

	return HelpStyle.Render(content)
}

func main() {
	m := Config{}

	listItems := []list.Item{
		item{title: "Invoices", desc: "Show, modify and create invoices."},
		item{title: "Receivers", desc: "Show, modify and create invoice receivers."},
		item{title: "Config", desc: "Misc configs related to application."},
	}
	m.viewList = list.New(listItems, list.NewDefaultDelegate(), 10, 5)

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

	m.invoicePreview = *flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.invoicePreview.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.invoicePreview.AddColumns(columns)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}

type InitEvent struct{}

func (m Config) Init() tea.Cmd {
	m.viewList = list.New(m.viewList.Items(), list.NewDefaultDelegate(), 0, 0)

	return func() tea.Msg {
		return InitEvent{}
	}
}
func (m Config) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case JumpMainView:
		m.view = ViewMain
	case InitEvent:
		m.view = ViewMain

		// viper
		appname := "invoice-maker"
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/" + appname)
		if err := viper.ReadInConfig(); err != nil {
			return m, tea.Quit
		}

		if err := viper.UnmarshalKey("issuer", &m.config.Issuer); err != nil {
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("receivers", &m.config.Receivers); err != nil {
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("invoices", &m.config.Invoices); err != nil {
			log.Fatal(err)
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("invoiceDirectory", &m.config.InvoiceDirectory); err != nil {
			log.Fatal(err)
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("font", &m.config.Font); err != nil {
			log.Fatal(err)
			return m, tea.Quit
		}

		//viper end

		// invoices
		rows := []table.Row{}
		for _, v := range m.config.Invoices {
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
		//m.invoiceView.SetRows(rows)
		// invoices end
	case tea.WindowSizeMsg:
		m.viewList.SetSize(msg.Width, msg.Height)

		m.invoicesTable.SetWidth(msg.Width)
		m.invoicesTable.SetHeight(msg.Height)

		m.invoicePreview.SetWidth(msg.Width)
		m.invoicePreview.SetHeight(msg.Height)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyCtrlD || msg.Type == tea.KeyCtrlQ || msg.String() == "q" {
			return m, tea.Quit
		}

		switch m.view {
		case ViewMain:
			switch msg.String() {
			case "enter", "l":
				item := m.viewList.GlobalIndex()

				switch item {
				case 0:
					m.view = ViewInvoices
				case 1:
					m.view = ViewReceivers
				case 2:
					m.view = ViewConfig
				}
			}
		case ViewInvoices:
			switch msg.String() {
			case "h":
				m.view = ViewMain
			case "j":
				m.invoicesTable.MoveDown(1)
			case "k":
				m.invoicesTable.MoveUp(1)
			case "p":
				m.view = ViewInvoicePreview
			}
			cmds = append(cmds, cmd)
		case ViewInvoicePreview:
			switch msg.String() {
			case "h", "p":
				m.view = ViewInvoices
			case "P":
				picked := m.invoicesTable.SelectedRow()
				inv, _, err := GetInvoice(m.config.Invoices, picked[2], picked[4])
				if err != nil {
					panic(err)
				}

				invContent, err := template.GetContent(inv)
				m.invoicePrintPath = m.printInvoice(invContent)
				m.view = ViewInvoicePrint
			}
		case ViewInvoicePrint:
			switch msg.String() {
			case "s":
				go func(file string) {
					cmd := exec.Command("xdg-open", file)
					_, err := cmd.Output()
					if err != nil {
						log.Fatal(err)
					}
				}(m.invoicePrintPath)
			}
		default:
			switch msg.String() {
			case "h":
				m.view = ViewMain
			}
		}
	default:
	}

	switch m.view {
	case ViewMain:
		m.viewList, cmd = m.viewList.Update(msg)
		cmds = append(cmds, cmd)
	case ViewInvoices:
		m.invoicesTable, cmd = m.invoicesTable.Update(msg)
		m.invoicePreview.GetColumn(0).GetCell(0).SetContent(m.invoicesTable.View())
		cmds = append(cmds, cmd)
	case ViewReceivers:
		m.receivers.Receivers = m.config.Receivers
		m.receivers, cmd = m.receivers.Update(msg)

	}

	return m, tea.Batch(cmds...)
}

func (m Config) View() string {
	switch m.view {
	case ViewMain:
		m.viewList.Title = "Choose action"
		return m.viewList.View()
	case ViewInvoices:
		availHeight := m.invoicePreview.GetHeight()

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

		m.invoicePreview.SetHeight(availHeight)
		content := m.invoicePreview.Render()

		var sections []string
		sections = append(sections, content)
		sections = append(sections, helpContent)
		return lipgloss.JoinVertical(lipgloss.Left, sections...)
	case ViewReceivers:
		return m.receivers.View()
	case ViewInvoicePreview:
		picked := m.invoicesTable.SelectedRow()
		if picked == nil {
			return ""
		}
		content, err := GetInvoiceContent(m.config.Invoices, picked[2], picked[4])
		if err != nil {
			panic(err)
		}
		return content
	case ViewInvoicePrint:
		if m.invoicePrintPath == "" {
			return "Missing invoice print path..."

		}
		content := fmt.Sprintf("Your invoice got print at %s", m.invoicePrintPath)
		return content
	default:
		return "invalid view"
	}
}

func GetInvoice(invoices []config.Invoice, invoiceNo string, netSum string) (*config.Invoice, int, error) {
	for i, v := range invoices {
		if v.InvoiceNo == invoiceNo && v.NetSum() == netSum {
			return &v, i, nil
		}
	}

	return nil, 0, errors.New("no such invoice")
}
func GetInvoiceContent(invoices []config.Invoice, invoiceNo string, netSum string) (string, error) {
	invoice, _, err := GetInvoice(invoices, invoiceNo, netSum)
	if err != nil {
		return "", err
	}

	content, err := template.GetContent(invoice)
	if err != nil {
		return "", err
	}

	return content, nil
}

func SaveFile(dirname string, filename string, content []byte) error {
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

func (m *Config) printInvoice(invContent string) string {
	dir, err := m.config.GetInvoiceDirectory()
	if err != nil {
		panic(err)
	}
	fonts, err := font.FindFonts(m.config.Font.Family, m.config.Font.Style)
	if err != nil {
		panic(err)
	}

	if len(fonts) == 0 {
		errMsg := fmt.Sprint(
			"font from the config could not be found in the system, font-family: ",
			m.config.Font.Family, "font-style: ", m.config.Font.Style)
		panic(errMsg)
	}

	htmlBytes, err := template.ToHTML(invContent)

	name := time.Now().Format("2006-01-02 15:04:05")
	mdName := name + ".md"
	htmlName := name + ".html"
	pdfName := name + ".pdf"

	if err := SaveFile(dir, mdName, []byte(invContent)); err != nil {
		panic("issue while writting markdown file: " + err.Error())
	}
	if err := SaveFile(dir, htmlName, htmlBytes); err != nil {
		panic("issue while writting html file: " + err.Error())
	}

	re := regexp.MustCompile(`<?.pre>`)
	pdfContent := re.ReplaceAllString(invContent, "")

	pdf.InitializePdf("")

	err = pdf.SetFont(
		m.config.Font.Family,
		m.config.Font.Style,
		m.config.Font.Filepath,
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
