package ui

import (
	"log"
	"strings"

	"github.com/rivo/tview"
)

// NewEditPage provides the form to edit an existing task.
func NewEditPage(p *BoardPage, listIdx, taskIdx int) *tview.Form {
	task, err := p.data.GetTask(listIdx, taskIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	form := tview.NewForm().
		AddInputField("Task", task[0], 20, nil, nil).
		AddInputField("Task Description", task[1], 20, nil, nil)
	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskName = strings.TrimSpace(taskName)
		if len(taskName) <= 0 {
			emptyTitleNameModal := EmptyTitleNameModal()
			pages.AddAndSwitchToPage("emptyTitle", emptyTitleNameModal, true)
			return
		}
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		taskDesc = strings.TrimSpace(taskDesc)
		activeListIdx := p.activeListIdx
		err := p.data.EditTask(activeListIdx, p.activeTaskIdxs[activeListIdx], taskName, taskDesc)
		// editTaskCommand := command.CreateEditTaskCommand(activeListIdx, p.activeTaskIdxs[activeListIdx], taskName, taskDesc)
		// p.command.Execute(editTaskCommand)
		if err != nil {
			app.Stop()
			panic(err)
		}
		p.redraw(activeListIdx)
		pages.SwitchToPage("board")
		app.SetFocus(p.lists[p.activeListIdx])
	}).
		AddButton("Cancel", func() {
			pages.RemovePage("edit")
			pages.SwitchToPage("board")
			app.SetFocus(p.lists[p.activeListIdx])
		})
	form.SetBorder(true).SetTitle("Edit Task").SetTitleAlign(tview.AlignCenter)
	return form
}
