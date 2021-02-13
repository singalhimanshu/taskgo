package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/command"
	"github.com/singalhimanshu/taskgo/parser"
)

// BoardPage contains all the information that the board page requires.
type BoardPage struct {
	lists          []*tview.List
	theme          *tview.Theme
	data           parser.Data
	command        *command.CommandManager
	activeListIdx  int
	activeTaskIdxs []int
}

// NewBoardPage adds the data to BoardPage structure.
func NewBoardPage(fileName string) *BoardPage {
	theme := defaultTheme()
	data := parser.Data{
		FileName: fileName,
	}
	if err := data.ParseData(); err != nil {
		log.Fatal(err)
	}
	command := command.CreateNewCommand(&data)
	listCount := len(data.GetListNames())
	return &BoardPage{
		lists:          make([]*tview.List, listCount),
		data:           data,
		command:        command,
		theme:          theme,
		activeListIdx:  0,
		activeTaskIdxs: make([]int, listCount),
	}
}

// Page displays the contents of the board.
func (p *BoardPage) Page() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	listNames := p.data.GetListNames()
	for i := 0; i < len(listNames); i++ {
		p.lists[i] = tview.NewList()
		p.lists[i].
			ShowSecondaryText(false).
			SetBorder(true)
		// Highlights the first list
		if i == 0 {
			p.lists[i].SetBorderColor(theme.ContrastBackgroundColor)
		}
		p.lists[i].SetTitle(listNames[i])
		p.setInputCapture(i)
		p.addTasksToList(i)
		flex.AddItem(p.lists[i], 0, 1, i == 0)
	}
	boardName := p.data.GetBoardName()
	boardName = "Board: " + boardName
	frame := tview.NewFrame(flex).
		SetBorders(0, 0, 1, 0, 1, 1).
		AddText(boardName, true, tview.AlignCenter, p.theme.TitleColor).
		AddText("?: help \t q:quit", false, tview.AlignCenter, p.theme.PrimaryTextColor)
	return frame
}

func (p *BoardPage) down() {
	activeList := p.lists[p.activeListIdx]
	curIdx := activeList.GetCurrentItem()
	listLen := activeList.GetItemCount()
	if listLen == 0 {
		return
	}
	newIdx := (curIdx + 1) % listLen
	p.activeTaskIdxs[p.activeListIdx] = newIdx
	p.lists[p.activeListIdx].SetCurrentItem(newIdx)
}

func (p *BoardPage) up() {
	activeList := p.lists[p.activeListIdx]
	curIdx := activeList.GetCurrentItem()
	listLen := activeList.GetItemCount()
	if listLen == 0 {
		return
	}
	newIdx := (curIdx - 1 + listLen) % listLen
	p.activeTaskIdxs[p.activeListIdx] = newIdx
	p.lists[p.activeListIdx].SetCurrentItem(newIdx)
}

func (p *BoardPage) left() {
	listCount := len(p.lists)
	p.lists[p.activeListIdx].SetBorderColor(theme.PrimitiveBackgroundColor)
	p.activeListIdx = (p.activeListIdx - 1 + listCount) % listCount
	p.lists[p.activeListIdx].SetBorderColor(theme.ContrastBackgroundColor)
	app.SetFocus(p.lists[p.activeListIdx])
}

func (p *BoardPage) right() {
	listCount := len(p.lists)
	p.lists[p.activeListIdx].SetBorderColor(theme.PrimitiveBackgroundColor)
	p.activeListIdx = (p.activeListIdx + 1) % listCount
	p.lists[p.activeListIdx].SetBorderColor(theme.ContrastBackgroundColor)
	app.SetFocus(p.lists[p.activeListIdx])
}

func (p *BoardPage) moveDown() {
	activeListIdx := p.activeListIdx
	taskCount, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	activeTaskIdx := p.activeTaskIdxs[activeListIdx]
	if activeTaskIdx+1 >= taskCount {
		return
	}
	p.swapListItem(activeListIdx, activeTaskIdx, activeTaskIdx+1)
	p.data.Save()
	p.redraw(p.activeListIdx)
	p.down()
}

func (p *BoardPage) moveUp() {
	activeListIdx := p.activeListIdx
	activeTaskIdx := p.activeTaskIdxs[activeListIdx]
	if activeTaskIdx == 0 {
		return
	}
	p.swapListItem(activeListIdx, activeTaskIdx, activeTaskIdx-1)
	p.data.Save()
	p.redraw(p.activeListIdx)
	p.up()
}

func (p *BoardPage) moveLeft() {
	activeListIdx := p.activeListIdx
	if activeListIdx == 0 {
		return
	}
	taskCount, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if taskCount == 0 {
		return
	}
	p.moveTask(p.activeTaskIdxs[activeListIdx], activeListIdx, activeListIdx-1)
	p.data.Save()
	if err := p.fixActiveTaskIdx(); err != nil {
		app.Stop()
		log.Fatal(err)
	}
	p.redraw(p.activeListIdx)
	p.left()
	lastIdx, err := p.data.GetTaskCount(p.activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if lastIdx > 0 {
		p.activeTaskIdxs[p.activeListIdx] = lastIdx - 1
	}
	p.redraw(p.activeListIdx)
}

func (p *BoardPage) moveRight() {
	activeListIdx := p.activeListIdx
	listCount := len(p.lists)
	if activeListIdx+1 >= listCount {
		return
	}
	taskCount, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if taskCount == 0 {
		return
	}
	p.moveTask(p.activeTaskIdxs[activeListIdx], activeListIdx, activeListIdx+1)
	p.data.Save()
	p.redraw(p.activeListIdx)
	if err := p.fixActiveTaskIdx(); err != nil {
		app.Stop()
		log.Fatal(err)
	}
	p.right()
	lastIdx, err := p.data.GetTaskCount(p.activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if lastIdx > 0 {
		p.activeTaskIdxs[p.activeListIdx] = lastIdx - 1
	}
	p.redraw(p.activeListIdx)
}

func (p *BoardPage) redraw(listIdx int) {
	p.lists[listIdx].Clear()
	tasks, err := p.data.GetTasks(listIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	for _, item := range tasks {
		p.lists[listIdx].AddItem(item, "", 0, nil)
	}
	activeListIdx := p.activeListIdx
	p.lists[activeListIdx].SetCurrentItem(p.activeTaskIdxs[activeListIdx])
}

func (p *BoardPage) removeTask() {
	activeListIdx := p.activeListIdx
	taskCount, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if taskCount < 1 {
		return
	}
	removeTaskIdx := p.activeTaskIdxs[activeListIdx]
	removeTaskCommand := command.CreateRemoveTaskCommand(activeListIdx, removeTaskIdx)
	if err := p.command.Execute(removeTaskCommand); err != nil {
		app.Stop()
		log.Fatal(err)
	}
	p.data.Save()
	p.redraw(activeListIdx)
}

func (p *BoardPage) taskCompleted() {
	activeListIdx := p.activeListIdx
	activeTaskIdx := p.activeTaskIdxs[activeListIdx]
	listCount := p.data.GetListCount()
	taskDoneIdx := listCount - 1
	if activeListIdx == taskDoneIdx {
		return
	}
	taskCount, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	if taskCount <= 0 {
		return
	}
	p.moveTask(activeTaskIdx, activeListIdx, taskDoneIdx)
	p.redraw(activeListIdx)
	p.redraw(taskDoneIdx)
	if err := p.fixActiveTaskIdx(); err != nil {
		app.Stop()
		log.Fatal(err)
	}
}

func (p *BoardPage) fixActiveTaskIdx() error {
	taskCount, err := p.data.GetTaskCount(p.activeListIdx)
	if err != nil {
		return err
	}
	if taskCount > 0 && p.activeTaskIdxs[p.activeListIdx] >= taskCount {
		p.activeTaskIdxs[p.activeListIdx]--
	}
	return nil
}

func (p *BoardPage) setInputCapture(i int) {
	p.lists[i].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			p.down()
			return nil
		case tcell.KeyUp:
			p.up()
			return nil
		case tcell.KeyRight:
			p.right()
			return nil
		case tcell.KeyLeft:
			p.left()
			return nil
		case tcell.KeyCtrlR:
			p.redo()
		}
		switch event.Rune() {
		case 'j', tcell.RuneDArrow:
			p.down()
		case 'k', tcell.RuneUArrow:
			p.up()
		case 'h', tcell.RuneLArrow:
			p.left()
		case 'l', tcell.RuneRArrow:
			p.right()
		case 'J':
			p.moveDown()
		case 'K':
			p.moveUp()
		case 'H':
			p.moveLeft()
		case 'L':
			p.moveRight()
		case 'a':
			taskPos := p.activeTaskIdxs[p.activeListIdx]
			pages.AddAndSwitchToPage("add", NewAddPage(p, taskPos), true)
		case 'A':
			lastTaskPos, err := p.data.GetTaskCount(p.activeListIdx)
			if err != nil {
				app.Stop()
				panic(err)
			}
			pages.AddAndSwitchToPage("add", NewAddPage(p, lastTaskPos-1), true)
		case 'D':
			p.removeTask()
		case 'd':
			p.taskCompleted()
		case 'e':
			pages.AddAndSwitchToPage("edit", NewEditPage(p, p.activeListIdx, p.activeTaskIdxs[p.activeListIdx]), true)
		case 'u':
			p.undo()
		case 'q':
			p.data.Save()
			app.Stop()
		case '?':
			pages.AddAndSwitchToPage("help", NewHelpPage(p), true)
		case rune(tcell.KeyEnter):
			pages.AddAndSwitchToPage("info", NewInfoPage(p, p.activeListIdx, p.activeTaskIdxs[p.activeListIdx]), true)
		case 'g':
			p.focusFirst()
		case 'G':
			p.focusLast()
		default:
		}
		return event
	})
}

func (p *BoardPage) addTasksToList(listIdx int) {
	tasks, err := p.data.GetTasks(listIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	for _, item := range tasks {
		p.lists[listIdx].AddItem(item, "", 0, nil)
	}
}

func (p *BoardPage) focusFirst() {
	activeListIdx := p.activeListIdx
	if p.activeTaskIdxs[activeListIdx] == 0 {
		return
	}
	p.activeTaskIdxs[activeListIdx] = 0
	p.redraw(activeListIdx)
}

func (p *BoardPage) focusLast() {
	activeListIdx := p.activeListIdx
	lastIdx, err := p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if p.activeTaskIdxs[activeListIdx] == lastIdx {
		return
	}
	if lastIdx != 0 {
		lastIdx--
	}
	p.activeTaskIdxs[activeListIdx] = lastIdx
	p.redraw(activeListIdx)
}

func (p *BoardPage) redrawAll() {
	listCount := len(p.lists)
	for listIdx := 0; listIdx < listCount; listIdx++ {
		p.redraw(listIdx)
	}
}

func (p *BoardPage) swapListItem(listIdx, taskIdxFirst, taskIdxSecond int) {
	swapListItemCommand := command.CreateSwapListItemCommand(listIdx, taskIdxFirst, taskIdxSecond)
	if err := p.command.Execute(swapListItemCommand); err != nil {
		app.Stop()
		log.Fatal(err)
	}
}

func (p *BoardPage) moveTask(prevTaskIdx, prevListIdx, newListIdx int) {
	moveTaskCommand := command.CreateMoveTaskCommand(prevTaskIdx, prevListIdx, newListIdx)
	if err := p.command.Execute(moveTaskCommand); err != nil {
		app.Stop()
		log.Fatal(err)
	}
}

func (p *BoardPage) undo() {
	if err := p.command.Undo(); err != nil {
		app.Stop()
		panic(err)
	}
	p.redrawAll()
}

func (p *BoardPage) redo() {
	if err := p.command.Redo(); err != nil {
		app.Stop()
		log.Fatal(err)
	}
	p.redrawAll()
}
