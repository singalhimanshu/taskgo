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
	fieldWidth := 20
	form := tview.NewForm().
		AddInputField("Task", "", fieldWidth, nil, nil).
		AddInputField("Task Description", "", fieldWidth, nil, nil)
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
		switch event.Rune() {
		case 'q':
			closeAddPage()
		}
		return event
	})
	form.SetBorder(true).SetTitle("Create Task").SetTitleAlign(tview.AlignCenter)
	width, height := GetSize()
	return GetCenteredModal(form, width/2, height/2)
}

func closeAddPage() {
	pages.RemovePage("add")
	pages.SwitchToPage("board")
}

func EmptyTitleNameModal() *tview.Modal {
	emptyModal := tview.NewModal().
		SetText("Empty title name not allowed").
		SetBackgroundColor(theme.PrimitiveBackgroundColor).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				pages.SwitchToPage("add")
			}
		})
	return emptyModal
}
