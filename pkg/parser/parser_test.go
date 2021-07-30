package parser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddNewTask(t *testing.T) {
	t.Run("Add new task", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## TODO"}); err != nil {
			t.Fatal(err)
		}
		if err := testData.AddNewTask(0, "Test Task", "Test Task Desc", 0); err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		if taskCount != 1 {
			t.Fatalf("Want: %v, Got: %v", 1, taskCount)
		}
	})

	t.Run("Add new task: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## TODO"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.AddNewTask(100, "Test Task", "Test Task Desc", 0)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestEditTask(t *testing.T) {
	t.Run("Edit task", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"}); err != nil {
			t.Fatal(err)
		}
		err := testData.EditTask(0, 0, "Edit Task", "")
		if err != nil {
			t.Fatal(err)
		}
		editedTask, err := testData.GetTask(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if editedTask.ItemName != "Edit Task" {
			t.Fatalf("Want: %v, Got: %v", "Edit Task", editedTask)
		}
	})

	t.Run("Edit task: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.EditTask(100, 0, "Edit Task", "")
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Edit task: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.EditTask(0, 100, "Edit Task", "")
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestMoveTask(t *testing.T) {
	t.Run("Move Task", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"}); err != nil {
			t.Fatal(err)
		}
		err := testData.MoveTask(0, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		sourceTaskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		if sourceTaskCount != 0 {
			t.Fatalf("Want: %v, Got: %v", 0, sourceTaskCount)
		}
		destTaskCount, err := testData.GetTaskCount(1)
		if err != nil {
			t.Fatal(err)
		}
		if destTaskCount != 1 {
			t.Fatalf("Want: %v, Got: %v", 1, destTaskCount)
		}
		movedTask, err := testData.GetTask(1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if movedTask.ItemName != "Task" {
			t.Fatalf("Want: %v, Got: %v", "Task", movedTask.ItemName)
		}
	})

	t.Run("Move Task: source list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.MoveTask(0, 100, 0)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 2)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Move Task: dest list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.MoveTask(0, 0, 100)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 2)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestRemoveTask(t *testing.T) {
	t.Run("Remove Task", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		_, err := testData.RemoveTask(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		if taskCount != 0 {
			t.Fatalf("Want: %v, Got: %v", 0, taskCount)
		}
	})

	t.Run("Remove Task: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		_, gotErr := testData.RemoveTask(100, 0)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Remove Task: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		_, gotErr := testData.RemoveTask(0, 100)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestSwapListItems(t *testing.T) {
	t.Run("Swap list items", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"}); err != nil {
			t.Fatal(err)
		}
		err := testData.SwapListItems(0, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		task1, err := testData.GetTask(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if task1.ItemName != "Task 2" {
			t.Fatalf("Want: %v, Got: %v", "Task 2", task1.ItemName)
		}
		task2, err := testData.GetTask(0, 1)
		if err != nil {
			t.Fatal(err)
		}
		if task2.ItemName != "Task 1" {
			t.Fatalf("Want: %v, Got: %v", "Task 1", task2.ItemName)
		}
	})

	t.Run("Swap list items: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.SwapListItems(100, 0, 1)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Swap list items: first task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.SwapListItems(0, 100, 0)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 2)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Swap list items: second task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"}); err != nil {
			t.Fatal(err)
		}
		gotErr := testData.SwapListItems(0, 0, 100)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 2)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestAddSubtask(t *testing.T) {
	t.Run("Add subtask", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		wantSubtasks := []string{"Subtask"}
		if err := testData.AddSubtask(0, 0, wantSubtasks[0]); err != nil {
			t.Fatal(err)
		}
		gotSubtasks, err := testData.GetSubtasks(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(wantSubtasks, gotSubtasks) {
			t.Errorf("want: %v, Got: %v", wantSubtasks, gotSubtasks)
		}
	})

	t.Run("Add subtask: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if gotErr := testData.AddSubtask(100, 0, "Subtask"); wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Add subtask: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, taskCount)
		if gotErr := testData.AddSubtask(0, 100, "Subtask"); wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestEditSubtask(t *testing.T) {
	t.Run("Edit Subtask", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task", "+ Subtask"}); err != nil {
			t.Fatal(err)
		}
		wantSubtask := "New Subtask"
		if err := testData.EditSubtask(0, 0, 0, wantSubtask); err != nil {
			t.Fatal(err)
		}
		gotSubtask, err := testData.GetSubtask(0, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if gotSubtask != wantSubtask {
			t.Errorf("Got: %v, Want: %v", gotSubtask, wantSubtask)
		}
	})

	t.Run("Edit Subtask: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, testData.GetListCount())
		gotErr := testData.EditSubtask(100, 0, 0, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Edit Subtask: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, taskCount)
		gotErr := testData.EditSubtask(0, 100, 0, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Edit Subtask: subtask idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		subtaskCount, err := testData.GetSubtaskCount(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, subtaskCount)
		gotErr := testData.EditSubtask(0, 0, 100, "")
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})
}

func TestRemoveSubtask(t *testing.T) {
	t.Run("Remove Subtask", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task", "+ Subtask"}); err != nil {
			t.Fatal(err)
		}
		if err := testData.RemoveSubtask(0, 0, 0); err != nil {
			t.Fatal(err)
		}
		subtaskCount, err := testData.GetSubtaskCount(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		if subtaskCount != 0 {
			t.Fatalf("Want: %v, Got: %v", 0, subtaskCount)
		}
	})

	t.Run("Remove Subtask: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, testData.GetListCount())
		gotErr := testData.RemoveSubtask(100, 0, 0)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Remove Subtask: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, taskCount)
		gotErr := testData.RemoveSubtask(0, 100, 0)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})

	t.Run("Remove Subtask: subtask idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		subtaskCount, err := testData.GetSubtaskCount(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, subtaskCount)
		gotErr := testData.RemoveSubtask(0, 0, 100)
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("Got: %v, Want: %v", gotErr, wantErr)
		}
	})
}

func TestSwapSubtask(t *testing.T) {
	t.Run("Swap Subtask", func(t *testing.T) {
		subtasks := []string{"First Subtask", "Second Subtask"}
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task",
			"+ " + subtasks[0], "+ " + subtasks[1]}); err != nil {
			t.Fatal(err)
		}
		firstSubtaskIdx := 0
		secondSubtaskIdx := 1
		if err := testData.SwapSubtask(0, 0, firstSubtaskIdx, secondSubtaskIdx); err != nil {
			t.Fatal(err)
		}
		wantFirstSubtask := subtasks[1]
		wantSecondSubtask := subtasks[0]
		gotFirstSubtask, err := testData.GetSubtask(0, 0, firstSubtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotFirstSubtask != wantFirstSubtask {
			t.Fatalf("Got: %v, Want: %v", gotFirstSubtask, wantFirstSubtask)
		}
		gotSecondSubtask, err := testData.GetSubtask(0, 0, secondSubtaskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if gotSecondSubtask != wantSecondSubtask {
			t.Fatalf("Got: %v, Want: %v", gotSecondSubtask, wantSecondSubtask)
		}
	})

	t.Run("Swap Subtask: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		listCount := testData.GetListCount()
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, listCount)
		gotErr := testData.SwapSubtask(100, 0, 0, 0)
		if wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Swap Subtask: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		taskCount, err := testData.GetTaskCount(0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, taskCount)
		gotErr := testData.SwapSubtask(0, 100, 0, 0)
		if wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Swap Subtask: first subtask idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task"}); err != nil {
			t.Fatal(err)
		}
		subtaskCount, err := testData.GetSubtaskCount(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, subtaskCount)
		gotErr := testData.SwapSubtask(0, 0, 100, 0)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Swap Subtask: second subtask idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		if err := testData.ParseData([]string{"# taskgo", "## List", "- Task", "+ Subtask"}); err != nil {
			t.Fatal(err)
		}
		subtaskCount, err := testData.GetSubtaskCount(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, subtaskCount)
		gotErr := testData.SwapSubtask(0, 0, 0, 100)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}
