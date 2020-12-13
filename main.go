package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/singalhimanshu/taskgo/files"
	"github.com/singalhimanshu/taskgo/ui"
)

func main() {
	fileName := "taskgo.md"
	args := os.Args
	if len(args) == 2 {
		fileName = args[1]
		if !strings.HasSuffix(fileName, ".md") {
			log.Fatal("Invalid file extension (make sure it is a .md(markdown) file)")
		}
	}
	checkFile := files.CheckFile("/" + fileName)
	if !checkFile {
		fmt.Printf("%q doesn't exist. Do you want to create it? (Y[es]/n[o]) ", fileName)

		var createFile string
		fmt.Scanln(&createFile)

		if createFile == "y" || createFile == "Y" || createFile == "Yes" {
			files.CreateFile("/" + fileName)
			files.WriteInitialContent("/" + fileName)
		} else {
			return
		}
	}

	err := ui.Start("/" + fileName)

	if err != nil {
		panic(err)
	}
}
