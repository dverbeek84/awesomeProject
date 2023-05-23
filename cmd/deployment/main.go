package main

// Sort packages based on stlibs, remote, local.
import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"awesomeProject/cmd/deployment/app"
)

func main() {
	// Setup Logger to default to info.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Write logs to stderr.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Main entry point of the application.
	if err := app.Command.Execute(); err != nil {
		os.Exit(1)
	}
}
