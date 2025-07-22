package main

import (
	"fmt"
	"github.com/urfave/cli/v3"
)

func RunCount(ctx *cli.Command, command *LmcatArgs, pipeData string) error {
	countTokens := GetTokenFunc(command.approx)
	var tokenCount int
	var lineCount int
	if pipeData != "" {
		tokenCount = countTokens([]byte(pipeData))
		lineCount = countLines(pipeData)
	} else {
		firstArg := ctx.Args().First()
		tokenCount = countTokens([]byte(firstArg))
		lineCount = countLines(firstArg)
	}
	fmt.Println("Tokens:", tokenCount)
	fmt.Println("Lines:", lineCount)
	return nil
}

package main

import (
	"fmt"
	"github.com/urfave/cli/v3"
)

func RunCount(command *cli.Command) error {
	countTokens := GetTokenFunc(command.Bool("approx"))

	// Check for a positional argument (first non-flag argument)
	firstArg := command.Args().First()
	if firstArg == "" {
		return fmt.Errorf("no input provided, please specify a file, directory, or string")
	}

	// If not a file or directory, treat it as a string
	lineCount := countLines(firstArg)
	tokenCount := countTokens([]byte(firstArg))
	fmt.Printf("Lines: %d\nTokens: %d\n", lineCount, tokenCount)
	return nil
}
