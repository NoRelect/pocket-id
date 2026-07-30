package job

import (
	"context"
	"log/slog"

	"github.com/go-co-op/gocron/v2"

	"github.com/pocket-id/pocket-id/backend/internal/service"
	"github.com/pocket-id/pocket-id/backend/internal/webauthn/mds"
)

func (s *Scheduler) RegisterMDSUpdateJob(ctx context.Context, mdsService *mds.Service) error {
	job := &mdsUpdateJob{mdsService: mdsService}

	return s.RegisterJob(
		ctx,
		"UpdateFIDOMDS",
		gocron.DurationJob(mds.CacheTTL),
		job.run,
		service.RegisterJobOpts{RunImmediately: true},
	)
}

type mdsUpdateJob struct {
	mdsService *mds.Service
}

func (j *mdsUpdateJob) run(ctx context.Context) error {
	err := j.mdsService.Refresh(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to refresh FIDO MDS cache",
			slog.Any("error", err),
			slog.Time("lastRefreshed", j.mdsService.FetchedAt()),
		)
		return nil
	}
	return nil
}
