package labelinput

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	label string
	Input textinput.Model

	isValid func(v string) (bool, string)
	valid   bool
	errMsg  string
}

func (m *Model) Blur() {
	m.Input.Blur()
}

func New(label string) Model {
	i := textinput.New()
	i.Width = 80
	i.TextStyle = greyedOut
	i.PromptStyle = greyedOut

	return Model{
		label:   label,
		Input:   i,
		isValid: func(v string) (bool, string) { return true, "" },
	}
}

func (m *Model) SetValidation(isValid func(val string) (bool, string)) {
	m.isValid = isValid
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

	m.valid, m.errMsg = m.isValid(m.Input.Value())
	if m.Input.Focused() {
		m.Input.TextStyle = modifyInput
	} else {
		m.Input.TextStyle = greyedOut
	}

	return m, tea.Batch(cmds...)
}

var bg = lipgloss.Color("#333")
var red = lipgloss.Color("#F54927")
var white = lipgloss.Color("#fff")
var grey = lipgloss.Color("#555")

var greyedOut = lipgloss.NewStyle().Foreground(grey).Background(bg)
var modifyInput = lipgloss.NewStyle().Foreground(white).Background(bg)
var invalidInput = lipgloss.NewStyle().Foreground(red)

func (m Model) View() string {
	b := strings.Builder{}
	b.WriteString(m.label)
	b.WriteString(":\n")
	b.WriteString(m.Input.View())

	if m.Input.Focused() {
		b.WriteString(m.Input.Cursor.View())
	}

	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(invalidInput.Render(m.errMsg))
	}

	return b.String()
}
