package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"strings"
)

type createAccountRequest struct {
	Owner    string `validate:"required"`
	Currency string `validate:"required,oneof=EUR CAD USD"`
}

type getAccountRequest struct {
	ID int64 `validate:"required,min=1"`
}

type getListAccountsRequest struct {
	Page    int32 `validate:"min=1"`
	PerPage int32 `query:"per_page" validate:"min=1"`
}

func (server *Server) createAccount(ctx *fiber.Ctx) error {
	req := createAccountRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "there is an error in the type of provided variables!",
		}
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	account, err := server.db.CreateAccount(ctx.Context(), sql.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error!")
	}
	err = ctx.JSON(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}

func (server *Server) getAccount(ctx *fiber.Ctx) error {
	req := getAccountRequest{}
	err := ctx.ParamsParser(&req)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "there is an error in the type of provided variables!",
		}
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	account, err := server.db.GetAccountById(ctx.Context(), req.ID)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "Record was not found.",
		}
	}
	err = ctx.JSON(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}

func (server *Server) listAccounts(ctx *fiber.Ctx) error {
	req := getListAccountsRequest{
		Page:    1,
		PerPage: 10,
	}
	err := ctx.QueryParser(&req)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "there is an error in the type of provided variables!",
		}
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	if req.PerPage > 10 {
		req.PerPage = 10
	}
	accounts, err := server.db.GetAllAccountsPaginated(ctx.Context(), sql.GetAllAccountsPaginatedParams{
		Offset: (req.Page - 1) * req.PerPage,
		Limit:  req.PerPage,
	})
	err = ctx.JSON(accounts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}
