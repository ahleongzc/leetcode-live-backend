package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"

	"github.com/justinas/alice"
)

func (a *Application) StartHTTPServer(errChan chan error) *http.Server {
	httpServerConfig := config.LoadHTTPServerConfig()

	httpServer := &http.Server{
		Addr:         httpServerConfig.Address,
		Handler:      a.HTTPHandler(),
		IdleTimeout:  httpServerConfig.IdleTimeout,
		ReadTimeout:  httpServerConfig.ReadTimeout,
		WriteTimeout: httpServerConfig.WriteTimeout,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	a.logger.Info().Msg(fmt.Sprintf("http server has started at %s", time.Now().Format("2006-01-02 15:04:05")))

	return httpServer
}

func (a *Application) HTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// --- These routes are public and don't require authentication (if they are needed, they are handled at the service layer)
	mux.HandleFunc("GET /v1/health", a.healthHandler.HealthCheck)

	mux.HandleFunc("POST /v1/auth/login", a.authHandler.Login)
	mux.HandleFunc("POST /v1/user/register", a.userHandler.Register)

	mux.HandleFunc("GET /v1/interview/join", a.interviewHandler.JoinInterview)
	// ---

	// --- These routes require X-Session-Token to be in the headers
	protected := alice.New(a.middleware.Authenticate, a.middleware.SetUserID, a.middleware.SetSessionTokenInResponseHeader)
	mux.Handle("GET /v1/user", protected.ThenFunc(a.userHandler.GetUserProfile))
	mux.Handle("POST /v1/auth/logout", protected.ThenFunc(a.authHandler.Logout))

	mux.Handle("POST /v1/interview/set-up-new", protected.ThenFunc(a.interviewHandler.SetUpNewInterview))
	mux.Handle("POST /v1/interview/set-up-unfinished", protected.ThenFunc(a.interviewHandler.SetUpUnfinishedInterview))
	mux.Handle("POST /v1/interview/abandon-unfinished", protected.ThenFunc(a.interviewHandler.AbandonUnfinishedInterview))

	mux.Handle("GET /v1/interview/ongoing", protected.ThenFunc(a.interviewHandler.GetOngoingInterview))
	mux.Handle("GET /v1/interview/history", protected.ThenFunc(a.interviewHandler.GetInterviewHistory))
	mux.Handle("GET /v1/interview/unfinished", protected.ThenFunc(a.interviewHandler.GetUnfinishedInterview))
	// ---

	return alice.New(
		a.middleware.RecoverPanic,
		a.middleware.CORS,
		a.middleware.RecordRequestTimestampMS,
		a.middleware.Log,
	).Then(mux)
}
