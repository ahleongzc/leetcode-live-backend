package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"

	"github.com/justinas/alice"
)

type Application struct {
	authHandler           *handler.AuthHandler
	userHandler           *handler.UserHandler
	healthHandler         *handler.HealthHandler
	interviewHandler      *handler.InterviewHandler
	middleware            *middleware.Middleware
	housekeeper           background.HouseKeeper
	workerPool            background.WorkerPool
	inMemoryCallbackQueue infra.InMemoryCallbackQueue
}

func NewApplication(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	interviewHandler *handler.InterviewHandler,
	middleware *middleware.Middleware,
	housekeeper background.HouseKeeper,
	workerPool background.WorkerPool,
	inMemoryCallbackQueue infra.InMemoryCallbackQueue,
) *Application {
	return &Application{
		authHandler:           authHandler,
		userHandler:           userHandler,
		healthHandler:         healthHandler,
		interviewHandler:      interviewHandler,
		middleware:            middleware,
		housekeeper:           housekeeper,
		workerPool:            workerPool,
		inMemoryCallbackQueue: inMemoryCallbackQueue,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)

	mux.HandleFunc("POST /v1/auth/login", a.authHandler.Login)
	mux.HandleFunc("POST /v1/auth/status", a.authHandler.GetAuthStatus)
	mux.HandleFunc("POST /v1/auth/logout", a.authHandler.Logout)

	mux.HandleFunc("POST /v1/interview/set-up", a.interviewHandler.SetUpInterview)
	mux.HandleFunc("GET /v1/interview/join", a.interviewHandler.JoinInterview)

	mux.HandleFunc("POST /v1/user/register", a.userHandler.Register)

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

func (a *Application) StartWorkerPool(ctx context.Context, poolSize uint) {
	a.workerPool.Init(poolSize)
	a.workerPool.Start(ctx)
}
