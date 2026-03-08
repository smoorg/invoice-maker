package config_page

import (
	"path/filepath"

	"invoice-maker/internal/gui/modal"
	"invoice-maker/internal/types"

	"invoice-maker/pkg/config"
	"invoice-maker/pkg/font"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var title = "General Configuration"

func Render(tui *types.TUI) {
	tui.AddAndSwitchToPage(
		types.PageConfig,
		modal.Modal(tui, types.PageConfig, types.PageDefault, configPage(tui), 50, 11, title),
	)
}

func configPage(tui *types.TUI) *tview.Form {
	ff, err := font.GetFontFamilies()
	if err != nil {
		return nil
	}
	page := tview.NewForm()
	tui.SetDefaultStyle(page.Box)

	page.AddInputField("Invoice Directory", tui.Config.Config.InvoiceDirectory, 80, nil, func(text string) {
		tui.Config.Config.InvoiceDirectory = filepath.Join(text)
	})

	setFontStyle := func(opt string, i int) {
		if i < 0 || opt == "" {
			return
		}
		tui.Config.Config.SetStyle(opt)

		if font, err := font.FindFonts(tui.Config.Config.Family, tui.Config.Config.Style); err == nil {
			tui.Config.Config.SetPath(font[0].Filepath)
		}
	}

	setFontFamily := func(option string, optionIndex int) {
		// do nothing as nothing has been picked. Fixes issue when its triggered before "Font Style" exists in the form.
		if optionIndex < 0 {
			return
		}
		fonts, err := font.FindFonts(option, "")
		if err != nil || len(fonts) == 0 {
			return
		}
		tui.Config.Config.SetFamily(fonts[0].Family)

		styleDropdown := page.GetFormItem(page.GetFormItemIndex("Font Style"))
		if len(fonts) > 0 {
			styles, err := font.GetFontStyles(fonts[0].Family)

			if err != nil {
				return
			}

			switch v := styleDropdown.(type) {
			case *tview.DropDown:
				v.SetCurrentOption(-1)
				v.SetOptions(styles, setFontStyle)
			}
		}
	}

	page.AddDropDown("Font Family", ff, -1, setFontFamily)
	page.AddDropDown("Font Style", []string{}, -1, setFontStyle)

	pickedIndex := -1
	for i, v := range ff {
		if v == tui.Config.Config.Family {
			pickedIndex = i
		}
	}
	fontFamilyDropdown := page.GetFormItem(page.GetFormItemIndex("Font Family"))
	fontFamilyDropdown.(*tview.DropDown).SetCurrentOption(pickedIndex)

	pickedStyle := -1
	s, err := font.GetFontStyles(tui.Config.Config.Family)
	for i, v := range s {
		if v == tui.Config.Config.Style {
			pickedStyle = i
			break
		}
	}
	fontVariantDropdown := page.GetFormItem(page.GetFormItemIndex("Font Style"))
	fontVariantDropdown.(*tview.DropDown).SetOptions(s, setFontStyle)
	fontVariantDropdown.(*tview.DropDown).SetCurrentOption(pickedStyle)

	page.SetTitle(title)
	page.SetBorder(true)

	page.AddButton("Save", func() {
		if valid := config.IsValidInvoiceDirectory(tui.Config.Config.InvoiceDirectory); valid == false {
			msg := "Invoice directory provided is not a valid directory or no privileges to modify it. Please modify it accordingly and ensure its absolute path."
			modal.Error(tui, msg, types.PageConfig, 40, 5, "Error", func() {
				tui.RefreshConfig()
				tui.Rerender()
			})
			return
		}

		if err := tui.Config.WriteConfig(); err != nil {
			modal.Error(tui, err.Error(), types.PageConfig, 40, 5, "Error", func() { Render(tui) })
		}
		tui.SwitchToPage(types.PageDefault)
	})

	return page
}

func HandleEvents(eventKey *tcell.EventKey, tui *types.TUI) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEsc {
		tui.Pages.RemovePage(types.PageConfig)
		tui.SwitchToPrevious()
		return nil
	}

	return eventKey
}
