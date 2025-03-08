package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// High-level args
type HiArgs struct {
	regexContent  *regexp.Regexp
	regexFilepath *regexp.Regexp
	approx        bool
	debug         bool
	sequential    bool
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
}

func WalkFiles(bufSize int) chan *gocodewalker.File {
	fileListQueue := make(chan *gocodewalker.File, bufSize)

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

	return fileListQueue
}

func ConvertToHiArgs(command *cli.Command) *HiArgs {
	hiArgs := &HiArgs{
		approx:     command.Bool("approx"),
		debug:      command.Bool("debug"),
		sequential: command.Bool("sequential"),
	}
	if r := command.String("regex-content"); r != "" {
		hiArgs.regexContent = regexp.MustCompile(r)
	}
	if r := command.String("regex-filepath"); r != "" {
		hiArgs.regexFilepath = regexp.MustCompile(r)
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
