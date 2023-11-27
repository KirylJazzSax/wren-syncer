package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newRootCmd(run func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Short: "wren task syncer",
		Long:  "Sync wren tasks from polcode link app",
		Run:   run,
	}
}

func Execute(ctx context.Context) error {
	rootRun, err := newHandler("root", ctx)
	if err != nil {
		return err
	}

	syncRun, err := newHandler("sync", ctx)
	if err != nil {
		return err
	}

	rootCmd := newRootCmd(rootRun)
	syncCmd := newSyncCmd(syncRun)
	rootCmd.AddCommand(syncCmd)

	return rootCmd.Execute()
}
