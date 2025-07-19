package app

import (
	"context"
	"sync"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
)

type Application struct {
	HTTPServer *HTTPServer
	RPCServer  *RPCServer

	reviewConsumer *consumer.ReviewConsumer
	housekeeper    background.HouseKeeper
	workerPool     background.WorkerPool

	wg *sync.WaitGroup
}

func NewApplication(
	httpServer *HTTPServer,
	rpcServer *RPCServer,

	reviewConsumer *consumer.ReviewConsumer,
	housekeeper background.HouseKeeper,
	workerPool background.WorkerPool,
) *Application {
	return &Application{
		HTTPServer: httpServer,
		RPCServer:  rpcServer,

		housekeeper:    housekeeper,
		reviewConsumer: reviewConsumer,
		workerPool:     workerPool,

		wg: &sync.WaitGroup{},
	}
}

func (a *Application) StartHouseKeeping(ctx context.Context, interval time.Duration) {
	go a.housekeeper.Housekeep(ctx, interval)
}

func (a *Application) StartConsumers(ctx context.Context, workerCount uint) {
	go a.reviewConsumer.ConsumeAndProcess(ctx, workerCount)
}
