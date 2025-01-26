package main

import (
	"fmt"
	"github.com/boyter/gocodewalker"
)

func main2() {
	fmt.Println("using Gocodewalker!")
	fileListQueue := WalkFiles(100)
	for f := range fileListQueue {
		fmt.Println(f.Location, "::", f.Filename)
	}
}

func WalkFiles(bufSize int) chan *gocodewalker.File {
	fileListQueue := make(chan *gocodewalker.File, bufSize)

	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// restrict to only process files that have the .go extension
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")
	fileWalker.ExcludeListExtensions = append(fileWalker.ExcludeListExtensions, "mod", "sum")
	//fileWalker.ExcludeFilenameRegex = append(fileWalker.ExcludeFilenameRegex, regexp.MustCompile("LICENSE"))
	fileWalker.ExcludeFilename = append(fileWalker.ExcludeFilename, "LICENSE")

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		fmt.Println("ERR", e.Error())
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			fmt.Println("ERR", err.Error())
		}
	}()

	return fileListQueue
}
