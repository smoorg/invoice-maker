package datemodel

import (
	"fmt"
	_ "strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	labelinput "invoice-maker/cmd/v2/controls/label_input"
	_ "invoice-maker/cmd/v2/styles"
)

type DateInput struct {
	time   time.Time
	input  labelinput.Model
	errMsg error
	layout string
}

func (m DateInput) Focus() tea.Cmd {
	return m.input.Focus()
}

func (m DateInput) Blur() tea.Cmd {
	m.input.Blur()
	return nil
}

func (m *DateInput) Value() string {
	return m.Value()
	//return m.time.Format(m.layout)
}

func (m DateInput) SetValue(date string) {
	m.input.SetValue(date)
	m.errMsg = m.Validate(date)

}

func (m DateInput) Validate(date string) error {
	_, err := time.Parse(m.layout, date)
	if err != nil {
		return fmt.Errorf("%s %s", "Your date is invalid according to layout", m.layout)
	}

	return nil
}

func NewDateInput(label string, limit int, layout string) DateInput {
	m := DateInput{
		layout: layout,
		input:  labelinput.New(label),
	}

	return m
}

func (m DateInput) Init() {}

func (m DateInput) Update(msg tea.Msg) (DateInput, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// switch msg := msg.(type) {
	// case tea.KeyMsg:
		// runes without ctrl/meta/alt
		// if msg.Type == tea.KeyBackspace || (msg.Type == tea.KeyRunes && msg.Alt == false) {
		m.input, cmd = m.input.Update(msg)
	 	cmds = append(cmds, cmd)
	// 	// m.errMsg = m.Validate(m.input.Value())
	// 	// }
	// }

	return m, tea.Batch(cmds...)
}
func (m DateInput) View() string {
	return m.input.View()
	//value := m.input.View()

	//content := strings.Builder{}
	//content.WriteString(value)

	//if m.errMsg != nil {
	//	content.WriteString("\n")
	//	for range m.input.LabelOffset() {
	//		content.WriteString(" ")
	//	}

	//	content.WriteString(
	//		styles.InvalidInput.Render(
	//			m.errMsg.Error(),
	//		),
	//	)
	//}
	//return content.String()
}
