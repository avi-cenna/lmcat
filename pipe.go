package main

import (
	"fmt"
	"strings"
)

// RunPipe handles the --count flag functionality
func RunPipe(data []byte, hiArgs *HiArgs) error {
	countTokens := GetTokenFunc(hiArgs.approx)

	if hiArgs.count {
		lineCount := countLines(string(data))
		tokenCount := countTokens(data)
		fmt.Printf("Lines: %d\nTokens: %d\n", lineCount, tokenCount)
		return nil
	}

	// Split the data into lines
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	// Filter out empty lines
	var filePaths []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			filePaths = append(filePaths, trimmedLine)
		}
	}

	if len(filePaths) == 0 {
		return fmt.Errorf("no valid file paths found in input")
	}

	// Create a channel for results
	resultQueue := make(chan *ConcatFileResult, 100)

	// Use the shared ProcessResults function
	done := ProcessResults(resultQueue)

	// Process each line as a filepath using the function from concat.go
	processFileListConcat(filePaths, resultQueue, countTokens)
	<-done

	return nil
}

//// runCountString counts tokens and lines in a string
//func runCountString(str string, countTokens TokenFunc) error {
//	lineCount := countLines(str)
//	tokenCount := countTokens([]byte(str))
//	fmt.Printf("Lines: %d\nTokens: %d\n", lineCount, tokenCount)
//	return nil
//}
