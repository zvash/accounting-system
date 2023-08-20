package worker

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zvash/accounting-system/internal/sql"
)

type RedisTaskProcessor struct {
	server *asynq.Server
	db     sql.Store
}

func NewRedisTaskProcessor(redisOptions asynq.RedisClientOpt, store sql.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOptions,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				fmt.Printf("error: %v, type: %v, payload: %v -> process task failed",
					err,
					task.Type(),
					task.Payload(),
				)
			}),
		},
	)
	return &RedisTaskProcessor{
		server: server,
		db:     store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
