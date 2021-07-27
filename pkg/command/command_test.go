package command

import (
	"fmt"
	"log"
	"testing"

	"github.com/singalhimanshu/taskgo/pkg/parser"
)

func TestExecute(t *testing.T) {
	t.Run("Add task command", testAddTaskCommand(t))
	t.Run("Remove task Command", testRemoveTaskCommand(t))
	t.Run("Swap List Item Command", testSwapListItemCommand(t))
	t.Run("Move Task Command", testMoveTaskCommand(t))
	t.Run("Edit Task Command", testEditTaskCommand(t))
}

func testAddTaskCommand(t *testing.T) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := getNewData()
		if err != nil {
			t.Error(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			t.Error(err)
		}
		listIdx, taskIdx := 0, 0
		actualTaskTitle := "test1"
		actualTaskDesc := "test1 desc"
		addTaskCommand := CreateAddTaskCommand(listIdx, actualTaskTitle, actualTaskDesc, taskIdx)
		// Test Execute
		if err := testCommand.Execute(addTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		gotTaskData, err := testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData.ItemName, gotTaskData.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Undo
		if err := testCommand.Undo(); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		_, err = testData.GetTask(listIdx, taskIdx)
		if err == nil {
			t.Errorf("Expected error, but didn't get one: %v", err)
		}
		// Test Redo
		if err := testCommand.Redo(); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		gotTaskData, err = testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData.ItemName, gotTaskData.ItemDescription); err != nil {
			t.Error(err)
		}
	}
}

func testRemoveTaskCommand(t *testing.T) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := getNewData()
		if err != nil {
			t.Error(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			t.Error(err)
		}
		listIdx, taskIdx := 0, 0
		actualTaskTitle := "test1"
		actualTaskDesc := "test1 desc"
		addTaskCommand := CreateAddTaskCommand(listIdx, actualTaskTitle, actualTaskDesc, taskIdx)
		removeTaskCommand := CreateRemoveTaskCommand(listIdx, taskIdx)
		// Test Execute
		if err := testCommand.Execute(addTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := testCommand.Execute(removeTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		_, err = testData.GetTask(listIdx, taskIdx)
		if err == nil {
			t.Error("Expected error, but didn't get one")
		}
		// Test Undo
		if err := testCommand.Undo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		gotTaskData, err := testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData.ItemName, gotTaskData.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Redo
		if err := testCommand.Redo(); err != nil {
			t.Error(err)
		}
		_, err = testData.GetTask(listIdx, taskIdx)
		if err == nil {
			t.Error("Expected error, but didn't get one")
		}
	}
}

func testSwapListItemCommand(t *testing.T) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := getNewData()
		if err != nil {
			t.Error(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			t.Error(err)
		}
		firstListIdx, firstTaskIdx := 0, 0
		firstTaskTitle := "test1"
		firstTaskDesc := "test1 desc"
		firstAddTaskCommand := CreateAddTaskCommand(firstListIdx, firstTaskTitle, firstTaskDesc, firstTaskIdx)
		if err := testCommand.Execute(firstAddTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		secondListIdx, secondTaskIdx := 0, 1
		secondTaskTitle := "test2"
		secondTaskDesc := "test2 desc"
		secondAddTaskCommand := CreateAddTaskCommand(secondListIdx, secondTaskTitle, secondTaskDesc, secondTaskIdx)
		if err := testCommand.Execute(secondAddTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		swapListItemCommand := CreateSwapListItemCommand(firstListIdx, firstTaskIdx, secondTaskIdx)
		// Test Execute
		if err := testCommand.Execute(swapListItemCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		swappedTask, err := testData.GetTask(secondListIdx, secondTaskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask.ItemName, swappedTask.ItemDescription); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		// Test Undo
		if err := testCommand.Undo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		swappedTask, err = testData.GetTask(firstListIdx, firstTaskIdx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask.ItemName, swappedTask.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Redo
		if err := testCommand.Redo(); err != nil {
			t.Error(err)
		}
		swappedTask, err = testData.GetTask(secondListIdx, secondTaskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask.ItemName, swappedTask.ItemDescription); err != nil {
			t.Error(err)
		}
	}
}

func testMoveTaskCommand(t *testing.T) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := getNewData()
		if err != nil {
			t.Error(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			t.Error(err)
		}
		firstListIdx, firstTaskIdx := 0, 0
		firstTaskTitle := "test1"
		firstTaskDesc := "test1 desc"
		firstAddTaskCommand := CreateAddTaskCommand(firstListIdx, firstTaskTitle, firstTaskDesc, firstTaskIdx)
		secondListIdx := 1
		if err := testCommand.Execute(firstAddTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		moveTaskCommand := CreateMoveTaskCommand(firstTaskIdx, firstListIdx, secondListIdx)
		// Test Execute
		if err := testCommand.Execute(moveTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		movedTask, err := testData.GetTask(secondListIdx, firstTaskIdx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, movedTask.ItemName, movedTask.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Undo
		if err := testCommand.Undo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		movedTask, err = testData.GetTask(firstListIdx, firstTaskIdx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, movedTask.ItemName, movedTask.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Redo
		if err := testCommand.Redo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		movedTask, err = testData.GetTask(secondListIdx, firstListIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, movedTask.ItemName, movedTask.ItemDescription); err != nil {
			t.Error(err)
		}
	}
}

func testEditTaskCommand(t *testing.T) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := getNewData()
		if err != nil {
			t.Error(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			t.Error(err)
		}
		listIdx, taskIdx := 0, 0
		actualTaskTitle := "test1"
		actualTaskDesc := "test1 desc"
		addTaskCommand := CreateAddTaskCommand(listIdx, actualTaskTitle, actualTaskDesc, taskIdx)
		if err := testCommand.Execute(addTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		editedTaskTitle, editedTaskDesc := "test1 edited", "test1 desc edited"
		// Test Execute
		editTaskCommand := CreateEditTaskCommand(listIdx, taskIdx, editedTaskTitle, editedTaskDesc)
		if err := testCommand.Execute(editTaskCommand); err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		editedTask, err := testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(editedTaskTitle, editedTaskDesc, editedTask.ItemName, editedTask.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Undo
		if err := testCommand.Undo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		originalTask, err := testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, originalTask.ItemName, originalTask.ItemDescription); err != nil {
			t.Error(err)
		}
		// Test Redo
		if err := testCommand.Redo(); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		editedTask, err = testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
		if err := compareTaskData(editedTaskTitle, editedTaskDesc, editedTask.ItemName, editedTask.ItemDescription); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkAddTaskCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testData, err := getNewData()
		if err != nil {
			log.Fatal(err)
		}
		testCommand, err := getNewCommand(testData)
		if err != nil {
			log.Fatal(err)
		}
		listIdx, taskIdx := 0, 0
		taskTitle := "test"
		taskDesc := "test desc"
		for i := 0; i < 100; i++ {
			addTaskCommand := CreateAddTaskCommand(listIdx, taskTitle, taskDesc, taskIdx)
			_ = testCommand.Execute(addTaskCommand)
			taskIdx++
		}
	}
}

func compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskTitle, gotTaskDesc string) error {
	if gotTaskTitle != actualTaskTitle {
		return fmt.Errorf("Expect title: %q, Got: %q", actualTaskTitle, gotTaskTitle)
	}
	if gotTaskDesc != actualTaskDesc {
		return fmt.Errorf("Expect description: %q, Got: %q", actualTaskDesc, gotTaskDesc)
	}
	return nil
}

func getNewData() (*parser.Data, error) {
	tempFileContent := []string{"# taskgo", "## TODO", "## DOING", "## DONE"}
	newData := &parser.Data{}
	if err := newData.ParseData(tempFileContent); err != nil {
		return nil, fmt.Errorf("Unexpected Error: %v", err)
	}
	return newData, nil
}

func getNewCommand(testData *parser.Data) (newCommand *CommandManager, err error) {
	if err != nil {
		return nil, err
	}
	newCommand = CreateNewCommand(testData)
	return newCommand, nil
}
