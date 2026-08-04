package styles

import "github.com/charmbracelet/lipgloss"

var (
	Bg       = lipgloss.Color("#333")
	Red      = lipgloss.Color("#F54927")
	White    = lipgloss.Color("#fff")
	Grey     = lipgloss.Color("#555")
	DarkGrey = lipgloss.Color("#111")
)

var (
	GreyedOut    = lipgloss.NewStyle().Foreground(Grey).Background(DarkGrey)
	ModifyInput  = lipgloss.NewStyle().Foreground(White).Background(Bg)
	InvalidInput = lipgloss.NewStyle().Foreground(Red)
	BlurInput    = lipgloss.NewStyle().Background(Bg)
	Cursor       = lipgloss.NewStyle().Foreground(Grey).Background(ModifyInput.GetBackground())
)
