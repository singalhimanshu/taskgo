package parser

import (
	"fmt"
	"strings"

	"github.com/singalhimanshu/taskgo/files"
)

const fileName = "/taskgo.md"

type Data struct {
	boardName string
	lists     []List
}

type List struct {
	listTitle string
	listItems []ListItem
}

type ListItem struct {
	itemName        string
	itemDescription string
}

func (d *Data) ParseData() error {
	fileFound := files.CheckFile()
	if !fileFound {
		files.CreateFile()
	}
	fileContent := files.OpenFile(fileName)

	for lineNumber, line := range fileContent {
		line = strings.TrimSpace(line)

		// skip empty lines
		if len(line) < 1 {
			continue
		}

		if !files.CheckPrefix(line) {
			return fmt.Errorf("Error at line %v", lineNumber)
		}

		if strings.HasPrefix(line, "# ") {
			boardNameStartingIndex := strings.Index(line, " ") + 1
			boardName := line[boardNameStartingIndex:]

			d.boardName = boardName

		} else if strings.HasPrefix(line, "## ") {

			listNameStartIndex := strings.Index(line, " ") + 1
			listTitle := line[listNameStartIndex:]

			d.lists = append(d.lists, List{
				listTitle: listTitle,
			})

		} else if strings.HasPrefix(line, "- ") {
			listLen := len(d.lists)

			if listLen < 1 {
				return fmt.Errorf("Error at line %v", lineNumber)
			}

			currentList := d.lists[listLen-1]
			itemNameStartIndex := strings.Index(line, " ") + 1
			itemName := line[itemNameStartIndex:]

			currentList.listItems = append(currentList.listItems, ListItem{
				itemName: itemName,
			})

			d.lists[listLen-1] = currentList
		} else {
			return fmt.Errorf("Error at line %v", lineNumber)
		}
	}
	return nil
}

func (d *Data) GetBoardName() string {
	return d.boardName
}

func (d *Data) GetListNames() []string {

	var listNames []string

	for _, list := range d.lists {
		listNames = append(listNames, list.listTitle)
	}

	return listNames
}

func (d *Data) GetTasks(idx int) []string {
	var tasks []string

	for _, item := range d.lists[idx].listItems {
		tasks = append(tasks, item.itemName)
	}

	return tasks
}
