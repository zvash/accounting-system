package gapi

import (
	"context"
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	userParams := sql.CreateUserParams{
		Username: req.GetUsername(),
		Password: hashedPassword,
		Name:     req.GetName(),
		Email:    req.GetEmail(),
	}
	user, err := server.db.CreateUser(ctx, userParams)
	if err != nil {
		if sql.ErrorCode(err) == sql.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, "this user is already exists!")
		}
		return nil, status.Errorf(codes.Internal, "error while trying to create the new user.")
	}
	resp := &pb.CreateUserResponse{
		User: dbUserToProtobufUser(user),
	}
	return resp, nil
}
