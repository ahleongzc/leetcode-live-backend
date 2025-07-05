package consumer

import (
	"context"
	"encoding/json"
	"runtime/debug"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
	messagequeue "github.com/ahleongzc/leetcode-live-backend/internal/infra/message_queue"
	"github.com/ahleongzc/leetcode-live-backend/internal/model"
	"github.com/ahleongzc/leetcode-live-backend/internal/scenario"

	"github.com/rs/zerolog"
)

func NewReviewConsumer(
	reviewScenario scenario.ReviewScenario,
	consumer messagequeue.MessageQueueConsumer,
	logger *zerolog.Logger,
) *ReviewConsumer {
	return &ReviewConsumer{
		reviewScenario: reviewScenario,
		consumer:       consumer,
		logger:         logger,
	}
}

type ReviewConsumer struct {
	reviewScenario scenario.ReviewScenario
	consumer       messagequeue.MessageQueueConsumer
	logger         *zerolog.Logger
}

func (r *ReviewConsumer) ConsumeAndProcess(ctx context.Context, workerCount uint) {
	deliveryChan, err := r.consumer.StartConsuming(ctx, common.REVIEW_QUEUE)
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
					defer func() {
						if err := recover(); err != nil {
							stackTrace := debug.Stack()
							r.logger.Error().
								Interface("panic", err).
								Bytes("stackTrace", stackTrace).
								Uint("workerNumber", i).
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
						continue
					}

					if err := r.reviewScenario.ReviewInterviewPerformance(ctx, reviewMessage.InterviewID); err != nil {
						r.logger.Error().Err(err).Msg("unable to review interview performance")
						delivery.Nack(true)
						continue
					}

					delivery.Ack()
				}
			}
		}()
	}
}
