package main

import (
	"fmt"
	"invoice-maker/pkg/config"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type View uint64

const (
	ViewMain View = iota
	ViewInvoices
	ViewReceivers
	ViewConfig
)

type Config struct {
	//Font             config.FontCfg   `yaml:"font"`
	Issuer           config.Issuer    `yaml:"issuer"`
	Receivers        []config.Company `yaml:"receivers"`
	Invoices         []config.Invoice `yaml:"invoices"`
	InvoiceDirectory string           `yaml:"invoiceDirectory"`
	viewList         list.Model
	view             View
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func main() {
	c := Config{}

	listItems := []list.Item{
		item{title: "Invoices", desc: "Show, modify and create invoices."},
		item{title: "Receivers", desc: "Show, modify and create invoice receivers."},
		item{title: "Config", desc: "Misc configs related to application."},
	}
	c.viewList = list.New(listItems, list.NewDefaultDelegate(), 10, 5)

	p := tea.NewProgram(c, tea.WithAltScreen())
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

	switch msg := msg.(type) {
	case InitEvent:
		appname := "invoice-maker"
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/" + appname)

		if err := viper.ReadInConfig(); err != nil {
			return m, tea.Quit
		}

		if err := viper.UnmarshalKey("issuer", &m.Issuer); err != nil {
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("receivers", &m.Receivers); err != nil {
			return m, tea.Quit
		}
		if err := viper.UnmarshalKey("invoices", &m.Invoices); err != nil {
			log.Fatal(err)
			return m, tea.Quit
		}

		m.view = ViewMain
	case tea.WindowSizeMsg:
		m.viewList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC ||
			msg.Type == tea.KeyCtrlD ||
			msg.Type == tea.KeyCtrlQ ||
			msg.String() == "q" {
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
			}
		case ViewReceivers:
			switch msg.String() {
			case "h":
				m.view = ViewMain
			}
		default:
			log.Fatal("invalid view", m.view)
		}
	default:
	}

	m.viewList, cmd = m.viewList.Update(msg)

	return m, cmd
}

func (m Config) View() string {
	switch m.view {
	case ViewMain:
		m.viewList.Title = "Choose action"
		return m.viewList.View()
	case ViewInvoices:
		msg := "invoices\n"
		for _, v := range m.Invoices {
			msg += fmt.Sprintf("Number: %s/n" + v.InvoiceNo)
		}
		return msg
	case ViewReceivers:
		msg := "receivers\n"
		for _, v := range m.Receivers {
			msg += fmt.Sprintf("%s\n", v.Name)
			msg += fmt.Sprintf("    %s\n", v.Address)
			msg += fmt.Sprintf("    %s\n", v.TaxID)
		}
		return msg
	default:
		return "invalid view"
	}
}
