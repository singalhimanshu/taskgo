package ui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/command"
)

// NewAddPage provides the form to create a new task.
func NewAddPage(p *BoardPage, pos int) tview.Primitive {
	width, height := GetSize()
	form := tview.NewForm().
		AddInputField("Task", "", width/4, nil, nil).
		AddInputField("Task Description", "", width/4, nil, nil)
	form = form.AddButton("Save", func() {
		taskName := form.GetFormItemByLabel("Task").(*tview.InputField).GetText()
		taskName = strings.TrimSpace(taskName)
		if len(taskName) <= 0 {
			return
		}
		taskDesc := form.GetFormItemByLabel("Task Description").(*tview.InputField).GetText()
		taskDesc = strings.TrimSpace(taskDesc)
		addTaskCommand := command.CreateAddTaskCommand(p.activeListIdx, taskName, taskDesc, pos)
		if err := p.command.Execute(addTaskCommand); err != nil {
			app.Stop()
			log.Fatal(err)
		}
		p.redraw(p.activeListIdx)
		pages.SwitchToPage("board")
	}).AddButton("Cancel", func() {
		closeAddPage()
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			closeAddPage()
		}
		return event
	})
	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)
	return GetCenteredModal(form, width/2, height/2)
}

func closeAddPage() {
	pages.RemovePage("add")
	pages.SwitchToPage("board")
}
