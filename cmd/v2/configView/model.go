package configview

import (
	"invoice-maker/pkg"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfigModel struct {
	InvoiceDirectory string
	FontFamily       string
	FontStyle        string
	keys             keyMap
}

type keyMap struct {
	NextField key.Binding
	Back      key.Binding
}

func NewConfigModel() *ConfigModel {
	return &ConfigModel{
		keys: keyMap{
			NextField: key.NewBinding(
				key.WithKeys(tea.KeyTab.String()),
				key.WithHelp("tab", "select next field"),
			),
			Back: key.NewBinding(
				key.WithKeys(
					tea.KeyCtrlQ.String(),
					tea.KeyCtrlC.String(),
					tea.KeyCtrlD.String(),
					"q",
				),
				key.WithHelp("tab", "select next field"),
			),
		},
	}
}

func (m ConfigModel) Init() tea.Cmd { return nil }

func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.NextField):
		case key.Matches(msg, m.keys.Back):
			return m, pkg.GoMain()
		}
	}
	return m, cmd
}

func (m ConfigModel) View() string {
	return "config view"
}
