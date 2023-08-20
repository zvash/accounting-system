package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"github.com/zvash/accounting-system/internal/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	dto := protobufCreateUserToDTOCreateUser(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	createUserTransactionParams := sql.CreateUserTransactionParams{
		CreateUserParams: sql.CreateUserParams{
			Username: req.GetUsername(),
			Password: hashedPassword,
			Name:     req.GetName(),
			Email:    req.GetEmail(),
		},
		AfterCreate: func(user sql.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}
	transactionResult, err := server.db.CreateUserTransaction(ctx, createUserTransactionParams)
	if err != nil {
		if sql.ErrorCode(err) == sql.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "this user is already exists!")
		}
		return nil, status.Errorf(codes.Internal, "error while trying to create the new user.")
	}

	resp := &pb.CreateUserResponse{
		User: dbUserToProtobufUser(transactionResult.User),
	}
	return resp, nil
}
