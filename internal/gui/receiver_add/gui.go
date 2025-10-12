package receiver_add

import (
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"invoice-maker/pkg/config"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AddOrEditReceiver(tui *types.TUI, company *config.Company, title string, save func(), cancel func()) *tview.Form {
	form := tview.NewForm()
	form.AddInputField("Name", company.Name, 0, nil, func(text string) { company.Name = text })
	form.AddInputField("Address", company.Address, 0, nil, func(text string) { company.Address = text })
	form.AddInputField("TaxID", company.TaxID, 0, nil, func(text string) { company.TaxID = text })

	form.AddButton("Save", save)
	form.AddButton("Cancel", cancel)

	tui.SetDefaultStyle(form.Box)
	form.SetBorder(true)

	form.SetTitle(title)

	return form
}

func Render(tui *types.TUI) {
	tui.AddAndSwitchToPage(
		types.PageReceiverAdd,
		modal.Modal(tui, types.PageReceiverAdd, types.PageReceiverList, AddReceiver(tui), 50, 11, "Add Receiver"),
	)
}

func AddReceiver(tui *types.TUI) tview.Primitive {
	company := &config.Company{}

	save := func() {
		tui.Config.Receivers = append(tui.Config.Receivers, *company)
		err := tui.Config.WriteConfig()
		if err != nil {
			modal.Error(tui, err.Error(), types.PageReceiverAdd, 40, 5, "Error", func() { Render(tui) })
		}
		tui.SwitchToPage(types.PageReceiverList)
	}

	cancel := func() {
		company = &config.Company{}
		tui.SwitchToPage(types.PageReceiverList)
	}

	return AddOrEditReceiver(tui, company, "Add Receiver", save, cancel)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.SwitchToPrevious()
		return nil
	}

	return eventKey
}
