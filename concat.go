package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type ConcatFileResult struct {
	Location   string
	Content    []byte
	TokenCount int
}

//// ProcessConcatResults handles displaying file content and statistics from the result channel
//func ProcessConcatResults(resultQueue chan *ConcatFileResult) chan struct{} {
//	done := make(chan struct{})
//
//	go func() {
//		totalTokens := 0
//		fileCount := 0
//
//		for f := range resultQueue {
//			fmt.Printf("// BEGIN FILE: %s\n", f.Location)
//			fmt.Println(strings.TrimSpace(string(f.Content)))
//			fmt.Printf("// END FILE: %s\n", f.Location)
//			fmt.Println()
//
//			totalTokens += f.TokenCount
//			fileCount++
//		}
//
//		eprintln("Total tokens:", totalTokens)
//		eprintln("Total files:", fileCount)
//		close(done)
//	}()
//
//	return done
//}

func RunConcat(command *cli.Command) error {
	countTokens := GetTokenFunc(command.Bool("approx"))
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *ConcatFileResult, 100)

	wg := sync.WaitGroup{}
	for filepath := range fileQueue {
		wg.Add(1)
		go func(filepath string) {
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
		}(filepath)
	}

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

	wg.Wait()
	close(resultQueue)
	<-done
	return nil
}

func processFileListConcat(
	fileList []string,
	resultQueue chan *ConcatFileResult,
	countTokens TokenFunc) {

	wg := sync.WaitGroup{}
	for _, filepath := range fileList {
		wg.Add(1)
		go func(filepath string) {
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
		}(filepath)
	}

	wg.Wait()
	close(resultQueue)
}
