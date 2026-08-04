package edit

import (
	labelinput "invoice-maker/cmd/v2/controls/label_input"
	"invoice-maker/cmd/v2/styles"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	pkg_help "invoice-maker/pkg/help"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"strings"
)

func New() ReceiverEdit {
	m := ReceiverEdit{}

	m.keys = keyMap{
		Back: key.NewBinding(
			key.WithKeys("h", tea.KeyLeft.String()),
			key.WithHelp("←/h", "go back"),
		),
		Esc: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String()),
			key.WithHelp("esc", "blur focus/exit"),
		),
		Next: key.NewBinding(
			key.WithKeys("tab", tea.KeyTab.String()),
			key.WithHelp("tab", "focus next input"),
		),
		Prev: key.NewBinding(
			key.WithKeys(tea.KeyShiftTab.String()),
			key.WithHelp("shift+tab", "focus prev input"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save"),
		),
	}
	m.helpContent = help.New().ShortHelpView(pkg_help.MapToBindingsList(m.keys))

	m.inputName = labelinput.New("Name")
	m.inputAddress = labelinput.New("Address")
	m.inputTaxID = labelinput.New("Tax ID")

	return m
}

func (m ReceiverEdit) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m ReceiverEdit) SetModel(c *config.Company) {
	m.receiver = *c
	m.inputName.SetValue(m.receiver.Name)
	m.inputAddress.SetValue(m.receiver.Address)
	m.inputTaxID.SetValue(m.receiver.TaxID)
}

func (m ReceiverEdit) Update(msg tea.Msg) (ReceiverEdit, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case pkg.JumpReceiverEdit:
		m.receiver = msg.Receiver
		m.inputName.SetValue(msg.Receiver.Name)
		m.inputAddress.SetValue(msg.Receiver.Address)
		m.inputTaxID.SetValue(msg.Receiver.TaxID)
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			cmds = append(cmds, pkg.GoReceivers())
		case key.Matches(msg, m.keys.Esc):
			if m.focus == 0 {
				cmds = append(cmds, pkg.GoReceivers())
				break
			}
			m.Blur()
			m.focus = 0
		case key.Matches(msg, m.keys.Next):
			m.Blur()
			m.focus = IncrementTillLimit(m.focus, FocusTaxID, 1, FocusNone)
			m.Focus(m.focus)
		case key.Matches(msg, m.keys.Prev):
			m.Blur()
			m.focus = IncrementTillLimit(m.focus, FocusNone, -1, FocusTaxID)
			m.Focus(m.focus)
		case key.Matches(msg, m.keys.Save):
			c := config.Company{
				Name:    m.inputName.Value(),
				Address: m.inputAddress.Value(),
				TaxID:   m.inputTaxID.Value(),
			}
			return m, pkg.UpdateReceiver(m.receiver, c)
		}
	}

	switch m.focus {
	case FocusName:
		m.inputName, cmd = m.inputName.Update(msg)
		cmds = append(cmds, cmd)
	case FocusAddress:
		m.inputAddress, cmd = m.inputAddress.Update(msg)
		cmds = append(cmds, cmd)
	case FocusTaxID:
		m.inputTaxID, cmd = m.inputTaxID.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ReceiverEdit) View() string {
	content := strings.Builder{}
	s := styles.MarginStyle

	ss := []string{
		s.Render(m.inputName.View()),
		s.Render(m.inputAddress.View()),
		s.Render(m.inputTaxID.View(), "\n"),
	}

	log.Println("height", m.height)
	// lines to fill so help lands on the last line
	// Margin style takes additional line for each field
	linesLeft := m.height - len(ss)*2 - 1
	log.Println("lines left", linesLeft)
	for range linesLeft {
		ss = append(ss, "\n")
	}

	for _, v := range ss {
		content.WriteString(v)
	}
	content.WriteString(m.helpContent)

	return content.String()
}

func IncrementTillLimit(curVal, limit, inc, defaultVal int) int {
	if curVal == limit {
		return defaultVal
	} else {
		return curVal + inc
	}
}
