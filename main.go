package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var converter = md.NewConverter("", true, nil)

func main() {
	app := &cli.App{
		Name:  "lmcat",
		Usage: "Process and concatenate files",
		Flags: []cli.Flag{
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
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cliCtx *cli.Context) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println("Current working directory:", cwd)

	// Initialize a map to store counts
	counts := make(map[string]int)

	// Create output file
	outputFile, err := os.Create("F:/output.txt")
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// Walk through the directory
	files, err := getRipgrepFiles(cliCtx.String("regex"), cliCtx.String("glob"))
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	for _, f := range files {
		fmt.Println("Visiting file:", f)
		// Get the file extension
		ext := strings.ToLower(filepath.Ext(f))
		if ext != "" {
			counts[ext]++
		}
		if ext == ".htm" || true {
			err := appendFileContent(f, writer)
			if err != nil {
				log.Printf("Error appending file %s: %v", f, err)
			}
		}
	}

	// Print the results
	for ext, count := range counts {
		fmt.Printf("%s: %d\n", ext, count)
	}

	fmt.Println("Text files have been concatenated into output.txt")

	// Count tokens in the output file
	tokenCount, err := countTokens("F:/output.txt")
	if err != nil {
		return fmt.Errorf("error counting tokens: %w", err)
	}

	fmt.Printf("Total tokens in output file: %d\n", tokenCount)

	if tokenCount > 1000 {
		log.Printf("WARNING: Token count (%d) exceeds 1000\n", tokenCount)
	}

	return nil
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
