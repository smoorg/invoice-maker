package configview

import (
	"log"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	labelinput "invoice-maker/cmd/v2/controls/label_input"
	singleselect "invoice-maker/cmd/v2/controls/single_select"
	"invoice-maker/pkg"
	"invoice-maker/pkg/config"
	"invoice-maker/pkg/font"
	pkg_help "invoice-maker/pkg/help"
)

const (
	InputFontFamily = iota
	InputInvoiceDir
)

type ConfigModel struct {
	FontConfig       config.FontCfg
	keys             keyMap
	focusID          int
	helpContent      string
	FontFamily       singleselect.Model
	InvoiceDirectory labelinput.Model
}

func (m ConfigModel) SetSize(width int, height int) {
}

type keyMap struct {
	NextField key.Binding
	Esc       key.Binding
	Back      key.Binding
}

func New(cfg config.Config) ConfigModel {
	// Build font family select items
	fontFamilies, err := font.GetFontFamilies()
	if err != nil {
		fontFamilies = []string{}
	}

	// Move current font to top of list
	currentIdx := slices.Index(fontFamilies, cfg.Config.Family)
	log.Printf("Font from config: %q, found at index: %d", cfg.Config.Family, currentIdx)
	if currentIdx > 0 {
		current := fontFamilies[currentIdx]
		fontFamilies = slices.Delete(fontFamilies, currentIdx, currentIdx+1)
		fontFamilies = slices.Insert(fontFamilies, 0, current)
	}
	log.Printf("First font in list after reorder: %q", fontFamilies[0])

	items := make([]singleselect.Item, len(fontFamilies))
	for i, f := range fontFamilies {
		items[i] = singleselect.Item{Label: f, Value: f}
	}

	fontFamily := singleselect.New("Font Family", 20, items)

	// Invoice directory input
	invoiceDir := labelinput.New("Invoice Directory", 20)
	invoiceDir.SetValue(cfg.Config.InvoiceDirectory)

	keymap := keyMap{
		NextField: key.NewBinding(
			key.WithKeys(tea.KeyTab.String()),
			key.WithHelp("tab", "select next field"),
		),
		Esc: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String()),
			key.WithHelp("Esc", "blur focus"),
		),
		Back: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "go to main view (when nothing focused)"),
		),
	}

	helpBubble := help.New()
	helpView := helpBubble.ShortHelpView(pkg_help.MapToBindingsList(keymap))

	return ConfigModel{
		FontConfig:       cfg.Config,
		keys:             keymap,
		focusID:          InputFontFamily,
		helpContent:      helpView,
		FontFamily:       fontFamily,
		InvoiceDirectory: invoiceDir,
	}
}

func (m ConfigModel) Init() tea.Cmd {
	return nil
}

func (m ConfigModel) Update(msg tea.Msg) (ConfigModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("updating config model", msg.String())
		switch {
		case key.Matches(msg, m.keys.NextField):
			log.Println("updating config next field", m.focusID)
			switch m.focusID {
			case 0:
				m.FontFamily.Focus()
				m.focusID++
			case 1:
				m.focusID++
				m.FontFamily.Blur()
				cmd = m.InvoiceDirectory.Focus()
				cmds = append(cmds, cmd)
			case 2:
				m.InvoiceDirectory.Blur()

				// last field, goes back to nothing
				m.focusID = 0
			default:
			}
			case key.Matches(msg, m.keys.Esc):
				m.focusID = 0
				m.FontFamily.Blur()
				m.InvoiceDirectory.Blur()
		case key.Matches(msg, m.keys.Back):
			// TODO: think of a better way to avoid key event fallback
			// I wanted to avoid triggering esc event when we have
			// not picked anything with select.
			if m.focusID == 0 {
				cmd = pkg.GoMain()
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}
	}
	m.FontFamily, cmd = m.FontFamily.Update(msg)
	cmds = append(cmds, cmd)
	m.InvoiceDirectory, cmd = m.InvoiceDirectory.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ConfigModel) View() string {
	b := strings.Builder{}

	b.WriteString(m.FontFamily.View())
	b.WriteString("\n\n")
	b.WriteString(m.InvoiceDirectory.View())
	b.WriteString("\n\n")

	return b.String()
}
