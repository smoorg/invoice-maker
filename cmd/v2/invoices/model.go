package invoices

import (
	"errors"
	"fmt"
	"log"
	"os"
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
	ViewMain InvoiceView = iota
	ViewPreview
	ViewPrint
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
	directory string
	font      config.FontCfg

	table table.Model
	flex  flexbox.HorizontalFlexBox

	helpContent  string
	printContent string
	printPath    string
	view         InvoiceView
	keys         KeyMap
}

func (m *InvoicesModel) SetConfig(cfg config.Config) {
	m.invoices = cfg.Invoices
	m.directory = cfg.InvoiceDirectory
	m.font = cfg.Font
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

	m.flex = *flexbox.NewHorizontal(0, 0)
	columns := []*flexbox.Column{
		m.flex.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.flex.AddColumns(columns)
	m.directory = config.InvoiceDirectory
	m.font = config.Font

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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			switch m.view {
			case ViewMain:
				cmd = pkg.GoMain()
				cmds = append(cmds, cmd)
			case ViewPreview:
				m.view = ViewMain
			case ViewPrint:
				m.view = ViewMain
			}
		case key.Matches(msg, m.keys.Down):
			m.table.MoveDown(1)
		case key.Matches(msg, m.keys.Up):
			m.table.MoveUp(1)
		case key.Matches(msg, m.keys.Next):
			switch m.view {
			case ViewMain:
				m.view = ViewPreview
			}
		case key.Matches(msg, m.keys.Print):
			m.view = ViewPrint

			picked := m.table.SelectedRow()
			invContent, err := getInvoiceContent(m.invoices, picked[2], picked[4])
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
		case key.Matches(msg, m.keys.Edit):
			// TODO: implement
			panic("unimplemented")
		}
	}

	switch m.view {
	case ViewMain:
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
			}
		}
	case ViewPreview:
		m.invoice.Update(msg)
	case ViewPrint:
	}

	return m, tea.Batch(cmds...)
}
func (m InvoicesModel) View() string {
	content := ""
	switch m.view {
	case ViewPrint:
		if m.printPath == "" {
			return "Missing invoice print path..."

		}
		content = fmt.Sprintf("Your invoice got print at:\n%s", m.printPath)

		buttonStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

		activeButtonStyle := buttonStyle.
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#F25D94")).
			MarginRight(2).
			Underline(true)
		//normal := lipgloss.Color("#EEEEEE")
		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(content)
		okButton := activeButtonStyle.Render("Yes")
		cancelButton := buttonStyle.Render("Maybe")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
		dialogBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		//base := lipgloss.NewStyle().Foreground(normal)
		dialog := lipgloss.Place(150, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		)
		content = dialog
	case ViewPreview:
		return m.invoice.View()
	case ViewMain:
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
