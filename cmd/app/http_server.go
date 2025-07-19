package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	httphandler "github.com/ahleongzc/leetcode-live-backend/internal/handler/http_handler"
	httpmiddleware "github.com/ahleongzc/leetcode-live-backend/internal/handler/http_handler/middleware"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	srv        *http.Server
	logger     *zerolog.Logger
	middleware *httpmiddleware.Middleware

	authHandler      *httphandler.AuthHandler
	userHandler      *httphandler.UserHandler
	healthHandler    *httphandler.HealthHandler
	interviewHandler *httphandler.InterviewHandler
}

func NewHTTPServer(
	logger *zerolog.Logger,
	middleware *httpmiddleware.Middleware,

	authHandler *httphandler.AuthHandler,
	userHandler *httphandler.UserHandler,
	healthHandler *httphandler.HealthHandler,
	interviewHandler *httphandler.InterviewHandler,
) (*HTTPServer, error) {
	httpServerConfig := config.LoadHTTPServerConfig()

	srv := &http.Server{
		Addr:         httpServerConfig.Address,
		IdleTimeout:  httpServerConfig.IdleTimeout,
		ReadTimeout:  httpServerConfig.ReadTimeout,
		WriteTimeout: httpServerConfig.WriteTimeout,
	}

	return &HTTPServer{
		srv:              srv,
		logger:           logger,
		middleware:       middleware,
		authHandler:      authHandler,
		userHandler:      userHandler,
		healthHandler:    healthHandler,
		interviewHandler: interviewHandler,
	}, nil
}

func (hs *HTTPServer) Serve(errChan chan error) *HTTPServer {
	hs.registerHandlers()

	go func() {
		err := hs.srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	hs.logger.Info().Msg(fmt.Sprintf("http server has started at %s", time.Now().Format("2006-01-02 15:04:05")))
	return hs
}

func (hs *HTTPServer) GracefullyTerminate(ctx context.Context) {
	if hs == nil || hs.srv == nil {
		return
	}
	hs.srv.Shutdown(ctx)
	hs.logger.Info().Msg(fmt.Sprintf("http server has gracefully terminated at %s", time.Now().Format("2006-01-02 15:04:05")))
}

func (hs *HTTPServer) registerHandlers() {
	if hs == nil || hs.srv == nil {
		return
	}
	hs.srv.Handler = hs.setUpHandlers()
}

func (hs *HTTPServer) setUpHandlers() http.Handler {
	mux := http.NewServeMux()

	// --- These routes are public and don't require authentication (if they are needed, they are handled at the service layer)
	mux.HandleFunc("GET /v1/health", hs.healthHandler.HealthCheck)

	mux.HandleFunc("POST /v1/auth/login", hs.authHandler.Login)
	mux.HandleFunc("POST /v1/user/register", hs.userHandler.Register)

	mux.HandleFunc("GET /v1/interview/join", hs.interviewHandler.JoinInterview)
	// ---

	// --- These routes require X-Session-Token to be in the headers
	protected := alice.New(hs.middleware.Authenticate, hs.middleware.SetUserID, hs.middleware.SetSessionTokenInResponseHeader)
	mux.Handle("GET /v1/user", protected.ThenFunc(hs.userHandler.GetUserProfile))
	mux.Handle("POST /v1/auth/logout", protected.ThenFunc(hs.authHandler.Logout))

	mux.Handle("POST /v1/interview/set-up-new", protected.ThenFunc(hs.interviewHandler.SetUpNewInterview))
	mux.Handle("POST /v1/interview/set-up-unfinished", protected.ThenFunc(hs.interviewHandler.SetUpUnfinishedInterview))
	mux.Handle("POST /v1/interview/abandon-unfinished", protected.ThenFunc(hs.interviewHandler.AbandonUnfinishedInterview))

	mux.Handle("GET /v1/interview/ongoing", protected.ThenFunc(hs.interviewHandler.GetOngoingInterview))
	mux.Handle("GET /v1/interview/history", protected.ThenFunc(hs.interviewHandler.GetInterviewHistory))
	mux.Handle("GET /v1/interview/unfinished", protected.ThenFunc(hs.interviewHandler.GetUnfinishedInterview))
	// ---

	return alice.New(
		hs.middleware.RecoverPanic,
		hs.middleware.CORS,
		hs.middleware.RecordRequestTimestampMS,
		hs.middleware.Log,
	).Then(mux)
}
