package app

// Sort packages based on stlibs, remote, local.
import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"awesomeProject/internal/database"
)

// Command is the root command for the Order service.
var Command = &cobra.Command{
	Use: "order",
}

// serviceCommand is used to start the Order service.
var serviceCommand = &cobra.Command{
	Use:   "run",
	Short: "Run Order service",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if debug is enabled before anything else otherwise you won't see everything in debug mode.
		if OrderServiceConfig.Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Warn().Msg("Debug enabled, do not use this in production!")
		}

		// Connect to database, so we can use the pointer in the interal database package
		db, err := database.Connect(OrderServiceConfig.Database.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot connect to database")
		}

		// Close connection after return statement from database.Connect().
		defer db.Close()

		// Stare GRPC server in a goroutine, so it's not blocking.
		go func() {
			startGRPCServer()
		}()

		startRESTServer()
	},
}

// migrationCommand does the database migration.
// In production, you should use another solution like https://github.com/bytebase/bytebase
var migrationCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Connect(OrderServiceConfig.Database.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot connect to database")
		}

		// Close connection after return statement from database.Connect().
		defer db.Close()

		// Migrate database, off course you need another solution in production.
		// GORM auto migrate won´t remove used columns. This is safe but not for data integrity.
		log.Info().Msg("Migrating Order database")
		err = database.MigrateOrderScheme()
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot migrate database")
		}
	},
}

// seedCommand seeds the database with dummy data.
// Don´t use this in production this is only for demonstration purpose.
var seedCommand = &cobra.Command{
	Use:   "seed",
	Short: "Run seed database",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.Connect(OrderServiceConfig.Database.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot connect to database")
		}

		// Close connection after return statement from database.Connect().
		defer db.Close()

		// Seed database with dummy data.
		log.Info().Msg("Seeding database")
		err = database.DummyOrderData()
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot seed database, database already seeded")
		}
	},
}
