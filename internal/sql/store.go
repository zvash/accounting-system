package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, arg TransferTransactionParams) (TransferTransactionResult, error)
	CreateUserTransaction(ctx context.Context, arg CreateUserTransactionParams) (CreateUserTransactionResult, error)
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
