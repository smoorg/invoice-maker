package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewButton(name, label string) ButtonModel {
	return ButtonModel{
		Name:  name,
		label: label,
		focus: false,
	}
}

type ButtonModel struct {
	Name  string
	label string
	focus bool
}

type ButtonSubmitEvent struct {
	ButtonName string
}

func (m *ButtonModel) Focus() {
	m.focus = true
}

func (m *ButtonModel) Blur() {
	m.focus = false
}

func (m ButtonModel) Init() {}
func (m ButtonModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			cmd = func() tea.Msg {
				return ButtonSubmitEvent{ButtonName: m.Name}
			}
		}
	}

	return cmd

}

var btnStyle = lipgloss.NewStyle().Background(lipgloss.Color("#FF6600")).Width(10)
var btnFocusStyle = lipgloss.NewStyle().Background(lipgloss.Color("#666666")).Width(10)

func (m ButtonModel) View() string {
	if m.focus {
		return btnFocusStyle.Render("Save")
	}
	return btnStyle.Render("Save")
}
