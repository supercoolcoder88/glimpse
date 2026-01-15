package main

import (
	"fmt"
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

	activeFilter := make(map[string]any)

	fieldFilterSidebar := tview.NewForm()

	for _, field := range logs.CommonFields {
		fieldFilterSidebar.AddInputField(field, "", 30, nil, func(text string) {
			activeFilter[field] = text
		})
	}

	grid.AddItem(fieldFilterSidebar, 0, 0, 1, 1, 0, 0, true)

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
