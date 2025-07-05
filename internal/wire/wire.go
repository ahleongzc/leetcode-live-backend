//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/ahleongzc/leetcode-live-backend/cmd/app"
	"github.com/ahleongzc/leetcode-live-backend/internal/background"
	"github.com/ahleongzc/leetcode-live-backend/internal/config"
	"github.com/ahleongzc/leetcode-live-backend/internal/consumer"
	"github.com/ahleongzc/leetcode-live-backend/internal/handler"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/llm"
	messagequeue "github.com/ahleongzc/leetcode-live-backend/internal/infra/message_queue"
	"github.com/ahleongzc/leetcode-live-backend/internal/infra/tts"
	"github.com/ahleongzc/leetcode-live-backend/internal/middleware"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"

	"github.com/google/wire"
)

func InitializeApplication() (*app.Application, error) {
	wire.Build(
		// Consumer
		consumer.NewReviewConsumer,

		// Handler
		handler.NewAuthHandler,
		handler.NewHealthHandler,
		handler.NewInterviewHandler,
		handler.NewUserHandler,

		// Service
		service.NewUserService,
		service.NewAuthService,
		service.NewInterviewService,

		// Scenario
		scenario.NewAuthScenario,
		scenario.NewQuestionScenario,
		scenario.NewInterviewScenario,
		scenario.NewTranscriptManager,
		scenario.NewIntentClassifier,
		scenario.NewReviewScenario,

		// Repo
		repo.NewReviewRepo,
		repo.NewQuestionRepo,
		repo.NewSessionRepo,
		repo.NewUserRepo,
		repo.NewInterviewRepo,
		repo.NewTranscriptRepo,
		repo.NewFileRepo,

		// Infra
		tts.NewTTS,
		llm.NewLLM,
		infra.NewPostgresDatabase,
		infra.NewZerologLogger,
		infra.NewCloudflareR2ObjectStorageClient,
		infra.NewInMemoryCallbackQueue,
		infra.NewHTTPCLient,
		wire.NewSet(
			messagequeue.NewMessageQueue,
			wire.Bind(new(messagequeue.MessageQueueProducer), new(messagequeue.MessageQueue)),
			wire.Bind(new(messagequeue.MessageQueueConsumer), new(messagequeue.MessageQueue)),
		),

		// Config
		config.LoadLLMConfig,
		config.LoadDatabaseConfig,
		config.LoadObjectStorageConfig,
		config.LoadTTSConfig,
		config.LoadWebsocketConfig,
		config.LoadInMemoryQueueConfig,
		config.LoadMessageQueueConfig,

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
