package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})
}

func dummyLog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen})

	log.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("Fibonacci is everywhere")

	log.Debug().
		Str("Name", "Tom").
		Send()
}
