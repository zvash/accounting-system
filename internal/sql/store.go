package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error)
}

type DBStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &DBStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (store *DBStore) executeTransaction(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTransactionParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *DBStore) TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error) {
	var result TransferTransactionResult
	err := store.executeTransaction(ctx, func(queries *Queries) error {
		var err error
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		accountsToUpdate, err := queries.GetTwoAccountsInvolvedInTransfer(ctx, GetTwoAccountsInvolvedInTransferParams{
			FromAccount: arg.FromAccountID,
			ToAccount:   arg.ToAccountID,
		})
		if err != nil {
			return err
		}
		if len(accountsToUpdate) != 2 {
			return errors.New("exactly two distinct accounts are required for a successful transfer")
		}
		accountIdToAccountMap := make(map[int64]Account)
		for _, account := range accountsToUpdate {
			accountIdToAccountMap[account.ID] = account
		}
		result.FromAccount, err = queries.UpdateAccountById(ctx, UpdateAccountByIdParams{
			ID:      arg.FromAccountID,
			Balance: pgtype.Int8{Int64: accountIdToAccountMap[arg.FromAccountID].Balance - arg.Amount, Valid: true},
		})
		result.ToAccount, err = queries.UpdateAccountById(ctx, UpdateAccountByIdParams{
			ID:      arg.ToAccountID,
			Balance: pgtype.Int8{Int64: accountIdToAccountMap[arg.ToAccountID].Balance + arg.Amount, Valid: true},
		})
		return nil
	})

	return result, err
}
