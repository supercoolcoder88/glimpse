package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewSearchButton() *tview.Button {
	searchButton := tview.NewButton("Search")

	searchButton.
		SetBackgroundColor(tcell.NewHexColor(0x3a3a3a)).
		SetBorder(true).
		SetBorderColor(tcell.ColorGray)

	return searchButton
}
