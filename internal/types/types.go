package types

import (
	"invoice-maker/internal/config"

	"github.com/rivo/tview"
)

const (
	PageDefault    string = "default"     // Menu page
	PageHelp       string = "help"        // Help modal
	PageIssuerEdit string = "issuer_edit" // Edit issuer details. There are just one available for the time being.

	PageInvoiceAdd  string = "invoice_add"
	PageInvoiceEdit string = "invoice_edit"
	PageInvoiceList string = "invoice_list"

	PageReceiverAdd  string = "receiver_add"
	PageReceiverEdit string = "receiver_edit"
	PageReceiverList string = "receiver_list"

	PageConfig string = "config" // General configuration page

	PageModal             string = "modal" // We expect just single modal at once
	PagePrintModal        string = "print_modal"
	PagePrintFailureModal string = "print_failure_modal"
)

type TUI struct {
	App          *tview.Application
	Pages        *tview.Pages
	Config       *config.Config
	PreviousPage string // Previously active page (i.e. to show modal on top of it or go back)
	ActivePage   string // Active page,. It also impacts which HandleEvents function is used currently.
	Rerender     func()
}

func (tui *TUI) RefreshConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal("config refresh error: ", err)
	}
	tui.Config = cfg
}

func (tui *TUI) SwitchToNext(nextPage string) {
	tui.PreviousPage = tui.ActivePage
	tui.ActivePage = nextPage
}

func (tui *TUI) SwitchToPrevious() {
	tui.ActivePage, tui.PreviousPage = tui.PreviousPage, tui.ActivePage
}

func (tui *TUI) SwitchToPage(page string) {
	tui.RefreshConfig()
	tui.SwitchToNext(page)
	tui.Rerender()
	tui.Pages.SwitchToPage(page)
}

func (tui *TUI) AddAndSwitchToPage(page string, item tview.Primitive) {
	tui.SwitchToNext(page)
	tui.RefreshConfig()

	tui.Pages.AddPage(page, item, true, true)

}
