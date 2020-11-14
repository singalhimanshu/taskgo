package parser

import (
	"reflect"
	"testing"
)

// For this test you would have to manually create a file taskgo.md
// Inside the parser folder with content: # Board Name
func TestGetBoardName(t *testing.T) {
	got := GetBoardName()
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
	got := GetListNames()
	want := []string{
		"TODO",
		"DOING",
		"DONE",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetTaskFromListName(t *testing.T) {
	got := GetTaskFromListName("TODO")
	want := []string{
		"Task 1",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
