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

	"github.com/gdamore/tcell/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rivo/tview"
)

type UIMessage interface{}

type AppendLog struct {
	Entry logs.Entry
}

type ShowSearchResults struct {
	Results []logs.Entry
}

type ShowError struct {
	Err error
}

type Quit struct{}

type UIState struct {
	mode string // live | search
}

func main() {
	defer os.Remove("glimpse_temp.db")

	tview.Styles.PrimitiveBackgroundColor = tcell.NewHexColor(0x1e1e1e)
	tview.Styles.ContrastBackgroundColor = tcell.NewHexColor(0x2a2a2a)
	tview.Styles.MoreContrastBackgroundColor = tcell.NewHexColor(0x333333)
	tview.Styles.BorderColor = tcell.ColorGray
	tview.Styles.TitleColor = tcell.ColorWhite
	tview.Styles.GraphicsColor = tcell.ColorWhite
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorLightGray

	sqlite, _ := db.Initialise()
	defer sqlite.Close()

	readCh := make(chan logs.Entry)
	sigs := make(chan os.Signal, 1)
	uiCh := make(chan UIMessage)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := tview.NewApplication()

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(35, 0)

	sidebar := components.NewSidebar(logs.CommonFields)
	searchbar := components.NewSearchBar()

	logDisplay := components.NewDisplay(app)

	// Search button
	search := func(filterSidebar *tview.Form, db *sqlx.DB) ([]logs.Entry, error) {
		inputs := filterSidebar.GetFormItemCount()

		rules := []logs.Rule{}

		for i := range inputs {
			f := filterSidebar.GetFormItem(i).(*tview.InputField)
			text := strings.TrimSpace(f.GetText())

			if text == "" {
				continue
			}

			s := strings.Fields(text)

			if len(s) != 2 {
				panic("incorrect filter format should be 'operator value'") // TODO: Handle error
			}
			r, err := logs.NewRule(f.GetLabel(), s[1], s[0])
			if err != nil {
				// TODO: Handle error here
			}
			rules = append(rules, *r)
		}

		filter := logs.NewFilter(db)
		res, err := filter.Apply(rules)
		return res, err
	}

	searchButton := components.NewSearchButton()
	searchButton.SetSelectedFunc(func() {
		res, err := search(sidebar, sqlite)
		if err != nil {
			uiCh <- ShowError{Err: err}
			return
		}

		uiCh <- ShowSearchResults{Results: res}
	})

	searchRow := tview.NewFlex().
		AddItem(searchbar, 0, 9, true).
		AddItem(searchButton, 0, 1, false)

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

	// UI state management
	go func() {
		state := UIState{
			mode: "live",
		}

		for msg := range uiCh {
			switch m := msg.(type) {

			case AppendLog:
				if state.mode == "live" {
					app.QueueUpdateDraw(func() {
						fmt.Fprint(logDisplay, formatLog(m.Entry))
						logDisplay.ScrollToEnd()
					})
				}

			case ShowSearchResults:
				state.mode = "search"

				app.QueueUpdateDraw(func() {
					logDisplay.SetText("")
					logDisplay.ScrollToBeginning()

					if len(m.Results) == 0 {
						fmt.Fprint(logDisplay, "\n [yellow]No results found[-]")
					} else {
						for _, r := range m.Results {
							fmt.Fprint(logDisplay, formatLog(r))
						}
					}
				})

			case ShowError:
				app.QueueUpdateDraw(func() {
					fmt.Fprintf(logDisplay, "\n [red]ERR: %v[-]\n", m.Err)
					logDisplay.ScrollToEnd()
				})

			case Quit:
				app.Stop()
				return
			}
		}
	}()

	// Read routine
	go func() {
		if err := logs.Read(os.Stdin, readCh, sqlite); err != nil {
			uiCh <- ShowError{Err: err}
			uiCh <- Quit{}
		}
	}()

	// Update app routine
	go func() {
		for entry := range readCh {
			uiCh <- AppendLog{Entry: entry}
		}
	}()

	if err := app.EnableMouse(true).SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

}

func formatLog(e logs.Entry) string {
	content := tview.Escape(e.Raw)
	separator := "[#444444]" + strings.Repeat("─", 80) + "[-]"

	return fmt.Sprintf("\n%s\n  %s\n", separator, content)
}
