package ui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
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
	width, height := GetSize()
	form := tview.NewForm().
		AddInputField("Task", task.ItemName, width/4, nil, nil).
		AddInputField("Task Description", task.ItemDescription, width/4, nil, nil)
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
	}).
		AddButton("Cancel", func() {
			closeRemovePage()
		})
	form.SetBorder(true).SetTitle("Edit Task").SetTitleAlign(tview.AlignCenter)
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeRemovePage()
		}
		return event
	})
	return GetCenteredModal(form, width/2, height/2)
}

func closeRemovePage() {
	pages.RemovePage("edit")
	pages.SwitchToPage("board")
}
