package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewSearchBar() *tview.InputField {
	searchbar := tview.NewInputField()

	searchbar.SetBorder(true).
		SetTitle(" Search ").
		SetBackgroundColor(tcell.NewHexColor(0x2a2a2a)).
		SetBorderColor(tcell.ColorGray)
	return searchbar
}
