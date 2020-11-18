package ui

import "github.com/rivo/tview"

func NewAddPage() *tview.Form {
	form := tview.NewForm().
		AddInputField("Task", "", 20, nil, nil).
		AddInputField("Task Description", "", 20, nil, nil).
		AddButton("Save", nil).
		AddButton("Cancel", func() {
			pages.RemovePage("add")
			pages.SwitchToPage("board")
		})

	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)

	return form
}
