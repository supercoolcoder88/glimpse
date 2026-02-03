package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewSearchBar() *tview.InputField {
	searchbar := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldBackgroundColor(tcell.ColorDarkSlateGrey)

	return searchbar
}
