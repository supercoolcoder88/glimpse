package main

import (
	"fmt"
	"glimpse/logs"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
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
		SetColumns(35, 0)

	// Filter area
	fieldFilterSidebar := tview.NewForm()
	fieldFilterSidebar.SetItemPadding(1).SetBorder(true).SetTitle(" Filters ")

	activeFilters := make(map[string]string)
	for _, field := range logs.CommonFields {
		fieldFilterSidebar.AddInputField(field, "", 0, nil, func(text string) {
			activeFilters[field] = text
		})
	}

	// Search bar component
	searchbar := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldBackgroundColor(tcell.ColorDarkSlateGrey)

	// Search button
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

	// Search row
	searchRow := tview.NewFlex().
		AddItem(searchbar, 0, 1, true).
		AddItem(searchButton, 12, 0, false)

	// Draw items
	// Create a placeholder for where the logs will go
	logDisplay := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	logDisplay.SetBorder(true).SetTitle(" Log Output ")

	// searchrow on top, logs on bot
	rightSide := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchRow, 3, 0, false).
		AddItem(logDisplay, 0, 1, false)

	// Draw grid
	grid.Clear()
	grid.SetRows(0).SetColumns(35, 0)
	grid.AddItem(fieldFilterSidebar, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(rightSide, 0, 1, 1, 1, 0, 0, false)

	// Log display area
	// textArea := tview.NewTextView()
	if err := app.EnableMouse(true).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
