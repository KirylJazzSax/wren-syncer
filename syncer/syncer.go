package syncer

import (
	"context"

	"wren-time-syncer/repository"
)

type Syncer interface {
	SyncIssues(ctx context.Context, tasks []repository.Issue, force bool) error
}
