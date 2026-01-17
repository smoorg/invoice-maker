package main

import (
	"fmt"
	"invoice-maker/cmd/v2/invoices"
	"invoice-maker/cmd/v2/receiver"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	pkg_help "invoice-maker/pkg/help"
	"invoice-maker/pkg/view"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/viper"
)

type SwitchViewEvent struct {
	item int
}
type Config struct {
	config   config.Config
	viewList list.Model
	view     view.View

	invoiceModel invoices.InvoicesModel
	receivers    receiver.ReceiversModel
	keys         keymap
	helpContent  string
}

type keymap struct {
	Quit key.Binding
	Next key.Binding
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var HelpStyle lipgloss.Style = lipgloss.NewStyle().Padding(1, 0, 0, 2).Foreground(lipgloss.Color("#666"))

func main() {
	m := Config{}

	listItems := []list.Item{
		item{title: "Invoices", desc: "Show, modify and create invoices."},
		item{title: "Receivers", desc: "Show, modify and create invoice receivers."},
		item{title: "Config", desc: "Misc configs related to application."},
	}
	m.viewList = list.New(listItems, list.NewDefaultDelegate(), 10, 5)
	m.viewList.SetShowHelp(false)
	m.viewList.SetShowFilter(false)

	m.invoiceModel = invoices.New(m.config)
	m.receivers = receiver.New()

	m.keys = keymap{
		Quit: key.NewBinding(
			key.WithKeys(
				tea.KeyCtrlQ.String(),
				tea.KeyCtrlC.String(),
				tea.KeyCtrlD.String(),
				"q",
			),
			key.WithHelp("^C/^D/^Q/q", "quit"),
		),
		Next: key.NewBinding(
			key.WithKeys(
				"l",
				tea.KeyEnter.String(),
				tea.KeyRight.String(),
			),
			key.WithHelp("→/l/enter", "next view"),
		),
	}

	helpBubble := help.New()
	m.helpContent = helpBubble.ShortHelpView(pkg_help.MapToBindingsList(m.keys))

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
	case SwitchViewEvent:
		switch msg.item {
		case 0:
			m.view = view.ViewInvoices
		case 1:
			m.view = view.ViewReceivers
		case 2:
			m.view = view.ViewConfig
		}
	case pkg.JumpMainView:
		m.view = view.ViewMain
	case InitEvent:
		m.view = view.ViewMain

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
		m.invoiceModel.SetConfig(m.config)
		m.invoiceModel.SetRows(rows)

		m.receivers.SetConfig(m.config.Receivers)

	case tea.WindowSizeMsg:
		m.viewList.SetSize(msg.Width, msg.Height-1)
		m.invoiceModel.SetSize(msg.Width, msg.Height)
		m.receivers.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

		switch m.view {
		case view.ViewMain:
			switch {
			case key.Matches(msg, m.keys.Next):
				item := m.viewList.GlobalIndex()
				cmd = func() tea.Msg {
					return SwitchViewEvent{item: item}
				}
				cmds = append(cmds, cmd)
			}
		default:
		}
	default:
	}

	// second switch for moved out stuff
	switch m.view {
	case view.ViewMain:
		m.viewList, cmd = m.viewList.Update(msg)
		cmds = append(cmds, cmd)
	case view.ViewInvoices:
		m.invoiceModel, cmd = m.invoiceModel.Update(msg)
		cmds = append(cmds, cmd)
	case view.ViewReceivers:
		m.receivers, cmd = m.receivers.Update(msg)
		cmds = append(cmds, cmd)

	}

	return m, tea.Batch(cmds...)
}

func (m Config) View() string {
	switch m.view {
	case view.ViewMain:
		m.viewList.Title = "Choose action"
		content := m.viewList.View() + "\n" + m.helpContent
		return content
	case view.ViewInvoices:
		return m.invoiceModel.View()
	case view.ViewReceivers:
		return m.receivers.View()
	default:
		return "invalid view"
	}
}
