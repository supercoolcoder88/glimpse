package components

import (
	"github.com/rivo/tview"
)

func NewSearchButton() *tview.Button {
	searchButton := tview.NewButton("Search")

	return searchButton
}
