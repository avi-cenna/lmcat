package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"github.com/urfave/cli"
	"github.com/urfave/cli/v3"
	"io"
	"os"
)

func main() {
	app := &cli.Command{
		Name:  "lmcat",
		Usage: "Process and concatenate files",
		Flags: []cli.Flag{
			// TODO (AA) #51: add a flag here that will force sequential running
			//            of stats and concatenation commands
			&cli.StringFlag{
				Name:     "regex-content",
				Required: false,
				Usage:    "Regex pattern for ripgrep",
			},
			&cli.StringFlag{
				Name:     "regex-filepath",
				Required: false,
				Usage:    "Regex pattern for ripgrep",
			},
			&cli.BoolFlag{
				Name:     "stats",
				Required: false,
				Usage:    "Run file stats",
			},
			&cli.BoolFlag{
				Name:     "approx",
				Required: false,
				Usage:    "Use approximate token count",
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

	// Check if there's data being piped in
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped in
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		
		countTokens := GetTokenFunc(cliCtx.Bool("approx"))
		tokenCount := countTokens(data)
		
		_, err = fmt.Fprintf(os.Stdout, "Token count: %d\n", tokenCount)
		return err
	}

	hiArgs := ConvertToHiArgs(cliCtx)
	if hiArgs.debug {
		globalLevel = zerolog.DebugLevel
	} else {
		globalLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(globalLevel)

	if cliCtx.Bool("stats") {
		return RunStats(hiArgs)
	} else if cliCtx.Bool("gcw") {
		return nil
	}
	return RunConcat(cliCtx)
}

//func getRipgrepFiles(regex, glob string) ([]string, error) {
//	fmt.Println("glob pattern:", glob)
//	args := []string{"--files-with-matches", regex}
//	if glob != "" {
//		args = append(args, "--glob", glob)
//	}
//	cmd := exec.Command("rg", args...)
//	output, err := cmd.Output()
//	if err != nil {
//		return nil, fmt.Errorf("error running ripgrep: %w", err)
//	}
//	files := strings.Split(strings.TrimSpace(string(output)), "\n")
//	return files, nil
//}
