package api

import (
	"github.com/stretchr/testify/require"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/util"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store sql.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
