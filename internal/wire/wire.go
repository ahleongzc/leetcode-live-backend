//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ahleongzc/leetcode-live-backend/cmd/app"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"

	"github.com/google/wire"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// Handler
		handler.NewAuthHandler,
		handler.NewHealthHandler,

		// Service
		service.NewAuthService,

		// Repo
		repo.NewSessionRepo,
		repo.NewFileRepo,

		// Infra
		infra.NewTTS,
		infra.NewPostgresDatabase,
		infra.NewCloudflareR2ObjectStorageClient,

		// Config
		config.LoadDatabaseConfig,
		config.LoadObjectStorageConfig,
		config.LoadTTSConfig,

		// Application
		app.NewApplication,
	)
	return &app.Application{}, nil
}
