package ui

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/singalhimanshu/taskgo/parser"
)

type BoardPage struct {
	mainArea []*tview.List
	theme    *tview.Theme
}

func NewBoardPage() *BoardPage {
	theme := defaultTheme()

	return &BoardPage{
		mainArea: make([]*tview.List, 3),
		theme:    theme,
	}
}

func (p *BoardPage) Page() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	for i := 0; i < 3; i++ {
		p.mainArea[i] = tview.NewList()
		p.mainArea[i].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			globalInputCapture(event)

			switch event.Key() {
			case tcell.KeyUp:
				fmt.Println("key up pressed")
			}
			return event
		})

		p.mainArea[i].AddItem("testing", "", 0, nil).SetBorder(true)

		flex.AddItem(p.mainArea[i], 0, 1, i == 0)
	}

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		globalInputCapture(event)
		return event
	})
	boardName := parser.GetBoardName()
	boardName = "Board: " + boardName

	frame := tview.NewFrame(flex).
		SetBorders(0, 0, 1, 0, 1, 1).
		AddText(boardName, true, tview.AlignCenter, p.theme.TitleColor).
		AddText("?: help \t q:quit", false, tview.AlignCenter, p.theme.PrimaryTextColor)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		globalInputCapture(event)
		log.Println(event.Rune())
		return event
	})

	return frame
}
