package singleselect

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	Label string
	Value string
}

type Model struct {
	label      string
	labelLimit int
	items      []Item
	filtered   []int // indices into items that match filter
	cursor     int   // position in filtered list
	selected   int   // index in original items list
	expanded   bool
	focused    bool
	maxVisible int

	filter textinput.Model
}

func New(label string, limit int, items []Item) Model {
	f := textinput.New()
	f.Placeholder = "Type to filter..."
	f.Prompt = ""
	f.Width = 80

	// Initialize filtered to all items
	filtered := make([]int, len(items))
	for i := range items {
		filtered[i] = i
	}

	return Model{
		label:      label,
		labelLimit: limit,
		items:      items,
		filtered:   filtered,
		cursor:     0,
		selected:   0,
		expanded:   false,
		focused:    false,
		maxVisible: 10,
		filter:     f,
	}
}

func (m *Model) SetMaxVisible(n int) {
	m.maxVisible = n
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return m.filter.Focus()
}

func (m *Model) Focused() bool {
	return m.filter.Focused() && m.focused && m.expanded
}

func (m *Model) Blur() {
	m.focused = false
	m.expanded = false
	m.filter.Blur()
	m.applyFilter()
}

// itemLabels implements fuzzy.Source for fuzzy matching
type itemLabels []Item

func (il itemLabels) String(i int) string {
	return il[i].Label
}

func (il itemLabels) Len() int {
	return len(il)
}

func (m *Model) applyFilter() {
	query := m.filter.Value()
	m.filtered = m.filtered[:0]

	if query == "" {
		// No filter - show all items
		for i := range m.items {
			m.filtered = append(m.filtered, i)
		}
	} else {
		// Fuzzy match
		matches := fuzzy.FindFrom(query, itemLabels(m.items))
		for _, match := range matches {
			m.filtered = append(m.filtered, match.Index)
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
	Cancel  key.Binding
	Toggle  key.Binding
	Exit    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
	),
	Cancel: key.NewBinding(
		key.WithKeys(tea.KeyEsc.String()),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" "),
	),
}

type BlurEvent struct{}

func SendBlurEvent() tea.Cmd {
	return func() tea.Msg {
		return BlurEvent{}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case BlurEvent:
		if m.expanded {
			m.expanded = false
			m.filter.Blur()
			m.applyFilter()
		}
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.expanded && m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case key.Matches(msg, keys.Down):
			if m.expanded && m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil
		case key.Matches(msg, keys.Confirm):
			if m.expanded {
				if len(m.filtered) > 0 {
					m.selected = m.filtered[m.cursor]
				}
				m.expanded = false
				m.filter.Blur()
				m.filter.SetValue("")
				m.applyFilter()
			} else {
				m.expanded = true
				m.filter.Focus()
				m.cursor = 0
				// Position cursor at currently selected item if visible
				for i, idx := range m.filtered {
					if idx == m.selected {
						m.cursor = i
						break
					}
				}
			}
			return m, nil
		case key.Matches(msg, keys.Toggle):
			// Only toggle when collapsed - let spacebar pass through when expanded
			if !m.expanded {
				m.expanded = true
				m.filter.Focus()
				m.cursor = 0
				return m, cmd
			}
		case key.Matches(msg, keys.Cancel):
			if m.focused && m.Focused() {
				cmd = SendBlurEvent()
				return m, cmd
			}
		}
	}

	// Update filter input when expanded
	if m.expanded {
		prevValue := m.filter.Value()
		m.filter, cmd = m.filter.Update(msg)
		if m.filter.Value() != prevValue {
			m.applyFilter()
			m.cursor = 0
		}
	}

	return m, cmd
}

var (
	bg    = lipgloss.Color("#333")
	white = lipgloss.Color("#fff")
	grey  = lipgloss.Color("#555")

	labelStyle       = lipgloss.NewStyle().Foreground(white)
	selectedStyle    = lipgloss.NewStyle().Foreground(white).Background(bg)
	itemStyle        = lipgloss.NewStyle().Foreground(grey)
	cursorStyle      = lipgloss.NewStyle().Foreground(white).Background(lipgloss.Color("#555"))
	collapsedStyle   = lipgloss.NewStyle().Foreground(grey).Background(bg)
	collapsedFocused = lipgloss.NewStyle().Foreground(white).Background(bg)
	filterStyle      = lipgloss.NewStyle().Foreground(grey)
	matchCountStyle  = lipgloss.NewStyle().Foreground(grey).Italic(true)
)

func (m Model) View() string {
	b := strings.Builder{}

	// Label
	b.WriteString(labelStyle.Render(m.label))
	lettersLeft := m.labelLimit - len(m.label)
	for range lettersLeft {
		b.WriteString(" ")
	}
	b.WriteString(": ")

	if len(m.items) == 0 {
		b.WriteString(itemStyle.Render("(no items)"))
		return b.String()
	}

	if !m.expanded {
		// Collapsed view - show selected item
		selectedLabel := m.items[m.selected].Label
		if m.focused {
			b.WriteString(collapsedFocused.Render("[ " + selectedLabel + " ]"))
		} else {
			b.WriteString(collapsedStyle.Render("[ " + selectedLabel + " ]"))
		}
	} else {
		// Filter input
		b.WriteString(m.filter.View())
		b.WriteString("\n")

		if len(m.filtered) == 0 {
			b.WriteString(itemStyle.Render("  (no matches)"))
		} else {
			// Calculate visible window
			start, end := m.visibleRange()

			// Show scroll indicator at top
			if start > 0 {
				b.WriteString(matchCountStyle.Render("  ... more above"))
				b.WriteString("\n")
			}

			// Render visible items
			for i := start; i < end; i++ {
				if i > start || start > 0 {
					b.WriteString("\n")
				}

				itemIdx := m.filtered[i]
				item := m.items[itemIdx]

				prefix := "  "
				if itemIdx == m.selected {
					prefix = "* "
				}

				if i == m.cursor {
					b.WriteString(cursorStyle.Render(prefix + item.Label))
				} else if itemIdx == m.selected {
					b.WriteString(selectedStyle.Render(prefix + item.Label))
				} else {
					b.WriteString(itemStyle.Render(prefix + item.Label))
				}
			}

			// Show scroll indicator at bottom
			if end < len(m.filtered) {
				b.WriteString("\n")
				b.WriteString(matchCountStyle.Render("  ... more below"))
			}
		}
	}

	return b.String()
}

func (m Model) visibleRange() (start, end int) {
	total := len(m.filtered)
	if total <= m.maxVisible {
		return 0, total
	}

	// Center the cursor in the visible window
	half := m.maxVisible / 2
	start = max(m.cursor-half, 0)

	end = start + m.maxVisible
	if end > total {
		end = total
		start = end - m.maxVisible
	}

	return start, end
}
