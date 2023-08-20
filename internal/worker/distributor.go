package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

// TaskDistributor pushes the tasks to the redis
type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}
