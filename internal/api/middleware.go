package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/token"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func getUsernameFromAuthPayload(ctx *fiber.Ctx) string {
	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	return authPayload.Username
}

func unauthorizedAccess() error {
	return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
}

func getHeader(ctx *fiber.Ctx, key string, defaultValue ...string) (string, error) {
	headers := ctx.GetReqHeaders()
	for k, v := range headers {
		if strings.ToLower(k) == strings.ToLower(key) {
			return v, nil
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}
	return "", errors.New("could not find the value for the provided key")
}

func (server *Server) authMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header, err := getHeader(ctx, authorizationHeaderKey)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "authorization header is not provided")
		}
		fields := strings.Fields(header)
		if len(fields) < 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization token type")
		}
		accessToken := fields[1]
		payload, err := server.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		ctx.Locals(authorizationPayloadKey, payload)
		return ctx.Next()
	}
}
