package app

import (
	"github.com/spf13/cobra"
)

// Set defaults, load config en setup CLI commands and flags before anything else.
func init() {
	cobra.OnInitialize(setDefaults)
	cobra.OnInitialize(loadConfig)

	Command.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is $HOME/deployment.yaml)")

	Command.AddCommand(serviceCommand)
	Command.AddCommand(migrationCommand)
	Command.AddCommand(seedCommand)
}
