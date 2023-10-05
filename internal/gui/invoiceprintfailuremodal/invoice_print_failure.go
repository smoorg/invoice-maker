package invoiceprintfailuremodal

import (
	"fmt"
	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var pageName = types.PagePrintFailureModal

func Render(tui *types.TUI, failureMsg string) {
	content := fmt.Sprintf("There was an error during printing an invoice.\nUsually it happens due to not enough space for field content to render. Error message:\n\n%s", failureMsg)

	modal := tview.NewModal().
		SetText(content).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "OK":
				tui.RefreshConfig()
				tui.Pages.RemovePage(pageName)
				tui.SwitchToPrevious()
			}
		})
	tui.SetDefaultStyle(modal.Box)
	tui.AddAndSwitchToPage(pageName, modal)
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if vimkeys.Back(eventKey) {
		tui.SwitchToPage(tui.PreviousPage)
		return nil
	}
	if vimkeys.Down(eventKey) {
		return tcell.NewEventKey(tcell.KeyDown, tcell.RuneDArrow, tcell.ModNone)
	}
	if vimkeys.Up(eventKey) {
		return tcell.NewEventKey(tcell.KeyUp, tcell.RuneUArrow, tcell.ModNone)
	}
	if vimkeys.Forward(eventKey) || eventKey.Rune() == 'e' {
		return tcell.NewEventKey(tcell.KeyEnter, rune(tcell.KeyEnter), tcell.ModNone)
	}

	return eventKey
}
