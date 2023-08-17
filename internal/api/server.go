package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"github.com/zvash/accounting-system/internal/token"
	"github.com/zvash/accounting-system/internal/util"
)

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Server struct {
	config     util.Config
	db         sql.Store
	router     *fiber.App
	validator  *XValidator
	tokenMaker token.Maker
}

func NewServer(config util.Config, db sql.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	requestValidator := NewValidator()
	server := &Server{
		config:     config,
		db:         db,
		tokenMaker: tokenMaker,
	}
	router := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	router.Post("/accounts", server.createAccount)
	router.Get("/accounts/:id", server.getAccount)
	router.Get("/accounts", server.listAccounts)

	router.Post("/transfers", server.createTransfer)

	router.Post("/users", server.createUser)

	server.router = router
	server.validator = requestValidator
	return server, nil
}

func (server *Server) Start(address string) error {
	err := server.router.Listen(address)
	if err != nil {
		return err
	}
	return nil
}
