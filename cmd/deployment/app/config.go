package app

// Sort packages base on stlibs, remote, local.
import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"awesomeProject/internal/config"
)

var ConfigFile string
var DeploymentServiceConfig deploymentServiceConfig

type deploymentServiceConfig struct {
	Debug bool

	// Import config for the internal package config to share code as much as possible.
	Application config.Application
	Database    config.Database
	GRPC        config.GRPC
	Queue       config.Queue
}

// SetDefaults sets default configuration to be production ready.
// Never set development config as default to reduce security mis configurations.
func setDefaults() {
	// Explicit set debug mode to false.
	DeploymentServiceConfig.Debug = false

	// By default, listen on all interfaces on port 8080.
	DeploymentServiceConfig.Application.Address = "0.0.0.0"
	DeploymentServiceConfig.Application.Port = 8080

	DeploymentServiceConfig.Database.Name = "/tmp/deployment.sql"
}

func loadConfig() {
	// Check if --config-file flag is parsed otherwise use the default filename.
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		// Search for a config file name called deployment.yaml in the current or home directory.
		viper.SetConfigName("deployment")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/")
	}

	// Read config file, if error continue with default configuration.
	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Msg("Cannot load configuration file, using defaults.")
	}

	// Parse config file, always use the pointer otherwise you will get an empty copy of the config struct.
	if err := viper.Unmarshal(&DeploymentServiceConfig); err != nil {
		log.Err(err).Msg("Cannot parse configuration file")
		os.Exit(1)
	}
}
