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

func (d *Data) GetList(listIdx int) (*List, error) {
	listCount := d.GetListCount()
	if err := checkBounds(listIdx, listCount); err != nil {
		return nil, err
	}
	return &d.lists[listIdx], nil
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
func (d *Data) GetTask(listIdx, taskIdx int) (*ListItem, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return nil, err
	}
	taskCount, err := d.GetTaskCount(listIdx)
	if err != nil {
		return nil, err
	}
	if err := checkBounds(taskIdx, taskCount); err != nil {
		return nil, err
	}
	return &list.listItems[taskIdx], nil
}

func (d *Data) GetSubtask(listIdx, taskIdx, subtaskIdx int) (string, error) {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return "", err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return "", nil
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return "", err
	}
	return task.Subtasks[subtaskIdx], nil
}

// GetTasks returns a list of all the tasks of a particular list.
// Example: ["Task 1", "Task 2"]
func (d *Data) GetTasks(listIdx int) ([]string, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return nil, err
	}
	var tasks []string
	for _, item := range list.listItems {
		tasks = append(tasks, item.ItemName)
	}
	return tasks, nil
}

// GetSubtasks return a list of strings and an error that
// contains the name of the subtask. Example: ["Subtask 1",
// Subtask 2"]. If the listIdx or taskIdx is out of bounds,
// nil and an error is returned.
func (d *Data) GetSubtasks(listIdx, taskIdx int) ([]string, error) {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return nil, err
	}
	return task.Subtasks, nil
}

// AddNewTask adds a new task to a list provided the list index and the title of that task.
// It returns an error if the index is out of bounds.
func (d *Data) AddNewTask(listIdx int, taskTitle, taskDesc string, taskIdx int) error {
	if err := checkBounds(listIdx, d.GetListCount()); err != nil {
		return err
	}
	newTask := ListItem{
		ItemName:        taskTitle,
		ItemDescription: taskDesc,
	}
	err := d.insertTask(listIdx, newTask, taskIdx)
	if err != nil {
		return err
	}
	d.Save()
	return nil
}

func (d *Data) AddSubtask(listIdx, taskIdx int, subtask string) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return err
	}
	task.Subtasks = append(task.Subtasks, subtask)
	return nil
}

// EditTask edits a task title and description given the index of (list
// and task), task title and description. It returns an error if the
// index are out of bounds.
func (d *Data) EditTask(listIdx, taskIdx int, taskTitle, taskDesc string) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return err
	}
	task.ItemName = taskTitle
	task.ItemDescription = taskDesc
	d.Save()
	return nil
}

func (d *Data) EditSubtask(listIdx, taskIdx, subtaskIdx int, newSubtask string) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return err
	}
	task.Subtasks[subtaskIdx] = newSubtask
	return nil
}

// MoveTask moves a task from one list to another.
// It returns an error if any of the index is out of bounds.
func (d *Data) MoveTask(taskIdx, sourceListIdx, destListIdx int) error {
	sourceTask, err := d.GetTask(sourceListIdx, taskIdx)
	if err != nil {
		return err
	}
	newListTaskCount, err := d.GetTaskCount(destListIdx)
	if err != nil {
		return err
	}
	taskTitle := sourceTask.ItemName
	taskDesc := sourceTask.ItemDescription
	err = d.AddNewTask(destListIdx, taskTitle, taskDesc, newListTaskCount)
	if err != nil {
		return err
	}
	_, err = d.RemoveTask(sourceListIdx, taskIdx)
	return err
}

// RemoveTask removes a task given the index of list and the task.
// It returns an error if any of the index is out of bounds.
func (d *Data) RemoveTask(listIdx, taskIdx int) (ListItem, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return ListItem{}, err
	}
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return ListItem{}, err
	}
	taskData := *task
	list.listItems = append(list.listItems[:taskIdx], list.listItems[taskIdx+1:]...)
	d.Save()
	return taskData, nil
}

func (d *Data) RemoveSubtask(listIdx, taskIdx, subtaskIdx int) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return err
	}
	subtaskCount, err := d.GetSubtaskCount(listIdx, taskIdx)
	if err != nil {
		return err
	}
	if err := checkBounds(subtaskIdx, subtaskCount); err != nil {
		return err
	}
	task.Subtasks = append(task.Subtasks[:subtaskIdx], task.Subtasks[subtaskIdx+1:]...)
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
func (d *Data) SwapListItems(listIdx, firstTaskIdx, secondTaskIdx int) error {
	firstTask, err := d.GetTask(listIdx, firstTaskIdx)
	if err != nil {
		return err
	}
	secondTask, err := d.GetTask(listIdx, secondTaskIdx)
	if err != nil {
		return err
	}
	*firstTask, *secondTask = *secondTask, *firstTask
	d.Save()
	return nil
}

func (d *Data) SwapSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx int) error {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
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
	task.Subtasks[firstSubtaskIdx], task.Subtasks[secondSubtaskIdx] =
		task.Subtasks[secondSubtaskIdx], task.Subtasks[firstSubtaskIdx]
	return nil
}

// GetTaskCount gives the task count of a particular list.
// It returns an int(number of tasks) and an error(if the list index is out of bounds).
func (d *Data) GetTaskCount(listIdx int) (int, error) {
	list, err := d.GetList(listIdx)
	if err != nil {
		return 0, err
	}
	return len(list.listItems), nil
}

func (d *Data) GetSubtaskCount(listIdx, taskIdx int) (int, error) {
	task, err := d.GetTask(listIdx, taskIdx)
	if err != nil {
		return 0, err
	}
	return len(task.Subtasks), nil
}

// GetListCount returns the count of lists.
func (d *Data) GetListCount() int {
	return len(d.lists)
}

func (d *Data) insertTask(listIdx int, task ListItem, taskIdx int) error {
	list, err := d.GetList(listIdx)
	if err != nil {
		return err
	}
	if len(list.listItems) < 1 {
		list.listItems = append(list.listItems, task)
		return nil
	}

	list.listItems = append(list.listItems, ListItem{})
	copy(list.listItems[(taskIdx+1):], list.listItems[taskIdx:])
	list.listItems[taskIdx] = task
	return nil
}

func checkBounds(idx, boundary int) error {
	if idx < 0 || idx >= boundary {
		return fmt.Errorf("Index Out of Bounds: got %v, length: %v", idx, boundary)
	}
	return nil
}
