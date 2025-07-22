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
