package main

import (
	"fmt"
	_ "github.com/h2non/filetype/matchers"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	// Initialize a map to store extCounts
	extCounts := make(map[string]int)
	extTokens := make(map[string]int)

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		tokenCount, err := CountTokens(f)
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

func CountTokens(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		return 0, err
	}
	tokens := tke.Encode(string(content), nil, nil)
	return len(tokens), nil
}

func printCountsAndTokens(extCounts map[string]int, extTokens map[string]int) {
	sortedKeys := make([]string, 0, len(extCounts))
	for ext := range extCounts {
		sortedKeys = append(sortedKeys, ext)
	}
	sort.Strings(sortedKeys)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Extension", "Count", "Tokens"})
	for _, ext := range sortedKeys {
		t.AppendRow(table.Row{ext, extCounts[ext], extTokens[ext]})
	}
	t.Render()
}
