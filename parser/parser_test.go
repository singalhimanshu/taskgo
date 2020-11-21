package parser

import (
	"reflect"
	"testing"
)

// For this test you would have to manually create a file taskgo.md
// Inside the parser folder with content: # Board Name
func TestGetBoardName(t *testing.T) {
	d, err := getNewData()
	if err != nil {
		t.Fatal(err)
	}

	got := d.GetBoardName()
	want := "Board Name"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

// For this test you would have to manually create a file taskgo.md
// Inside the parser folder with content:
// ## TODO
// ## DOING
// ## DONE
func TestGetListNames(t *testing.T) {
	d, err := getNewData()
	if err != nil {
		t.Fatal(err)
	}

	got := d.GetListNames()
	want := []string{
		"TODO",
		"DOING",
		"DONE",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetTasks(t *testing.T) {
	d, err := getNewData()
	if err != nil {
		t.Fatal(err)
	}

	got := d.GetTasks(0)
	want := []string{"Task 1"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestAddNewTask(t *testing.T) {
	d, err := getNewData()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("successfully add a new task", func(t *testing.T) {
		err = d.AddNewTask(0, "Task 2")
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}

		got := d.GetTasks(0)
		want := []string{
			"Task 1",
			"Task 2",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("throw error when index is out of bounds", func(t *testing.T) {

	})
}

func getNewData() (Data, error) {
	d := Data{}
	err := d.ParseData()
	return d, err
}
