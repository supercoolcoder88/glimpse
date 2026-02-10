package components

import "github.com/rivo/tview"

func NewDisplay(app *tview.Application) *tview.TextView {
	display := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	display.SetBorder(true).SetTitle(" Log Output ")
	return display
}
