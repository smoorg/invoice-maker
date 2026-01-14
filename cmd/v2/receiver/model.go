package receiver

import (
	"fmt"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

type ReceiversModel struct {
	Receivers []config.Company `yaml:"receivers"`
}

func (m ReceiversModel) Update(msg tea.Msg) (ReceiversModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			return m, pkg.GoMain()
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

