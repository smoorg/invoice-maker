package styles

import "github.com/charmbracelet/lipgloss"

var Bg = lipgloss.Color("#333")
var Red = lipgloss.Color("#F54927")
var White = lipgloss.Color("#fff")
var Grey = lipgloss.Color("#555")

var GreyedOut = lipgloss.NewStyle().Foreground(Grey).Background(lipgloss.NoColor{})
var ModifyInput = lipgloss.NewStyle().Foreground(White).Background(Bg)
var InvalidInput = lipgloss.NewStyle().Foreground(Red)
