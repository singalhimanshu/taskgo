package ui

import "github.com/rivo/tview"

const helpText = `j: down
k: up
h: left
l: right
a: add task
D: delete a task
d: mark a task as done
e: change/edit task
L: move task right
H: move task left
J: move task down
K: move task up
Enter: view task info
q: quit
`

// NewHelpPage displays the help page that contains all the keybinds of the application
func NewHelpPage(p *BoardPage) *tview.Modal {
	help := tview.NewModal().
		SetText(helpText).
		SetBackgroundColor(theme.PrimitiveBackgroundColor).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "OK" {
				pages.HidePage("help")
				pages.SwitchToPage("board")
				app.SetFocus(p.lists[p.activeListIdx])
			}
		})
	return help
}
