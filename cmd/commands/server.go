package commands

import (
	"github.com/spf13/cobra"
	"go-rest-api-boilerplate/config"
	"go-rest-api-boilerplate/internal/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the API server",
	Long:  "Run the API server",
	Run: func(c *cobra.Command, args []string) {
		config.Init()
		server.NewServer().Run()
	},
}
