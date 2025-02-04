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

func RunConcat(command *cli.Command) error {
	countTokens := GetTokenFunc(command.Bool("approx"))
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *ConcatFileResult, 100)
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
			tokenCount := countTokens(fileBytes)
			log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
			resultQueue <- &ConcatFileResult{
				Location:   f.Location,
				Content:    fileBytes,
				TokenCount: tokenCount}
		}(f)
	}

	wg.Wait()
	close(resultQueue)
}

// func RunConcatSequentially(cliCtx *cli.Command) error {
// 	// Get current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return fmt.Errorf("error getting current working directory: %w", err)
// 	}
// 	fmt.Println("Current working directory:", cwd)

// 	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
// 	if err != nil {
// 		return fmt.Errorf("error getting files: %w", err)
// 	}

// 	for _, f := range files {

// 	go func() {
// 		totalTokens := 0
// 		fileCount := 0

// 		for f := range resultQueue {
// 			fmt.Printf("// BEGIN FILE: %s\n", f.Location)
// 			fmt.Println(strings.TrimSpace(string(f.Content)))
// 			fmt.Printf("// END FILE: %s\n", f.Location)
// 			fmt.Println()

// 			totalTokens += f.TokenCount
// 			fileCount++
// 		}

// 		eprintln("Total tokens:", totalTokens)
// 		eprintln("Total files:", fileCount)
// 		close(done)
// 	}()

// 	processFilesConcat(fileQueue, resultQueue, countTokens)
// 	<-done
// 	return nil
// }
// }

// func processFilesConcatTmp(
// 	fileQueue chan *gocodewalker.File,
// 	resultQueue chan *ConcatFileResult,
// 	countTokens TokenFunc) {

// 	wg := sync.WaitGroup{}
// 	for f := range fileQueue {

// 		wg.Add(1)
// 		go func(f *gocodewalker.File) {
// 			log.Debug().Str("file", f.Location).Msg("Processing file")
// 			defer wg.Done()
// 			if !IsLikelyTextFile(f.Location) {
// 				return
// 			}

// 		tokenCount, err := CountTokensInText(byteArr)
// 		if err != nil {
// 			log.Printf("Error counting extTokens in file %s: %v", f, err)
// 		}

// 		totalTokens += tokenCount

// 		fmt.Println("// BEGIN FILE:", f)
// 		fmt.Println(strings.TrimSpace(string(byteArr)))
// 		fmt.Println("// END FILE:", f)
// 		fmt.Println()
// 		}(f)
// 	}

// 	// fmt.Println("Total tokens:", totalTokens)
// 	// fmt.Println("Total files:", len(files))

// }
