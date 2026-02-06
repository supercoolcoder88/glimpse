package main

import (
	"fmt"
	"glimpse/components"
	"glimpse/db"
	"glimpse/logs"
	"os"
	"os/signal"
	"syscall"

	"github.com/rivo/tview"
)

func main() {
	defer os.Remove("glimpse_temp.db")

	// Cleanup function
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(35, 0)

	sidebar := components.NewSidebar(logs.CommonFields)

	// Search bar component
	searchbar := components.NewSearchBar()

	// Search button
	// searchButton := components.NewSearchButton(sidebar)

	// Search row
	searchRow := tview.NewFlex().
		AddItem(searchbar, 0, 1, true)

	// Draw items
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
	grid.AddItem(sidebar, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(rightSide, 0, 1, 1, 1, 0, 0, false)

	// Reading the logs
	readCh := make(chan logs.Entry)

	// Shutdown goroutine
	go func() {
		<-sigs
		app.Stop()
	}()

	go func() {
		sqlite, _ := db.Initialise()
		defer sqlite.Close()

		if err := logs.Read(os.Stdin, readCh, sqlite); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		for entry := range readCh {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(logDisplay, "%s \n", entry.Raw)
			})
		}
	}()

	defer close(readCh)

	if err := app.EnableMouse(true).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

}
