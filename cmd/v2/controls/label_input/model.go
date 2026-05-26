package labelinput

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	label      string
	labelLimit int
	Input      textinput.Model

	isValid func(v string) (bool, error)
	valid   bool
	errMsg  error
}

func (m *Model) Blur() {
	m.Input.Blur()
	m.Input.TextStyle = greyedOut
}

func (m Model) Focused() bool {
	return m.Input.Focused()
}

func New(label string, limit int) Model {
	i := textinput.New()
	i.Width = 80
	i.TextStyle = greyedOut
	i.PromptStyle = greyedOut
	i.Prompt = ""

	return Model{
		label:      label,
		labelLimit: limit,
		Input:      i,
		isValid:    func(v string) (bool, error) { return true, nil },
		valid:      false,
		errMsg:     nil,
	}
}

func (m *Model) SetValidation(isValid func(val string) (bool, error)) {
	m.isValid = isValid
}

func (m *Model) Focus() tea.Cmd {
	m.Input.TextStyle = modifyInput
	return m.Input.Focus()
}

func (m Model) Init() {}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	m.valid, m.errMsg = m.isValid(m.Input.Value())
	if m.errMsg != nil {
		panic(m.errMsg)
	}
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

var greyedOut = lipgloss.NewStyle().Foreground(grey).Background(lipgloss.NoColor{})
var modifyInput = lipgloss.NewStyle().Foreground(white).Background(bg)
var invalidInput = lipgloss.NewStyle().Foreground(red)

func (m Model) View() string {
	b := strings.Builder{}
	b.WriteString(m.label)

	limitLeft := m.labelLimit - len(m.label)
	for range limitLeft {
		b.WriteString(" ")
	}

	b.WriteString(": ")
	b.WriteString(m.Input.View())

	if m.Input.Focused() {
		b.WriteString(m.Input.Cursor.View())
	}

	if m.errMsg != nil {
		b.WriteString("\n")
		b.WriteString(invalidInput.Render(m.errMsg.Error()))
	}

	return b.String()
}
