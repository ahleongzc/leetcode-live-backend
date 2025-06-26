package app

import (
	"net/http"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Handler() http.Handler {
	mux := http.NewServeMux()
	return mux
}
