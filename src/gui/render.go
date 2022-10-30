package gui

import (
	"invoice-maker/config"
	"invoice-maker/gui/config_page"
	"invoice-maker/gui/invoice_add"
	"invoice-maker/gui/invoice_edit"
	"invoice-maker/gui/invoice_list"
	"invoice-maker/gui/issuer_edit"
	"invoice-maker/gui/menu"
	"invoice-maker/gui/receiver_add"
	"invoice-maker/gui/receiver_edit"
	"invoice-maker/gui/receiver_list"
	"invoice-maker/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/* Initializes application. Run it once in main only. */
func Run(tui *types.TUI) error {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal(err)
	}
	tui.Config = cfg

	tui.App = tview.NewApplication()
	tui.Pages = tview.NewPages()
	tui.ActivePage = types.PageDefault
	tui.Rerender = func() { Render(tui) }

	Render(tui)

	return tui.App.Run()
}

func Render(tui *types.TUI) {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal(err)
	}
	tui.Config = cfg
	setup(tui)
}

func setup(tui *types.TUI) {
	menu.Render(tui)

	receiver_list.Render(tui)
	issuer_edit.Render(tui)
	invoice_list.Render(tui)

	tui.App.
		SetRoot(tui.Pages, true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return eventHandler(event, tui)
		})
}

func eventHandler(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyCtrlC ||
		eventKey.Key() == tcell.KeyCtrlD {
		tui.App.Stop()
	}

	switch tui.ActivePage {
	case types.PageDefault:
		return menu.HandleEvents(eventKey, tui)
	case types.PageIssuerEdit:
		return issuer_edit.HandleEvents(eventKey, tui)

	case types.PageReceiverList:
		return receiver_list.HandleEvents(eventKey, tui)
	case types.PageReceiverAdd:
		return receiver_add.HandleEvents(eventKey, tui)
	case types.PageReceiverEdit:
		return receiver_edit.HandleEvents(eventKey, tui)

	case types.PageInvoiceList:
		return invoice_list.HandleEvents(eventKey, tui)
	case types.PageInvoiceAdd:
		return invoice_add.HandleEvents(eventKey, tui)
	case types.PageInvoiceEdit:
		return invoice_edit.HandleEvents(eventKey, tui)

	case types.PageConfig:
		return config_page.HandleEvents(eventKey, tui)
	}

	return eventKey
}
