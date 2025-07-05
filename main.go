package main

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "lmcat",
		Usage: "Process and concatenate files",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "stats",
				Required: false,
				Aliases:  []string{"s"},
				Usage:    "Run file stats",
			},
			&cli.BoolFlag{
				Name:     "approx",
				Required: false,
				Usage:    "Use approximate token count",
			},
			//TODO (for AI):  implement this flag:
			//   - if there is piped data, count the lines and tokens
			//   - if there is no piped data, then get the first positional argument, and
			//     if it's a file, count tokens in the file, if it's a directory, count tokens in the directory,
			//     if it's a string, count tokens and lines in the string
			&cli.BoolFlag{
				Name: "count",
				//Aliases:  []string{"count"},
				Required: false,
				Usage:    "TODO",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Aliases:  []string{"d", "verbose", "v"},
				Required: false,
				Usage:    "Enable debug logging",
			},
		},
		Action: run,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("Error running app")
	}
}

func run(ctx context.Context, cliCtx *cli.Command) error {
	var globalLevel zerolog.Level
	hiArgs := ParseArgs(cliCtx)
	if hiArgs.debug {
		globalLevel = zerolog.DebugLevel
	} else {
		globalLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(globalLevel)

	// Check if there's data being piped in
	pipedData, err := getPipedData()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading piped data")
	}
	//if pipedData != "" {
	//	lineCount := countLines(string(data))
	//	tokenCount := GetTokenFunc(hiArgs.approx)(data)
	//	_, err = fmt.Fprintf(os.Stdout, "Lines: %d\nTokens: %d", lineCount, tokenCount)
	//	return err
	//}

	if hiArgs.count {
		return RunCount(cliCtx, hiArgs, pipedData)
	} else if hiArgs.stats {
		return RunStats(hiArgs, pipedData)
	} else {
		return RunConcat(cliCtx, pipedData)
	}
}

func getPipedData() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}
