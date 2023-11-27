package repository

import (
	"context"
	"time"
)

type WorklogDTO struct {
	TimeSpent uint64
	Comment   string
	Started   time.Time
}

type IssueRepository interface {
	GetIssues(ctx context.Context, d time.Time) ([]Issue, error)
}

type WorklogRepository interface {
	Add(ctx context.Context, taskKey string, worklog *WorklogDTO) error
}
