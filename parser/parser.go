package parser

import (
	"strings"

	"github.com/singalhimanshu/taskgo/files"
)

const fileName = "/taskgo.md"

func GetBoardName() string {

	var boardName string

	fileContent := files.OpenFile(fileName)

	for _, line := range fileContent {
		// ignore empty lines
		if len(line) < 1 {
			continue
		}
		if strings.HasPrefix(line, "# ") {
			line = strings.TrimSpace(line)
			boardNamestartingIndex := strings.Index(line, " ") + 1
			boardName = line[boardNamestartingIndex:]
		}
	}

	return boardName
}

func GetListNames() []string {

	var listNames []string

	fileContent := files.OpenFile(fileName)

	for _, line := range fileContent {
		// ignore empty lines
		if len(line) < 1 {
			continue
		}
		if strings.HasPrefix(line, "## ") {
			line = strings.TrimSpace(line)
			liststartingIndex := strings.Index(line, " ") + 1
			listNames = append(listNames, line[liststartingIndex:])
		}
	}

	return listNames
}

func GetTaskFromListName(listName string) []string {
	searchText := "## " + listName
	var taskNames []string

	fileContent := files.OpenFile(fileName)

	taskStartIndex := -1

	for i, line := range fileContent {
		if strings.HasPrefix(line, searchText) {
			taskStartIndex = i
			break
		}
	}

	if taskStartIndex == -1 {
		return []string{}
	}

	for i := taskStartIndex + 1; i < len(fileContent); i++ {
		if len(fileContent[i]) < 1 {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(fileContent[i]), "- ") {
			tempLine := strings.TrimSpace(fileContent[i])
			taskNameIdx := strings.Index(tempLine, " ") + 1
			taskNames = append(taskNames, tempLine[taskNameIdx:])
		} else if strings.HasPrefix(strings.TrimSpace(fileContent[i]), "## ") {
			break
		}
	}

	return taskNames
}
