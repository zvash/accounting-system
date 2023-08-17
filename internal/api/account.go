package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"strings"
)

type createAccountRequest struct {
	Currency string `validate:"required,currency"`
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
			Code:    fiber.StatusBadRequest,
			Message: "there is an error in the type of provided variables!",
		}
	}
	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	account, err := server.db.CreateAccount(ctx.Context(), sql.CreateAccountParams{
		Owner:    getUsernameFromAuthPayload(ctx),
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		errCode := sql.ErrorCode(err)
		if errCode == sql.ForeignKeyViolation || errCode == sql.UniqueViolation {
			return fiber.NewError(fiber.StatusForbidden, "Creating a new account with provided data is forbidden!")
		}
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
			Code:    fiber.StatusUnprocessableEntity,
			Message: strings.Join(errorsBag, " and "),
		}
	}
	owner := getUsernameFromAuthPayload(ctx)
	account, err := server.db.GetUserAccountById(ctx.Context(), sql.GetUserAccountByIdParams{
		ID:    req.ID,
		Owner: owner,
	})
	if err != nil {
		if errors.Is(err, sql.ErrRecordNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "Record was not found.",
			}
		}
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error!",
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
	owner := getUsernameFromAuthPayload(ctx)
	accounts, err := server.db.GetAllUserAccountsPaginated(ctx.Context(), sql.GetAllUserAccountsPaginatedParams{
		Offset: (req.Page - 1) * req.PerPage,
		Limit:  req.PerPage,
		Owner:  owner,
	})
	err = ctx.JSON(accounts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}
	return nil
}
