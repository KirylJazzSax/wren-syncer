package syncer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	mockprogress "wren-time-syncer/mocks/progress"
	mockrepository "wren-time-syncer/mocks/repository"
	"wren-time-syncer/repository"
	"wren-time-syncer/utils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var MockedResponseErr = errors.New("worklog hasn't been added.")

func TestSyncIssue(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	config := &utils.Config{}
	repo := mockrepository.NewMockWorklogRepository(ctrl)
	prog := mockprogress.NewMockWriter(ctrl)

	jiraSyncer := NewJiraSyncer(repo, config, prog)

	t.Run("Synced successfully", func(t *testing.T) {
		issues := makeIssues(2, false, false)

		configureRepository(repo, ctx, &issues[1], nil)
		configureProgress(prog, 1)

		err := jiraSyncer.SyncIssues(ctx, issues[1:], false)
		assert.NoError(t, err)
		assert.False(t, issues[0].IsSynced)
		assert.False(t, issues[0].HasSyncError)

		assert.True(t, issues[1].IsSynced)
		assert.False(t, issues[1].HasSyncError)
	})

	t.Run("Not synced", func(t *testing.T) {
		issues := makeIssues(2, false, false)

		configureRepository(repo, ctx, &issues[0], MockedResponseErr)
		configureProgress(prog, 1)

		err := jiraSyncer.SyncIssues(ctx, issues[:1], false)
		assert.NoError(t, err)
		assert.False(t, issues[0].IsSynced)
		assert.True(t, issues[0].HasSyncError)

		assert.False(t, issues[1].IsSynced)
		assert.False(t, issues[1].HasSyncError)
	})
}

func TestSyncIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	config := &utils.Config{}
	ctx := context.TODO()

	repo := mockrepository.NewMockWorklogRepository(ctrl)
	prog := mockprogress.NewMockWriter(ctrl)

	s := NewJiraSyncer(repo, config, prog)

	t.Run("Force sync", func(t *testing.T) {
		possiblySyncedIssues := makeIssues(2, false, true)

		configureProgress(prog, 2)

		for _, issue := range possiblySyncedIssues {
			configureRepository(repo, ctx, &issue, nil)
		}

		err := s.SyncIssues(ctx, possiblySyncedIssues, true)
		assert.NoError(t, err)
		for _, issue := range possiblySyncedIssues {
			assert.True(t, issue.IsSynced)
			assert.False(t, issue.HasSyncError)
		}
	})

	t.Run("Force sync already synced", func(t *testing.T) {
		alreadySynced := makeIssues(3, true, true)

		configureProgress(prog, 3)

		err := s.SyncIssues(ctx, alreadySynced, true)
		assert.NoError(t, err)
		for _, issue := range alreadySynced {
			assert.True(t, issue.IsSynced)
			assert.False(t, issue.HasSyncError)
		}
	})

	t.Run("Issues to sync", func(t *testing.T) {
		possiblySyncedIssues := makeIssues(4, false, false)

		configureProgress(prog, 4)

		for _, issue := range possiblySyncedIssues {
			configureRepository(repo, ctx, &issue, nil)
		}

		err := s.SyncIssues(ctx, possiblySyncedIssues, true)
		assert.NoError(t, err)
		for _, issue := range possiblySyncedIssues {
			assert.True(t, issue.IsSynced)
			assert.False(t, issue.HasSyncError)
		}
	})
}

func makeIssues(length int, synced bool, possiblySynced bool) []repository.Issue {
	issues := make([]repository.Issue, length)

	for i := range issues {
		add := time.Duration(rand.Intn(100))
		issues[i].Comment = fmt.Sprintf("TRADE-%d:Dev:Com", i)
		issues[i].Minutes = rand.Uint64()
		issues[i].From = time.Now().Add(-add * time.Minute)
		issues[i].IsSynced = synced
		issues[i].PossiblySynced = possiblySynced
	}

	return issues
}

func configureRepository(repo *mockrepository.MockWorklogRepository, ctx context.Context, issue *repository.Issue, err error) {
	repo.EXPECT().Add(ctx, issue.GetKey(), &repository.WorklogDTO{
		TimeSpent: issue.Minutes,
		Started:   issue.From,
		Comment:   issue.GetComment(),
	}).Return(err)
}

func configureProgress(prog *mockprogress.MockWriter, issuesLength int) {
	prog.EXPECT().Start(int64(issuesLength), "Sync tasks").Times(1)
	prog.EXPECT().Inc(int64(1)).Times(issuesLength)
	prog.EXPECT().Stop().Times(1)
}
