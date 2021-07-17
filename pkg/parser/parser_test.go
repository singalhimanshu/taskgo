package parser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddSubtask(t *testing.T) {
	t.Run("Test Add subtask is successful", func(t *testing.T) {
		testData, err := getTestDataWithSubtasks(nil)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		wantSubtasks := []string{"Subtask 1"}
		if err := testData.AddSubtask(listIdx, taskIdx, wantSubtasks[0]); err != nil {
			t.Fatal(err)
		}
		gotSubtasks, err := testData.GetSubtasks(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(wantSubtasks, gotSubtasks) {
			t.Errorf("want: %v, Got: %v", wantSubtasks, gotSubtasks)
		}
	})

	t.Run("Test Add subtask out of bounds", func(t *testing.T) {
		subtask := "Subtask"
		testData, err := getTestDataWithSubtasks(nil)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 10
		taskIdx := 10
		taskCount := 1
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, taskCount)
		if gotErr := testData.AddSubtask(listIdx, taskIdx, subtask); wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestGetSubtasks(t *testing.T) {
	wantSubtasks := []string{"Subtask 1", "Subtask 2"}
	testData, err := getTestDataWithSubtasks(wantSubtasks)
	if err != nil {
		t.Fatal(err)
	}
	listIdx := 0
	taskIdx := 0
	gotSubtasks, err := testData.GetSubtasks(listIdx, taskIdx)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(wantSubtasks, gotSubtasks) {
		t.Errorf("want: %v, Got: %v", wantSubtasks, gotSubtasks)
	}
}

func TestGetSubtaskCount(t *testing.T) {
	t.Run("Get Subtask Count successful", func(t *testing.T) {
		subtasks := []string{"Subtask 1", "Subtask 2"}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		wantSubtaskCount := len(subtasks)
		gotSubtaskCount, err := testData.GetSubtaskCount(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotSubtaskCount != wantSubtaskCount {
			t.Errorf("Want: %v, Got: %v", wantSubtaskCount, gotSubtaskCount)
		}
	})

	t.Run("Get Subtask Count list index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 100
		taskIdx := 0
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, testData.GetListCount())
		_, gotErr := testData.GetSubtaskCount(listIdx, taskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Get Subtask Count task index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 100
		taskCount, err := testData.GetTaskCount(listIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", taskIdx, taskCount)
		_, gotErr := testData.GetSubtaskCount(listIdx, taskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestGetSubtask(t *testing.T) {
	t.Run("Get Subtask successful", func(t *testing.T) {
		wantSubtasks := []string{"Subtask 1"}
		testData, err := getTestDataWithSubtasks(wantSubtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 0
		gotSubtask, err := testData.GetSubtask(listIdx, taskIdx, subtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if wantSubtasks[0] != gotSubtask {
			t.Errorf("want: %v, Got: %v", wantSubtasks[0], gotSubtask)
		}
	})

	t.Run("Get Subtask out of bounds", func(t *testing.T) {
		testData, err := getTestDataWithSubtasks(nil)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 0
		subtaskCount := 0
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", subtaskIdx, subtaskCount)
		_, gotErr := testData.GetSubtask(listIdx, taskIdx, subtaskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestEditSubtask(t *testing.T) {
	t.Run("Edit Subtask successful", func(t *testing.T) {
		oldSubtasks := []string{"Subtask"}
		testData, err := getTestDataWithSubtasks(oldSubtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 0
		wantSubtask := "New Subtask"
		if err := testData.EditSubtask(listIdx, taskIdx, subtaskIdx, wantSubtask); err != nil {
			t.Fatal(err)
		}
		gotSubtask, err := testData.GetSubtask(listIdx, taskIdx, subtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotSubtask != wantSubtask {
			t.Errorf("Got: %v, Want: %v", gotSubtask, wantSubtask)
		}
	})

	t.Run("Edit Subtask: list index out of bounds", func(t *testing.T) {
		oldSubtasks := []string{}
		testData, err := getTestDataWithSubtasks(oldSubtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 100
		taskIdx := 0
		subtaskIdx := 0
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, testData.GetListCount())
		gotErr := testData.EditSubtask(listIdx, taskIdx, subtaskIdx, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Edit Subtask: task index out of bounds", func(t *testing.T) {
		oldSubtasks := []string{}
		testData, err := getTestDataWithSubtasks(oldSubtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 100
		subtaskIdx := 0
		taskCount, err := testData.GetTaskCount(listIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", taskIdx, taskCount)
		gotErr := testData.EditSubtask(listIdx, taskIdx, subtaskIdx, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Edit Subtask: subtask index out of bounds", func(t *testing.T) {
		oldSubtasks := []string{}
		testData, err := getTestDataWithSubtasks(oldSubtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 100
		subtaskCount, err := testData.GetSubtaskCount(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", subtaskIdx, subtaskCount)
		gotErr := testData.EditSubtask(listIdx, taskIdx, subtaskIdx, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})
}

func TestRemoveSubtask(t *testing.T) {
	t.Run("Remove Subtask successful", func(t *testing.T) {
		subtasks := []string{"First Subtask"}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 0
		if err := testData.RemoveSubtask(listIdx, taskIdx, subtaskIdx); err != nil {
			t.Fatal(err)
		}
		_, gotErr := testData.GetSubtask(listIdx, taskIdx, subtaskIdx)
		if err != nil {
			t.Fatalf("Expected no error, Got: %v", gotErr)
		}
	})

	t.Run("Remove Subtask: list index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 100
		taskIdx := 0
		subtaskIdx := 0
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, testData.GetListCount())
		gotErr := testData.RemoveSubtask(listIdx, taskIdx, subtaskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Remove Subtask: task index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 100
		subtaskIdx := 0
		taskCount, err := testData.GetTaskCount(listIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", taskIdx, taskCount)
		gotErr := testData.RemoveSubtask(listIdx, taskIdx, subtaskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Remove Subtask: subtask index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		subtaskIdx := 100
		subtaskCount, err := testData.GetSubtaskCount(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", subtaskIdx, subtaskCount)
		gotErr := testData.RemoveSubtask(listIdx, taskIdx, subtaskIdx)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})
}

func TestMoveSubtask(t *testing.T) {
	t.Run("Move Subtask successful", func(t *testing.T) {
		subtasks := []string{"First Subtask", "Second Subtask"}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		firstSubtaskIdx := 0
		secondSubtaskIdx := 1
		if err := testData.MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx); err != nil {
			t.Fatal(err)
		}
		wantFirstSubtask := subtasks[1]
		wantSecondSubtask := subtasks[0]
		gotFirstSubtask, err := testData.GetSubtask(listIdx, taskIdx, firstSubtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotFirstSubtask != wantFirstSubtask {
			t.Fatalf("Got: %v, Want: %v", gotFirstSubtask, wantFirstSubtask)
		}
		gotSecondSubtask, err := testData.GetSubtask(listIdx, taskIdx, secondSubtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotSecondSubtask != wantSecondSubtask {
			t.Fatalf("Got: %v, Want: %v", gotSecondSubtask, wantSecondSubtask)
		}
	})

	t.Run("Move Subtask: list index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 100
		taskIdx := 0
		firstSubtaskIdx := 0
		secondSubtaskIdx := 1
		listCount := testData.GetListCount()
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, listCount)
		gotErr := testData.MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx)
		if wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Move Subtask: task index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 100
		firstSubtaskIdx := 0
		secondSubtaskIdx := 1
		taskCount, err := testData.GetTaskCount(listIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", taskIdx, taskCount)
		gotErr := testData.MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx)
		if wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Move Subtask: first subtask index out of bounds", func(t *testing.T) {
		subtasks := []string{}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		firstSubtaskIdx := 100
		secondSubtaskIdx := 0
		subtaskCount, err := testData.GetSubtaskCount(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", firstSubtaskIdx, subtaskCount)
		gotErr := testData.MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Move Subtask: second subtask index out of bounds", func(t *testing.T) {
		subtasks := []string{"First Subtask"}
		testData, err := getTestDataWithSubtasks(subtasks)
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskIdx := 0
		firstSubtaskIdx := 0
		secondSubtaskIdx := 100
		subtaskCount, err := testData.GetSubtaskCount(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", secondSubtaskIdx, subtaskCount)
		gotErr := testData.MoveSubtask(listIdx, taskIdx, firstSubtaskIdx, secondSubtaskIdx)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func getTestDataWithSubtasks(subtasks []string) (*Data, error) {
	tempFileContent := []string{"# taskgo", "## List", "- Task"}
	for _, task := range subtasks {
		tempFileContent = append(tempFileContent, "+ "+task)
	}
	newData := &Data{}
	if err := newData.ParseData(tempFileContent); err != nil {
		return nil, fmt.Errorf("Unwant Error: %v", err)
	}
	return newData, nil
}
