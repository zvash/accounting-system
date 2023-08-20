package gapi

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	dto := protobufUpdateUserToDTOUpdateUser(req)
	if errs := server.validator.Validate(dto); errs != nil {
		return nil, errorResponsesToStatusErrors(errs)
	}
	user, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized.")
	}
	if user.Username != req.GetUsername() {
		return nil, status.Error(codes.Unauthenticated, "unauthorized.")
	}
	updatedUser, err := server.db.UpdateUserByUsername(ctx, sql.UpdateUserByUsernameParams{
		Username: req.GetUsername(),
		Name:     pgtype.Text{Valid: req.Name != nil, String: req.GetName()},
		Email:    pgtype.Text{Valid: req.Email != nil, String: req.GetEmail()},
	})
	if err != nil {
		if sql.ErrorCode(err) == sql.UniqueViolation {
			return nil, status.Error(codes.InvalidArgument, "cannot change the email to the provided value.")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	resp := &pb.UpdateUserResponse{
		User: dbUserToProtobufUser(updatedUser),
	}
	return resp, nil
}
