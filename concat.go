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

func RunConcat(command *cli.Command, pipeData string) error {
	countTokens := GetTokenFunc(command.Bool("approx"))

	// Create the file queue which will be used to process files
	var fileQueue chan string
	if pipeData != "" {
		fileQueue = ConvertPipeDataToChannel(pipeData, bufSize)
	} else {
		fileQueue = WalkFiles(bufSize)
	}

	log.Debug().Msg("processFilesConcat: Starting file processing")
	resultQueue := make(chan *ConcatFileResult, bufSize)
	done := make(chan struct{})

	// Start a goroutine to process the result queue
	go func() {
		log.Debug().Msg("RunConcat: Result processing goroutine started")
		totalTokens := 0
		fileCount := 0

		for f := range resultQueue {
			log.Debug().Str("file", f.Location).Int("tokens", f.TokenCount).Msg("RunConcat: Processing result")
			fmt.Printf("// BEGIN FILE: %s\n", f.Location)
			fmt.Println(strings.TrimSpace(string(f.Content)))
			fmt.Printf("// END FILE: %s\n", f.Location)
			fmt.Println()

			totalTokens += f.TokenCount
			fileCount++
		}

		log.Debug().Int("totalFiles", fileCount).Int("totalTokens", totalTokens).Msg("RunConcat: Finished processing all results")
		errprintln("Total tokens:", totalTokens)
		errprintln("Total files:", fileCount)
		log.Debug().Msg("RunConcat: Result processing goroutine finished")
		close(done)
	}()

	// Start processing files in parallel
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
			resultQueue <- &ConcatFileResult{
				Location:   filePath,
				Content:    fileBytes,
				TokenCount: tokenCount}
		}(filePath)
	}

	log.Debug().Msg("processFilesConcat: Waiting for all goroutines to complete")
	wg.Wait()
	log.Debug().Msg("processFilesConcat: All goroutines completed, closing resultQueue")
	close(resultQueue)

	<-done

	return nil
}

//func processFilesConcat(
//	fileQueue chan string,
//	resultQueue chan *ConcatFileResult,
//	countTokens TokenFunc) {
//
//	log.Debug().Msg("processFilesConcat: Starting file processing")
//	wg := sync.WaitGroup{}
//	for filePath := range fileQueue {
//
//		wg.Add(1)
//		go func(filePath string) {
//			log.Debug().Str("file", filePath).Msg("Processing file")
//			defer wg.Done()
//			if !IsLikelyTextFile(filePath) {
//				return
//			}
//			fileBytes, err := os.ReadFile(filePath)
//			if err != nil {
//				log.Err(err).Str("file", filePath).Msg("Error reading file")
//				return
//			}
//			tokenCount := countTokens(fileBytes)
//			log.Debug().Str("file", filePath).Int("tokenCount", tokenCount).Msg("Counted tokens")
//			resultQueue <- &ConcatFileResult{
//				Location:   filePath,
//				Content:    fileBytes,
//				TokenCount: tokenCount}
//		}(filePath)
//	}
//
//	log.Debug().Msg("processFilesConcat: Waiting for all goroutines to complete")
//	wg.Wait()
//	log.Debug().Msg("processFilesConcat: All goroutines completed, closing resultQueue")
//	close(resultQueue)
//}
