package ui

import (
	"github.com/rivo/tview"
)

// NewAddPage provides the form to create a new task.
func NewAddPage(p *BoardPage) *tview.Form {
	form := tview.NewForm().
		AddInputField("Task", "", 20, nil, nil).
		AddInputField("Task Description", "", 20, nil, nil)

	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		err := p.data.AddNewTask(p.activeListIdx, taskName, taskDesc)
		if err != nil {
			panic(err)
		}

		p.data.Save()
		p.redraw()
		pages.SwitchToPage("board")
		app.SetFocus(p.lists[p.activeListIdx])
	}).
		AddButton("Cancel", func() {
			pages.RemovePage("add")
			pages.SwitchToPage("board")
			app.SetFocus(p.lists[p.activeListIdx])
		})

	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)

	return form
}
