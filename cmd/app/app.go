package app

import (
	"context"
	"sync"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
	handler "github.com/ahleongzc/leetcode-live-backend/internal/http_handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/http_handler/middleware"

	"github.com/rs/zerolog"
)

type Application struct {
	logger           *zerolog.Logger
	authHandler      *handler.AuthHandler
	userHandler      *handler.UserHandler
	healthHandler    *handler.HealthHandler
	interviewHandler *handler.InterviewHandler
	middleware       *middleware.Middleware
	reviewConsumer   *consumer.ReviewConsumer
	wg               *sync.WaitGroup
	housekeeper      background.HouseKeeper
	workerPool       background.WorkerPool
}

func NewApplication(
	logger *zerolog.Logger,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	middleware *middleware.Middleware,
	healthHandler *handler.HealthHandler,
	reviewConsumer *consumer.ReviewConsumer,
	interviewHandler *handler.InterviewHandler,
	housekeeper background.HouseKeeper,
	workerPool background.WorkerPool,
) *Application {
	return &Application{
		logger:           logger,
		authHandler:      authHandler,
		userHandler:      userHandler,
		healthHandler:    healthHandler,
		interviewHandler: interviewHandler,
		middleware:       middleware,
		housekeeper:      housekeeper,
		workerPool:       workerPool,
		reviewConsumer:   reviewConsumer,
		wg:               &sync.WaitGroup{},
	}
}

func (a *Application) StartHouseKeeping(ctx context.Context, interval time.Duration) {
	go a.housekeeper.Housekeep(ctx, interval)
}

func (a *Application) StartConsumers(ctx context.Context, workerCount uint) {
	go a.reviewConsumer.ConsumeAndProcess(ctx, workerCount)
}
