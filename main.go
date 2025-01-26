package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	//"github.com/urfave/cli"
	"github.com/urfave/cli/v3"

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
				//Aliases: []string{"g"},
				Usage: "Glob pattern for ripgrep",
			},
			&cli.StringFlag{
				Name:     "regex",
				Required: false,
				Usage:    "Regex pattern for ripgrep",
			},
			&cli.BoolFlag{
				Name:     "stats",
				Required: false,
				Usage:    "Run file stats",
			},
			&cli.BoolFlag{
				Name:     "gcw",
				Required: false,
				Usage:    "Gocodewalker",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Aliases:  []string{"d", "verbose", "v"},
				Required: false,
				Usage:    "Enable debug logging",
			},
			&cli.BoolFlag{
				Name:     "approx",
				Required: false,
				Usage:    "Use approximate token count",
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
	if cliCtx.Bool("debug") {
		globalLevel = zerolog.DebugLevel
	} else {
		globalLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(globalLevel)
	//initTokenizer(cliCtx.Bool("approx"))

	if cliCtx.Bool("stats") {
		return RunStats(cliCtx)
	} else if cliCtx.Bool("gcw") {
		main2()
		return nil
	}
	return RunConcatOld(cliCtx)
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
