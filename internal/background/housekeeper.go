package background

import (
	"context"
	"fmt"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/rs/zerolog"
)

type HouseKeeper interface {
	Housekeep(interval time.Duration, done chan bool)
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

func (h *HousekeeperImpl) Housekeep(interval time.Duration, done chan bool) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {

		case t := <-ticker.C:
			h.logger.Info().Msg(
				fmt.Sprintf("housekeeping at %s", t.Format("2006-01-02 15:04:05")),
			)

			ctx := context.Background()

			err := h.sessionRepo.DeleteExpired(ctx)
			if err != nil {
				h.logger.Log().Err(err)
			}

		case <-done:
			h.logger.Log().Msg("gracefully terminating housekeeping")
			return
		}
	}
}
