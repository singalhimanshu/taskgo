package files

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestCreateFile(t *testing.T) {
	CreateFile()
	got := CheckFile()
	want := true

	if got != want {
		t.Errorf("file wasn't created successfully")
	}
}

func TestCheckFile(t *testing.T) {
	got := CheckFile()
	want := true

	if got != want {
		t.Errorf("file expected but not present")
	}
}

func TestWriteInitialContent(t *testing.T) {
	WriteInitialContent()
	want := fmt.Sprintf(initialFileContent, "files", "TODO", "DOING", "DONE")

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	got := string(content)

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetBoardName(t *testing.T) {
	got := GetDirectoryName()
	want := "files"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCheckFileSyntax(t *testing.T) {
	got := CheckFileSyntax()
	want := true

	if got != want {
		t.Errorf("file syntax correct but got false")
	}
}

func TestCheckPrefix(t *testing.T) {
	testCases := []struct {
		TestName string
		Prefix   string
		Want     bool
	}{
		{
			"Testing # Board Name should be parsed correctly",
			"# Board Name",
			true,
		},
		{
			"Testing # TODO Name should be parsed correctly",
			"## TODO",
			true,
		},
		{
			"Testing - Task should be parsed correctly",
			"- Task",
			true,
		},
		{
			"Testing > Task Description should be parsed correctly",
			"> Task Description",
			true,
		},
		{
			"Testing * Subtask should be parsed correctly",
			"* Subtask",
			true,
		},
		{
			"#BoardNameWithoutSpace should not parse",
			"#BoardNameWithoutSpace",
			false,
		},
		{
			"##ListNameWithoutSpace should not parse",
			"##ListNameWithoutSpace",
			false,
		},
		{
			"-TaskNameWithoutSpace should not parse",
			"-TaskNameWithoutSpace",
			false,
		},
		{
			">TaskDescWithoutSpace should not parse",
			">TaskDescWithoutSpace",
			false,
		},
		{
			"*SubtaskWithoutSpace should not parse",
			"*SubtaskWithoutSpace",
			false,
		},
		{
			"SimpleText should not parse",
			"SimpleText",
			false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			got := CheckPrefix(testCase.Prefix)

			if got != testCase.Want {
				t.Errorf("got %v want %v", got, testCase.Want)
			}
		})
	}
}
