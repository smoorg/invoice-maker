package edit

import (
	labelinput "invoice-maker/cmd/v2/controls/label_input"
	"invoice-maker/pkg/config"

	"github.com/charmbracelet/bubbles/key"
)

const (
	FocusNone = iota
	FocusName
	FocusAddress
	FocusTaxID
)

type keyMap struct {
	Back key.Binding
	Esc  key.Binding
	Next key.Binding
	Prev key.Binding
	Save key.Binding
}

type ReceiverEdit struct {
	receiver config.Company
	keys     keyMap

	inputName    labelinput.Model
	inputAddress labelinput.Model
	inputTaxID   labelinput.Model

	focus       int
	height      int
	helpContent string
}
