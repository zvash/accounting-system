package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
)

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Server struct {
	db        sql.Store
	router    *fiber.App
	validator *XValidator
}

func NewServer(db sql.Store) *Server {
	requestValidator := NewValidator()
	server := &Server{db: db}
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
	return server
}

func (server *Server) Start(address string) error {
	err := server.router.Listen(address)
	if err != nil {
		return err
	}
	return nil
}
