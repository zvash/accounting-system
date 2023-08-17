package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *fiber.Ctx) error {
	var req renewAccessTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "there is an error in the type of provided variables!")
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	session, err := server.db.GetSession(ctx.Context(), refreshPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "error while handling the token renewal procedure")
	}

	if session.IsBlocked {
		return fiber.NewError(fiber.StatusUnauthorized, "blocked session")
	}

	if session.Username != refreshPayload.Username {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	if session.RefreshToken != req.RefreshToken {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	if time.Now().After(session.ExpiresAt) {
		return fiber.NewError(fiber.StatusUnauthorized, "expired session")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error while handling the token renewal procedure")
	}

	resp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	if err := ctx.JSON(resp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}
