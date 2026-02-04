package edit

import (
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Back key.Binding
}

type ReceiverEdit struct {
	receiver config.Company
	keys     keyMap
}

func (m *ReceiverEdit) SetReceiver(v config.Company) {
	m.receiver = v
}

func New() ReceiverEdit {
	m := ReceiverEdit{}

	m.keys = keyMap{
		Back: key.NewBinding(
			key.WithKeys("h", tea.KeyLeft.String()),
			key.WithHelp("←/h", "go back"),
		),
	}

	return m
}

func (m ReceiverEdit) Init() {}

func (m ReceiverEdit) Update(msg tea.Msg) (ReceiverEdit, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			cmds = append(cmds, pkg.GoReceivers())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ReceiverEdit) View() string {
	content := "receiver edit" + m.receiver.Name

	return content
}
