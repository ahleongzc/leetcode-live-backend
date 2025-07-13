package consumer

import (
	"context"
	"encoding/json"
	"runtime/debug"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	"github.com/ahleongzc/leetcode-live-backend/internal/domain/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/repo"
	"github.com/ahleongzc/leetcode-live-backend/internal/service"
	"github.com/rs/zerolog"
)

func NewReviewConsumer(
	reviewService service.ReviewService,
	consumerRepo repo.MessageQueueConsumerRepo,
	logger *zerolog.Logger,
) *ReviewConsumer {
	return &ReviewConsumer{
		reviewService: reviewService,
		consumerRepo:  consumerRepo,
		logger:        logger,
	}
}

type ReviewConsumer struct {
	reviewService service.ReviewService
	consumerRepo  repo.MessageQueueConsumerRepo
	logger        *zerolog.Logger
}

func (r *ReviewConsumer) ConsumeAndProcess(ctx context.Context, workerCount uint) {
	deliveryChan, err := r.consumerRepo.StartConsuming(ctx, common.REVIEW_QUEUE)
	if err != nil {
		panic(err)
	}

	for i := range workerCount {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					stackTrace := debug.Stack()
					r.logger.Error().
						Interface("panic", err).
						Bytes("stackTrace", stackTrace).
						Uint("workerNumber", i).
						Msg("panic recovered in review consumer")
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case delivery := <-deliveryChan:
					r.processMessage(ctx, delivery)
				}
			}
		}()
	}
}

func (r *ReviewConsumer) processMessage(ctx context.Context, delivery *model.Delivery) {
	defer func() {
		if err := recover(); err != nil {
			stackTrace := debug.Stack()
			r.logger.Error().
				Interface("panic", err).
				Bytes("stackTrace", stackTrace).
				Msg("panic recovered when consuming messages")

			if err := delivery.Nack(true); err != nil {
				r.logger.Error().Err(err).Msg("failed to nack message after panic")
			}
		}
	}()

	reviewMessage := &model.ReviewMessage{}
	if err := json.Unmarshal(delivery.Body, reviewMessage); err != nil {
		r.logger.Error().Err(err).Msg("unable to marshal review message")
		delivery.Nack(true)
		return
	}

	if err := r.reviewService.ReviewInterviewPerformance(ctx, reviewMessage.InterviewID); err != nil {
		r.logger.Error().Err(err).Msg("unable to review interview performance")
		delivery.Nack(true)
		return
	}

	delivery.Ack()
}
