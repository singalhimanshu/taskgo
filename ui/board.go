package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/parser"
)

// BoardPage contains all the information that the board page requires.
type BoardPage struct {
	lists          []*tview.List
	theme          *tview.Theme
	data           parser.Data
	activeListIdx  int
	activeTaskIdxs []int
	fileName       string
}

// NewBoardPage adds the data to BoardPage structure.
func NewBoardPage(fileName string) *BoardPage {
	theme := defaultTheme()

	data := parser.Data{}
	if err := data.ParseData(fileName); err != nil {
		log.Fatal(err)
	}

	listCount := len(data.GetListNames())

	return &BoardPage{
		lists:          make([]*tview.List, listCount),
		data:           data,
		theme:          theme,
		activeListIdx:  0,
		activeTaskIdxs: make([]int, listCount),
		fileName:       fileName,
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
		if i == 0 {
			p.lists[i].SetBorderColor(theme.ContrastBackgroundColor)
		}
		p.lists[i].SetTitle(listNames[i])

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
			}
			switch event.Rune() {
			case 'j':
				p.down()
			case 'k':
				p.up()
			case 'h':
				p.left()
			case 'l':
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
				pages.AddAndSwitchToPage("add", NewAddPage(p), true)
			case 'D':
				p.removeTask()
			case 'C':
				pages.AddAndSwitchToPage("edit", NewEditPage(p, p.activeListIdx, p.activeTaskIdxs[p.activeListIdx]), true)
			case 'q':
				p.data.Save(p.fileName)
				app.Stop()
			case '?':
				pages.AddAndSwitchToPage("help", NewHelpPage(p), true)
			case ' ':
				pages.AddAndSwitchToPage("info", NewInfoPage(p, p.activeListIdx, p.activeTaskIdxs[p.activeListIdx]), true)
			default:
			}
			return event
		})

		for _, item := range p.data.GetTasks(i) {
			p.lists[i].AddItem(item, "", 0, nil)
		}

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

	err = p.data.SwapListItems(activeListIdx,
		activeTaskIdx,
		activeTaskIdx+1)

	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	p.data.Save(p.fileName)
	p.redraw()
	p.down()
}

func (p *BoardPage) moveUp() {
	activeListIdx := p.activeListIdx
	activeTaskIdx := p.activeTaskIdxs[activeListIdx]
	if activeTaskIdx == 0 {
		return
	}

	err := p.data.SwapListItems(activeListIdx,
		activeTaskIdx,
		activeTaskIdx-1)

	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	p.data.Save(p.fileName)
	p.redraw()
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
	err = p.data.MoveTask(p.activeTaskIdxs[activeListIdx],
		activeListIdx, activeListIdx-1)

	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	p.data.Save(p.fileName)
	taskCount, err = p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if taskCount > 0 && p.activeTaskIdxs[p.activeListIdx] >= taskCount {
		p.activeTaskIdxs[p.activeListIdx]--
	}

	p.redraw()
	p.left()
	lastIdx, err := p.data.GetTaskCount(p.activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if lastIdx > 0 {
		p.activeTaskIdxs[p.activeListIdx] = lastIdx - 1
	}
	p.redraw()
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

	err = p.data.MoveTask(p.activeTaskIdxs[activeListIdx],
		activeListIdx, activeListIdx+1)

	if err != nil {
		app.Stop()
		log.Fatal(err)
	}

	p.data.Save(p.fileName)
	p.redraw()
	taskCount, err = p.data.GetTaskCount(activeListIdx)
	if err != nil {
		app.Stop()
		panic(err)
	}
	if taskCount > 0 && p.activeTaskIdxs[p.activeListIdx] >= taskCount {
		p.activeTaskIdxs[p.activeListIdx]--
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
	p.redraw()
}

func (p *BoardPage) redraw() {
	activeListIdx := p.activeListIdx
	p.lists[activeListIdx].Clear()
	tasks := p.data.GetTasks(activeListIdx)
	for _, item := range tasks {
		p.lists[activeListIdx].AddItem(item, "", 0, nil)
	}
	p.lists[activeListIdx].SetCurrentItem(p.activeTaskIdxs[activeListIdx])
}

func (p *BoardPage) removeTask() {
	activeListIdx := p.activeListIdx
	removeTaskIdx := p.activeTaskIdxs[activeListIdx]
	err := p.data.RemoveTask(activeListIdx, removeTaskIdx)
	if err != nil {
		app.Stop()
		log.Fatal(err)
	}
	p.data.Save(p.fileName)
	p.redraw()
}
