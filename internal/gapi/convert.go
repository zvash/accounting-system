package gapi

import (
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(dbUser sql.User) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		Name:              dbUser.Name,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreateAt:          timestamppb.New(dbUser.CreatedAt),
	}
}
