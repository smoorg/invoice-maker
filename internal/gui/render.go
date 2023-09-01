package gui

import (
	"invoice-maker/internal/config"
	"invoice-maker/internal/gui/config_page"
	"invoice-maker/internal/gui/help"
	"invoice-maker/internal/gui/invoice_add"
	"invoice-maker/internal/gui/invoice_edit"
	"invoice-maker/internal/gui/invoice_list"
	"invoice-maker/internal/gui/invoice_print_modal"
	"invoice-maker/internal/gui/invoiceprintfailuremodal"
	"invoice-maker/internal/gui/issuer_edit"
	"invoice-maker/internal/gui/menu"
	"invoice-maker/internal/gui/modal"

	"invoice-maker/internal/gui/receiver_add"
	"invoice-maker/internal/gui/receiver_edit"
	"invoice-maker/internal/gui/receiver_list"
	"invoice-maker/internal/types"

	"github.com/gdamore/tcell/v2"
)

/* Initializes application. Run it once in main only. */
func Run(tui *types.TUI) error {
	tui.ActivePage = types.PageDefault
	tui.Rerender = func() { Render(tui) }

	Render(tui)

	return tui.App.Run()
}

func Render(tui *types.TUI) {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal("render error:", err)
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

	if eventKey.Key() == tcell.KeyHelp || eventKey.Key() == tcell.KeyCtrlH {
		help.Render(tui)
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

	case types.PageModal:
		return modal.HandleEvents(eventKey, tui)
	case types.PagePrintModal:
		return invoiceprintmodal.HandleEvents(eventKey, tui)
	case types.PagePrintFailureModal:
		return invoiceprintfailuremodal.HandleEvents(eventKey, tui)

	case types.PageHelp:
		return help.HandleEvents(eventKey, tui)

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
