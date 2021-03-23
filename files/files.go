package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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

// CheckFile checks if the file is present in the current directory or not.
func CheckFile(fileName string) bool {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filePath := dir + fileName
	_, err = os.Stat(filePath)
	return err == nil
}

// CreateFile Creates the file.
func CreateFile(fileName string) {
	f, err := OpenFileWriteOnly(fileName)
	defer f.Close()
	if err != nil {
		log.Fatalf("Cannot create file %q, ERR: %v", fileName, err)
	}
}

// WriteInitialContent Writes initial content to the file.
func WriteInitialContent(fileName string) {
	f, err := OpenFileWriteOnly(fileName)
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

// OpenFileWriteOnly opens file in write only mode.
func OpenFileWriteOnly(fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Cannot open file: %v", err)
	}
	return os.OpenFile(dir+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// GetDirectoryName returns the name of current working directory.
func GetDirectoryName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get directory name: %v", err)
	}
	dirs := strings.Split(dir, "/")
	dirName := dirs[len(dirs)-1]
	return dirName
}

// CheckPrefix checks prefix if it matches to the given set of prefix.
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
