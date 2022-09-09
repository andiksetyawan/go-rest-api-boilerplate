package commands

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "go-rest-api-boilerplate",
		Short: "Run service",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}
	command.AddCommand(serverCmd, NewMigrateCmd())
	return command
}
