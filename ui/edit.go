package ui

import (
	"log"
	"strings"

	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/command"
)

// NewEditPage provides the form to edit an existing task.
func NewEditPage(p *BoardPage, listIdx, taskIdx int) tview.Primitive {
	task, err := p.data.GetTask(listIdx, taskIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	fieldWidth := 20
	form := tview.NewForm().
		AddInputField("Task", task.ItemName, fieldWidth, nil, nil).
		AddInputField("Task Description", task.ItemDescription, fieldWidth, nil, nil)
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
		editTaskCommand := command.CreateEditTaskCommand(activeListIdx, p.activeTaskIdxs[activeListIdx], taskName, taskDesc)
		if err := p.command.Execute(editTaskCommand); err != nil {
			app.Stop()
			log.Fatal(err)
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
	width, height := GetSize()
	return GetCenteredModal(form, width/2, height/2)
}
