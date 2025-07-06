package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"

	"github.com/justinas/alice"
)

type Application struct {
	authHandler      *handler.AuthHandler
	userHandler      *handler.UserHandler
	healthHandler    *handler.HealthHandler
	interviewHandler *handler.InterviewHandler
	middleware       *middleware.Middleware
	reviewConsumer   *consumer.ReviewConsumer
	housekeeper      background.HouseKeeper
	workerPool       background.WorkerPool
}

func NewApplication(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	interviewHandler *handler.InterviewHandler,
	middleware *middleware.Middleware,
	reviewConsumer *consumer.ReviewConsumer,
	housekeeper background.HouseKeeper,
	workerPool background.WorkerPool,
) *Application {
	return &Application{
		authHandler:      authHandler,
		userHandler:      userHandler,
		healthHandler:    healthHandler,
		interviewHandler: interviewHandler,
		middleware:       middleware,
		housekeeper:      housekeeper,
		workerPool:       workerPool,
		reviewConsumer:   reviewConsumer,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)

	mux.HandleFunc("POST /v1/auth/login", a.authHandler.Login)
	mux.HandleFunc("POST /v1/user/register", a.userHandler.Register)

	mux.HandleFunc("GET /v1/interview/join", a.interviewHandler.JoinInterview)

	// These routes require X-Session-Token to be in the headers
	protected := alice.New(a.middleware.Authenticate, a.middleware.SetUserID, a.middleware.SetSessionTokenInResponseHeader)

	mux.Handle("POST /v1/auth/status", protected.ThenFunc(a.authHandler.GetStatus))
	mux.Handle("POST /v1/auth/logout", protected.ThenFunc(a.authHandler.Logout))

	mux.Handle("POST /v1/interview/set-up", protected.ThenFunc(a.interviewHandler.SetUpInterview))

	return alice.New(
		a.middleware.RecoverPanic,
		a.middleware.CORS,
		a.middleware.RecordRequestTimestampMS,
		a.middleware.Log,
	).Then(mux)
}

func (a *Application) StartHouseKeeping(ctx context.Context, interval time.Duration) {
	a.housekeeper.Housekeep(ctx, interval)
}

func (a *Application) StartConsumers(ctx context.Context, workerCount uint) {
	a.reviewConsumer.ConsumeAndProcess(ctx, workerCount)
}
