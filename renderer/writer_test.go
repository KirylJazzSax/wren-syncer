package renderer

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	mocktable "wren-time-syncer/mocks/go-pretty/table"
	mockprompt "wren-time-syncer/mocks/prompt"
	mocksyncer "wren-time-syncer/mocks/syncer"
	"wren-time-syncer/repository"
	"wren-time-syncer/syncer"

	"github.com/golang/mock/gomock"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()

	tab := mocktable.NewMockWriter(ctrl)
	sp := mockprompt.NewMockSelectPromptRunner(ctrl)
	cp := mockprompt.NewMockConfirmPromptRunner(ctrl)
	s := mocksyncer.NewMockSyncer(ctrl)

	do.Provide(nil, func(i *do.Injector) (syncer.Syncer, error) {
		return s, nil
	})

	do.Provide(nil, func(i *do.Injector) (io.Writer, error) {
		return os.Stdout, nil
	})

	t.Run("Basic run", func(t *testing.T) {
		tab.EXPECT().ResetRows().Times(2)
		tab.EXPECT().ResetHeaders().Times(2)
		tab.EXPECT().AppendRows(gomock.Any()).Times(2)
		tab.EXPECT().Render().Times(2)
		tab.EXPECT().AppendHeader(table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}).Times(2)

		cp.EXPECT().SetLabel("You sure").Times(1)
		sp.EXPECT().Run().Return(2, "Sync All", nil).Times(1)
		issues := makeIssues(false)

		items := preparePromptItems(issues)
		sp.EXPECT().SetItems(items).Times(2)

		s.EXPECT().SyncIssues(ctx, issues, false).Times(1)

		wr := newConsoleWriterWithMocks(tab, sp, cp)

		err := wr.Render(ctx, issues)
		assert.NoError(t, err)
	})

	t.Run("Sync 2 issues and exit", func(t *testing.T) {
		tab.EXPECT().ResetRows().Times(3)
		tab.EXPECT().ResetHeaders().Times(3)
		tab.EXPECT().AppendRows(gomock.Any()).Times(3)
		tab.EXPECT().Render().Times(3)
		tab.EXPECT().AppendHeader(table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}).Times(3)

		cp.EXPECT().SetLabel("You sure").Times(1)
		sp.EXPECT().Run().Return(1, "TRADE-2", nil).Times(1)
		sp.EXPECT().Run().Return(0, "TRADE-1", nil).Times(1)
		sp.EXPECT().Run().Return(4, "Exit", nil).Times(1)

		issues := makeIssues(false)

		items := preparePromptItems(issues)
		sp.EXPECT().SetItems(items).Times(3)

		s.EXPECT().SyncIssues(ctx, issues[1:2], false).Times(1)
		s.EXPECT().SyncIssues(ctx, issues[0:1], false).Times(1)

		wr := newConsoleWriterWithMocks(tab, sp, cp)

		err := wr.Render(ctx, issues)
		assert.NoError(t, err)
	})

	t.Run("Sync forced", func(t *testing.T) {
		tab.EXPECT().ResetRows().Times(2)
		tab.EXPECT().ResetHeaders().Times(2)
		tab.EXPECT().AppendRows(gomock.Any()).Times(2)
		tab.EXPECT().Render().Times(2)
		tab.EXPECT().AppendHeader(table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}).Times(2)

		cp.EXPECT().SetLabel("You sure").Times(1)
		cp.EXPECT().Run().Return("y", nil).Times(1)
		sp.EXPECT().Run().Return(4, "Sync All no matter what", nil).Times(1)

		issues := makeIssues(true)

		items := preparePromptItems(issues)
		sp.EXPECT().SetItems(items).Times(2)

		s.EXPECT().SyncIssues(ctx, issues, true).Times(1)

		wr := newConsoleWriterWithMocks(tab, sp, cp)

		err := wr.Render(ctx, issues)
		assert.NoError(t, err)
	})

	t.Run("Do not sync forced", func(t *testing.T) {
		tab.EXPECT().ResetRows().Times(2)
		tab.EXPECT().ResetHeaders().Times(2)
		tab.EXPECT().AppendRows(gomock.Any()).Times(2)
		tab.EXPECT().Render().Times(2)
		tab.EXPECT().AppendHeader(table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}).Times(2)

		cp.EXPECT().SetLabel("You sure").Times(1)
		cp.EXPECT().Run().Return("n", nil).Times(1)
		sp.EXPECT().Run().Return(4, "Sync All no matter what", nil).Times(1)

		issues := makeIssues(true)

		items := preparePromptItems(issues)
		sp.EXPECT().SetItems(items).Times(2)

		wr := newConsoleWriterWithMocks(tab, sp, cp)

		err := wr.Render(ctx, issues)
		assert.NoError(t, err)
	})

}

func makeIssues(ps bool) []repository.Issue {
	return []repository.Issue{
		{
			Comment:        "TRADE-1:Dev:Comment",
			From:           time.Now().Add(-20 * time.Minute),
			To:             time.Now().Add(-10 * time.Minute),
			Minutes:        10,
			PossiblySynced: ps,
		},
		{
			Comment:        "TRADE-2:Dev:Comment",
			From:           time.Now().Add(-50 * time.Minute),
			To:             time.Now().Add(-40 * time.Minute),
			Minutes:        10,
			PossiblySynced: ps,
		},
		{
			Comment:        "TRADE-3:Dev:Comment",
			From:           time.Now().Add(-50 * time.Minute),
			To:             time.Now().Add(-40 * time.Minute),
			Minutes:        10,
			PossiblySynced: ps,
		},
	}
}
