package cmd

import (
	"github.com/spf13/cobra"
)

func newSyncCmd(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync wren jira tasks",
		Long:  "will display tasks what we need to sync then you will choose do you need to sync it or not",
		Run:   run,
	}

	cmd.Flags().String("date", "", "Day you want sync tasks in format 2006-01-02")
	return cmd
}
