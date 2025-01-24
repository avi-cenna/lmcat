package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := &cli.App{
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
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cliCtx *cli.Context) error {
	if cliCtx.Bool("stats") {
		return RunStats(cliCtx)
	} else if cliCtx.Bool("gcw") {
		main2()
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
