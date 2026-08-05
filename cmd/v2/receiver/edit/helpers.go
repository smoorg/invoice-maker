package edit

import (
	"fmt"
	"invoice-maker/pkg/config"
)

func (m *ReceiverEdit) SetReceiver(v config.Company) {
	m.receiver = v
}

func (m *ReceiverEdit) Focus(focus int) {
	switch focus {
	case FocusNone:
		break
	case FocusName:
		m.inputName.Focus()
	case FocusAddress:
		m.inputAddress.Focus()
	case FocusTaxID:
		m.inputTaxID.Focus()
	default:
		panic("unimplemented focus type " + fmt.Sprint(focus))
	}
}

func (m *ReceiverEdit) Blur() {
	m.inputAddress.Blur()
	m.inputName.Blur()
	m.inputTaxID.Blur()
}

