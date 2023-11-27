package cmd

import (
	"context"
	"fmt"
	"time"

	"wren-time-syncer/renderer"
	"wren-time-syncer/repository"
	"wren-time-syncer/utils"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func newHandler(cmdName string, ctx context.Context) (func(cmd *cobra.Command, args []string), error) {
	switch cmdName {
	case "root":
		return emptyHandler(), nil
	case "sync":
		return newSyncCommand(ctx), nil
	}

	return nil, fmt.Errorf("not supported cmdName %s", cmdName)
}

func emptyHandler() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {}
}

func newSyncCommand(ctx context.Context) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		dateString := utils.Must(cmd.Flags().GetString("date"))
		repo := do.MustInvoke[repository.IssueRepository](nil)

		dayToSync := time.Now()
		if dateString != "" {
			dayToSync = utils.Must(time.Parse("2006-01-02", dateString))
		}

		issues := utils.Must(repo.GetIssues(ctx, dayToSync))
		writer := do.MustInvoke[renderer.Writer](nil)
		writer.Render(ctx, issues)
	}
}
