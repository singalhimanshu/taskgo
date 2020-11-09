package files

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const FILE_NAME = "taskgo.md"

const INITIAL_FILE_CONTENT = `# %s

## %s


## %s


## %s


`

func CheckFile() bool {
	filePath := "./" + FILE_NAME
	_, err := os.Stat(filePath)
	return err == nil
}

func CreateFile() {
	f, err := OpenFileWriteOnly()
	defer f.Close()

	if err != nil {
		log.Fatalf("Cannot create file %q, ERR: %v", FILE_NAME, err)
	}
}

func WriteInitialContent() {
	f, err := OpenFileWriteOnly()
	defer f.Close()

	if err != nil {
		log.Fatalf("Cannot Open file %q, ERR: %v", FILE_NAME, err)
	}

	// TODO: Make these customizable
	_, err = f.WriteString(fmt.Sprintf(INITIAL_FILE_CONTENT, GetBoardName(), "TODO", "DOING", "DONE"))

	if err != nil {
		log.Fatalf("Cannot write contents to file (%v): %v", FILE_NAME, err)
	}
}

func OpenFileWriteOnly() (*os.File, error) {
	return os.OpenFile(FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func GetBoardName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get directory name: %v", err)
	}

	dirs := strings.Split(dir, "/")
	dirName := dirs[len(dirs)-1]

	return dirName
}
