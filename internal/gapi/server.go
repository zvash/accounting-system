package gapi

import (
	"fmt"
	"github.com/zvash/accounting-system/internal/pb"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/token"
	"github.com/zvash/accounting-system/internal/util"
	"github.com/zvash/accounting-system/internal/val"
	"github.com/zvash/accounting-system/internal/worker"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedAccountingSystemServer
	config          util.Config
	db              sql.Store
	validator       *val.XValidator
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store sql.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		db:              store,
		validator:       val.NewValidator(),
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
