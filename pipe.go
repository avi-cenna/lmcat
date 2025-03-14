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

	// Extract file paths from input data
	filePaths := extractFilePaths(data)

	if len(filePaths) == 0 {
		return fmt.Errorf("no valid file paths found in input")
	}

	var done chan struct{}
	if hiArgs.stats {
		resultQueue := make(chan *StatsFileResult, 100)
		done = processStatsResults(resultQueue)
		processFileListForStats(filePaths, resultQueue, countTokens)
		<-done
	} else {
		resultQueue := make(chan *ConcatFileResult, 100)
		done = ProcessConcatResults(resultQueue)
		processFileListConcat(filePaths, resultQueue, countTokens)
		<-done
	}

	return nil
}

// extractFilePaths returns a list of file paths from the input data
func extractFilePaths(data []byte) []string {
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

	return filePaths
}

//// runCountString counts tokens and lines in a string
//func runCountString(str string, countTokens TokenFunc) error {
//	lineCount := countLines(str)
//	tokenCount := countTokens([]byte(str))
//	fmt.Printf("Lines: %d\nTokens: %d\n", lineCount, tokenCount)
//	return nil
//}
