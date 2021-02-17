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

func compareTaskData(actualTaskTitle, actualTaskDesc, gotTaskTitle, gotTaskDesc string) error {
	if gotTaskTitle != actualTaskTitle {
		return fmt.Errorf("Expect title: %q, Got: %q", actualTaskTitle, gotTaskTitle)
	}
	if gotTaskDesc != actualTaskDesc {
		return fmt.Errorf("Expect description: %q, Got: %q", actualTaskDesc, gotTaskDesc)
	}
	return nil
}
