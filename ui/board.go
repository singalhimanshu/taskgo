package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/parser"
)

type BoardPage struct {
	lists          []*tview.List
	theme          *tview.Theme
	data           parser.Data
	activeListIdx  int
	activeTaskIdxs []int
}

func NewBoardPage() *BoardPage {
	theme := defaultTheme()

	data := parser.Data{}
	if err := data.ParseData(); err != nil {
		log.Fatal(err)
	}

	listCount := len(data.GetListNames())

	return &BoardPage{
		lists:          make([]*tview.List, listCount),
		data:           data,
		theme:          theme,
		activeListIdx:  0,
		activeTaskIdxs: make([]int, listCount),
	}
}

func (p *BoardPage) Page() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	listNames := p.data.GetListNames()

	for i := 0; i < len(listNames); i++ {
		p.lists[i] = tview.NewList()

		p.lists[i].
			ShowSecondaryText(false).
			SetBorder(true).
			SetBorderColor(theme.BorderColor)

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
			case 'a':
				pages.AddAndSwitchToPage("add", NewAddPage(p), true)
			case 'D':
				p.removeTask()
			case 'C':
			case 'q':
				p.data.Save()
				app.Stop()
			case '?':
				pages.SwitchToPage("help")
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
	newIdx := (curIdx + 1) % listLen
	p.activeTaskIdxs[p.activeListIdx] = newIdx
	p.lists[p.activeListIdx].SetCurrentItem(newIdx)
}

func (p *BoardPage) up() {
	activeList := p.lists[p.activeListIdx]
	curIdx := activeList.GetCurrentItem()
	listLen := activeList.GetItemCount()
	newIdx := (curIdx - 1 + listLen) % listLen
	p.activeTaskIdxs[p.activeListIdx] = newIdx
	p.lists[p.activeListIdx].SetCurrentItem(newIdx)
}

func (p *BoardPage) left() {
	listCount := len(p.lists)
	p.activeListIdx = (p.activeListIdx - 1 + listCount) % listCount
	app.SetFocus(p.lists[p.activeListIdx])
}

func (p *BoardPage) right() {
	listCount := len(p.lists)
	p.activeListIdx = (p.activeListIdx + 1) % listCount
	app.SetFocus(p.lists[p.activeListIdx])
}

func (p *BoardPage) redraw() {
	activeListIdx := p.activeListIdx
	p.lists[activeListIdx].Clear()
	tasks := p.data.GetTasks(activeListIdx)
	for _, item := range tasks {
		p.lists[activeListIdx].AddItem(item, "", 0, nil)
	}
}

func (p *BoardPage) removeTask() {
	activeListIdx := p.activeListIdx
	removeTaskIdx := p.activeTaskIdxs[activeListIdx]
	err := p.data.RemoveTask(activeListIdx, removeTaskIdx)
	if err != nil {
		log.Fatal(err)
	}
	p.data.Save()
	p.redraw()
}
