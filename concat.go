package main

import (
	"fmt"
	_ "github.com/h2non/filetype/matchers"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"strings"
	"sync"
)

func RunConcat(cliCtx *cli.Command) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	var totalTokens int

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
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

			tokenCount, err := CountTokensInText(byteArr)
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

func RunConcatSequentially(cliCtx *cli.Command) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	var totalTokens int

	files, err := FindFiles(cliCtx.String("glob"), cliCtx.String("regex"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		byteArr, err := ReadFile(f)
		if err != nil {
			log.Fatalln("Error reading file:", err)
		}

		tokenCount, err := CountTokensInText(byteArr)
		if err != nil {
			log.Printf("Error counting extTokens in file %s: %v", f, err)
		}

		totalTokens += tokenCount

		fmt.Println("// BEGIN FILE:", f)
		fmt.Println(strings.TrimSpace(string(byteArr)))
		fmt.Println("// END FILE:", f)
		fmt.Println()
	}

	log.Println("Total tokens:", totalTokens)
	log.Println("Total files:", len(files))

	return nil
}
