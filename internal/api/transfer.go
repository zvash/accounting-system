package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/zvash/accounting-system/internal/sql"
	"strings"
	"time"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,currency"`
}

type destinationAccount struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
}

type transferResponse struct {
	FromAccount sql.Account        `json:"from_account"`
	ToAccount   destinationAccount `json:"to_account"`
	Amount      int64              `json:"amount"`
	Date        time.Time          `json:"date"`
}

func (server *Server) createTransfer(ctx *fiber.Ctx) error {
	req := transferRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not parse the input data")
	}

	if errs := server.validator.Validate(req); errs != nil {
		errorsBag := server.validator.makeErrorBag(errs)
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: strings.Join(errorsBag, " and "),
		}
	}

	owner := getUsernameFromAuthPayload(ctx)
	fromAccount, err := server.checkIfAccountCurrencyIsAMatch(ctx, req.FromAccountID, req.Currency)
	if fromAccount.Owner != owner {
		return unauthorizedAccess()
	}
	if err != nil {
		return err
	}
	if _, err := server.checkIfAccountCurrencyIsAMatch(ctx, req.ToAccountID, req.Currency); err != nil {
		return err
	}

	result, err := server.db.TransferTransaction(ctx.Context(), sql.TransferTransactionParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		if errors.Is(err, sql.InsufficientFundsError{}) {
			return fiber.NewError(fiber.StatusBadRequest, "Insufficient funds")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Fund transfer was unsuccessful!")
	}
	response := transferResponse{
		FromAccount: result.FromAccount,
		ToAccount: destinationAccount{
			ID:    result.ToAccount.ID,
			Owner: result.ToAccount.Owner,
		},
		Amount: result.Transfer.Amount,
		Date:   result.Transfer.CreatedAt,
	}
	err = ctx.JSON(response)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error while creating the response!")
	}

	return nil
}

func (server *Server) checkIfAccountCurrencyIsAMatch(ctx *fiber.Ctx, accountId int64, currency string) (sql.Account, error) {
	account, err := server.db.GetAccountById(ctx.Context(), accountId)
	if err != nil {
		if errors.Is(err, sql.ErrRecordNotFound) {
			return account, fiber.NewError(fiber.StatusUnprocessableEntity, "Could not find the account.")
		}
		return account, fiber.NewError(fiber.StatusInternalServerError, "Error while trying to access the account.")
	}
	if account.Currency != currency {
		return account, fiber.NewError(fiber.StatusUnprocessableEntity, "Currency mismatched.")
	}
	return account, nil
}
