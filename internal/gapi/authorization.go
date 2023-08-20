package gapi

import (
	"context"
	"fmt"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}

func (server *Server) getAuthenticatedUser(ctx context.Context) (sql.User, error) {
	user := sql.User{}
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return user, err
	}
	user, err = server.db.GetUserByUserName(ctx, payload.Username)
	if err != nil {
		return user, nil
	}
	return user, nil
}
