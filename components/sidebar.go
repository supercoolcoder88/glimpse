package components

import (
	"github.com/rivo/tview"
)

func NewSidebar(fields []string) *tview.Form {
	form := tview.NewForm()
	form.SetItemPadding(1).SetBorder(true).SetTitle(" Filters ")

	for _, field := range fields {
		form.AddInputField(field, "", 0, nil, nil)
	}

	return form
}
