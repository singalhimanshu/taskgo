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

// Checks if the file(taskgo.md) is present in the current directory or not.
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

// Creates the file(taskgo.md).
func CreateFile() {
	f, err := OpenFileWriteOnly()
	defer f.Close()

	if err != nil {
		log.Fatalf("Cannot create file %q, ERR: %v", fileName, err)
	}
}

// Writes initial content to the file.
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

// Opens file in write only mode.
func OpenFileWriteOnly() (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", err)
	}

	return os.OpenFile(dir+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// Returns the name of current working directory.
func GetDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get directory name: %v", err)
	}

	dirs := strings.Split(dir, "/")
	dirName := dirs[len(dirs)-1]

	return dirName
}

// Helper function to check prefix if it matches to the given set of prefix.
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

// Checks the file syntax for any errors.
func CheckFileSyntax() bool {
	fileContent := OpenFile(fileName)

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

// FilePath returns the complete path of file given the fileName.
func FilePath(fileName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filePath := dir + fileName

	return filePath, nil
}

// OpenFile opens the given file and returns the content of file line by line
// as a slice of string.
func OpenFile(fileName string) []string {
	filePath, err := FilePath(fileName)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileContent []string

	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	return fileContent
}

// WriteFile writes a slice of string to a file line by line.
// It returns an error if file cannot be opened or the content can't be written.
func WriteFile(fileContent []string, fileName string) error {
	filePath, err := FilePath(fileName)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range fileContent {
		fmt.Fprintln(w, line)
	}

	return w.Flush()
}
