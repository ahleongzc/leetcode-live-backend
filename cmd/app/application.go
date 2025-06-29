package app

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"
	"github.com/justinas/alice"
)

type Application struct {
	authHandler      *handler.AuthHandler
	healthHandler    *handler.HealthHandler
	interviewHandler *handler.InterviewHandler
	middleware       *middleware.Middleware
	housekeeper      background.HouseKeeper
}

func NewApplication(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	interviewHandler *handler.InterviewHandler,
	middleware *middleware.Middleware,
	housekeeper background.HouseKeeper,
) *Application {
	return &Application{
		authHandler:      authHandler,
		healthHandler:    healthHandler,
		interviewHandler: interviewHandler,
		middleware:       middleware,
		housekeeper:      housekeeper,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)

	mux.HandleFunc("POST /v1/auth/login", a.authHandler.Login)
	mux.HandleFunc("POST /v1/auth/status", a.authHandler.GetAuthStatus)
	mux.HandleFunc("POST /v1/auth/logout", a.authHandler.Logout)

	mux.HandleFunc("POST /v1/interview/set-up-interview", a.interviewHandler.SetUpInterview)
	mux.HandleFunc("GET /v1/interview/join-interview", a.interviewHandler.JoinInterview)

	return alice.New(
		a.middleware.RecoverPanic,
		a.middleware.CORS,
		a.middleware.RecordRequestTimestampMS,
		a.middleware.Log,
	).Then(mux)
}

// TODO: Add a done channel here for graceful termination
func (a *Application) StartBackgroundTasks() {
	go a.housekeeper.Housekeep(common.HOUSEKEEPING_INTERVAL, nil)
}
