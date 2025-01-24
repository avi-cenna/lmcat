package main

import (
	"fmt"
	"github.com/boyter/gocodewalker"
	"log"
	"sync"

	//_ "github.com/h2non/filetype/matchers"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tiktoken-go/tokenizer"
	"github.com/urfave/cli"
	_ "log"
	"os"
	"path/filepath"
	"slices"
	_ "sort"
	"strings"
)

//type TokenCounter struct {
//	encoder *tokenizer.Codec
//}

//func NewTokenCounter() (*TokenCounter, error) {
//	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
//	if err != nil {
//		return nil, fmt.Errorf("failed to initialize tokenizer: %w", err)
//	}
//	return &TokenCounter{encoder: enc}, nil
//}

type FileResult struct {
	// Full path to the file
	Location string
	// File contents
	Content []byte
	// Count of tokens in the file
	TokenCount int
}

func RunStats(cliCtx *cli.Context) error {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return err
	}

	fileQueue := WalkFiles()
	resultQueue := make(chan *FileResult, 100)
	fmt.Println(resultQueue)

	extCounts := make(map[string]int)
	extTokens := make(map[string]int)
	for f := range fileQueue {

		fmt.Println(f.Location)
		fileBytes, err := os.ReadFile(f.Location)
		if err != nil {
			log.Println("Error reading file:", err)
		}
		tokenCount, err := CountTokensInText(enc, fileBytes)
		if err != nil {
			log.Printf("Error counting tokens in file %s: %v", f.Location, err)
		}
		ext := extensionOrBase(f.Location)
		extCounts[ext]++
		extTokens[ext] += tokenCount
	}

	printCountsAndTokens(extCounts, extTokens)

	return nil
}

func processFiles(
	fileQueue chan *gocodewalker.File,
	resultQueue chan *FileResult,
	codec tokenizer.Codec) {
	wg := sync.WaitGroup{}
	for f := range fileQueue {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fileBytes, err := os.ReadFile(f.Location)
			if err != nil {
				log.Println("Error reading file:", err)
				return
			}
			tokenCount, err := CountTokensInText(codec, fileBytes)
			if err != nil {
				log.Printf("Error counting tokens in file %s: %v", f.Location, err)
				return
			}
			resultQueue <- &FileResult{Location: f.Location, Content: fileBytes, TokenCount: tokenCount}
		}()
	}
	wg.Wait()
	close(resultQueue)
}

// ..
// Note: tokenizer.Codec is an interface type

func CountTokensInText(codec tokenizer.Codec, text []byte) (int, error) {
	ids, _, err := codec.Encode(string(text))
	if err != nil {
		return 0, err
	}
	return len(ids), nil
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

func extensionOrBase(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != "" {
		return ext
	}
	return filepath.Base(filePath)
}
