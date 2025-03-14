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

func RunStats(command *HiArgs) error {

	countTokens := GetTokenFunc(command.approx)
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *StatsFileResult, 100)
	done := processStatsResults(resultQueue)

	if command.sequential {
		processFilesSequential(command, fileQueue, resultQueue, countTokens)
	} else {
		processFileQueueForStats(fileQueue, resultQueue, countTokens)
	}
	<-done
	return nil
}

// processStatsResults handles processing and displaying statistics from the result channel
func processStatsResults(resultQueue chan *StatsFileResult) chan struct{} {
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

	return done
}

func processFileQueueForStats(
	fileQueue chan string,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for filepath := range fileQueue {
		wg.Add(1)
		go processStatsFileWorker(&wg, filepath, resultQueue, countTokens)
	}

	wg.Wait()
	close(resultQueue)
}

func processFileListForStats(
	fileList []string,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for _, filepath := range fileList {
		wg.Add(1)
		go processStatsFileWorker(&wg, filepath, resultQueue, countTokens)
	}

	wg.Wait()
	close(resultQueue)
}

func processFilesSequential(
	command *HiArgs,
	fileQueue chan string,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	for filepath := range fileQueue {
		log.Debug().Str("file", filepath).Msg("Processing file")

		if !IsLikelyTextFile(filepath) {
			continue
		}

		fileBytes, err := os.ReadFile(filepath)
		if err != nil {
			log.Err(err).Str("file", filepath).Msg("Error reading file")
			continue
		}

		tokenCount := countTokens(fileBytes)
		log.Debug().Str("file", filepath).Int("tokenCount", tokenCount).Msg("Counted tokens")
		resultQueue <- &StatsFileResult{Location: filepath, TokenCount: tokenCount}
	}
	close(resultQueue)
}

// processStatsFileWorker handles the processing of a single file for stats
// collection
func processStatsFileWorker(
	wg *sync.WaitGroup,
	filepath string,
	resultQueue chan *StatsFileResult,
	countTokens TokenFunc) {

	defer wg.Done()
	log.Debug().Str("file", filepath).Msg("Processing file")

	if !IsLikelyTextFile(filepath) {
		return
	}

	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		log.Err(err).Str("file", filepath).Msg("Error reading file")
		return
	}

	tokenCount := countTokens(fileBytes)
	log.Debug().Str("file", filepath).Int("tokenCount", tokenCount).Msg("Counted tokens")
	resultQueue <- &StatsFileResult{Location: filepath, TokenCount: tokenCount}
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
