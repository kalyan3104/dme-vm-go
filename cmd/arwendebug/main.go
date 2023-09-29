package main

import (
	"os"

	logger "github.com/kalyan3104/dme-logger-go"
	"github.com/kalyan3104/dme-vm-go/arwendebug"
)

var log = logger.GetOrCreate("arwendebug")

const (
	// ErrCodeSuccess signals success
	ErrCodeSuccess = iota
	// ErrCodeCriticalError signals a critical error
	ErrCodeCriticalError
)

func main() {
	logger.ToggleLoggerName(true)
	_ = logger.SetLogLevel("*:TRACE")

	facade := arwendebug.NewDebugFacade()
	app := initializeCLI(facade)

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(ErrCodeCriticalError)
	}

	os.Exit(ErrCodeSuccess)
}
