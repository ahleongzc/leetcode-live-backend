//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ahleongzc/leetcode-live-backend/cmd/app"
	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"

	"github.com/google/wire"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// Handler
		handler.NewAuthHandler,
		handler.NewHealthHandler,
		handler.NewInterviewHandler,

		// Service
		service.NewAuthService,
		service.NewInterviewService,

		// Scenario
		scenario.NewAuthScenario,
		scenario.NewTranscriptManager,
		scenario.NewIntentClassifier,

		// Repo
		repo.NewSessionRepo,
		repo.NewUserRepo,
		repo.NewInterviewRepo,
		repo.NewTranscriptRepo,
		repo.NewFileRepo,

		// Infra
		infra.NewTTS,
		infra.NewLLM,
		infra.NewPostgresDatabase,
		infra.NewZerologLogger,
		infra.NewCloudflareR2ObjectStorageClient,

		// Config
		config.LoadDatabaseConfig,
		config.LoadObjectStorageConfig,
		config.LoadTTSConfig,
		config.LoadWebsocketConfig,

		// Middleware
		middleware.NewMiddleware,

		// Housekeeping
		background.NewHouseKeeper,

		// Application
		app.NewApplication,
	)
	return &app.Application{}, nil
}
