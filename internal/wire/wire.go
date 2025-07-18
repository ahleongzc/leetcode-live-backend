//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ahleongzc/leetcode-live-backend/cmd/app"
	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
	httpHandler "github.com/ahleongzc/leetcode-live-backend/internal/http_handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/http_handler/middleware"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/cloudflare"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/fasttext"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/http"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/postgres"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo/zerolog"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"

	"github.com/google/wire"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// Consumer
		consumer.NewReviewConsumer,

		// Handler
		httpHandler.NewAuthHandler,
		httpHandler.NewHealthHandler,
		httpHandler.NewInterviewHandler,
		httpHandler.NewUserHandler,

		// Service
		service.NewUserService,
		service.NewAuthService,
		service.NewInterviewService,
		service.NewQuestionService,
		service.NewReviewService,
		service.NewTranscriptManager,

		// Use case
		service.NewAIUseCase,

		// Repo
		repo.NewReviewRepo,
		repo.NewSettingRepo,
		repo.NewQuestionRepo,
		repo.NewSessionRepo,
		repo.NewUserRepo,
		repo.NewInterviewRepo,
		repo.NewTranscriptRepo,
		repo.NewFileRepo,
		repo.NewLLMRepo,
		repo.NewTTSRepo,
		repo.NewInMemoryCallbackQueueRepo,
		repo.NewIntentClassificationRepo,
		wire.NewSet(
			repo.NewMessageQueueRepo,
			wire.Bind(new(repo.MessageQueueProducerRepo), new(repo.MessageQueueRepo)),
			wire.Bind(new(repo.MessageQueueConsumerRepo), new(repo.MessageQueueRepo)),
		),

		// HTTP
		http.NewHTTPCLient,

		// Fasttext
		fasttext.NewFastTextPool,

		// Postgres
		postgres.NewPostgresDatabase,

		// Zerolog
		zerolog.NewZerologLogger,

		// Cloudflare
		cloudflare.NewCloudflareR2ObjectStorageClient,

		// Config
		config.LoadLLMConfig,
		config.LoadDatabaseConfig,
		config.LoadObjectStorageConfig,
		config.LoadTTSConfig,
		config.LoadWebsocketConfig,
		config.LoadInMemoryQueueConfig,
		config.LoadMessageQueueConfig,
		config.LoadIntentClassificationConfig,

		// Middleware
		middleware.NewMiddleware,

		// Housekeeping
		background.NewHouseKeeper,
		background.NewWorkerPool,

		// Application
		app.NewApplication,
	)
	return &app.Application{}, nil
}
