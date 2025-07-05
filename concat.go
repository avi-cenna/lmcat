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

	var fileQueue chan string
	if pipeData != "" {
		fileQueue = ConvertPipeDataToChannel(pipeData, bufSize)
	} else {
		fileQueue = WalkFiles(bufSize)
	}

	resultQueue := make(chan *ConcatFileResult, bufSize)
	processFilesConcat(fileQueue, resultQueue, countTokens)

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

		errprintln("Total tokens:", totalTokens)
		errprintln("Total files:", fileCount)
		close(done)
	}()

	<-done

	return nil
}

func processFilesConcat(
	fileQueue chan string,
	resultQueue chan *ConcatFileResult,
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
			resultQueue <- &ConcatFileResult{
				Location:   filePath,
				Content:    fileBytes,
				TokenCount: tokenCount}
		}(filePath)
	}

	wg.Wait()
	close(resultQueue)
}
