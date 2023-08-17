// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: accounts.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addToAccountBalanceById = `-- name: AddToAccountBalanceById :one
UPDATE accounts
SET balance = balance + $1
WHERE id = $2
RETURNING id, owner, balance, currency, created_at
`

type AddToAccountBalanceByIdParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddToAccountBalanceById(ctx context.Context, arg AddToAccountBalanceByIdParams) (Account, error) {
	row := q.db.QueryRow(ctx, addToAccountBalanceById, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING id, owner, balance, currency, created_at
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccountById = `-- name: DeleteAccountById :execrows
DELETE
FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccountById(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAccountById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, owner, balance, currency, created_at
FROM accounts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAccountById(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountByIdForUpdate = `-- name: GetAccountByIdForUpdate :one
SELECT id, owner, balance, currency, created_at
FROM accounts
WHERE id = $1
LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetAccountByIdForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountByIdForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const getAllAccounts = `-- name: GetAllAccounts :many
SELECT id, owner, balance, currency, created_at
FROM accounts
ORDER BY created_at
`

func (q *Queries) GetAllAccounts(ctx context.Context) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAllAccounts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllAccountsPaginated = `-- name: GetAllAccountsPaginated :many
SELECT id, owner, balance, currency, created_at
FROM accounts
ORDER BY created_at
OFFSET $1 LIMIT $2
`

type GetAllAccountsPaginatedParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) GetAllAccountsPaginated(ctx context.Context, arg GetAllAccountsPaginatedParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAllAccountsPaginated, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUserAccountsPaginated = `-- name: GetAllUserAccountsPaginated :many
SELECT id, owner, balance, currency, created_at
FROM accounts
WHERE owner = $1
ORDER BY created_at
OFFSET $2 LIMIT $3
`

type GetAllUserAccountsPaginatedParams struct {
	Owner  string `json:"owner"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

func (q *Queries) GetAllUserAccountsPaginated(ctx context.Context, arg GetAllUserAccountsPaginatedParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAllUserAccountsPaginated, arg.Owner, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTwoAccountsInvolvedInTransfer = `-- name: GetTwoAccountsInvolvedInTransfer :many
SELECT id, owner, balance, currency, created_at
FROM accounts
WHERE id = $1
   OR id = $2
ORDER BY id
LIMIT 2 FOR NO KEY UPDATE
`

type GetTwoAccountsInvolvedInTransferParams struct {
	FromAccount int64 `json:"from_account"`
	ToAccount   int64 `json:"to_account"`
}

func (q *Queries) GetTwoAccountsInvolvedInTransfer(ctx context.Context, arg GetTwoAccountsInvolvedInTransferParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, getTwoAccountsInvolvedInTransfer, arg.FromAccount, arg.ToAccount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserAccountById = `-- name: GetUserAccountById :one
SELECT id, owner, balance, currency, created_at
FROM accounts
WHERE owner = $1
  AND id = $2
`

type GetUserAccountByIdParams struct {
	Owner string `json:"owner"`
	ID    int64  `json:"id"`
}

func (q *Queries) GetUserAccountById(ctx context.Context, arg GetUserAccountByIdParams) (Account, error) {
	row := q.db.QueryRow(ctx, getUserAccountById, arg.Owner, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateAccountBalanceById = `-- name: UpdateAccountBalanceById :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING id, owner, balance, currency, created_at
`

type UpdateAccountBalanceByIdParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateAccountBalanceById(ctx context.Context, arg UpdateAccountBalanceByIdParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccountBalanceById, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}

const updateAccountById = `-- name: UpdateAccountById :one
UPDATE accounts
SET owner    = COALESCE($1::varchar, owner),
    balance  = COALESCE($2::bigint, balance),
    currency = COALESCE($3::varchar, currency)
WHERE id = $4::bigint
RETURNING id, owner, balance, currency, created_at
`

type UpdateAccountByIdParams struct {
	Owner    pgtype.Text `json:"owner"`
	Balance  pgtype.Int8 `json:"balance"`
	Currency pgtype.Text `json:"currency"`
	ID       int64       `json:"id"`
}

func (q *Queries) UpdateAccountById(ctx context.Context, arg UpdateAccountByIdParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccountById,
		arg.Owner,
		arg.Balance,
		arg.Currency,
		arg.ID,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)
	return i, err
}
