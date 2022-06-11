package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

var App = tview.NewApplication()

func main() {
	if len(os.Args) < 2 {
		log.Panic("provide path argument")
	}
	absPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Panic("could not get abs path:", err)
	}
	tree, err := GetFileTreeView(absPath)
	if err != nil {
		log.Panic("could not get file tree view:", err)
	}

	App.SetRoot(tree, true).EnableMouse(true)
	err = App.Run()
	if err != nil {
		log.Panic("could not run:", err)
	}
}
