package main

import (
	"fmt"
	"glimpse/components"
	"glimpse/db"
	"glimpse/logs"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/rivo/tview"
)

func main() {
	defer os.Remove("glimpse_temp.db")

	sqlite, _ := db.Initialise()
	defer sqlite.Close()

	readCh := make(chan logs.Entry)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(35, 0)

	sidebar := components.NewSidebar(logs.CommonFields)
	searchbar := components.NewSearchBar()

	searchRow := tview.NewFlex().
		AddItem(searchbar, 0, 1, true)
	logDisplay := components.NewDisplay(app)

		// Search button

	search := func(filterSidebar *tview.Form, logDisplay *tview.TextView, ch chan logs.Entry, db *sqlx.DB) {
		inputs := filterSidebar.GetFormItemCount()

		rules := []logs.Rule{}

		for i := range inputs {
			f := filterSidebar.GetFormItem(i).(*tview.InputField)
			s := strings.Split(strings.TrimSpace(f.GetText()), "")
			if len(s) != 2 {
				panic("incorrect filter format should be *field operator value*") // TODO: Handle error
			}
			r, err := logs.NewRule(f.GetLabel(), s[1], s[0])
			if err != nil {
				// TODO: Handle error here
			}
			rules = append(rules, *r)
		}

		filter := logs.NewFilter(db)
		res, err := filter.Apply(rules)

		if err != nil {
			// TODO: Handle error
		}

		logDisplay.SetText("")
		logDisplay.ScrollToBeginning()

		for _, r := range res {
			ch <- r
		}
	}

	searchButton := components.NewSearchButton()
	searchButton.SetSelectedFunc(func () {
		search(sidebar, logDisplay, readCh, sqlite)
	})

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

	// Shutdown goroutine
	go func() {
		<-sigs
		app.Stop()
	}()

	// Read routine
	go func() {
		if err := logs.Read(os.Stdin, readCh, sqlite); err != nil {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(logDisplay, "[red]Error: %v\n", err)
				app.Stop()
			})
		}
	}()

	go func() {
		for entry := range readCh {
			app.QueueUpdateDraw(func() {
				fmt.Fprintf(logDisplay, "%s \n", entry.Raw)
			})
		}
	}()

	if err := app.EnableMouse(true).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

}
