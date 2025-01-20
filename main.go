package main

import (
	"bufio"
	"bytes"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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
	}
	return RunConcat(cliCtx)
}

func getRipgrepFiles(regex, glob string) ([]string, error) {
	fmt.Println("glob pattern:", glob)
	args := []string{"--files-with-matches", regex}
	if glob != "" {
		args = append(args, "--glob", glob)
	}
	cmd := exec.Command("rg", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running ripgrep: %w", err)
	}
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	return files, nil
}

func appendFileContent(filePath string, writer *bufio.Writer) error {
	fmt.Println("Trying to append file content:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = writer.WriteString(fmt.Sprintf("\n// Contents of %s \n", filePath))
	if err != nil {
		return err
	}

	markdown, err := convertHTMLToMarkdown(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, &markdown)
	if err != nil {
		return err
	}

	_, err = writer.WriteString("\n")
	return err
}

func convertHTMLToMarkdown(filePath string) (bytes.Buffer, error) {
	emptyBuffer := bytes.Buffer{}
	file, err := os.Open(filePath)
	if err != nil {
		return emptyBuffer, err
	}
	defer file.Close()

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertReader(file)
	if err != nil {
		return emptyBuffer, err
	}

	return markdown, nil
}

func countTokens(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	encoding := "cl100k_base"
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, err
	}

	tokens := tke.Encode(string(content), nil, nil)
	return len(tokens), nil
}
