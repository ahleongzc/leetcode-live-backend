package background

import (
	"context"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/repo"

	"github.com/rs/zerolog"
)

type HouseKeeper interface {
	Housekeep(ctx context.Context, interval time.Duration)
}

type HousekeeperImpl struct {
	sessionRepo repo.SessionRepo
	logger      *zerolog.Logger
}

func NewHouseKeeper(
	sessionRepo repo.SessionRepo,
	logger *zerolog.Logger,
) HouseKeeper {
	return &HousekeeperImpl{
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (h *HousekeeperImpl) Housekeep(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx := context.Background()
				h.deleteExpiredSession(ctx)
			case <-ctx.Done():
				h.logger.Log().Msg("gracefully terminating housekeeping")
				return
			}
		}
	}()
}

func (h *HousekeeperImpl) deleteExpiredSession(ctx context.Context) {
	start := time.Now()
	deletedCount, err := h.sessionRepo.DeleteExpired(ctx)
	duration := time.Since(start)
	if err != nil {
		h.logger.Error().
			Err(err).
			Dur("duration", time.Duration(duration.Seconds())).
			Msg("failed to delete expired sessions")
		return
	}

	if deletedCount == 0 {
		return
	}

	h.logger.Info().
		Int("sessionDeleted", int(deletedCount)).
		Dur("duration", time.Duration(duration.Seconds())).
		Msg("deleted expired session successfully")
}
