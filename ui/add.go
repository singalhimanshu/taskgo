package ui

import (
	"github.com/rivo/tview"
)

func NewAddPage(p *BoardPage) *tview.Form {
	form := tview.NewForm().
		AddInputField("Task", "", 20, nil, nil).
		AddInputField("Task Description", "", 20, nil, nil)

	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		p.data.AddNewTask(p.activeListIdx, taskName, taskDesc)
		p.data.Save()
		p.redraw()
		pages.SwitchToPage("board")
	}).
		AddButton("Cancel", func() {
			pages.RemovePage("add")
			pages.SwitchToPage("board")
		})

	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)

	return form
}
