package command

import (
	"fmt"
	"testing"

	"github.com/singalhimanshu/taskgo/parser"
)

func TestExecute(t *testing.T) {
	testData := &parser.Data{}
	tempFileContent := []string{"# taskgo", "## TODO", "## DOING", "## DONE"}
	if err := testData.ParseData(tempFileContent); err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	testCommand := CreateNewCommand(testData)
	t.Run("Add task command", testAddTaskCommand(t, testCommand, testData))
	if err := testData.ParseData(tempFileContent); err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	testData = &parser.Data{}
	if err := testData.ParseData(tempFileContent); err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	testCommand = CreateNewCommand(testData)
	t.Run("Remove task Command", testRemoveTaskCommand(t, testCommand, testData))
	testData = &parser.Data{}
	if err := testData.ParseData(tempFileContent); err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}
	testCommand = CreateNewCommand(testData)
	t.Run("Swap List Item Command", testSwapListItemCommand(t, testCommand, testData))
}

func testAddTaskCommand(t *testing.T, testCommand *CommandManager, testData *parser.Data) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
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
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData[0], gotTaskData[1]); err != nil {
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
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData[0], gotTaskData[1]); err != nil {
			t.Error(err)
		}
	}
}

func testRemoveTaskCommand(t *testing.T, testCommand *CommandManager, testData *parser.Data) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
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
		_, err := testData.GetTask(listIdx, taskIdx)
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
		if err := compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskData[0], gotTaskData[1]); err != nil {
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

func testSwapListItemCommand(t *testing.T, testCommand *CommandManager, testData *parser.Data) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
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
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask[0], swappedTask[1]); err != nil {
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
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask[0], swappedTask[1]); err != nil {
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
		if err := compareTaskData(firstTaskTitle, firstTaskDesc, swappedTask[0], swappedTask[1]); err != nil {
			t.Error(err)
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
