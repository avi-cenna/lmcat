package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"io"
	"os"
)

func main() {
	app := &cli.Command{
		Name:  "lmcat",
		Usage: "Process and concatenate files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "regex-content",
				Required: false,
				Usage:    "Regex pattern for ripgrep",
			},
			&cli.StringFlag{
				Name:     "regex-filepath",
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
				Name:     "approx",
				Required: false,
				Usage:    "Use approximate token count",
			},
			&cli.BoolFlag{
				Name:     "sequential",
				Aliases:  []string{"seq"},
				Required: false,
				Usage:    "Force sequential processing",
			},
			&cli.BoolFlag{
				Name:     "count",
				Aliases:  []string{"c"},
				Required: false,
				Usage:    "Count lines and tokens in input (piped data, file, directory, or string)",
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
	hiArgs := ConvertToHiArgs(cliCtx)
	if hiArgs.debug {
		globalLevel = zerolog.DebugLevel
	} else {
		globalLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(globalLevel)

	// Check if there's data being piped in
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		return RunPipe(data, hiArgs)
	}

	// Handle different commands based on flags
	if cliCtx.Bool("count") {
		return RunCount(cliCtx)
	} else if cliCtx.Bool("stats") {
		return RunStats(hiArgs)
	} else {
		return RunConcat(cliCtx)
	}
}
