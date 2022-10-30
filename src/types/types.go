package types

import (
	"invoice-maker/config"

	"github.com/rivo/tview"
)

const (
	PageDefault    string = "default"     // Menu page
	PageIssuerEdit string = "issuer_edit" // Edit issuer details. There are just one available for the time being.

	PageInvoiceAdd  string = "invoice_add"
	PageInvoiceEdit string = "invoice_edit"
	PageInvoiceList string = "invoice_list"

	PageReceiverAdd  string = "receiver_add"
	PageReceiverEdit string = "receiver_edit"
	PageReceiverList string = "receiver_list"

	PageConfig   string = "config"    // General configuration page
	PageModal    string = "modal"     // we expect just single modal at once
	PageErrModal string = "err_modal" // we expect just single error modal at once, can overlap existing modal
)

type TUI struct {
	App        *tview.Application
	Pages      *tview.Pages
	Config     *config.Config
	ActivePage string
	Rerender   func()
}

func (tui *TUI) SwitchToPage(page string) {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal(err)
	}
	tui.Config = cfg

	tui.ActivePage = page
	tui.Rerender()
	tui.Pages.SwitchToPage(page)
}
func (tui *TUI) AddAndSwitchToPage(page string, item tview.Primitive) {
	cfg, err := config.GetConfig()
	if err != nil {
		tui.Fatal(err)
	}
	tui.Config = cfg

	tui.ActivePage = page
	tui.Pages.AddAndSwitchToPage(page, item, true)
}
