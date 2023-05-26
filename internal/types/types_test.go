package types_test

import (
	"invoice-maker/internal/types"
	"testing"

	"github.com/rivo/tview"
)

func TestSwitchToNextPage(t *testing.T) {
	tui := &types.TUI{}

	tui.ActivePage = "active_page"
	tui.PreviousPage = "inactive_page"

	tui.SwitchToNext("test")

	if tui.ActivePage != "test" {
		t.Error("switch to next should replace active page")
	}

	if tui.PreviousPage != "active_page" {
		t.Error("switch to next should replace previous page to last active")
	}
}

func TestAddAndSwitchToPage(t *testing.T) {
	tui := &types.TUI{}

	tui.ActivePage = "active_page"
	tui.PreviousPage = "inactive_page"

	tui.App = tview.NewApplication()
	tui.Pages = tview.NewPages()
	tui.App.SetRoot(tui.Pages, true)
	form := tview.NewForm()
	tui.AddAndSwitchToPage("test", form)

	if tui.ActivePage != "test" {
		t.Error("switch to next should replace active page")
	}

	if tui.PreviousPage != "active_page" {
		t.Error("switch to next should replace previous page to last active")
	}
}
