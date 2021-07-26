package parser

import (
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	t.Run("Get List successful", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## List"})
		if err != nil {
			t.Fatal(err)
		}
		_, err = testData.GetList(0)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Get List out of bounds", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## List"})
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 100
		_, gotErr := testData.GetList(listIdx)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", listIdx, testData.GetListCount())
		if wantErr.Error() != gotErr.Error() {
			t.Errorf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestGetTask(t *testing.T) {
	t.Run("Get Task: in bounds", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## List"})
		if err != nil {
			t.Fatal(err)
		}
		listIdx := 0
		taskTitle := "Test Task"
		taskDesc := "Test Task Desc"
		taskIdx := 0
		err = testData.AddNewTask(listIdx, taskTitle, taskDesc, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		gotTask, err := testData.GetTask(listIdx, taskIdx)
		if err != nil {
			t.Fatal(err)
		}
		if taskTitle != gotTask.ItemName {
			t.Fatalf("Want: %v, Got: %v", taskTitle, gotTask.ItemName)
		}
		if taskDesc != gotTask.ItemDescription {
			t.Fatalf("Want: %v, Got %v", taskDesc, gotTask.ItemDescription)
		}
	})

	t.Run("Get Task: list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## List"})
		if err != nil {
			t.Fatal(err)
		}
		_, gotErr := testData.GetTask(100, 0)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 1)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})

	t.Run("Get Task: task idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## List"})
		if err != nil {
			t.Fatal(err)
		}
		_, gotErr := testData.GetTask(0, 100)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 0)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}

func TestAddNewTask(t *testing.T) {
	t.Run("Add new task", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## TODO"})
		if err != nil {
			t.Fatal(err)
		}
		if err = testData.AddNewTask(0, "Test Task", "Test Task Desc", 0); err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## TODO"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"})
		if err != nil {
			t.Fatal(err)
		}
		err = testData.EditTask(0, 0, "Edit Task", "")
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
		err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## TODO", "- Task"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"})
		if err != nil {
			t.Fatal(err)
		}
		err = testData.MoveTask(0, 0, 1)
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
		if movedTask.ItemName != "Task" {
			t.Fatalf("Want: %v, Got: %v", "Task", movedTask.ItemName)
		}
	})

	t.Run("Move Task: source list idx out of bounds", func(t *testing.T) {
		testData := &Data{}
		err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## Source List", "- Task", "## Dest List"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task"})
		if err != nil {
			t.Fatal(err)
		}
		_, err = testData.RemoveTask(0, 0)
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"})
		if err != nil {
			t.Fatal(err)
		}
		err = testData.SwapListItems(0, 0, 1)
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"})
		if err != nil {
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
		err := testData.ParseData([]string{"# taskgo", "## List", "- Task 1", "- Task 2"})
		if err != nil {
			t.Fatal(err)
		}
		gotErr := testData.SwapListItems(0, 0, 100)
		wantErr := fmt.Errorf("Index Out of Bounds: got %v, length: %v", 100, 2)
		if wantErr.Error() != gotErr.Error() {
			t.Fatalf("Want: %v, Got: %v", wantErr, gotErr)
		}
	})
}
