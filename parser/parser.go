package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const fileName = "/taskgo.md"

func GetBoardName() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := dir + fileName

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var boardName string

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
		if strings.HasPrefix(line, "# ") {
			line = strings.TrimSpace(line)
			boardNamestartingIndex := strings.Index(line, " ") + 1
			boardName = line[boardNamestartingIndex:]
		}
	}

	return boardName
}

func GetListNames() []string {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filePath := dir + fileName

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var listNames []string

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
		if strings.HasPrefix(line, "## ") {
			line = strings.TrimSpace(line)
			liststartingIndex := strings.Index(line, " ") + 1
			listNames = append(listNames, line[liststartingIndex:])
		}
	}

	return listNames
}
