package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	_ "sort"
	"strings"
	"sync"

	_ "github.com/h2non/filetype/matchers"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli"
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

func RunStats(cliCtx *cli.Context) error {
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
		tokenCount, err := CountTokensFilePath(f)
		if err != nil {
			log.Printf("Error counting extTokens in file %s: %v", f, err)
		}

		ext := strings.ToLower(filepath.Ext(f))
		if ext != "" {
			extCounts[ext]++
			extTokens[ext] += tokenCount
		} else {
			extCounts[filepath.Base(f)]++
			extTokens[filepath.Base(f)] += tokenCount
		}
	}

	printCountsAndTokens(extCounts, extTokens)
	return nil
}

func RunStatsParallel(cliCtx *cli.Context) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	extCounts := make(map[string]int)
	extTokens := make(map[string]int)
	var mu sync.Mutex // Mutex to protect concurrent map access

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	// Create a channel to receive results
	results := make(chan struct {
		ext   string
		count int
	}, len(files))

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Process files concurrently
	for _, f := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			tokenCount, err := CountTokensFilePath(filePath)
			if err != nil {
				log.Printf("Error counting tokens in file %s: %v", filePath, err)
				return
			}

			ext := strings.ToLower(filepath.Ext(filePath))
			if ext == "" {
				ext = filepath.Base(filePath)
			}

			results <- struct {
				ext   string
				count int
			}{ext, tokenCount}
		}(f)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		mu.Lock()
		extCounts[result.ext]++
		extTokens[result.ext] += result.count
		mu.Unlock()
	}

	printCountsAndTokens(extCounts, extTokens)
	return nil
}

func CountTokensFilePath(filePath string) (int, error) {
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

// /
func extensionOrBase(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != "" {
		return ext
	}
	return filepath.Base(filePath)
}
