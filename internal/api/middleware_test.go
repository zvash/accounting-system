package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/token"
	"net/http"
	"testing"
	"time"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	createdToken, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, createdToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, resp *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "UnsupportedTokenType",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderKey, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
		{
			name: "NoAuthHearer",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, resp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.router.Get(
				authPath,
				server.authMiddleware(),
				func(ctx *fiber.Ctx) error {
					return ctx.JSON(struct{}{})
				},
			)
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			resp, err := server.router.Test(request)
			require.NoError(t, err)
			tc.checkResponse(t, resp)
		})
	}
}
