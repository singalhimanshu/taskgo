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
	want := fmt.Sprintf(INITIAL_FILE_CONTENT, "files", "TODO", "DOING", "DONE")

	content, err := ioutil.ReadFile(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	got := string(content)

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetBoardName(t *testing.T) {
	got := GetBoardName()
	want := "files"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
