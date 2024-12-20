package main

import (
	"fmt"
	_ "github.com/h2non/filetype/matchers"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

var (
	tke *tiktoken.Tiktoken
)

func init() {
	encoding := "cl100k_base"
	var err error
	tke, err = tiktoken.GetEncoding(encoding)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO (AA): refactor as necessary to be in more idiomatic Go, and DRY

func RunStats(cliCtx *cli.Command) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	extCounts := make(map[string]int)
	extTokens := make(map[string]int)

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, f := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tokenCount, err := countTokensFilePath(f)
			if err != nil {
				log.Printf("Error counting extTokens in file %s: %v", f, err)
			}

			ext := extensionOrBase(f)
			mu.Lock()
			extCounts[ext]++
			extTokens[ext] += tokenCount
			mu.Unlock()
		}()
	}
	wg.Wait()

	printCountsAndTokens(extCounts, extTokens)
	return nil
}

func RunStatsSequential(cliCtx *cli.Command) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	extCounts := make(map[string]int)
	extTokens := make(map[string]int)

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		tokenCount, err := countTokensFilePath(f)
		if err != nil {
			log.Printf("Error counting extTokens in file %s: %v", f, err)
		}

		ext := extensionOrBase(f)
		extCounts[ext]++
		extTokens[ext] += tokenCount
	}

	printCountsAndTokens(extCounts, extTokens)
	return nil
}

func countTokensFilePath(filePath string) (int, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return CountTokensInText(content)
}

func CountTokensInText(text []byte) (int, error) {
	tokens := tke.Encode(string(text), nil, nil)
	return len(tokens), nil
}

func printCountsAndTokens(extCounts map[string]int, extTokens map[string]int) {
	sortedKeys := make([]string, 0, len(extCounts))
	for ext := range extCounts {
		sortedKeys = append(sortedKeys, ext)
	}
	slices.Sort(sortedKeys)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Extension", "Count", "Tokens"})
	for _, ext := range sortedKeys {
		t.AppendRow(table.Row{ext, extCounts[ext], extTokens[ext]})
	}
	t.Render()
}

// extensionOrBase returns the extension of a file path if it exists, otherwise
// it returns the base name of the file path.
func extensionOrBase(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != "" {
		return ext
	}
	return filepath.Base(filePath)
}
