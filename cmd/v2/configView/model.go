package configview

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	labelinput "invoice-maker/cmd/v2/controls/label_input"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/font"
	pkg_help "invoice-maker/pkg/help"
)

type FocID int

func (m *ConfigModel) Increment() {
	if m.focID >= (len(m.inputs) - 1) {
		m.focID = 0
	} else {
		m.focID++
	}
}

type ConfigModel struct {
	FontConfig config.FontCfg
	keys       keyMap
	inputs     []*labelinput.Model
	focID      int
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
	inputs := make([]*labelinput.Model, 0, 3)

	family := labelinput.New("Font Family")
	family.Input.SetValue(config.Config.Family)
	family.Focus()
	family.SetValidation(func(val string) (bool, string) {
		fonts, err := font.GetFontFamilies()
		if err != nil {
			return false, "unable to get fonts"
		}

		if val == "" {
			return false, "Font Family has to be set!"
		}

		fontExists := font.HasFontFamily(fonts, val)
		if !fontExists {
			return false, "No such font!"
		}

		return true, ""
	})

	//style := labelinput.New("Font Style")
	//style.Input.SetValue(config.Config.Style)

	dir := labelinput.New("Invoice Directory")
	dir.Input.SetValue(config.Config.InvoiceDirectory)

	//inputs = append(inputs, family, style, dir)
	inputs = append(inputs, &family, &dir)

	keymap := keyMap{
		NextField: key.NewBinding(
			key.WithKeys(tea.KeyTab.String()),
			key.WithHelp("tab", "select next field"),
		),
		Back: key.NewBinding(
			key.WithKeys(
				tea.KeyEsc.String(),
			),
			key.WithHelp("esc", "go back"),
		),
	}

	return ConfigModel{
		FontConfig: config.Config,
		keys:       keymap,
		inputs:     inputs,
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return nil
}

func (m *ConfigModel) SwitchInput() tea.Cmd {
	m.inputs[m.focID].Blur()
	m.Increment()
	return m.inputs[m.focID].Focus()
}

func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("updating config model", msg.String())
		switch {
		case key.Matches(msg, m.keys.NextField):
			log.Println("updating config next field", m.focID)
			cmd = m.SwitchInput()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keys.Back):
			m.inputs[m.focID].Blur()
			cmd = pkg.GoMain()
			cmds = append(cmds, cmd)
		default:
		}
	}
	for _, input := range m.inputs {
		*input, cmd = input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m ConfigModel) View() string {
	b := strings.Builder{}
	const nl string = "\n"

	for _, v := range m.inputs {
		b.WriteString(v.View())
		b.WriteString(nl)
	}

	return b.String()
}
