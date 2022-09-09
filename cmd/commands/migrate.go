package commands

import (
	"github.com/spf13/cobra"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/db"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run the Up db migration",
	Long:  "Run the Up db migration",
	Run: func(c *cobra.Command, args []string) {
		config.Init()
		pg := db.NewPostgreeDb(config.App.DbHost, config.App.DbPort, config.App.DbName, config.App.DbUser, config.App.DbPass)
		pg.Connect()
		pg.MigrateUp()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Run the down db migration",
	Long:  "Run the down db migration",
	Run: func(c *cobra.Command, args []string) {
		config.Init()
		pg := db.NewPostgreeDb(config.App.DbHost, config.App.DbPort, config.App.DbName, config.App.DbUser, config.App.DbPass)
		pg.Connect()
		pg.MigrateDown()
	},
}

func NewMigrateCmd() *cobra.Command {
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run the db migrations",
		Long:  "Run the db migrations",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)
	return migrateCmd
}
