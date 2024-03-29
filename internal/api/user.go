package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zvash/accounting-system/internal/dto"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"strings"
	"time"
)

type userResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func mapUserModelToUserResponse(dbUser *sql.User) userResponse {
	return userResponse{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Name:     dbUser.Name,
	}
}

func (server *Server) createUser(ctx *fiber.Ctx) error {
	req := dto.CreateUserRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "there is an error in the type of provided variables!")
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.MakeErrorBag(errs)
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

func (server *Server) loginUser(ctx *fiber.Ctx) error {
	req := dto.UserLoginRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "there is an error in the type of provided variables!")
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.MakeErrorBag(errs)
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

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}

	userAgent, _ := getHeader(ctx, "User-Agent", "")
	ipAddress := ctx.IP()
	session, err := server.db.CreateSession(ctx.Context(), sql.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     ipAddress,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	loginResponse := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  mapUserModelToUserResponse(&user),
	}
	if err := ctx.JSON(loginResponse); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}

func invalidCredentials() error {
	return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials.")
}
