package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/boyter/gocodewalker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const bufSize = 100

type LmcatArgs struct {
	approx bool
	count  bool
	debug  bool
	stats  bool
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
}

func ConvertPipeDataToChannel(pipeData string, bufSize int) chan string {
	fileQueue := make(chan string, bufSize)
	go func() {
		defer close(fileQueue)
		lines := strings.Split(pipeData, "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" {
				fileQueue <- trimmedLine
			}
		}
	}()
	return fileQueue
}

// WalkFiles returns a channel of file paths by walking the current directory
func WalkFiles(bufSize int) chan string {
	log.Debug().Msg("WalkFiles: Starting file walking")
	// Original channel for gocodewalker.File objects
	fileListQueue := make(chan *gocodewalker.File, bufSize)
	// New channel for file paths
	filePaths := make(chan string, bufSize)

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
		log.Debug().Msg("WalkFiles: File walker goroutine started")
		err := fileWalker.Start()
		if err != nil {
			log.Err(err).Msg("Error starting file walker")
		}
		log.Debug().Msg("WalkFiles: File walker goroutine finished")
	}()

	// Transform File objects to file paths
	go func() {
		log.Debug().Msg("WalkFiles: Transformation goroutine started")
		defer close(filePaths)
		fileCount := 0
		for file := range fileListQueue {
			fileCount++
			log.Debug().
				Str("file", file.Location).
				Int("count", fileCount).
				Msg("WalkFiles: Processing file")
			filePaths <- file.Location
		}
		log.Debug().Msg("WalkFiles: fileListQueue closed, closing filePaths channel")
		log.Debug().Int("totalFiles", fileCount).Msg("WalkFiles: Transformation goroutine finished")
	}()

	return filePaths
}

func ParseArgs(command *cli.Command) *LmcatArgs {
	hiArgs := &LmcatArgs{
		approx: command.Bool("approx"),
		count:  command.Bool("count"),
		debug:  command.Bool("debug"),
		stats:  command.Bool("stats"),
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

//func println(a ...interface{}) {
//	if _, err := fmt.Println(a...); err != nil {
//		log.Err(err).Msg("Error writing to stderr")
//	}
//}

func errprintln(a ...interface{}) {
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
