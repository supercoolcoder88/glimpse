package main

import (
	"fmt"
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

	grid := tview.NewGrid().
		SetRows(9, 1).
		SetColumns(3, 7)

	if err := tview.NewApplication().SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
