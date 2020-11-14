package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const fileName = "/taskgo.md"

const initialFileContent = `# %s

## %s


## %s


## %s


`

var validPrefixes = [...]string{
	// Board Name
	"# ",
	// List Name
	"## ",
	// Task
	"- ",
	// Task Description
	"> ",
	// Subtask
	"* ",
}

func CheckFile() bool {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := dir + fileName
	log.Println(filePath)
	_, err = os.Stat(filePath)
	return err == nil
}

func CreateFile() {
	f, err := OpenFileWriteOnly()
	defer f.Close()

	if err != nil {
		log.Fatalf("Cannot create file %q, ERR: %v", fileName, err)
	}
}

func WriteInitialContent() {
	f, err := OpenFileWriteOnly()
	defer f.Close()

	if err != nil {
		log.Fatalf("Cannot Open file %q, ERR: %v", fileName, err)
	}

	// TODO: Make these customizable
	_, err = f.WriteString(fmt.Sprintf(initialFileContent, GetDirectoryName(), "TODO", "DOING", "DONE"))

	if err != nil {
		log.Fatalf("Cannot write contents to file (%v): %v", fileName, err)
	}
}

func OpenFileWriteOnly() (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return os.OpenFile(dir+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func GetDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get directory name: %v", err)
	}

	dirs := strings.Split(dir, "/")
	dirName := dirs[len(dirs)-1]

	return dirName
}

func CheckPrefix(line string) bool {
	result := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			result = true
			break
		}
	}
	return result
}

func CheckFileSyntax() bool {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(dir + fileName)
	defer file.Close()

	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileContent []string

	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	for _, line := range fileContent {
		// ignore empty lines
		if len(line) < 1 {
			continue
		}
		if !CheckPrefix(line) {
			return false
		}
	}
	return true
}
