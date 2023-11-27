package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	pw "wren-time-syncer/progress"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
)

type Issue struct {
	CategoryId     uint      `json:"categoryId"`
	CategoryName   string    `json:"categoryName"`
	From           time.Time `json:"fromDate"`
	To             time.Time `json:"toDate"`
	Comment        string    `json:"comment"`
	Minutes        uint64    `json:"minutes"`
	IsSynced       bool      `json:"synced"`
	PossiblySynced bool      `json:"possiblySynced"`
	HasFetchError  bool      `json:"hasFetchError"`
	HasSyncError   bool      `json:"hasSyncError"`
}

type JiraIssueRepository struct {
	waitGroup      *sync.WaitGroup
	client         *http.Client
	jiraClient     *jira.Client
	config         *utils.Config
	progressWriter pw.Writer
}

func (i *Issue) GetKey() string {
	return strings.Split(i.Comment, ":")[0]
}

func (issue *Issue) GetComment() string {
	return strings.Join(strings.Split(issue.Comment, ":")[1:], ":")
}

func (r *JiraIssueRepository) GetIssues(ctx context.Context, d time.Time) ([]Issue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%sapi/time/%s", r.config.LinkHost, d.Format("2006-01-02")), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", r.config.LinkAuthHeader)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.LogError(err.Error())
		}
	}()

	tasks := make([]Issue, 0)
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}

	r.waitGroup.Add(len(tasks))

	if err := r.progressWriter.Start(int64(len(tasks)), "Fetch tasks"); err != nil {
		return nil, err
	}

	for i := range tasks {
		go func(idx int) {
			defer r.waitGroup.Done()
			possiblySynced, err := r.isPossiblySynced(ctx, &tasks[idx], d)

			if err != nil {
				utils.LogError(fmt.Sprintf("%s: %s", tasks[idx].GetKey(), err.Error()))
				tasks[idx].HasFetchError = true
			}

			tasks[idx].PossiblySynced = possiblySynced
			if err := r.progressWriter.Inc(1); err != nil {
				utils.LogError(err.Error())
			}
		}(i)
	}
	r.waitGroup.Wait()

	return tasks, nil
}

func (r *JiraIssueRepository) isPossiblySynced(ctx context.Context, i *Issue, d time.Time) (bool, error) {
	unixStr := strconv.FormatInt(d.UnixMilli(), 10)
	jqlStr := fmt.Sprintf(
		"worklogDate = %s AND worklogAuthor = currentUser() AND issuekey = %s AND worklogComment ~ %s",
		unixStr,
		i.GetKey(),
		fmt.Sprintf("\"%s\"", i.GetComment()),
	)

	_, _, err := r.jiraClient.Issue.GetWithContext(ctx, i.GetKey(), &jira.GetQueryOptions{})
	issues := []jira.Issue{}
	if err == nil {
		issues, _, _ = r.jiraClient.Issue.SearchWithContext(ctx, jqlStr, &jira.SearchOptions{})
	}

	return len(issues) > 0, err
}

func NewJiraIssueRepository(jiraClient *jira.Client, client *http.Client, config *utils.Config, w pw.Writer) *JiraIssueRepository {
	return &JiraIssueRepository{
		waitGroup:      &sync.WaitGroup{},
		client:         client,
		jiraClient:     jiraClient,
		config:         config,
		progressWriter: w,
	}
}
