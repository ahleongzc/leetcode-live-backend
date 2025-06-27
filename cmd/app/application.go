package app

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
)

type Application struct {
	authHandler   *handler.AuthHandler
	healthHandler *handler.HealthHandler
}

func NewApplication(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
) *Application {
	return &Application{
		authHandler:   authHandler,
		healthHandler: healthHandler,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)

	// Auth
	mux.HandleFunc("POST /v1/login", a.authHandler.Login)

	return mux
}
