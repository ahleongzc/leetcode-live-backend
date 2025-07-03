package background

import (
	"context"
	"fmt"
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
		for {
			select {
			case t := <-ticker.C:
				h.logger.Info().Msg(
					fmt.Sprintf("housekeeping at %s", t.Format("2006-01-02 15:04:05")),
				)
				h.deleteExpiredSession(ctx)
			case <-ctx.Done():
				h.logger.Log().Msg("gracefully terminating housekeeping")
				ticker.Stop()
				return
			}
		}
	}()
}

func (h *HousekeeperImpl) deleteExpiredSession(ctx context.Context) {
	err := h.sessionRepo.DeleteExpired(ctx)
	if err != nil {
		h.logger.Log().Err(err)
	}
}
