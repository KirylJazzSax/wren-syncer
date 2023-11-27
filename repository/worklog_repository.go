package repository

import (
	"context"
	"strconv"

	"github.com/andygrunwald/go-jira"
)

type JiraWorklogRepository struct {
	jiraClient *jira.Client
}

func (r *JiraWorklogRepository) Add(ctx context.Context, taskKey string, worklogDto *WorklogDTO) error {
	started := jira.Time(worklogDto.Started)
	worklog := &jira.WorklogRecord{
		TimeSpent: strconv.FormatUint(worklogDto.TimeSpent, 10) + "m",
		Comment:   worklogDto.Comment,
		Started:   &started,
	}
	if _, _, err := r.jiraClient.Issue.AddWorklogRecordWithContext(ctx, taskKey, worklog); err != nil {
		return err
	}

	return nil
}

func NewWorklogRepository(client *jira.Client) *JiraWorklogRepository {
	return &JiraWorklogRepository{
		jiraClient: client,
	}
}
