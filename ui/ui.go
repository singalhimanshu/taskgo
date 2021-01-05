package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	theme *tview.Theme
	pages *tview.Pages
)

var globalInputCapture = func(event *tcell.EventKey) *tcell.EventKey {
	s := string(event.Rune())
	switch s {
	case "?":
		pages.SwitchToPage("help")
	case "q":
		app.Stop()
	}
	return event
}

func defaultTheme() *tview.Theme {
	return &tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorBlack,          // Main background color for primitives.
		ContrastBackgroundColor:     tcell.ColorBlue,           // Background color for contrasting elements.
		MoreContrastBackgroundColor: tcell.ColorGreen,          // Background color for even more contrasting elements.
		BorderColor:                 tcell.ColorGrey,           // Box borders.
		TitleColor:                  tcell.ColorCoral,          // Box titles.
		GraphicsColor:               tcell.ColorFuchsia,        // Graphics.
		PrimaryTextColor:            tcell.ColorWhite,          // Primary text.
		SecondaryTextColor:          tcell.ColorAqua,           // Secondary text (e.g. labels).
		TertiaryTextColor:           tcell.ColorMediumSeaGreen, // Tertiary text (e.g. subtitles, notes).
		InverseTextColor:            tcell.ColorBlue,           // Text on primary-colored backgrounds.
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,       // Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}

// Start runs the application
func Start(fileName string) error {
	app = tview.NewApplication()
	start(fileName)
	if err := app.Run(); err != nil {
		return fmt.Errorf("Error running app: %s", err)
	}
	return nil
}

func start(fileName string) {
	theme = defaultTheme()
	boardPage := NewBoardPage(fileName)
	boardPageFrame := boardPage.Page()
	pages = tview.NewPages().
		AddPage("board", boardPageFrame, true, true)
	app.SetRoot(pages, true).SetFocus(boardPageFrame)
}
