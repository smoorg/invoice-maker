package issuer_edit

import (
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI) {
	tui.Pages.AddPage(
		types.PageIssuerEdit,
		modal.Modal(tui, types.PageIssuerEdit, types.PageDefault, EditIssuer(tui), 50, 17, "Edit issuer"),
		true,
		false)
}

func goBack(tui *types.TUI) {
	tui.SwitchToPage(types.PageDefault)
}

func EditIssuer(tui *types.TUI) *tview.Form {
	issuerForm := tview.NewForm().
		AddInputField("Name", tui.Config.Issuer.Name, 50, nil, func(text string) { tui.Config.Issuer.Name = text }).
		AddInputField("Address", tui.Config.Issuer.Address, 50, nil, func(text string) { tui.Config.Issuer.Address = text }).
		AddInputField("TaxID", tui.Config.Issuer.TaxID, 20, nil, func(text string) { tui.Config.Issuer.TaxID = text }).
		AddInputField("Account", tui.Config.Issuer.Account, 32, nil, func(text string) { tui.Config.Issuer.Account = text }).
		AddInputField("BankName", tui.Config.Issuer.BankName, 40, nil, func(text string) { tui.Config.Issuer.BankName = text }).
		AddInputField("BIC", tui.Config.Issuer.BIC, 12, nil, func(text string) { tui.Config.Issuer.BIC = text }).
		AddButton("Save", func() {
			if err := tui.Config.WriteConfig(); err != nil {
				modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui) })
			}
			goBack(tui)
		}).
		AddButton("Cancel", func() {
			goBack(tui)
		})

	issuerForm.SetTitle("Issuer details").SetBorder(true)
	issuerForm.SetBorder(true)

	return issuerForm
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		goBack(tui)
	}

	return eventKey
}
