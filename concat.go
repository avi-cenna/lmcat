package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type ConcatFileResult struct {
	Location   string
	Content    []byte
	TokenCount int
}

// ProcessResults handles displaying file content and statistics from the result channel
func ProcessResults(resultQueue chan *ConcatFileResult) chan struct{} {
	done := make(chan struct{})

	go func() {
		totalTokens := 0
		fileCount := 0

		for f := range resultQueue {
			fmt.Printf("// BEGIN FILE: %s\n", f.Location)
			fmt.Println(strings.TrimSpace(string(f.Content)))
			fmt.Printf("// END FILE: %s\n", f.Location)
			fmt.Println()

			totalTokens += f.TokenCount
			fileCount++
		}

		eprintln("Total tokens:", totalTokens)
		eprintln("Total files:", fileCount)
		close(done)
	}()

	return done
}

func RunConcat(command *cli.Command) error {
	countTokens := GetTokenFunc(command.Bool("approx"))
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *ConcatFileResult, 100)

	done := ProcessResults(resultQueue)

	processFilesConcat(fileQueue, resultQueue, countTokens)
	<-done
	return nil
}

func processFilesConcat(
	fileQueue chan *gocodewalker.File,
	resultQueue chan *ConcatFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for f := range fileQueue {
		wg.Add(1)
		go processFileWorker(&wg, f.Location, resultQueue, countTokens)
	}

	wg.Wait()
	close(resultQueue)
}

func processFileListConcat(
	fileList []string,
	resultQueue chan *ConcatFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for _, f := range fileList {
		wg.Add(1)
		go processFileWorker(&wg, f, resultQueue, countTokens)
	}

	wg.Wait()
	close(resultQueue)
}

// processFileWorker handles the processing of a single file
func processFileWorker(
	wg *sync.WaitGroup,
	filepath string,
	resultQueue chan *ConcatFileResult,
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

	resultQueue <- &ConcatFileResult{
		Location:   filepath,
		Content:    fileBytes,
		TokenCount: tokenCount,
	}
}

// func processFilesConcatSequential(
// 	fileQueue chan *gocodewalker.File,
// 	resultQueue chan *ConcatFileResult,
// 	countTokens TokenFunc) {

// 	for f := range fileQueue {
// 		log.Debug().Str("file", f.Location).Msg("Processing file")
// 		if !IsLikelyTextFile(f.Location) {
// 			return
// 		}
// 		fileBytes, err := os.ReadFile(f.Location)
// 		if err != nil {
// 			log.Err(err).Str("file", f.Location).Msg("Error reading file")
// 			return
// 		}
// 		tokenCount := countTokens(fileBytes)
// 		log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
// 		resultQueue <- &ConcatFileResult{
// 			Location:   f.Location,
// 			Content:    fileBytes,
// 			TokenCount: tokenCount}
// 	}

// 	close(resultQueue)
// }
