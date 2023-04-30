package config_page

import (
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var title = "General Configuration"

func Render(tui *types.TUI) {
	tui.AddAndSwitchToPage(
		types.PageConfig,
		modal.Modal(tui, types.PageDefault, configPage(tui), 50, 11, title),
	)
}

func configPage(tui *types.TUI) *tview.Form {
	page := tview.NewForm()

	page.AddInputField("Invoice Directory", tui.Config.InvoiceDirectory, 80, nil, func(text string) {
		tui.Config.InvoiceDirectory = filepath.Join(text)
	})

	page.SetTitle(title)
	page.SetBorder(true)

	page.AddButton("Save", func() {
		if valid := config.IsValidInvoiceDirectory(tui.Config.InvoiceDirectory); valid == false {
			msg := "Invoice directory provided is not a valid directory or no privileges to modify it. Please modify it accordingly and ensure its absolute path."
			modal.Error(tui, msg, types.PageConfig, 40, 5, "Error", func() {
				tui.Config.ReloadConfig()
				tui.Rerender()
			})
			return
		}

		if err := tui.Config.WriteConfig(); err != nil {
			modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui) })
		}
		tui.Pages.RemovePage(types.PageConfig)
		tui.SwitchToPage(types.PageDefault)
	})

	return page
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.Pages.RemovePage(types.PageConfig)
		tui.SwitchToPage(types.PageDefault)
		return nil
	}

	return eventKey
}
