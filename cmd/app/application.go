package app

import (
	"net/http"

	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
)

type Application struct {
	authHandler *handler.AuthHandler
}

func NewApplication(
	authHandler *handler.AuthHandler,
) *Application {
	return &Application{
		authHandler: authHandler,
	}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/login", a.authHandler.Login)

	return mux
}
