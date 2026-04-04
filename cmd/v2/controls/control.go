package controls

import tea "github.com/charmbracelet/bubbletea"

// Control is an interface for form controls that can be focused and rendered.
type Control interface {
	Focus() tea.Cmd
	Blur()
	Focused() bool
	View() string
}
