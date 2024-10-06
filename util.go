package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

// FindFiles uses ripgrep to find files that match a glob pattern
// and/or contain a regex pattern. It filters out files that are
// likely not text files, both via ripgrep and a list of known
// non-text file extensions.
func FindFiles(glob, regex string) ([]string, error) {
	log.Println("glob pattern:", glob)
	args := []string{"--files-with-matches", regex}
	if glob != "" {
		args = append(args, "--glob", glob)
	}

	log.Println("rg args:", args)
	cmd := exec.Command("rg", args...)
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("error running ripgrep: %w", err)
	}
	rgFiles := strings.Split(strings.TrimSpace(string(output)), "\n")

	var textFiles []string
	for _, f := range rgFiles {
		if IsLikelyTextFile(f) {
			textFiles = append(textFiles, f)
		}
	}

	return textFiles, nil
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}
