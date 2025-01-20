package main

import (
	_ "github.com/h2non/filetype/matchers"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tiktoken-go/tokenizer"
	"github.com/urfave/cli"
	_ "log"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	_ "sort"
	"strings"
)

type TokenCounter struct {
	encoder *tokenizer.Codec
}

func NewTokenCounter() (*TokenCounter, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tokenizer: %w", err)
	}
	return &TokenCounter{encoder: enc}, nil
}

func RunStats(cliCtx *cli.Context) error {
	return nil
}

func (tc *TokenCounter) CountTokensFilePath(filePath string) (int, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return tc.CountTokensInText(content)
}

// TODO: adjust this so that it accepts a pointer to the Codec tokenizer instead of creating it each time.
func CountTokensInText(text []byte) (int, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		return 0, err
	}

	ids, _, err := enc.Encode(string(text))
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
