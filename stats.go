package main

import (
	"bufio"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	_ "github.com/h2non/filetype/matchers"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	tke *tiktoken.Tiktoken
)

func runStats(cliCtx *cli.Context) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	// Initialize a map to store counts
	counts := make(map[string]int)

	// Create output file
	outputFile, err := os.Create("F:/output.txt")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// Walk through the directory
	files, err := getRipgrepFiles(cliCtx.String("regex"), cliCtx.String("glob"))
	foo := filetype.GetType("F:/output.txt")
	bar := matchers.Application
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		fmt.Println("Visiting file:", f)
		// Get the file extension
		ext := strings.ToLower(filepath.Ext(f))
		if ext != "" {
			counts[ext]++
		}
		if ext == ".htm" || true {
			err := appendFileContent(f, writer)
			if err != nil {
				log.Printf("Error appending file %s: %v", f, err)
			}
		}
	}

	// Print the results
	for ext, count := range counts {
		fmt.Printf("%s: %d\n", ext, count)
	}

	fmt.Println("Text files have been concatenated into output.txt")

	// Count tokens in the output file
	tokenCount, err := countTokens("F:/output.txt")
	if err != nil {
		return fmt.Errorf("error counting tokens: %w", err)
	}

	fmt.Printf("Total tokens in output file: %d\n", tokenCount)

	if tokenCount > 1000 {
		log.Printf("WARNING: Token count (%d) exceeds 1000\n", tokenCount)
	}

	return nil
}

func countTokens2(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// print log warning
		log.Println("Error reading file:", err)
		return 0, err
	}

	tokens := tke.Encode(string(content), nil, nil)
	return len(tokens), nil
}

func init() {
	encoding := "cl100k_base"
	var err error
	tke, err = tiktoken.GetEncoding(encoding)
	if err != nil {
		log.Fatal(err)
	}
}
