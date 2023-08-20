package gapi

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/zvash/accounting-system/internal/dto"
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protobufCreateUserToDTOCreateUser(pcu *pb.CreateUserRequest) dto.CreateUserRequest {
	return dto.CreateUserRequest{
		Username:             pcu.GetUsername(),
		Password:             pcu.GetPassword(),
		PasswordConfirmation: pcu.GetPasswordConfirmation(),
		Email:                pcu.GetEmail(),
		Name:                 pcu.GetEmail(),
	}
}

func errorResponseToErrorDetailsBadRequestFieldViolation(er val.ErrorResponse) *errdetails.BadRequest_FieldViolation {
	fieldName := strcase.ToSnake(er.FailedField)
	return &errdetails.BadRequest_FieldViolation{
		Field: fieldName,
		Description: fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			fieldName,
			er.Value,
			er.Tag,
		),
	}
}

func errorResponsesToErrorDetailsBadRequestFieldViolations(ers []val.ErrorResponse) (violations []*errdetails.BadRequest_FieldViolation) {
	for _, er := range ers {
		violations = append(violations, errorResponseToErrorDetailsBadRequestFieldViolation(er))
	}
	return
}

func errorResponsesToStatusErrors(errs []val.ErrorResponse) error {
	violations := errorResponsesToErrorDetailsBadRequestFieldViolations(errs)
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}

func protobufLoginRequestToDTOLoginRequest(lr *pb.LoginRequest) dto.UserLoginRequest {
	return dto.UserLoginRequest{
		Username: lr.GetUsername(),
		Password: lr.GetPassword(),
	}
}

func protobufUpdateUserToDTOUpdateUser(pcu *pb.UpdateUserRequest) dto.UpdateUserRequest {
	return dto.UpdateUserRequest{
		Username: pcu.GetUsername(),
		Email:    pcu.GetEmail(),
		Name:     pcu.GetName(),
	}
}

func protobufChangePasswordToDTOChangePassword(pcu *pb.ChangePasswordRequest) dto.ChangePasswordRequest {
	return dto.ChangePasswordRequest{
		Username:             pcu.GetUsername(),
		Password:             pcu.GetPassword(),
		PasswordConfirmation: pcu.GetPasswordConfirmation(),
	}
}

func dbUserToProtobufUser(dbUser sql.User) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		Name:              dbUser.Name,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreateAt:          timestamppb.New(dbUser.CreatedAt),
	}
}
