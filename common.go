package main

import (
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
	"regexp"
	"strings"
	"time"

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
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
}

func WalkFiles(bufSize int) chan *gocodewalker.File {
	fileListQueue := make(chan *gocodewalker.File, bufSize)

	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// restrict to only process files that have the .go extension
	//fileWalker.AllowListExtensions = append(fileWalker.AllowListExtensions, "go")
	fileWalker.ExcludeListExtensions = append(fileWalker.ExcludeListExtensions, "mod", "sum")
	//fileWalker.ExcludeFilenameRegex = append(fileWalker.ExcludeFilenameRegex, regexp.MustCompile("LICENSE"))
	fileWalker.ExcludeFilename = append(fileWalker.ExcludeFilename, "LICENSE")

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		fmt.Println("ERR", e.Error())
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			fmt.Println("ERR", err.Error())
		}
	}()

	return fileListQueue
}

// func readFile(filePath string) []byte {
// 	content, err := os.ReadFile(filePath)
// 	if err != nil {
// 		log.Err(err).Str("file", filePath).Msg("Error reading file")
// 		return []byte{}
// 	}
// 	return content
// }
//

func ConvertToHiArgs(command *cli.Command) *HiArgs {
	hiArgs := &HiArgs{
		approx: command.Bool("approx"),
		debug:  command.Bool("debug"),
	}
	if r := command.String("regex-content"); r != "" {
		hiArgs.regexContent = regexp.MustCompile(r)
	}
	if r := command.String("regex-filepath"); r != "" {
		hiArgs.regexFilepath = regexp.MustCompile(r)
	}
	return hiArgs
}

func filepathMatches(f *gocodewalker.File, regex *regexp.Regexp) bool {
	if regex == nil {
		return true
	}
	return regex.MatchString(f.Location)
}

func fileContentMatches(fc []byte, regex *regexp.Regexp) bool {
	if regex == nil {
		return true
	}
	return regex.MatchString(string(fc))
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

//func _dummyLog() {
//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
//	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
//
//	log.Debug().
//		Str("Scale", "833 cents").
//		Float64("Interval", 833.09).
//		Msg("Fibonacci is everywhere")
//
//	log.Debug().
//		Str("Name", "Tom").
//		Send()
//}
