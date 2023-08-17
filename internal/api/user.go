package api

import (
	"errors"
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

type loginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	userResponse
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

type userLoginRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

func (server *Server) loginUser(ctx *fiber.Ctx) error {
	req := userLoginRequest{}
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
	user, err := server.db.GetUserByUserName(ctx.Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrRecordNotFound) {
			return invalidCredentials()
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return invalidCredentials()
	}

	accessToken, _, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}
	refreshToken, _, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}

	loginResponse := loginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		userResponse: mapUserModelToUserResponse(&user),
	}
	if err := ctx.JSON(loginResponse); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}

func invalidCredentials() error {
	return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials.")
}
