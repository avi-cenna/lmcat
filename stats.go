package main

import (
	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog/log"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	tiktoken "github.com/tiktoken-go/tokenizer"
	"github.com/urfave/cli/v3"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type StatsFileResult struct {
	Location string
	//Content    []byte
	TokenCount int
}

func RunStats(command *cli.Command) error {
	enc, err := tiktoken.Get(tiktoken.Cl100kBase)
	if err != nil {
		return err
	}

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

	processFiles(fileQueue, resultQueue, enc)
	<-done
	return nil
}

func processFiles(
	fileQueue chan *gocodewalker.File,
	resultQueue chan *StatsFileResult,
	codec tiktoken.Codec) {

	wg := sync.WaitGroup{}
	for f := range fileQueue {

		wg.Add(1)
		go func(f *gocodewalker.File) {
			log.Debug().Str("file", f.Location).Msg("Processing file")
			defer wg.Done()
			if !IsLikelyTextFile(f.Location) {
				return
			}
			fileBytes, err := os.ReadFile(f.Location)
			if err != nil {
				log.Err(err).Str("file", f.Location).Msg("Error reading file")
				return
			}
			tokenCount, err := CountTokensInText(codec, fileBytes)
			log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
			if err != nil {
				log.Err(err).Str("file", f.Location).Msg("Error counting tokens")
				return
			}
			resultQueue <- &StatsFileResult{Location: f.Location, TokenCount: tokenCount}
		}(f)
	}

	wg.Wait()
	close(resultQueue)
}

func CountTokensInText(codec tiktoken.Codec, text []byte) (int, error) {
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
