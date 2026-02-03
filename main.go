package main

import (
	"fmt"
	"glimpse/components"
	"glimpse/logs"
	"os"
	"os/signal"
	"syscall"

	"github.com/rivo/tview"
)

func main() {
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
	// searchRow := tview.NewFlex().
	// 	AddItem(searchbar, 0, 1, true).
	// 	AddItem(searchButton, 12, 0, false)

	searchRow := tview.NewFlex().
		AddItem(searchbar, 0, 1, true)

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
	grid.AddItem(sidebar, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(rightSide, 0, 1, 1, 1, 0, 0, false)

	// Reading the logs
	readCh := make(chan logs.Entry)
	go func() {
		<-sigs
		app.Stop()
	}()

	go func() {
		err := logs.Read(os.Stdin, readCh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		for entry := range readCh {
			// app.QueueUpdateDraw here
		}
	}()

	// Log display area
	// textArea := tview.NewTextView()
	if err := app.EnableMouse(true).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

	fmt.Println("\nShutting down... cleaning up temporary database.")
	defer os.Remove("glimpse_temp.db")
}
