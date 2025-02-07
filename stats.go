package main

import (
	"sync"

	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog/log"

	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type StatsFileResult struct {
	Location string
	//Content    []byte
	TokenCount int
}

func RunStats(command *HiArgs) error {

	countTokens := GetTokenFunc(command.approx)
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *StatsFileResult, 100)
	done := make(chan struct{})

	go func() {
		extCounts := make(map[string]int)
		extTokens := make(map[string]int)
		for f := range resultQueue {
			ext := extensionOrBase(f.Location)
			extCounts[ext]++
			extTokens[ext] += f.TokenCount
		}
		printCountsAndTokens(extCounts, extTokens)
		close(done)
	}()

	processFiles(command, fileQueue, resultQueue, countTokens)
	<-done
	return nil
}

func processFiles(
	command *HiArgs,
	fileQueue chan *gocodewalker.File,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for f := range fileQueue {

		wg.Add(1)
		go func(f *gocodewalker.File) {
			log.Debug().Str("file", f.Location).Msg("Processing file")
			defer wg.Done()
			if !IsLikelyTextFile(f.Location) {
				return
			}
			if command.regexFilepath != nil && !command.regexFilepath.MatchString(f.Location) {
				return
			}
			fileBytes, err := os.ReadFile(f.Location)
			if err != nil {
				log.Err(err).Str("file", f.Location).Msg("Error reading file")
				return
			}
			if command.regexContent != nil && !command.regexContent.Match(fileBytes) {
				return
			}
			tokenCount := countTokens(fileBytes)
			log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
			resultQueue <- &StatsFileResult{Location: f.Location, TokenCount: tokenCount}
		}(f)
	}

	wg.Wait()
	close(resultQueue)
}

func processFilesSequential(
	command *HiArgs,
	fileQueue chan *gocodewalker.File,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	for f := range fileQueue {
		log.Debug().Str("file", f.Location).Msg("Processing file")
		
		if !IsLikelyTextFile(f.Location) {
			continue
		}
		if command.regexFilepath != nil && !command.regexFilepath.MatchString(f.Location) {
			continue
		}
		
		fileBytes, err := os.ReadFile(f.Location)
		if err != nil {
			log.Err(err).Str("file", f.Location).Msg("Error reading file")
			continue
		}
		
		if command.regexContent != nil && !command.regexContent.Match(fileBytes) {
			continue
		}
		
		tokenCount := countTokens(fileBytes)
		log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
		resultQueue <- &StatsFileResult{Location: f.Location, TokenCount: tokenCount}
	}
	close(resultQueue)

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
