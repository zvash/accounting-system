// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: transfers.sql

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRow(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransferById = `-- name: DeleteTransferById :execrows
DELETE
FROM transfers
WHERE id = $1
`

func (q *Queries) DeleteTransferById(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.Exec(ctx, deleteTransferById, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getAllTransfers = `-- name: GetAllTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
ORDER BY created_at
`

func (q *Queries) GetAllTransfers(ctx context.Context) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, getAllTransfers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
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

const getAllTransfersPaginated = `-- name: GetAllTransfersPaginated :many
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
ORDER BY created_at
OFFSET $1 LIMIT $2
`

type GetAllTransfersPaginatedParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) GetAllTransfersPaginated(ctx context.Context, arg GetAllTransfersPaginatedParams) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, getAllTransfersPaginated, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
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

const getTransferById = `-- name: GetTransferById :one
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTransferById(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRow(ctx, getTransferById, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const updateTransferById = `-- name: UpdateTransferById :one
UPDATE transfers
SET from_account_id = coalesce($2, from_account_id),
    to_account_id   = coalesce($3, to_account_id),
    amount          = coalesce($4, amount)
WHERE id = $1
RETURNING id, from_account_id, to_account_id, amount, created_at
`

type UpdateTransferByIdParams struct {
	ID            int64       `json:"id"`
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	Amount        pgtype.Int8 `json:"amount"`
}

func (q *Queries) UpdateTransferById(ctx context.Context, arg UpdateTransferByIdParams) (Transfer, error) {
	row := q.db.QueryRow(ctx, updateTransferById,
		arg.ID,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Amount,
	)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
