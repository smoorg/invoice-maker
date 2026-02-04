package receiver

import (
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	pkg_help "invoice-maker/pkg/help"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up   key.Binding
	Down key.Binding
	Back key.Binding
	Edit key.Binding
}

type view uint8

type ReceiversModel struct {
	receivers []config.Company `yaml:"receivers"`
	table     table.Model

	keyBindings keyMap

	flex        flexbox.HorizontalFlexBox
	helpContent string
}

func New() ReceiversModel {
	m := ReceiversModel{}

	m.table = table.New(table.WithColumns([]table.Column{
		{Title: "Name", Width: 30},
		{Title: "Address", Width: 45},
		{Title: "Tax", Width: 18},
	}))

	m.flex = *flexbox.NewHorizontal(0, 0)

	columns := []*flexbox.Column{
		m.flex.NewColumn().AddCells(
			flexbox.NewCell(1, 1),
		),
	}
	m.flex.AddColumns(columns)

	m.keyBindings = keyMap{
		Back: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "go back"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", tea.KeyUp.String()),
			key.WithHelp("k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", tea.KeyDown.String()),
			key.WithHelp("j", "down"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
	}

	helpBubble := help.New()
	m.helpContent = helpBubble.ShortHelpView(pkg_help.MapToBindingsList(m.keyBindings))

	return m
}

func (m *ReceiversModel) SetSize(width, height int) {
	m.flex.SetWidth(width)
	m.flex.SetHeight(height)

	m.table.SetWidth(width)
	m.table.SetHeight(height)
}

func (m *ReceiversModel) SetConfig(receivers []config.Company) {
	m.receivers = receivers

	rows := []table.Row{}
	for _, v := range m.receivers {
		rows = append(rows, table.Row{
			v.Name,
			v.Address,
			v.TaxID,
		})
	}
	m.table.SetRows(rows)
}

func (m ReceiversModel) Init() tea.Cmd {
	return nil
}

func (m ReceiversModel) Update(msg tea.Msg) (ReceiversModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height)
		m.flex.SetWidth(msg.Width)
		m.flex.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyBindings.Back):
			return m, pkg.GoMain()
		case key.Matches(msg, m.keyBindings.Down):
			m.table.MoveDown(1)
		case key.Matches(msg, m.keyBindings.Up):
			m.table.MoveUp(1)
		case key.Matches(msg, m.keyBindings.Edit):
			row := m.table.SelectedRow()
			for _, v := range m.receivers {
				if v.Name == row[0] && v.Address == row[1] &&
					v.TaxID == row[2] {
					cmds = append(cmds, pkg.GoReceiverEdit(v))
				}
			}
		}

	}
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	m.flex.GetColumn(0).GetCell(0).SetContent(m.table.View())

	return m, tea.Batch(cmds...)
}

func (m ReceiversModel) View() string {
	var content string
	var sections []string
	availHeight := m.flex.GetHeight()
	availHeight -= lipgloss.Height(m.helpContent)
	m.flex.SetHeight(availHeight)

	content = m.flex.Render()
	sections = append(sections, content)
	sections = append(sections, m.helpContent)
	content = lipgloss.JoinVertical(lipgloss.Left, sections...)

	return content
}
