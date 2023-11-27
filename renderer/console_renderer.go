package renderer

import (
	"context"
	"fmt"
	"io"
	"strings"

	"wren-time-syncer/prompt"
	"wren-time-syncer/repository"
	"wren-time-syncer/syncer"
	"wren-time-syncer/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/samber/do"
)

type ConsoleWriter struct {
	tableWriter   table.Writer
	selectPrompt  prompt.SelectPromptRunner
	confirmPrompt prompt.ConfirmPromptRunner
}

func (t *ConsoleWriter) Render(ctx context.Context, issues []repository.Issue) error {
	return t.runPrompt(ctx, issues)
}

// Here I pass io.Writer and not a table.Writer because it gives us more flexibility. (I think so)
// table.Writer is sort of a helper, which I can replace with a custom writer, which can be anything that helps render issues beautifully.
func NewConsoleWriter(writer io.Writer) *ConsoleWriter {
	return &ConsoleWriter{
		tableWriter:   newTable(writer),
		confirmPrompt: prompt.NewConfirmPrompt(),
		selectPrompt:  prompt.NewSelectPrompt(),
	}
}

func preparePromptItems(issues []repository.Issue) []string {
	hasMightSynced := false
	items := make([]string, 0)
	var sb strings.Builder

	for i, task := range issues {
		syncJiraError := issues[i].HasSyncError
		isPossiblySynced := issues[i].PossiblySynced
		isSynced := task.IsSynced

		sb.WriteString(fmt.Sprintf("%s ", task.GetKey()))

		if isSynced {
			sb.WriteString(fmt.Sprintf("- %s", utils.SyncStatusSynced))
		}

		if isPossiblySynced && !isSynced && !syncJiraError {
			hasMightSynced = true
			sb.WriteString("- task could be synced")
		}

		if syncJiraError {
			sb.WriteString("- had error while sync")
		}

		items = append(items, sb.String())
		sb.Reset()
	}

	items = append(items, utils.SyncAllMessage)
	if hasMightSynced {
		items = append(items, utils.SyncAllForceMessage)
	}
	items = append(items, utils.ExitMessage)
	return items
}

func (t *ConsoleWriter) runPrompt(ctx context.Context, issues []repository.Issue) error {

	sync, err := do.Invoke[syncer.Syncer](nil)
	if err != nil {
		return err
	}

	t.confirmPrompt.SetLabel("You sure")

	for {
		t.render(issues)

		items := preparePromptItems(issues)
		t.selectPrompt.SetItems(items)

		i, r, err := t.selectPrompt.Run()
		if err != nil {
			utils.LogError(err.Error())
			return err
		}

		if r == utils.ExitMessage {
			fmt.Println("Bye")
			return nil
		}

		if r == utils.SyncAllMessage || r == utils.SyncOnlyNotSyncedMessage {
			if err := sync.SyncIssues(ctx, issues, false); err != nil {
				return err
			}
			break
		}

		if r == utils.SyncAllForceMessage {
			s, err := t.confirmPrompt.Run()
			if err != nil {
				return err
			}

			if strings.EqualFold(s, "y") {
				if err := sync.SyncIssues(ctx, issues, true); err != nil {
					return err
				}
			}
			break
		}

		if issues[i].PossiblySynced {
			s, err := t.confirmPrompt.Run()
			if err != nil {
				return err
			}

			if strings.EqualFold(s, "y") {
				if err := sync.SyncIssues(ctx, issues[i:i+1], true); err != nil {
					utils.LogError(err.Error())
				}
			}
			continue
		}

		if err := sync.SyncIssues(ctx, issues[i:i+1], false); err != nil {
			utils.LogError(err.Error())
		}
	}

	t.render(issues)
	out, err := do.Invoke[io.Writer](nil)
	if err != nil {
		utils.LogError(err.Error())
		return err
	}

	_, err = out.Write([]byte("Bye\n"))
	return err
}

func (t *ConsoleWriter) render(issues []repository.Issue) {
	t.tableWriter.ResetRows()
	t.tableWriter.ResetHeaders()
	header := table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}
	t.tableWriter.AppendHeader(header)

	rows := make([]table.Row, len(issues))
	for i, t := range issues {
		syncStatus := utils.SyncStatusGoodToGo

		if issues[i].PossiblySynced {
			syncStatus = utils.SyncStatusPossiblySynced
		}

		if t.IsSynced {
			syncStatus = utils.SyncStatusSynced
		}

		if issues[i].HasFetchError {
			syncStatus = "Fetch error occured"
		}

		if issues[i].HasSyncError {
			syncStatus = "Sync error occured"
		}

		rows[i] = table.Row{
			t.GetKey(),
			t.GetComment(),
			t.From.Format("2006-01-02"),
			t.From.Format("15:04"),
			t.To.Format("15:04"),
			fmt.Sprintf("%dm", t.Minutes),
			syncStatus,
		}
	}

	t.tableWriter.AppendRows(rows)
	t.tableWriter.Render()
}

func newTable(writer io.Writer) table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(writer)
	tbl.SetStyle(table.StyleColoredMagentaWhiteOnBlack)
	tbl.SetRowPainter(func(row table.Row) text.Colors {
		if row[len(row)-1] == utils.SyncStatusPossiblySynced {
			return text.Colors{text.BgHiMagenta, text.FgHiBlack}
		}

		if row[len(row)-1] == utils.SyncStatusSynced {
			return text.Colors{text.BgHiGreen, text.FgHiWhite}
		}

		if row[len(row)-1] == utils.SyncStatusGoodToGo {
			return text.Colors{text.BgBlack, text.FgHiWhite}
		}

		return text.Colors{text.BgRed, text.FgHiWhite}
	})
	return tbl
}

// It is for teststing only; w io.WriteClosernot not sure if it good idea to create it here.
func newConsoleWriterWithMocks(t table.Writer, sp prompt.SelectPromptRunner, cp prompt.ConfirmPromptRunner) *ConsoleWriter {
	return &ConsoleWriter{
		tableWriter:   t,
		confirmPrompt: cp,
		selectPrompt:  sp,
	}
}
