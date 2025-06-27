package app

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"
	"github.com/justinas/alice"
)

type Application struct {
	authHandler      *handler.AuthHandler
	healthHandler    *handler.HealthHandler
	interviewHandler *handler.InterviewHandler
	middleware       *middleware.Middleware
}

func NewApplication(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	interviewHandler *handler.InterviewHandler,
	middleware *middleware.Middleware,
) *Application {
	return &Application{
		authHandler:      authHandler,
		healthHandler:    healthHandler,
		interviewHandler: interviewHandler,
		middleware:       middleware,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)
	mux.HandleFunc("POST /v1/login", a.authHandler.Login)
	mux.HandleFunc("GET /v1/start-interview", a.authHandler.Login)

	return alice.New(
		a.middleware.RecoverPanic,
		a.middleware.CORS,
		a.middleware.RecordRequestTimestampMS,
		a.middleware.Log,
	).Then(mux)
}
