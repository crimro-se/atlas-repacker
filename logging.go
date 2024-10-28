package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// If any errors exist, logs the (optional) msg and all errors.
func logErrors(errs []error, msg ...string) {
	if len(errs) == 0 {
		return
	}

	if len(msg) == 1 {
		log.Error().Msg(msg[0])
	}
	for _, err := range errs {
		log.Err(err).Send()
	}
}

// if err isn't nil, show it to the user and optionally a message, then os.exit
func errHandler(err error, msg ...string) {
	if err != nil {
		if len(msg) == 1 {
			log.Fatal().Err(err).Msg(msg[0])
		}
		log.Fatal().Err(err).Send()
	}
}

// convenience wrapper for any function returning one value and an error.
// If the error isn't nil, logs and dies. msg is optional.
func must1[T any](x T, err error, msg ...string) T {
	errHandler(err, msg...)
	return x
}

// convenience wrapper for any function returning two values and an error.
// If the error isn't nil, logs and dies. msg is optional.
func must2[T1, T2 any](x T1, y T2, err error, msg ...string) (T1, T2) {
	errHandler(err, msg...)
	return x, y
}
