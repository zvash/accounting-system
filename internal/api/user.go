package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"strings"
)

type createUserRequest struct {
	Username             string `json:"username" validate:"required,alphanum"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=6,eqfield=Password"`
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func mapUserModelToUserResponse(dbUser *sql.User) userResponse {
	return userResponse{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Name:     dbUser.Name,
	}
}

func (server *Server) createUser(ctx *fiber.Ctx) error {
	req := createUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "there is an error in the type of provided variables!")
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	password, err := util.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the user!")
	}
	user, err := server.db.CreateUser(ctx.Context(), sql.CreateUserParams{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	})
	if err != nil {
		if sql.ErrorCode(err) == sql.UniqueViolation {
			return fiber.NewError(fiber.StatusBadRequest, "this user is already exists!")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the user!")
	}
	userResponse := mapUserModelToUserResponse(&user)
	if err := ctx.JSON(userResponse); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}
