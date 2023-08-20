package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/zvash/accounting-system/internal/sql"
	"log"
)

const TaskSendVerifyEmail = "task:send-verify-email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Printf("type: %v, payload: %v, queue: %v, max_retry: %v -> enqueued task.",
		info.Type,
		info.Payload,
		info.Queue,
		info.MaxRetry,
	)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal TaskSendVerifyEmail: %w", asynq.SkipRetry)
	}
	user, err := processor.db.GetUserByUserName(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrRecordNotFound) {
			return fmt.Errorf("user doesn't exists: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	//TODO: send verification email to user
	log.Printf("type: %v, payload: %v, email: %v -> processed task.",
		task.Type,
		task.Payload,
		user.Email,
	)
	return nil
}
