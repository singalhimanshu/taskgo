package ui

import (
	"strings"

	"github.com/rivo/tview"
)

// NewAddPage provides the form to create a new task.
func NewAddPage(p *BoardPage, pos int) *tview.Form {
	form := tview.NewForm().
		AddInputField("Task", "", 20, nil, nil).
		AddInputField("Task Description", "", 20, nil, nil)
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
		err := p.data.AddNewTask(p.activeListIdx, taskName, taskDesc, pos)
		if err != nil {
			app.Stop()
			panic(err)
		}
		p.data.Save(p.fileName)
		p.redraw(p.activeListIdx)
		p.down()
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
