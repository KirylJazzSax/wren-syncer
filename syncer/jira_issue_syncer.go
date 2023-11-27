package syncer

import (
	"context"
	"sync"

	pw "wren-time-syncer/progress"
	"wren-time-syncer/repository"
	"wren-time-syncer/utils"
)

type JiraSyncer struct {
	waitGroup      *sync.WaitGroup
	config         *utils.Config
	progressWriter pw.Writer
	repo           repository.WorklogRepository
}

func (s *JiraSyncer) syncIssue(ctx context.Context, i *repository.Issue) error {
	worklogDto := &repository.WorklogDTO{
		TimeSpent: i.Minutes,
		Comment:   i.GetComment(),
		Started:   i.From,
	}

	if err := s.repo.Add(ctx, i.GetKey(), worklogDto); err != nil {
		i.IsSynced = false
		i.HasSyncError = true
		return err
	}

	i.HasSyncError = false
	i.IsSynced = true
	return nil
}

func (s *JiraSyncer) SyncIssues(ctx context.Context, issues []repository.Issue, force bool) error {
	if err := s.progressWriter.Start(int64(len(issues)), "Sync tasks"); err != nil {
		return err
	}

	s.waitGroup.Add(len(issues))
	defer s.waitGroup.Wait()

	for i := range issues {
		go func(idx int) {
			defer s.waitGroup.Done()
			if (force || (!force && !issues[idx].PossiblySynced)) && !issues[idx].IsSynced {
				if err := s.syncIssue(ctx, &issues[idx]); err != nil {
					utils.LogError(err.Error())
					issues[idx].IsSynced = false
					issues[idx].HasSyncError = true
				} else {
					issues[idx].IsSynced = true
					issues[idx].HasSyncError = false
				}
			}

			if err := s.progressWriter.Inc(1); err != nil {
				utils.LogError(err.Error())
			}
		}(i)
	}
	return nil
}

func NewJiraSyncer(repo repository.WorklogRepository, config *utils.Config, w pw.Writer) *JiraSyncer {
	return &JiraSyncer{
		waitGroup:      &sync.WaitGroup{},
		config:         config,
		progressWriter: w,
		repo:           repo,
	}
}
