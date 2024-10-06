package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GetRipgrepFiles(regex, glob string) ([]string, error) {
	log.Println("glob pattern:", glob)
	args := []string{"--files-with-matches", regex}
	if glob != "" {
		args = append(args, "--glob", glob)
	}
	//if regex != "" {
	//	args = append(args, regex)
	//}

	log.Println("rg args:", args)
	cmd := exec.Command("rg", args...)
	output, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("error running ripgrep: %w", err)
	}
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}
