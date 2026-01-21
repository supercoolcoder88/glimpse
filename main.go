package main

import (
	"fmt"
	"glimpse/logs"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rivo/tview"
)

func main() {
	// Cleanup function
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nShutting down... cleaning up temporary database.")
		os.Remove("glimpse_temp.db")
		os.Exit(0)
	}()

	// err := logs.Read(os.Stdin)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	os.Exit(1)
	// }

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(3, 7)

	// Filter area
	fieldFilterSidebar := tview.NewForm()

	activeFilters := make(map[string]string)
	for _, field := range logs.CommonFields {
		fieldFilterSidebar.AddInputField(field, "", 30, nil, func(text string) {
			activeFilters[field] = text
		})
	}

	searchButton := tview.NewButton("Search")

	searchButton.SetSelectedFunc(func() {
		rules := []*logs.FilterRule{}
		// We assume the filter is structured with a pattern like "operator value"
		for k, v := range activeFilters {
			s := strings.Split(strings.TrimSpace(v), "")
			if len(s) != 2 {
				panic("incorrect filter format should be *field operator value*")
			}
			r, err := logs.NewFilterRule(k, s[1], s[0])
			if err != nil {
				// TODO: Handle error here
			}
			rules = append(rules, r)
		}
	})

	grid.AddItem(fieldFilterSidebar, 0, 0, 1, 1, 0, 0, true)

	// Log display area
	// textArea := tview.NewTextView()
	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
