package repository

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	mockprogress "wren-time-syncer/mocks/progress"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIssueValueGetters(t *testing.T) {
	issue := &Issue{
		Comment: "TRADE-1:Dev:Com1",
	}

	assert.Equal(t, issue.GetKey(), "TRADE-1", "Should return string before first :")
	assert.Equal(t, issue.GetComment(), "Dev:Com1", "Should return string after first :")
}

func TestGetIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()

	t.Run("With possibly synced", func(t *testing.T) {
		issuesLength := 5
		mux := utils.MakeTestServeMux(issuesLength, true, 200)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		prog := mockprogress.NewMockWriter(ctrl)
		configureProgress(prog, issuesLength)

		tClient := ts.Client()
		jiraClient := utils.Must(jira.NewClient(tClient, utils.Host(ts)))
		repo := NewJiraIssueRepository(jiraClient, tClient, newTestConfig(ts), prog)

		issues, err := repo.GetIssues(ctx, time.Now())
		assert.NoError(t, err)

		for _, issue := range issues {
			assert.False(t, issue.HasFetchError)
			assert.True(t, issue.PossiblySynced)
		}
	})

	t.Run("All issues should not be possibly synced", func(t *testing.T) {
		issuesLength := 3
		mux := utils.MakeTestServeMux(issuesLength, false, 200)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		prog := mockprogress.NewMockWriter(ctrl)
		configureProgress(prog, issuesLength)

		tClient := ts.Client()
		jiraClient := utils.Must(jira.NewClient(tClient, utils.Host(ts)))
		repo := NewJiraIssueRepository(jiraClient, tClient, newTestConfig(ts), prog)

		issues, err := repo.GetIssues(ctx, time.Now())
		assert.NoError(t, err)

		for _, issue := range issues {
			assert.False(t, issue.HasFetchError)
			assert.False(t, issue.PossiblySynced)
		}
	})

	t.Run("With jira fetch errors", func(t *testing.T) {
		issuesLength := 4
		mux := utils.MakeTestServeMux(issuesLength, true, 404)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		prog := mockprogress.NewMockWriter(ctrl)
		configureProgress(prog, issuesLength)

		tClient := ts.Client()
		jiraClient := utils.Must(jira.NewClient(tClient, utils.Host(ts)))
		repo := NewJiraIssueRepository(jiraClient, tClient, newTestConfig(ts), prog)

		issues, err := repo.GetIssues(ctx, time.Now())
		assert.NoError(t, err)

		for _, issue := range issues {
			assert.True(t, issue.HasFetchError)
			assert.False(t, issue.PossiblySynced)
		}
	})
}

func TestAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	issuesLength := 1

	prog := mockprogress.NewMockWriter(ctrl)
	configureProgress(prog, issuesLength)

	mux := utils.MakeTestServeMux(issuesLength, true, 200)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	tClient := ts.Client()
	jiraClient := utils.Must(jira.NewClient(tClient, utils.Host(ts)))

	repo := NewWorklogRepository(jiraClient)
	issueRepo := NewJiraIssueRepository(jiraClient, tClient, newTestConfig(ts), prog)
	issues, _ := issueRepo.GetIssues(ctx, time.Now())

	err := repo.Add(ctx, issues[0].GetKey(), &WorklogDTO{
		TimeSpent: issues[0].Minutes,
		Comment:   issues[0].GetComment(),
		Started:   time.Now(),
	})
	assert.NoError(t, err)
}

func configureProgress(prog *mockprogress.MockWriter, issuesLength int) {
	prog.EXPECT().Start(int64(issuesLength), "Fetch tasks").Times(1)
	prog.EXPECT().Inc(int64(1)).Times(issuesLength)
	prog.EXPECT().Stop().Times(1)
}

func newTestConfig(ts *httptest.Server) *utils.Config {
	return &utils.Config{
		RequestTimeout: "1000",
		LinkHost:       utils.Host(ts),
	}
}
