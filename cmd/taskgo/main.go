package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/singalhimanshu/taskgo/internal/ui"
	"github.com/singalhimanshu/taskgo/pkg/files"
)

const defaultFileName = "taskgo.md"

func main() {
	fileName := flag.String("f", defaultFileName, "markdown file to use as task storage")
	flag.Parse()
	if !strings.HasSuffix(*fileName, ".md") {
		log.Fatal("Invalid file extension (make sure it is a .md(markdown) file)")
	}
	checkFile := files.CheckFile("/" + *fileName)
	if !checkFile {
		fmt.Printf("%q doesn't exist. Do you want to create it? (Y[es]/n[o]) ", *fileName)
		var createFile string
		fmt.Scanln(&createFile)
		if createFile == "y" || createFile == "Y" || createFile == "Yes" {
			files.CreateFile("/" + *fileName)
			files.WriteInitialContent("/" + *fileName)
		} else {
			return
		}
	}
	err := ui.Start("/" + *fileName)
	if err != nil {
		log.Fatal(err)
	}
}
