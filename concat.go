package main

import (
	"fmt"
	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog/log"
	"github.com/tiktoken-go/tokenizer"
	"github.com/urfave/cli/v3"
	"os"
	"strings"
	"sync"
)

type ConcatFileResult struct {
	Location   string
	Content    []byte
	TokenCount int
}

func RunConcat(command *cli.Command) error {
	fileQueue := WalkFiles(100)
	resultQueue := make(chan *ConcatFileResult, 100)
	done := make(chan struct{})

	go func() {
		//extCounts := make(map[string]int)
		//extTokens := make(map[string]int)
		//TODO: alter this so that for each file it does the following:
		//   - print BEGIN FILE: <filename>, to stdout
		//   - print the file content, to stdout
		//   - print END FILE: <filename>, to stdout
		//   - print a newline, to stdout
		//   - At the end, print the total number of tokens and # files scanned, make sure it's printed to stderr
		for f := range resultQueue {
			ext := extensionOrBase(f.Location)
			fmt.Println("ext:", ext, f.Location)
			//extCounts[ext]++
			//extTokens[ext] += f.TokenCount
		}
		//printCountsAndTokens(extCounts, extTokens)
		close(done)
	}()

	processFilesConcat(fileQueue, resultQueue)
	<-done
	return nil
}

func processFilesConcat(
	fileQueue chan *gocodewalker.File,
	resultQueue chan *ConcatFileResult) {

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
			tokenCount, err := tokenCounter.CountTokens(fileBytes)
			log.Debug().Str("file", f.Location).Int("tokenCount", tokenCount).Msg("Counted tokens")
			if err != nil {
				log.Err(err).Str("file", f.Location).Msg("Error counting tokens")
				return
			}
			resultQueue <- &ConcatFileResult{
				Location:   f.Location,
				Content:    fileBytes,
				TokenCount: tokenCount}
		}(f)
	}

	wg.Wait()
	close(resultQueue)
}

// TODO (AA): refactor as necessary to be in more idiomatic Go
func RunConcatOld(command *cli.Command) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	var totalTokens int

	files, err := FindFiles(command.String("glob"), command.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, f := range files {

		go func() {
			wg.Add(1)
			defer wg.Done()

			byteArr, err := ReadFile(f)
			if err != nil {
				log.Fatal().Err(err).Msg("Error reading file")
			}

			tokenCount, err := CountTokensInText(enc, byteArr)
			if err != nil {
				log.Printf("Error counting extTokens in file %s: %v", f, err)
			}

			mu.Lock()
			totalTokens += tokenCount
			mu.Unlock()

			mu.Lock()
			fmt.Println("// BEGIN FILE:", f)
			fmt.Println(strings.TrimSpace(string(byteArr)))
			fmt.Println("// END FILE:", f)
			fmt.Println()
			mu.Unlock()
		}()
	}

	wg.Wait()

	// error print ln here
	_, err = fmt.Fprintln(os.Stderr, "Total tokens:", totalTokens)
	_, err = fmt.Fprintln(os.Stderr, "Total files:", len(files))
	if err != nil {
		return err
	}

	return nil
}
