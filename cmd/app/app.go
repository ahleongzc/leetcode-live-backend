package app

import (
	"context"
	"sync"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
	httphandler "github.com/ahleongzc/leetcode-live-backend/internal/http_handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/http_handler/middleware"
	rpchandler "github.com/ahleongzc/leetcode-live-backend/internal/rpc_handler"

	"github.com/rs/zerolog"
)

type Application struct {
	logger *zerolog.Logger

	proxyHandler *rpchandler.ProxyHandler

	authHandler      *httphandler.AuthHandler
	userHandler      *httphandler.UserHandler
	healthHandler    *httphandler.HealthHandler
	interviewHandler *httphandler.InterviewHandler

	middleware *middleware.Middleware

	reviewConsumer *consumer.ReviewConsumer

	wg *sync.WaitGroup

	housekeeper background.HouseKeeper
	workerPool  background.WorkerPool
}

func NewApplication(
	logger *zerolog.Logger,

	proxyHandler *rpchandler.ProxyHandler,

	authHandler *httphandler.AuthHandler,
	userHandler *httphandler.UserHandler,
	healthHandler *httphandler.HealthHandler,
	interviewHandler *httphandler.InterviewHandler,

	middleware *middleware.Middleware,

	reviewConsumer *consumer.ReviewConsumer,
	housekeeper background.HouseKeeper,
	workerPool background.WorkerPool,
) *Application {
	return &Application{
		logger:           logger,
		proxyHandler:     proxyHandler,
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
