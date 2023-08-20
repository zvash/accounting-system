package sql

import (
	"context"
)

type CreateUserTransactionParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTransactionResult struct {
	User User
}

func (store *DBStore) CreateUserTransaction(ctx context.Context, arg CreateUserTransactionParams) (CreateUserTransactionResult, error) {
	var result CreateUserTransactionResult
	err := store.executeTransaction(ctx, func(queries *Queries) error {
		var err error
		result.User, err = queries.CreateUser(ctx, CreateUserParams{
			Username: arg.Username,
			Name:     arg.Name,
			Email:    arg.Email,
			Password: arg.Password,
		})
		if err != nil {
			return err
		}

		if err := arg.AfterCreate(result.User); err != nil {
			return err
		}

		return nil
	})

	return result, err
}
