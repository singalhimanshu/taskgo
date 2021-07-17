package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/singalhimanshu/taskgo/pkg/files"
)

// A Data represents the board name and a slice of list.
type Data struct {
	boardName string
	lists     []List
	fileName  string
}

// A List represents the title of list and a list of items inside it (i.e tasks).
type List struct {
	listTitle string
	listItems []ListItem
}

// A ListItem represents the name of item and it's description.
type ListItem struct {
	ItemName        string
	ItemDescription string
	Subtasks        []string
}

func (d *Data) SetFileName(fileName string) {
	d.fileName = fileName
}

func (d *Data) GetContentFromFile() []string {
	if !files.CheckFile(d.fileName) {
		files.CreateFile(d.fileName)
	}
	return files.OpenFile(d.fileName)
}

// ParseData parses the contents of the file to custom type Data
// It returns an error if the syntax of file is incorrect
func (d *Data) ParseData(fileContent []string) error {
	for lineNumber, line := range fileContent {
		line = strings.TrimSpace(line)
		// skip empty lines
		if len(line) < 1 {
			continue
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
			listCount := d.GetListCount()
			if listCount < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList := d.lists[listCount-1]
			itemNameStartIndex := strings.Index(line, " ") + 1
			itemName := line[itemNameStartIndex:]
			currentList.listItems = append(currentList.listItems, ListItem{
				ItemName: itemName,
			})
			d.lists[listCount-1] = currentList
		} else if strings.HasPrefix(line, "> ") {
			listCount := d.GetListCount()
			if listCount < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList := d.lists[listCount-1]
			itemDescStartIndex := strings.Index(line, " ") + 1
			itemDesc := line[itemDescStartIndex:]
			listItemLen := len(currentList.listItems)
			if listItemLen < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList.listItems[listItemLen-1].ItemDescription = itemDesc
			d.lists[listCount-1] = currentList
		} else if strings.HasPrefix(line, "+ ") {
			listCount := d.GetListCount()
			if listCount < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList := d.lists[listCount-1]
			subtaskStartIndex := strings.Index(line, " ") + 1
			subtask := line[subtaskStartIndex:]
			listItemLen := len(currentList.listItems)
			if listItemLen < 1 {
				return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
			}
			currentList.listItems[listItemLen-1].Subtasks = append(currentList.listItems[listItemLen-1].Subtasks, subtask)
		} else {
			return fmt.Errorf("Error at line %v of file %v\n Line: %v", lineNumber, d.fileName, line)
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

// GetTask gives the task title and description given the list index and
// task index. It returns an array of string and error if any of the
// index are out of bounds.
func (d *Data) GetTask(listIdx, taskIdx int) (ListItem, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return ListItem{}, err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return ListItem{}, err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return ListItem{}, err
	}
	return d.lists[listIdx].listItems[taskIdx], nil
}

func (d *Data) GetSubtask(listIdx, taskIdx, subtaskIdx int) (string, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return "", nil
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return "", err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return "", err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return "", nil
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return "", err
	}
	return d.lists[listIdx].listItems[taskIdx].Subtasks[subtaskIdx], nil
}

// GetTasks returns a list of all the tasks of a particular list.
// Example: ["Task 1", "Task 2"]
func (d *Data) GetTasks(listIdx int) ([]string, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return nil, err
	}
	var tasks []string
	for _, item := range d.lists[listIdx].listItems {
		tasks = append(tasks, item.ItemName)
	}
	return tasks, nil
}

// GetSubtasks return a list of strings and an error that
// contains the name of the subtask. Example: ["Subtask 1",
// Subtask 2"]. If the listIdx or taskIdx is out of bounds,
// nil and an error is returned.
func (d *Data) GetSubtasks(listIdx, taskIdx int) ([]string, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return nil, err
	}
	return d.lists[listIdx].listItems[taskIdx].Subtasks, nil
}

// AddNewTask adds a new task to a list provided the list index and the title of that task.
// It returns an error if the index is out of bounds.
func (d *Data) AddNewTask(listIdx int, taskTitle, taskDesc string, taskPos int) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	newTask := ListItem{
		ItemName:        taskTitle,
		ItemDescription: taskDesc,
	}
	d.insertTask(listIdx, newTask, taskPos)
	d.Save()
	return nil
}

func (d *Data) AddSubtask(listIdx, taskIdx int, subtask string) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return err
	}
	d.lists[listIdx].listItems[taskIdx].Subtasks = append(d.lists[listIdx].listItems[taskIdx].Subtasks, subtask)
	return nil
}

// EditTask edits a task title and description given the index of (list
// and task), task title and description. It returns an error if the
// index are out of bounds.
func (d *Data) EditTask(listIdx, taskIdx int, taskTitle, taskDesc string) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return err
	}
	d.lists[listIdx].listItems[taskIdx].ItemName = taskTitle
	d.lists[listIdx].listItems[taskIdx].ItemDescription = taskDesc
	d.Save()
	return nil
}

func (d *Data) EditSubtask(listIdx, taskIdx, subtaskIdx int, newSubtask string) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return err
	}
	d.lists[listIdx].listItems[taskIdx].Subtasks[subtaskIdx] = newSubtask
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
	newListTaskCount, err := d.GetTaskCount(newListIdx)
	if err != nil {
		return err
	}
	taskTitle := d.lists[prevListIdx].listItems[prevTaskIdx].ItemName
	taskDesc := d.lists[prevListIdx].listItems[prevTaskIdx].ItemDescription
	err = d.AddNewTask(newListIdx, taskTitle, taskDesc, newListTaskCount)
	if err != nil {
		return err
	}
	_, err = d.RemoveTask(prevListIdx, prevTaskIdx)
	return err
}

func (d *Data) MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx int) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(firstSubtaskIdx, subtaskCount); err != nil {
		return err
	}
	if err := checkBounds(secondSubtaskIdx, subtaskCount); err != nil {
		return err
	}
	// TODO(singalhimanshu): Refactor this (idea: make variables for listItem and so)
	d.lists[listIdx].listItems[taskIdx].Subtasks[firstSubtaskIdx],
		d.lists[listIdx].listItems[taskIdx].Subtasks[secondSubtaskIdx] =
		d.lists[listIdx].listItems[taskIdx].Subtasks[secondSubtaskIdx],
		d.lists[listIdx].listItems[taskIdx].Subtasks[firstSubtaskIdx]
	return nil
}

// RemoveTask removes a task given the index of list and the task.
// It returns an error if any of the index is out of bounds.
func (d *Data) RemoveTask(listIdx, taskIdx int) (ListItem, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return ListItem{}, err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return ListItem{}, err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return ListItem{}, fmt.Errorf("Index out of bounds(task): %v", taskIdx)
	}
	taskData := d.lists[listIdx].listItems[taskIdx]
	d.lists[listIdx].listItems = append(d.lists[listIdx].listItems[:taskIdx],
		d.lists[listIdx].listItems[taskIdx+1:]...)
	d.Save()
	return taskData, nil
}

func (d *Data) RemoveSubtask(listIdx, taskIdx, subtaskIdx int) error {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return err
	}
	d.lists[listIdx].listItems[taskIdx].Subtasks = append(
		d.lists[listIdx].listItems[taskIdx].Subtasks[:subtaskIdx],
		d.lists[listIdx].listItems[taskIdx].Subtasks[subtaskIdx+1:]...,
	)
	return nil
}

// Save saves the content of Data onto the file.
func (d *Data) Save() {
	if d.fileName == "" {
		return
	}
	var fileContent []string
	fileContent = append(fileContent, "# "+d.boardName+"\n")
	for _, list := range d.lists {
		fileContent = append(fileContent, "## "+list.listTitle)
		for _, listItem := range list.listItems {
			fileContent = append(fileContent, "\t- "+listItem.ItemName)
			if len(listItem.ItemDescription) > 0 {
				fileContent = append(fileContent, "\t\t> "+listItem.ItemDescription)
			}
		}
		fileContent = append(fileContent, "\n")
	}
	err := files.WriteFile(fileContent, d.fileName)
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
	d.Save()
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

func (d *Data) GetSubtaskCount(listIdx, taskIdx int) (int, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return 0, err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return 0, err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return 0, err
	}
	return len(d.lists[listIdx].listItems[taskIdx].Subtasks), nil
}

// GetListCount returns the count of lists.
func (d *Data) GetListCount() int {
	return len(d.lists)
}

func (d *Data) insertTask(listIdx int, task ListItem, taskPos int) {
	if len(d.lists[listIdx].listItems) < 1 {
		d.lists[listIdx].listItems = append(d.lists[listIdx].listItems, task)
		return
	}

	d.lists[listIdx].listItems = append(d.lists[listIdx].listItems, ListItem{})
	copy(d.lists[listIdx].listItems[(taskPos+1):],
		d.lists[listIdx].listItems[taskPos:])
	d.lists[listIdx].listItems[taskPos] = task
}

func swap(first, second *ListItem) {
	*second, *first = *first, *second
}

func checkBounds(idx, boundary int) error {
	if idx < 0 || idx >= boundary {
		return fmt.Errorf("Index Out of Bounds: got %v, length: %v", idx, boundary)
	}
	return nil
}
