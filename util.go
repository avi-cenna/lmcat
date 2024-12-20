package main

import (
	"fmt"
	"github.com/boyter/gocodewalker"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReadFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil, err
	}
	return content, nil
}

// WalkFiles walks the files in a directory recursively
func WalkFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error walking path:", err)
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking files: %w", err)
	}
	return files, nil
}

// WalkFilesChan walks the files in a directory recursively, and send results to channel
func WalkFilesChan(dir string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("Error walking path:", err)
				return err
			}
			if !info.IsDir() {
				out <- path
			}
			return nil
		})
		if err != nil {
			log.Println("Error walking files:", err)
		}
	}()
	return out
}

// WalkFilesChanGitignore WalkFilesChan walks the files in a directory recursively, and send results to channel.
// It will ignore files in the .gitignore file. It will also ignore the .git directory.
func WalkFilesChanGitignore(dir string) <-chan *gocodewalker.File {
	//fileListQueue := make(chan *gocodewalker.File, 100)
	fileListQueue := make(chan *gocodewalker.File)
	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// restrict to only process files that have the .go extension
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		fmt.Println("ERR", e.Error())
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	return fileListQueue
}

// FindFiles uses ripgrep to find files that match a glob pattern
// and/or contain a regex pattern. It filters out files that are
// likely not text files, both via ripgrep and a list of known
// non-text file extensions.
func FindFiles(glob, regex string) ([]string, error) {
	log.Println("glob pattern:", glob)
	args := []string{"--files-with-matches", regex}
	if glob != "" {
		// TODO (AA): there seems to be some issue here with the globbing.
		//            amd i doing it wrong?
		args = append(args, "--glob", glob)
	}

	log.Println("rg args:", args)
	cmd := exec.Command("rg", args...)
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("error running ripgrep: %w", err)
	}
	rgFiles := strings.Split(strings.TrimSpace(string(output)), "\n")

	exclude := mapset.NewSet("LICENSE", "go.mod", "go.sum")

	var textFiles []string
	for _, f := range rgFiles {
		if IsLikelyTextFile(f) && !exclude.Contains(f) {
			textFiles = append(textFiles, f)
		}
	}

	log.Println("text files:", textFiles)
	return textFiles, nil
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func Test1() {
	//fileListQueue := make(chan *gocodewalker.File, 100)
	fileListQueue := make(chan *gocodewalker.File)
	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// restrict to only process files that have the .go extension
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		fmt.Println("ERR", e.Error())
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for f := range fileListQueue {
		fmt.Println(f.Location, ":::", f.Filename)
	}
}
