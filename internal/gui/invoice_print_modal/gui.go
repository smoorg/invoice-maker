package invoiceprintmodal

import (
	"fmt"
	"log"
	"os/exec"

	"invoice-maker/internal/types"
	"invoice-maker/internal/vimkeys"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Render(tui *types.TUI, filepath string) {
	content := fmt.Sprintf("invoice saved at: \n %s", filepath)

	modal := tview.NewModal().
		SetText(content).
		AddButtons([]string{"OK", "Show"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "OK":
				tui.RefreshConfig()
				tui.Pages.RemovePage(types.PagePrintModal)
				tui.SwitchToPrevious()
			case "Show":
				go func(file string) {
					cmd := exec.Command("xdg-open", file)
					_, err := cmd.Output()
					if err != nil {
						log.Fatal(err)
					}
				}(filepath)
			}

		})
	tui.SetDefaultStyle(modal.Box)
	tui.AddAndSwitchToPage(types.PagePrintModal, modal)
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
