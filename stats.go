package main

import (
	"sync"

	"github.com/rs/zerolog/log"

	"os"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
)

type StatsFileResult struct {
	Location string
	//Content    []byte
	TokenCount int
}

func RunStats(command *LmcatArgs, pipeData string) error {
	countTokens := GetTokenFunc(command.approx)
	var fileQueue chan string
	if pipeData != "" {
		fileQueue = ConvertPipeDataToChannel(pipeData, bufSize)
	} else {
		fileQueue = WalkFiles(bufSize)
	}
	resultQueue := make(chan *StatsFileResult, 100)
	processFilesStats(command, fileQueue, resultQueue, countTokens)

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

	<-done
	return nil
}

func processFilesStats(
	command *LmcatArgs,
	fileQueue chan string,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for filePath := range fileQueue {

		wg.Add(1)
		go func(filePath string) {
			log.Debug().Str("file", filePath).Msg("Processing file")
			defer wg.Done()
			if !IsLikelyTextFile(filePath) {
				return
			}
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				log.Err(err).Str("file", filePath).Msg("Error reading file")
				return
			}
			tokenCount := countTokens(fileBytes)
			log.Debug().Str("file", filePath).Int("tokenCount", tokenCount).Msg("Counted tokens")
			resultQueue <- &StatsFileResult{Location: filePath, TokenCount: tokenCount}
		}(filePath)
	}

	wg.Wait()
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
