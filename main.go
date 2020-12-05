package main

import (
	"fmt"

	"github.com/singalhimanshu/taskgo/files"
	"github.com/singalhimanshu/taskgo/ui"
)

func main() {
	checkFile := files.CheckFile()
	if !checkFile {
		fmt.Print("taskgo.md doesn't exist. Do you want to create it? (Y[es]/n[o]) ")

		var createFile string
		fmt.Scanln(&createFile)

		if createFile == "y" || createFile == "Y" || createFile == "Yes" {
			files.CreateFile()
			files.WriteInitialContent()
		} else {
			return
		}
	}

	err := ui.Start()

	if err != nil {
		panic(err)
	}
}
