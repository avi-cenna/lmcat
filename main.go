package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	app := &cli.Command{
		Name:  "lmcat",
		Usage: "Process and concatenate files",
		Flags: []cli.Flag{
			// TODO (AA): add a flag here that will force sequential running
			//            of stats and concatenation commands
			&cli.StringFlag{
				Name:     "glob",
				Required: false,
				Aliases:  []string{"g"},
				Usage:    "Glob pattern for ripgrep",
			},
			&cli.StringFlag{
				Name:     "regex",
				Required: false,
				Aliases:  []string{"r"},
				Usage:    "Regex pattern for ripgrep",
			},
			&cli.BoolFlag{
				Name:     "stats",
				Required: false,
				Aliases:  []string{"s"},
				Usage:    "Run file stats",
			},
			&cli.BoolFlag{
				Name:     "seq",
				Required: false,
				Usage:    "Run sequentially",
			},
		},
		Action: run,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

	// if err := app.Run(os.Args); err != nil {
	// log.Fatal(err)
	// }
}

func run(cliCtx context.Context, cmd *cli.Command) error {

	Test1()
	return nil

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	foo := WalkFilesChan(cwd)
	for f := range foo {
		fmt.Println(f)
	}

	return nil

	if cmd.Bool("stats") {
		if cmd.Bool("seq") {
			return RunStatsSequential(cmd)
		} else {
			return RunStats(cmd)
		}
	}
	if cmd.Bool("seq") {
		return RunConcatSequentially(cmd)
	} else {
		return RunConcat(cmd)
	}
}
