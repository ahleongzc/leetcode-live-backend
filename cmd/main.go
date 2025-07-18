package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/wire"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app, err := wire.InitializeApplication()
	if err != nil {
		panic(err)
	}

	errChan := make(chan error)
	defer close(errChan)

	waitForTermination(errChan)

	ctx := context.Background()

	httpServer := app.StartHTTPServer(errChan)
	rpcServer := app.StartRPCServer(errChan)

	app.StartHouseKeeping(ctx, config.HOUSEKEEPING_INTERVAL)
	app.StartConsumers(ctx, config.CONSUMER_POOL_SIZE)

	<-errChan

	httpServer.Shutdown(ctx)
	rpcServer.GracefulStop()
}

func waitForTermination(errChan chan error) {
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		sig := <-signalChan

		errChan <- errors.New(sig.String())
	}()
}
