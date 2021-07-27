package ui

import (
	"github.com/rivo/tview"
	"golang.org/x/term"
)

func GetCenteredModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func GetSize() (width, height int) {
	width, height, err := term.GetSize(0)
	if err != nil {
		return 0, 0
	}
	return width, height
}
