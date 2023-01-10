package console

import (
	"github.com/irvankadhafi/go-boilerplate/internal/config"
	"github.com/irvankadhafi/go-boilerplate/internal/db"
	"github.com/irvankadhafi/go-boilerplate/utils"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Long:  `This subcommand used to migrate database`,
	Run:   processMigration,
}

func init() {
	migrateCmd.PersistentFlags().Int("step", 0, "maximum migration steps")
	migrateCmd.PersistentFlags().String("direction", "up", "migration direction")
	RootCmd.AddCommand(migrateCmd)
}

func processMigration(cmd *cobra.Command, args []string) {
	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		log.WithField("stepStr", stepStr).Fatal("Failed to parse step to int: ", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "./db/migration",
	}

	migrate.SetTable("schema_migrations")
	db.InitializePostgresConn()
	sqlDB, err := db.PostgreSQL.DB()
	if err != nil {
		log.WithField("DatabaseDSN", config.DatabaseDSN()).Fatal("Failed to connect database: ", err)
	}

	var n int
	if direction == "down" {
		n, err = migrate.ExecMax(sqlDB, "postgres", migrations, migrate.Down, step)
	} else {
		n, err = migrate.ExecMax(sqlDB, "postgres", migrations, migrate.Up, step)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"migrations": utils.Dump(migrations),
			"direction":  direction}).
			Fatal("Failed to migrate database: ", err)
	}

	log.Infof("Applied %d migrations!\n", n)

}
