package ui

import (
	"log"

	"github.com/rivo/tview"
)

// NewEditPage provides the form to edit an existing task.
func NewEditPage(p *BoardPage, listIdx, taskIdx int) *tview.Form {
	task, err := p.data.GetTask(listIdx, taskIdx)
	if err != nil {
		log.Fatal(err)
	}

	form := tview.NewForm().
		AddInputField("Task", task[0], 20, nil, nil).
		AddInputField("Task Description", task[1], 20, nil, nil)

	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		activeListIdx := p.activeListIdx
		err := p.data.EditTask(activeListIdx, p.activeTaskIdxs[activeListIdx], taskName, taskDesc)
		if err != nil {
			panic(err)
		}

		p.data.Save()
		p.redraw()
		pages.SwitchToPage("board")
	}).
		AddButton("Cancel", func() {
			pages.RemovePage("edit")
			pages.SwitchToPage("board")
		})

	form.SetBorder(true).SetTitle("Edit Task").SetTitleAlign(tview.AlignCenter)

	return form
}
