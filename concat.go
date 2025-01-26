package main

import (
	"fmt"
	"github.com/tiktoken-go/tokenizer"
	"github.com/urfave/cli/v3"
	"log"
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
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
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
				log.Fatalln("Error reading file:", err)
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
	log.Println("Total tokens:", totalTokens)
	log.Println("Total files:", len(files))

	return nil
}
