package labelinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	label string
	Input textinput.Model
}

func (m *Model) Blur() {
	m.Input.Blur()
}

func New(label string) Model {
	i := textinput.New()
	i.Width = 80
	i.TextStyle = gray
	i.PromptStyle = gray

	return Model{
		label: label,
		Input: i,
	}
}

func (m *Model) Focus() tea.Cmd {
	return m.Input.Focus()
}

func (m Model) Init() {}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)
	if m.Input.Focused() {
		m.Input.TextStyle = purple
		cmds = append(cmds, cmd)
	} else {
		m.Input.TextStyle = gray
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

var gray = lipgloss.NewStyle().Foreground(lipgloss.Color("#555")).Background(lipgloss.Color("#333"))
var purple = lipgloss.NewStyle().Foreground(lipgloss.Color("#fff")).Background(lipgloss.Color("#333"))

func (m Model) View() string {
	cursor := m.Input.Cursor.View()
	return fmt.Sprintf("%s:\n%s%s", m.label, m.Input.View(), cursor)
	//return m.Input.View()
}
