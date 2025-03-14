package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// High-level args
type HiArgs struct {
	//regexContent  *regexp.Regexp
	//regexFilepath *regexp.Regexp
	approx     bool
	count      bool
	debug      bool
	sequential bool
	stats      bool
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
}

func WalkFiles(bufSize int) chan string {
	fileListQueue := make(chan *gocodewalker.File, bufSize)
	filePathQueue := make(chan string, bufSize)

	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// restrict to only process files that have the .go extension
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")
	fileWalker.ExcludeListExtensions = append(
		fileWalker.ExcludeListExtensions,
		"mod",
		"sum",
		"lock",
		"bsp",
	)
	fileWalker.ExcludeFilename = append(fileWalker.ExcludeFilename, "LICENSE", "package-lock.json")

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		log.Err(e).Msg("Error walking files")
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			log.Err(err).Msg("Error starting file walker")
		}
	}()

	go func() {
		defer close(filePathQueue)
		for file := range fileListQueue {
			filePathQueue <- file.Location
		}
	}()

	return filePathQueue
}

func ConvertToHiArgs(command *cli.Command) *HiArgs {
	hiArgs := &HiArgs{
		approx:     command.Bool("approx"),
		count:      command.Bool("count"),
		debug:      command.Bool("debug"),
		sequential: command.Bool("sequential"),
		stats:      command.Bool("stats"),
	}
	return hiArgs
}

func countLines(text string) int {
	n := strings.Count(text, "\n")
	if len(text) > 0 && !strings.HasSuffix(text, "\n") {
		n++
	}
	return n
}

func eprintln(a ...interface{}) {
	if _, err := fmt.Fprintln(os.Stderr, a...); err != nil {
		log.Err(err).Msg("Error writing to stderr")
	}
}

func extensionOrBase(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != "" {
		return ext
	}
	return filepath.Base(filePath)
}
