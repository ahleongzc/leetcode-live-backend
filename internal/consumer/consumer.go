package consumer

import "context"

type Consumer interface {
	ConsumeAndProcess(ctx context.Context)
}
