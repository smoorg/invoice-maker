package configview

import (
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ConfigModel struct {
	InvoiceDirectory string
	FontFamily       string
	FontStyle        string
	keys             keyMap
	inputs           []textinput.Model
	focusedInput     int
}

func (m ConfigModel) SetSize(width int, height int) {
}

type keyMap struct {
	NextField key.Binding
	Back      key.Binding
}

const (
	InputFontStyle int = iota
	InputFontFamily
	InputDirectory
)

func New(config config.Config) ConfigModel {
	//inputs := make([]textinput.Model, 3)
	//inputs[InputFontStyle] = textinput.New()
	//inputs[InputFontStyle].Placeholder = "Font Style"
	//inputs[InputFontStyle].CharLimit = 40
	//inputs[InputFontFamily] = textinput.Model{Placeholder: "Font Family"}
	//inputs[InputDirectory] = textinput.Model{Placeholder: "Config directory"}

	keymap := keyMap{
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
	}

	return ConfigModel{
		FontFamily:       config.Font.Family,
		FontStyle:        config.Font.Style,
		InvoiceDirectory: config.InvoiceDirectory,
		keys:             keymap,
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return m.inputs[InputFontStyle].Focus()
}

func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.NextField):
		case key.Matches(msg, m.keys.Back):
			cmd = pkg.GoMain()
		}
	}
	return m, cmd
}

func (m ConfigModel) View() string {
	b := strings.Builder{}
	b.WriteString("Family: ")
	b.WriteString(m.FontFamily)
	b.WriteString("\n")

	b.WriteString("Style: ")
	b.WriteString(m.FontStyle)
	b.WriteString("\n")

	b.WriteString("Invoice' Directory: ")
	b.WriteString(m.InvoiceDirectory)
	b.WriteString("\n")

	return b.String()
}
