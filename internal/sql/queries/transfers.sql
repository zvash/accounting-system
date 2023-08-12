-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransferById :one
SELECT *
FROM transfers
WHERE id = $1
LIMIT 1;

-- name: GetAllTransfers :many
SELECT *
FROM transfers
ORDER BY created_at;

-- name: GetAllTransfersPaginated :many
SELECT *
FROM transfers
ORDER BY created_at
OFFSET $1 LIMIT $2;

-- name: UpdateTransferById :one
UPDATE transfers
SET from_account_id = coalesce(sqlc.narg('from_account_id'), from_account_id),
    to_account_id   = coalesce(sqlc.narg('to_account_id'), to_account_id),
    amount          = coalesce(sqlc.narg('amount'), amount)
WHERE id = $1
RETURNING *;

-- name: DeleteTransferById :execrows
DELETE
FROM transfers
WHERE id = $1;