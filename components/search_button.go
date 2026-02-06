package components

import (
	"glimpse/logs"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rivo/tview"
)

func NewSearchButton(filtersSidebar *tview.Form, search func(rules []logs.Rule, db *sqlx.DB) []logs.Entry) *tview.Button {
	searchButton := tview.NewButton("Search")

	inputs := filtersSidebar.GetFormItemCount()

	// form.GetFormItem(0).(*tview.InputField).GetText()
	searchButton.SetSelectedFunc(func() {
		rules := []logs.Rule{}

		// We assume the filter is structured with a pattern like "operator value"
		for i := range inputs {
			f := filtersSidebar.GetFormItem(i).(*tview.InputField)
			s := strings.Split(strings.TrimSpace(f.GetText()), "")
			if len(s) != 2 {
				panic("incorrect filter format should be *field operator value*")
			}
			r, err := logs.NewRule(f.GetLabel(), s[1], s[0])
			if err != nil {
				// TODO: Handle error here
			}
			rules = append(rules, *r)
		}

		// found, err := searchJSONLogs(rules, db)
	})

	return searchButton
}

func searchJSONLogs(rules []logs.Rule, db *sqlx.DB) ([]logs.Entry, error) {
	f := logs.NewFilter(db)
	res, err := f.Apply(rules)

	if err != nil {
		return nil, err
	}

	// TODO: Massage results
	return res, nil
}
