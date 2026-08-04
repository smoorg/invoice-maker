package labelinput

import (
	"invoice-maker/cmd/v2/styles"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	label      string
	labelLimit int
	input      textinput.Model

	valid   bool
	errMsg  error
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m *Model) SetValue(date string) {
	m.input.SetValue(date)
}

func (m *Model) Blur() {
	m.input.TextStyle = styles.GreyedOut
	m.input.Blur()
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func New(label string) Model {
	i := textinput.New()
	i.Width = 80
	i.TextStyle = styles.GreyedOut
	// i.PromptStyle = styles.GreyedOut
	// i.Cursor.Style = styles.Cursor
	// i.Cursor.TextStyle = styles.Cursor
	i.Prompt = ""

	return Model{
		label:      label,
		labelLimit: 20,
		input:      i,
		valid:      false,
		errMsg:     nil,
	}
}

func (m *Model) Focus() tea.Cmd {
	m.input.TextStyle = styles.ModifyInput
	return m.input.Focus()
}

func (m Model) Init() {}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	if m.input.Focused() {
		m.input.TextStyle = styles.ModifyInput
	} else {
		m.input.TextStyle = styles.GreyedOut
	}

	return m, tea.Batch(cmds...)
}

func (m Model) LabelOffset() int {
	return m.labelLimit
}

func (m Model) View() string {
	b := strings.Builder{}
	b.WriteString(m.label)

	limitLeft := m.labelLimit - len(m.label)
	for range limitLeft {
		b.WriteString(" ")
	}

	b.WriteString(": ")
	b.WriteString(m.input.View())

	if m.input.Focused() {
		b.WriteString(m.input.Cursor.View())
	}

	if m.errMsg != nil {
		b.WriteString("\n")
		b.WriteString(styles.InvalidInput.Render(m.errMsg.Error()))
	}

	return b.String()
}
