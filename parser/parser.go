package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/singalhimanshu/taskgo/files"
)

const fileName = "/taskgo.md"

// A Data represents the board name and a slice of list.
type Data struct {
	boardName string
	lists     []List
}

// A List represents the title of list and a list of items inside it (i.e tasks).
type List struct {
	listTitle string
	listItems []ListItem
}

// A ListItem represents the name of item and it's description.
type ListItem struct {
	itemName        string
	itemDescription string
}

// ParseData parses the contents of the file (taskgo.md) to custom type Data
// It returns an error if the syntax of file is incorrect
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
		} else if strings.HasPrefix(line, "> ") {
			listLen := len(d.lists)

			if listLen < 1 {
				return fmt.Errorf("Error at line %v", lineNumber)
			}

			currentList := d.lists[listLen-1]
			itemDescStartIndex := strings.Index(line, " ") + 1
			itemDesc := line[itemDescStartIndex:]

			listItemLen := len(currentList.listItems)

			if listItemLen < 1 {
				return fmt.Errorf("Error at line %v", lineNumber)
			}

			currentList.listItems[listItemLen-1].itemDescription = itemDesc

			d.lists[listLen-1] = currentList
		} else {
			return fmt.Errorf("Error at line %v", lineNumber)
		}
	}
	return nil
}

// GetBoardName returns the name of board.
func (d *Data) GetBoardName() string {
	return d.boardName
}

// GetListNames returns a list of all the list names.
// Example: ["TODO", "DOING", "DONE"]
func (d *Data) GetListNames() []string {

	var listNames []string

	for _, list := range d.lists {
		listNames = append(listNames, list.listTitle)
	}

	return listNames
}

// GetTasks returns a list of all the tasks of a particular list.
// Example: ["Task 1", "Task 2"]
func (d *Data) GetTasks(listIdx int) []string {
	var tasks []string

	for _, item := range d.lists[listIdx].listItems {
		tasks = append(tasks, item.itemName)
	}

	return tasks
}

// AddNewTask adds a new task to a list provided the list index and the title of that task.
// It returns an error if the index is out of bounds.
func (d *Data) AddNewTask(listIdx int, taskTitle, taskDesc string) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}

	d.lists[listIdx].listItems = append(d.lists[listIdx].listItems, ListItem{
		itemName:        taskTitle,
		itemDescription: taskDesc,
	})
	return nil
}

// MoveTask moves a task from one list to another.
// It returns an error if any of the index is out of bounds.
func (d *Data) MoveTask(prevTaskIdx, prevListIdx, newListIdx int) error {
	listCount := d.GetListCount()
	if err := checkBounds(prevListIdx, listCount); err != nil {
		return err
	}
	if err := checkBounds(newListIdx, listCount); err != nil {
		return err
	}

	taskCount, err := d.GetTaskCount(prevListIdx)
	if err != nil {
		return err
	}

	if err := checkBounds(prevTaskIdx, taskCount); err != nil {
		return err
	}

	taskTitle := d.lists[prevListIdx].listItems[prevTaskIdx].itemName
	taskDesc := d.lists[prevListIdx].listItems[prevTaskIdx].itemDescription
	d.AddNewTask(newListIdx, taskTitle, taskDesc)
	err = d.RemoveTask(prevListIdx, prevTaskIdx)
	return err
}

// RemoveTask removes a task given the index of list and the task.
// It returns an error if any of the index is out of bounds.
func (d *Data) RemoveTask(listIdx, taskIdx int) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}

	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}

	if err := checkBounds(taskIdx, taskCount); err != nil {
		return fmt.Errorf("Index out of bounds(task): %v", taskIdx)
	}

	d.lists[listIdx].listItems = append(d.lists[listIdx].listItems[:taskIdx],
		d.lists[listIdx].listItems[taskIdx+1:]...)
	return nil
}

// Save saves the content of Data onto the file (taskgo.md).
func (d *Data) Save() {
	var fileContent []string
	fileContent = append(fileContent, "# "+d.boardName+"\n")

	for _, list := range d.lists {
		fileContent = append(fileContent, "## "+list.listTitle)
		for _, listItem := range list.listItems {
			fileContent = append(fileContent, "\t- "+listItem.itemName)
			fileContent = append(fileContent, "\t\t> "+listItem.itemDescription)
		}
		fileContent = append(fileContent, "\n")
	}

	err := files.WriteFile(fileContent, fileName)
	if err != nil {
		log.Fatal(err)
	}
}

// SwapListItems swaps a task of one list with another list.
// It returns an error if any of the index is out of bounds
func (d *Data) SwapListItems(listIdx, taskIdxFirst, taskIdxSecond int) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return fmt.Errorf("Index out of bounds (list): %v", listIdx)
	}

	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	err = checkBounds(taskIdxFirst, taskCount)
	if err != nil {
		return err
	}
	err = checkBounds(taskIdxSecond, taskCount)
	if err != nil {
		return err
	}

	swap(&d.lists[listIdx].listItems[taskIdxFirst],
		&d.lists[listIdx].listItems[taskIdxSecond])

	return nil
}

// GetTaskCount gives the task count of a particular list.
// It returns an int(number of tasks) and an error(if the list index is out of bounds).
func (d *Data) GetTaskCount(listIdx int) (int, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return 0, err
	}
	return len(d.lists[listIdx].listItems), nil
}

func (d *Data) GetListCount() int {
	return len(d.lists)
}

func swap(first, second *ListItem) {
	*second, *first = *first, *second
}

func checkBounds(idx, boundary int) error {
	if idx < 0 || idx >= boundary {
		return fmt.Errorf("Index Out of Bounds: %v", idx)
	}
	return nil
}
